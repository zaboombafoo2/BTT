package corehttp

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	version "github.com/bittorrent/go-btfs"
	core "github.com/bittorrent/go-btfs/core"
	"github.com/bittorrent/go-btfs/core/coreapi"
	namesys "github.com/bittorrent/go-btfs/namesys"
	repo "github.com/bittorrent/go-btfs/repo"

	config "github.com/TRON-US/go-btfs-config"
	files "github.com/TRON-US/go-btfs-files"
	iface "github.com/TRON-US/interface-go-btfs-core"
	nsopts "github.com/TRON-US/interface-go-btfs-core/options/namesys"
	ipath "github.com/TRON-US/interface-go-btfs-core/path"
	datastore "github.com/ipfs/go-datastore"
	syncds "github.com/ipfs/go-datastore/sync"
	path "github.com/ipfs/go-path"
	ci "github.com/libp2p/go-libp2p/core/crypto"
	id "github.com/libp2p/go-libp2p/p2p/protocol/identify"
)

// `btfs object new unixfs-dir`
var emptyDir = "/btfs/QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn"

type mockNamesys map[string]path.Path

func (m mockNamesys) Resolve(ctx context.Context, name string, opts ...nsopts.ResolveOpt) (value path.Path, err error) {
	cfg := nsopts.DefaultResolveOpts()
	for _, o := range opts {
		o(&cfg)
	}
	depth := cfg.Depth
	if depth == nsopts.UnlimitedDepth {
		// max uint
		depth = ^uint(0)
	}
	for strings.HasPrefix(name, "/btns/") {
		if depth == 0 {
			return value, namesys.ErrResolveRecursion
		}
		depth--

		var ok bool
		value, ok = m[name]
		if !ok {
			return "", namesys.ErrResolveFailed
		}
		name = value.String()
	}
	return value, nil
}

func (m mockNamesys) ResolveAsync(ctx context.Context, name string, opts ...nsopts.ResolveOpt) <-chan namesys.Result {
	out := make(chan namesys.Result, 1)
	v, err := m.Resolve(ctx, name, opts...)
	out <- namesys.Result{Path: v, Err: err}
	close(out)
	return out
}

func (m mockNamesys) Publish(ctx context.Context, name ci.PrivKey, value path.Path) error {
	return errors.New("not implemented for mockNamesys")
}

func (m mockNamesys) PublishWithEOL(ctx context.Context, name ci.PrivKey, value path.Path, _ time.Time) error {
	return errors.New("not implemented for mockNamesys")
}

func (m mockNamesys) GetResolver(subs string) (namesys.Resolver, bool) {
	return nil, false
}

func newNodeWithMockNamesys(ns mockNamesys) (*core.IpfsNode, error) {
	c := config.Config{
		Identity: config.Identity{
			PeerID: "QmTFauExutTsy4XP6JbMFcw2Wa9645HJt2bTqL6qYDCKfe", // required by offline node
		},
	}
	r := &repo.Mock{
		C: c,
		D: syncds.MutexWrap(datastore.NewMapDatastore()),
	}
	n, err := core.NewNode(context.Background(), &core.BuildCfg{Repo: r})
	if err != nil {
		return nil, err
	}
	n.Namesys = ns
	return n, nil
}

type delegatedHandler struct {
	http.Handler
}

func (dh *delegatedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	dh.Handler.ServeHTTP(w, r)
}

func doWithoutRedirect(req *http.Request) (*http.Response, error) {
	tag := "without-redirect"
	c := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return errors.New(tag)
		},
	}
	res, err := c.Do(req)
	if err != nil && !strings.Contains(err.Error(), tag) {
		return nil, err
	}
	return res, nil
}

func newTestServerAndNode(t *testing.T, ns mockNamesys) (*httptest.Server, iface.CoreAPI, context.Context) {
	n, err := newNodeWithMockNamesys(ns)
	if err != nil {
		t.Fatal(err)
	}

	cfg, err := n.Repo.Config()
	if err != nil {
		t.Fatal(err)
	}
	cfg.Gateway.PathPrefixes = []string{"/good-prefix"}

	// need this variable here since we need to construct handler with
	// listener, and server with handler. yay cycles.
	dh := &delegatedHandler{}
	ts := httptest.NewServer(dh)
	t.Cleanup(func() { ts.Close() })

	dh.Handler, err = makeHandler(n,
		ts.Listener,
		HostnameOption(),
		GatewayOption(false, "/btfs", "/btns"),
		VersionOption(),
	)
	if err != nil {
		t.Fatal(err)
	}

	api, err := coreapi.NewCoreAPI(n)
	if err != nil {
		t.Fatal(err)
	}

	return ts, api, n.Context()
}

