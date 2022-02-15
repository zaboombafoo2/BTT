package settlement

import (
	cmds "github.com/bittorrent/go-btfs-cmds"
)

var SettlementCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Interact with chequebook services on BTFS.",
	},
	Subcommands: map[string]*cmds.Command{
		"list": ListSettlementCmd,
		"peer": PeerSettlementCmd,
	},
}
