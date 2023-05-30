package node

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	config "github.com/TRON-US/go-btfs-config"
	uio "github.com/TRON-US/go-unixfs/io"
	"github.com/bittorrent/go-btfs/core/node/libp2p"
	"github.com/bittorrent/go-btfs/p2p"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	offline "github.com/ipfs/go-ipfs-exchange-offline"
	util "github.com/ipfs/go-ipfs-util"
	log "github.com/ipfs/go-log"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	peer "github.com/libp2p/go-libp2p/core/peer"
	rcmgr "github.com/libp2p/go-libp2p/p2p/host/resource-manager"
	"github.com/tron-us/go-btfs-common/crypto"
	"go.uber.org/fx"
)

var logger = log.Logger("core:constructor")

var BaseLibP2P = fx.Options(
	fx.Provide(libp2p.UserAgent),
	fx.Provide(libp2p.PNet),
	fx.Provide(libp2p.ConnectionManager),
	fx.Provide(libp2p.Host),
	fx.Provide(libp2p.MultiaddrResolver),

	fx.Provide(libp2p.DiscoveryHandler),

	fx.Invoke(libp2p.PNetChecker),
)

func LibP2P(bcfg *BuildCfg, cfg *config.Config, userResourceOverrides rcmgr.PartialLimitConfig) fx.Option {
	// parse ConnMgr config
	var connmgr fx.Option

	// set connmgr based on Swarm.ConnMgr.Type
	connMgrType := cfg.Swarm.ConnMgr.Type.WithDefault(config.DefaultConnMgrType)
	switch connMgrType {
	case "none":
		connmgr = fx.Options() // noop
	case "", "basic":
		grace := cfg.Swarm.ConnMgr.GracePeriod.WithDefault(config.DefaultConnMgrGracePeriod)
		low := int(cfg.Swarm.ConnMgr.LowWater.WithDefault(config.DefaultConnMgrLowWater))
		high := int(cfg.Swarm.ConnMgr.HighWater.WithDefault(config.DefaultConnMgrHighWater))
		connmgr = fx.Provide(libp2p.ConnectionManager(low, high, grace))
	default:
		return fx.Error(fmt.Errorf("unrecognized Swarm.ConnMgr.Type: %q", connMgrType))
	}

	// parse PubSub config

	ps, disc := fx.Options(), fx.Options()
	if bcfg.getOpt("pubsub") || bcfg.getOpt("ipnsps") {
		disc = fx.Provide(libp2p.TopicDiscovery())

		var pubsubOptions []pubsub.Option
		pubsubOptions = append(
			pubsubOptions,
			pubsub.WithMessageSigning(!cfg.Pubsub.DisableSigning),
		)

		switch cfg.Pubsub.Router {
		case "":
			fallthrough
		case "gossipsub":
			ps = fx.Provide(libp2p.GossipSub(pubsubOptions...))
		case "floodsub":
			ps = fx.Provide(libp2p.FloodSub(pubsubOptions...))
		default:
			return fx.Error(fmt.Errorf("unknown pubsub router %s", cfg.Pubsub.Router))
		}
	}

	autonat := fx.Options()

	switch cfg.AutoNAT.ServiceMode {
	default:
		panic("BUG: unhandled autonat service mode")
	case config.AutoNATServiceDisabled:
	case config.AutoNATServiceUnset:
		// TODO
		//
		// We're enabling the AutoNAT service by default on _all_ nodes
		// for the moment.
		//
		// We should consider disabling it by default if the dht is set
		// to dhtclient.
		fallthrough
	case config.AutoNATServiceEnabled:
		autonat = fx.Provide(libp2p.AutoNATService(cfg.AutoNAT.Throttle))
	}

	enableRelayTransport := cfg.Swarm.Transports.Network.Relay.WithDefault(true) // nolint
	enableRelayService := cfg.Swarm.RelayService.Enabled.WithDefault(enableRelayTransport)
	enableRelayClient := cfg.Swarm.RelayClient.Enabled.WithDefault(enableRelayTransport)

	// Log error when relay subsystem could not be initialized due to missing dependency
	if !enableRelayTransport {
		if enableRelayService {
			logger.Warn("Failed to enable `Swarm.RelayService`, it requires `Swarm.Transports.Network.Relay` to be true.")
		}
		if enableRelayClient {
			logger.Warn("Failed to enable `Swarm.RelayClient`, it requires `Swarm.Transports.Network.Relay` to be true.")
		}
	}

	// TODO: Force users to migrate old config.
	// nolint
	if cfg.Swarm.DisableRelay {
		logger.Warn("The 'Swarm.DisableRelay' config field was removed." +
			"Use the 'Swarm.Transports.Network.Relay' instead.")
	}
	// nolint
	if cfg.Swarm.EnableAutoRelay {
		logger.Warn("The 'Swarm.EnableAutoRelay' config field was removed." +
			"Use the 'Swarm.RelayClient.Enabled' instead.")
	}
	// nolint
	if cfg.Swarm.EnableRelayHop {
		logger.Warn("The `Swarm.EnableRelayHop` config field was removed.\n" +
			"Use `Swarm.RelayService` to configure the circuit v2 relay.\n" +
			"If you want to continue running a circuit v1 relay, please use the standalone relay daemon: https://dist.ipfs.tech/#libp2p-relay-daemon (with RelayV1.Enabled: true)")
	}

	// Gather all the options
	opts := fx.Options(
		BaseLibP2P,

		fx.Provide(libp2p.ResourceManager(cfg.Swarm, userResourceOverrides)),
		fx.Provide(libp2p.AddrFilters(cfg.Swarm.AddrFilters)),
		fx.Provide(libp2p.AddrsFactory(cfg.Addresses.Announce, cfg.Addresses.NoAnnounce)),
		fx.Provide(libp2p.SmuxTransport(cfg.Swarm.Transports)),
		fx.Provide(libp2p.RelayTransport(enableRelayTransport)),
		fx.Provide(libp2p.RelayService(enableRelayService, cfg.Swarm.RelayService)),
		fx.Provide(libp2p.Transports(cfg.Swarm.Transports)),
		fx.Invoke(libp2p.StartListening(cfg.Addresses.Swarm)),
		fx.Invoke(libp2p.SetupDiscovery(cfg.Discovery.MDNS.Enabled, cfg.Discovery.MDNS.Interval)),
		fx.Provide(libp2p.ForceReachability(cfg.Internal.Libp2pForceReachability)),
		fx.Provide(libp2p.HolePunching(cfg.Swarm.EnableHolePunching, enableRelayClient)),

		fx.Provide(libp2p.Security(!bcfg.DisableEncryptedConnections, cfg.Swarm.Transports)),

		fx.Provide(libp2p.Routing),
		fx.Provide(libp2p.ContentRouting),

		fx.Provide(libp2p.BaseRouting(cfg.Experimental.AcceleratedDHTClient)),
		maybeProvide(libp2p.PubsubRouter, bcfg.getOpt("ipnsps")),

		maybeProvide(libp2p.BandwidthCounter, !cfg.Swarm.DisableBandwidthMetrics),
		maybeProvide(libp2p.NatPortMap, !cfg.Swarm.DisableNatPortMap),
		libp2p.MaybeAutoRelay(cfg.Swarm.RelayClient.StaticRelays, cfg.Peering, enableRelayClient),
		autonat,
		connmgr,
		ps,
		disc,
	)

	return opts
}