func matchPathOrBreadcrumbs(s string, expected string) bool {
	matched, _ := regexp.MatchString("Index of "+regexp.QuoteMeta(expected), s)
	return matched
}

func TestGatewayGet(t *testing.T) {
	ns := mockNamesys{}
	ts, api, ctx := newTestServerAndNode(t, ns)

	k, err := api.Unixfs().Add(ctx, files.NewBytesFile([]byte("fnord")))
	if err != nil {
		t.Fatal(err)
	}
	ns["/btns/example.com"] = path.FromString(k.String())
	ns["/btns/working.example.com"] = path.FromString(k.String())
	ns["/btns/double.example.com"] = path.FromString("/btns/working.example.com")
	ns["/btns/triple.example.com"] = path.FromString("/btns/double.example.com")
	ns["/btns/broken.example.com"] = path.FromString("/btns/" + k.Cid().String())
	// We picked .man because:
	// 1. It's a valid TLD.
	// 2. Go treats it as the file extension for "man" files (even though
	//    nobody actually *uses* this extension, AFAIK).
	//
	// Unfortunately, this may not work on all platforms as file type
	// detection is platform dependent.
	ns["/btns/example.man"] = path.FromString(k.String())

	t.Log(ts.URL)
	for i, test := range []struct {
		host   string
		path   string
		status int
		text   string
	}{
		{"127.0.0.1:8080", "/", http.StatusNotFound, "404 page not found\n"},
		{"127.0.0.1:8080", "/" + k.Cid().String(), http.StatusNotFound, "404 page not found\n"},
		{"127.0.0.1:8080", k.String(), http.StatusOK, "fnord"},
		{"127.0.0.1:8080", "/btns/nxdomain.example.com", http.StatusNotFound, "btfs resolve -r /btns/nxdomain.example.com: " + namesys.ErrResolveFailed.Error() + "\n"},
		{"127.0.0.1:8080", "/btns/%0D%0A%0D%0Ahello", http.StatusNotFound, "btfs resolve -r /btns/%0D%0A%0D%0Ahello: " + namesys.ErrResolveFailed.Error() + "\n"},
		{"127.0.0.1:8080", "/btns/example.com", http.StatusOK, "fnord"},
		{"example.com", "/", http.StatusOK, "fnord"},

		{"working.example.com", "/", http.StatusOK, "fnord"},
		{"double.example.com", "/", http.StatusOK, "fnord"},
		{"triple.example.com", "/", http.StatusOK, "fnord"},
		{"working.example.com", k.String(), http.StatusNotFound, "btfs resolve -r /btns/working.example.com" + k.String() + ": no link named \"btfs\" under " + k.Cid().String() + "\n"},
		{"broken.example.com", "/", http.StatusNotFound, "btfs resolve -r /btns/broken.example.com/: " + namesys.ErrResolveFailed.Error() + "\n"},
		{"broken.example.com", k.String(), http.StatusNotFound, "btfs resolve -r /btns/broken.example.com" + k.String() + ": " + namesys.ErrResolveFailed.Error() + "\n"},
		// This test case ensures we don't treat the TLD as a file extension.
		{"example.man", "/", http.StatusOK, "fnord"},
	} {
		var c http.Client
		r, err := http.NewRequest(http.MethodGet, ts.URL+test.path, nil)
		if err != nil {
			t.Fatal(err)
		}
		r.Host = test.host
		resp, err := c.Do(r)

		urlstr := "http://" + test.host + test.path
		if err != nil {
			t.Errorf("error requesting %s: %s", urlstr, err)
			continue
		}
		defer resp.Body.Close()
		contentType := resp.Header.Get("Content-Type")
		if contentType != "text/plain; charset=utf-8" {
			t.Errorf("expected content type to be text/plain, got %s", contentType)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if resp.StatusCode != test.status {
			t.Errorf("(%d) got %d, expected %d from %s", i, resp.StatusCode, test.status, urlstr)
			t.Errorf("Body: %s", body)
			continue
		}
		if err != nil {
			t.Fatalf("error reading response from %s: %s", urlstr, err)
		}
		if string(body) != test.text {
			t.Errorf("unexpected response body from %s: expected %q; got %q", urlstr, test.text, body)
			continue
		}
	}
}

func TestPretty404(t *testing.T) {
	ns := mockNamesys{}
	ts, api, ctx := newTestServerAndNode(t, ns)

	f1 := files.NewMapDirectory(map[string]files.Node{
		"ipfs-404.html": files.NewBytesFile([]byte("Custom 404")),
		"deeper": files.NewMapDirectory(map[string]files.Node{
			"ipfs-404.html": files.NewBytesFile([]byte("Deep custom 404")),
		}),
	})

	k, err := api.Unixfs().Add(ctx, f1)
	if err != nil {
		t.Fatal(err)
	}

	host := "example.net"
	ns["/btns/"+host] = path.FromString(k.String())

	for _, test := range []struct {
		path   string
		accept string
		status int
		text   string
	}{
		{"/ipfs-404.html", "text/html", http.StatusOK, "Custom 404"},
		{"/nope", "text/html", http.StatusNotFound, "Custom 404"},
		{"/nope", "text/*", http.StatusNotFound, "Custom 404"},
		{"/nope", "*/*", http.StatusNotFound, "Custom 404"},
		{"/nope", "application/json", http.StatusNotFound, "btfs resolve -r /btns/example.net/nope: no link named \"nope\" under QmcmnF7XG5G34RdqYErYDwCKNFQ6jb8oKVR21WAJgubiaj\n"},
		{"/deeper/nope", "text/html", http.StatusNotFound, "Deep custom 404"},
		{"/deeper/", "text/html", http.StatusOK, ""},
		{"/deeper", "text/html", http.StatusOK, ""},
		{"/nope/nope", "text/html", http.StatusNotFound, "Custom 404"},
	} {
		var c http.Client
		req, err := http.NewRequest("GET", ts.URL+test.path, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Add("Accept", test.accept)
		req.Host = host
		resp, err := c.Do(req)

		if err != nil {
			t.Fatalf("error requesting %s: %s", test.path, err)
		}

		defer resp.Body.Close()
		if resp.StatusCode != test.status {
			t.Fatalf("got %d, expected %d, from %s", resp.StatusCode, test.status, test.path)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("error reading response from %s: %s", test.path, err)
		}

		if test.text != "" && string(body) != test.text {
			t.Fatalf("unexpected response body from %s: got %q, expected %q", test.path, body, test.text)
		}
	}
}

func TestIPNSHostnameRedirect(t *testing.T) {
	ns := mockNamesys{}
	ts, api, ctx := newTestServerAndNode(t, ns)
	t.Logf("test server url: %s", ts.URL)

	// create /btns/example.net/foo/index.html

	f1 := files.NewMapDirectory(map[string]files.Node{
		"_": files.NewBytesFile([]byte("_")),
		"foo": files.NewMapDirectory(map[string]files.Node{
			"index.html": files.NewBytesFile([]byte("_")),
		}),
	})

	k, err := api.Unixfs().Add(ctx, f1)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("k: %s\n", k)
	ns["/btns/example.net"] = path.FromString(k.String())

	// make request to directory containing index.html
	req, err := http.NewRequest(http.MethodGet, ts.URL+"/foo", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Host = "example.net"

	res, err := doWithoutRedirect(req)
	if err != nil {
		t.Fatal(err)
	}

	// expect 302 redirect to same path, but with trailing slash
	if res.StatusCode != 302 {
		t.Errorf("status is %d, expected 302", res.StatusCode)
	}
	hdr := res.Header["Location"]
	if len(hdr) < 1 {
		t.Errorf("location header not present")
	} else if hdr[0] != "/foo/" {
		t.Errorf("location header is %v, expected /foo/", hdr[0])
	}

	// make request with prefix to directory containing index.html
	req, err = http.NewRequest(http.MethodGet, ts.URL+"/foo", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Host = "example.net"
	req.Header.Set("X-Btfs-Gateway-Prefix", "/good-prefix")

	res, err = doWithoutRedirect(req)
	if err != nil {
		t.Fatal(err)
	}

	// expect 302 redirect to same path, but with prefix and trailing slash
	if res.StatusCode != 302 {
		t.Errorf("status is %d, expected 302", res.StatusCode)
	}
	hdr = res.Header["Location"]
	if len(hdr) < 1 {
		t.Errorf("location header not present")
	} else if hdr[0] != "/good-prefix/foo/" {
		t.Errorf("location header is %v, expected /good-prefix/foo/", hdr[0])
	}

	// make sure /version isn't exposed
	req, err = http.NewRequest(http.MethodGet, ts.URL+"/version", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Host = "example.net"
	req.Header.Set("X-Btfs-Gateway-Prefix", "/good-prefix")

	res, err = doWithoutRedirect(req)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != 404 {
		t.Fatalf("expected a 404 error, got: %s", res.Status)
	}
}

func TestIPNSHostnameBacklinks(t *testing.T) {
	ns := mockNamesys{}
	ts, api, ctx := newTestServerAndNode(t, ns)
	t.Logf("test server url: %s", ts.URL)

	f1 := files.NewMapDirectory(map[string]files.Node{
		"file.txt": files.NewBytesFile([]byte("1")),
		"foo? #<'": files.NewMapDirectory(map[string]files.Node{
			"file.txt": files.NewBytesFile([]byte("2")),
			"bar": files.NewMapDirectory(map[string]files.Node{
				"file.txt": files.NewBytesFile([]byte("3")),
			}),
		}),
	})

	// create /btns/example.net/foo/
	k, err := api.Unixfs().Add(ctx, f1)
	if err != nil {
		t.Fatal(err)
	}

	k2, err := api.ResolvePath(ctx, ipath.Join(k, "foo? #<'"))
	if err != nil {
		t.Fatal(err)
	}

	k3, err := api.ResolvePath(ctx, ipath.Join(k, "foo? #<'/bar"))
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("k: %s\n", k)
	ns["/btns/example.net"] = path.FromString(k.String())

	// make request to directory listing
	req, err := http.NewRequest(http.MethodGet, ts.URL+"/foo%3F%20%23%3C%27/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Host = "example.net"

	res, err := doWithoutRedirect(req)
	if err != nil {
		t.Fatal(err)
	}

	// expect correct backlinks
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("error reading response: %s", err)
	}
	s := string(body)
	t.Logf("body: %s\n", string(body))

	//FIXME
	//if !matchPathOrBreadcrumbs(s, "/btns/example.net/foo? #&lt;&#39;/bar/") {
	//	t.Fatalf("expected a path in directory listing")
	//}
	if !strings.Contains(s, "<a href=\"/foo%3F%20%23%3C%27/./..\">") {
		t.Fatalf("expected backlink in directory listing")
	}
	if !strings.Contains(s, "<a href=\"/foo%3F%20%23%3C%27/file.txt\">") {
		t.Fatalf("expected file in directory listing")
	}
	if !strings.Contains(s, k2.Cid().String()) {
		t.Fatalf("expected hash in directory listing")
	}

	// make request to directory listing at root
	req, err = http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Host = "example.net"

	res, err = doWithoutRedirect(req)
	if err != nil {
		t.Fatal(err)
	}

	// expect correct backlinks at root
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("error reading response: %s", err)
	}
	s = string(body)
	t.Logf("body: %s\n", string(body))

	//TODO
	//if !matchPathOrBreadcrumbs(s, "/") {
	//	t.Fatalf("expected a path in directory listing")
	//}
	if !strings.Contains(s, "<a href=\"/\">") {
		t.Fatalf("expected backlink in directory listing")
	}
	if !strings.Contains(s, "<a href=\"/file.txt\">") {
		t.Fatalf("expected file in directory listing")
	}
	if !strings.Contains(s, k.Cid().String()) {
		t.Fatalf("expected hash in directory listing")
	}

	// make request to directory listing
	req, err = http.NewRequest(http.MethodGet, ts.URL+"/foo%3F%20%23%3C%27/bar/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Host = "example.net"

	res, err = doWithoutRedirect(req)
	if err != nil {
		t.Fatal(err)
	}

	// expect correct backlinks
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("error reading response: %s", err)
	}
	s = string(body)
	t.Logf("body: %s\n", string(body))

	//if !matchPathOrBreadcrumbs(s, "/btns/<a href=\"/btns/example.net\">example.net</a>/<a href=\"/btns/example.net/foo%3F%20%23%3C%27\">foo? #&lt;&#39;</a>/<a href=\"/btns/example.net/foo%3F%20%23%3C%27/bar\">bar</a>") {
	//	t.Fatalf("expected a path in directory listing")
	//}
	if !strings.Contains(s, "<a href=\"/foo%3F%20%23%3C%27/bar/./..\">") {
		t.Fatalf("expected backlink in directory listing")
	}
	if !strings.Contains(s, "<a href=\"/foo%3F%20%23%3C%27/bar/file.txt\">") {
		t.Fatalf("expected file in directory listing")
	}
	if !strings.Contains(s, k3.Cid().String()) {
		t.Fatalf("expected hash in directory listing")
	}

	// make request to directory listing with prefix
	req, err = http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Host = "example.net"
	req.Header.Set("X-Btfs-Gateway-Prefix", "/good-prefix")

	res, err = doWithoutRedirect(req)
	if err != nil {
		t.Fatal(err)
	}

	// expect correct backlinks with prefix
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("error reading response: %s", err)
	}
	s = string(body)
	t.Logf("body: %s\n", string(body))

	//if !matchPathOrBreadcrumbs(s, "/btns/<a href=\"/btns/example.net\">example.net</a>") {
	//	t.Fatalf("expected a path in directory listing")
	//}
	if !strings.Contains(s, "<a href=\"/good-prefix/\">") {
		t.Fatalf("expected backlink in directory listing")
	}
	if !strings.Contains(s, "<a href=\"/good-prefix/file.txt\">") {
		t.Fatalf("expected file in directory listing")
	}
	if !strings.Contains(s, k.Cid().String()) {
		t.Fatalf("expected hash in directory listing")
	}

	// make request to directory listing with illegal prefix
	req, err = http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Host = "example.net"
	req.Header.Set("X-Btfs-Gateway-Prefix", "/bad-prefix")

	// make request to directory listing with evil prefix
	req, err = http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Host = "example.net"
	req.Header.Set("X-Btfs-Gateway-Prefix", "//good-prefix/foo")

	res, err = doWithoutRedirect(req)
	if err != nil {
		t.Fatal(err)
	}

	// expect correct backlinks without illegal prefix
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("error reading response: %s", err)
	}
	s = string(body)
	t.Logf("body: %s\n", string(body))

	//TODO
	//if !matchPathOrBreadcrumbs(s, "/") {
	//	t.Fatalf("expected a path in directory listing")
	//}
	if !strings.Contains(s, "<a href=\"/\">") {
		t.Fatalf("expected backlink in directory listing")
	}
	if !strings.Contains(s, "<a href=\"/file.txt\">") {
		t.Fatalf("expected file in directory listing")
	}
	if !strings.Contains(s, k.Cid().String()) {
		t.Fatalf("expected hash in directory listing")
	}
}

func TestCacheControlImmutable(t *testing.T) {
	ts, _, _ := newTestServerAndNode(t, nil)
	t.Logf("test server url: %s", ts.URL)

	req, err := http.NewRequest(http.MethodGet, ts.URL+emptyDir+"/", nil)
	if err != nil {
		t.Fatal(err)
	}

	res, err := doWithoutRedirect(req)
	if err != nil {
		t.Fatal(err)
	}

	// check the immutable tag isn't set
	hdrs, ok := res.Header["Cache-Control"]
	if ok {
		for _, hdr := range hdrs {
			if strings.Contains(hdr, "immutable") {
				t.Fatalf("unexpected Cache-Control: immutable on directory listing: %s", hdr)
			}
		}
	}
}

func TestGoGetSupport(t *testing.T) {
	ts, _, _ := newTestServerAndNode(t, nil)
	t.Logf("test server url: %s", ts.URL)

	// mimic go-get
	req, err := http.NewRequest(http.MethodGet, ts.URL+emptyDir+"?go-get=1", nil)
	if err != nil {
		t.Fatal(err)
	}

	res, err := doWithoutRedirect(req)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != 200 {
		t.Errorf("status is %d, expected 200", res.StatusCode)
	}
}

func TestVersion(t *testing.T) {
	version.CurrentCommit = "theshortcommithash"

	ns := mockNamesys{}
	ts, _, _ := newTestServerAndNode(t, ns)
	t.Logf("test server url: %s", ts.URL)

	req, err := http.NewRequest(http.MethodGet, ts.URL+"/version", nil)
	if err != nil {
		t.Fatal(err)
	}

	res, err := doWithoutRedirect(req)
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("error reading response: %s", err)
	}
	s := string(body)

	if !strings.Contains(s, "Commit: theshortcommithash") {
		t.Fatalf("response doesn't contain commit:\n%s", s)
	}

	if !strings.Contains(s, "Client Version: "+version.UserAgent) {
		t.Fatalf("response doesn't contain client version:\n%s", s)
	}

	if !strings.Contains(s, "Protocol Version: "+id.DefaultProtocolVersion) {
		t.Fatalf("response doesn't contain protocol version:\n%s", s)
	}
}
