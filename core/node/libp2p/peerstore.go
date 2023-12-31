package libp2p

import (
	"context"

	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/p2p/host/peerstore/pstoremem"
	"go.uber.org/fx"
)

func Peerstore(lc fx.Lifecycle) peerstore.Peerstore {
	pstore, err := pstoremem.NewPeerstore()
	if err != nil {
		log.Errorln(err)
		return nil
	}
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return pstore.Close()
		},
	})

	return pstore
}