// Storage groups units which setup datastore based persistence and blockstore layers
func Storage(bcfg *BuildCfg, cfg *config.Config) fx.Option {
	cacheOpts := blockstore.DefaultCacheOpts()
	cacheOpts.HasBloomFilterSize = cfg.Datastore.BloomFilterSize
	if !bcfg.Permanent {
		cacheOpts.HasBloomFilterSize = 0
	}

	finalBstore := fx.Provide(GcBlockstoreCtor)
	if cfg.Experimental.FilestoreEnabled || cfg.Experimental.UrlstoreEnabled {
		finalBstore = fx.Provide(FilestoreBlockstoreCtor)
	}

	return fx.Options(
		fx.Provide(RepoConfig),
		fx.Provide(Datastore),
		fx.Provide(BaseBlockstoreCtor(cacheOpts, bcfg.NilRepo, cfg.Datastore.HashOnRead)),
		finalBstore,
	)
}

// Identity groups units providing cryptographic identity
func Identity(cfg *config.Config) fx.Option {
	// PeerID

	cid := cfg.Identity.PeerID
	if cid == "" {
		return fx.Error(errors.New("identity was not set in config (was 'btfs init' run?)"))
	}
	if len(cid) == 0 {
		return fx.Error(errors.New("no peer ID in config! (was 'btfs init' run?)"))
	}

	id, err := peer.Decode(cid)
	if err != nil {
		return fx.Error(fmt.Errorf("peer ID invalid: %s", err))
	}

	// Private Key
	// Use env override if available
	pk := os.Getenv("BTFS_PRIV_KEY")
	if pk != "" {
		// Override stored peer id
		cfg.Identity.PrivKey = pk
	}

	if cfg.Identity.PrivKey == "" {
		return fx.Options( // No PK (usually in tests)
			fx.Provide(PeerID(id)),
			fx.Provide(libp2p.Peerstore),
		)
	}

	sk, err := crypto.GetPrivKeyFromHexOrBase64(cfg.Identity.PrivKey)
	if err != nil {
		return fx.Error(err)
	}

	// Set correct peer id from overriden private key
	if pk != "" {
		pid, err := peer.IDFromPublicKey(sk.GetPublic())
		if err != nil {
			return fx.Error(err)
		}
		cfg.Identity.PeerID = pid.String()
		id = pid
	}

	return fx.Options( // Full identity
		fx.Provide(PeerID(id)),
		fx.Provide(PrivateKey(sk)),
		fx.Provide(libp2p.Peerstore),

		fx.Invoke(libp2p.PstoreAddSelfKeys),
	)
}

