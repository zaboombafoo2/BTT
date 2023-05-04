package cheque

import (
	"errors"
	"fmt"
	"github.com/bittorrent/go-btfs/chain/tokencfg"
	"github.com/bittorrent/go-btfs/utils"
	"io"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
)

type SendTotalCountRet struct {
	Count int `json:"count"`
}

var SendChequesCountCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "send cheque(s) count",
	},
	Options: []cmds.Option{
		cmds.StringOption(tokencfg.TokenTypeName, "tk", "file storage with token type,default WBTT, other TRX/USDD/USDT.").WithDefault("WBTT"),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		err := utils.CheckSimpleMode(env)
		if err != nil {
			return err
		}

		tokenStr := req.Options[tokencfg.TokenTypeName].(string)
		//fmt.Printf("... token:%+v\n", tokenStr)
		token, bl := tokencfg.MpTokenAddr[tokenStr]
		if !bl {
			return errors.New("your input token is none. ")
		}

		count, err := chain.SettleObject.VaultService.TotalIssuedCount(token)
		if err != nil {
			return err
		}

		return cmds.EmitOnce(res, &SendTotalCountRet{Count: count})
	},
	Type: SendTotalCountRet{},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, c *SendTotalCountRet) error {
			fmt.Println("send cheque(s) count: ", c.Count)

			return nil
		}),
	},
}