// IPNS groups namesys related units
var IPNS = fx.Options(
	fx.Provide(RecordValidator),
)

// Online groups online-only units
func Online(bcfg *BuildCfg, cfg *config.Config, userResourceOverrides rcmgr.PartialLimitConfig) fx.Option {

	// Namesys params

	ipnsCacheSize := cfg.Ipns.ResolveCacheSize
	if ipnsCacheSize == 0 {
		ipnsCacheSize = DefaultIpnsCacheSize
	}
	if ipnsCacheSize < 0 {
		return fx.Error(fmt.Errorf("cannot specify negative resolve cache size"))
	}

	// Republisher params

	var repubPeriod, recordLifetime time.Duration

	if cfg.Ipns.RepublishPeriod != "" {
		d, err := time.ParseDuration(cfg.Ipns.RepublishPeriod)
		if err != nil {
			return fx.Error(fmt.Errorf("failure to parse config setting BTNS.RepublishPeriod: %s", err))
		}

		if !util.Debug && (d < time.Minute || d > (time.Hour*24)) {
			return fx.Error(fmt.Errorf("config setting BTNS.RepublishPeriod is not between 1min and 1day: %s", d))
		}

		repubPeriod = d
	}

	if cfg.Ipns.RecordLifetime != "" {
		d, err := time.ParseDuration(cfg.Ipns.RecordLifetime)
		if err != nil {
			return fx.Error(fmt.Errorf("failure to parse config setting BTNS.RecordLifetime: %s", err))
		}

		recordLifetime = d
	}

	/* don't provide from bitswap when the strategic provider service is active */
	shouldBitswapProvide := !cfg.Experimental.StrategicProviding

	return fx.Options(
		fx.Provide(OnlineExchange(shouldBitswapProvide)),
		maybeProvide(Graphsync, cfg.Experimental.GraphsyncEnabled),
		fx.Provide(DNSResolver),
		fx.Provide(Namesys(ipnsCacheSize)),
		fx.Provide(Peering),
		PeerWith(cfg.Peering.Peers...),

		fx.Invoke(IpnsRepublisher(repubPeriod, recordLifetime)),

		fx.Provide(p2p.New),

		LibP2P(bcfg, cfg, userResourceOverrides),
		OnlineProviders(
			cfg.Experimental.StrategicProviding,
			cfg.Experimental.AcceleratedDHTClient,
			cfg.Reprovider.Strategy.WithDefault(config.DefaultReproviderStrategy),
			cfg.Reprovider.Interval.WithDefault(config.DefaultReproviderInterval),
		),
	)
}

// Offline groups offline alternatives to Online units
func Offline(cfg *config.Config) fx.Option {
	return fx.Options(
		fx.Provide(offline.Exchange),
		fx.Provide(DNSResolver),
		fx.Provide(Namesys(0)),
		fx.Provide(libp2p.Routing),
		fx.Provide(libp2p.ContentRouting),
		fx.Provide(libp2p.OfflineRouting),
		OfflineProviders(
			cfg.Experimental.StrategicProviding,
			cfg.Experimental.AcceleratedDHTClient,
			cfg.Reprovider.Strategy.WithDefault(config.DefaultReproviderStrategy),
			cfg.Reprovider.Interval.WithDefault(config.DefaultReproviderInterval),
		),
	)
}

// Core groups basic BTFS services
var Core = fx.Options(
	fx.Provide(BlockService),
	fx.Provide(Dag),
	fx.Provide(FetcherConfig),
	fx.Provide(Pinning),
	fx.Provide(Files),
)

func Networked(bcfg *BuildCfg, cfg *config.Config, userResourceOverrides rcmgr.PartialLimitConfig) fx.Option {
	if bcfg.Online {
		return Online(bcfg, cfg, userResourceOverrides)
	}
	return Offline(cfg)
}

// BTFS builds a group of fx Options based on the passed BuildCfg
func IPFS(ctx context.Context, bcfg *BuildCfg) fx.Option {
	if bcfg == nil {
		bcfg = new(BuildCfg)
	}

	bcfgOpts, cfg := bcfg.options(ctx)
	if cfg == nil {
		return bcfgOpts // error
	}
	userResourceOverrides, err := bcfg.Repo.UserResourceOverrides()
	if err != nil {
		return fx.Error(err)
	}
	// TEMP: setting global sharding switch here
	uio.UseHAMTSharding = cfg.Experimental.ShardingEnabled

	return fx.Options(
		bcfgOpts,

		fx.Provide(baseProcess),

		Storage(bcfg, cfg),
		Identity(cfg),
		IPNS,
		Networked(bcfg, cfg, userResourceOverrides),

		Core,
	)
}
