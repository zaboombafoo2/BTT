package settlement

import (
	"context"
	"fmt"
	"math/big"
	"time"

	cmds "github.com/bittorrent/go-btfs-cmds"
	"github.com/bittorrent/go-btfs/chain"
)

type settlementResponse struct {
	Peer               string   `json:"peer"`
	SettlementReceived *big.Int `json:"received"`
	SettlementSent     *big.Int `json:"sent"`
}

type settlementsResponse struct {
	TotalSettlementReceived  *big.Int             `json:"totalReceived"`
	TotalSettlementSent      *big.Int             `json:"totalSent"`
	SettlementReceivedCashed *big.Int             `json:"settlement_received_cashed"`
	SettlementSentCashed     *big.Int             `json:"settlement_sent_cashed"`
	Settlements              []settlementResponse `json:"settlements"`
}

var ListSettlementCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "list all settlements.",
	},
	RunTimeout: 5 * time.Minute,
	Arguments: []cmds.Argument{
		cmds.StringArg("token", true, false, "token"),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		token := req.Arguments[0]
		fmt.Printf("... token:%+v\n", token)

		settlementsSent, err := chain.SettleObject.SwapService.SettlementsSent(token)
		if err != nil {
			return err
		}
		settlementsReceived, err := chain.SettleObject.SwapService.SettlementsReceived(token)
		if err != nil {
			return err
		}

		totalReceived := big.NewInt(0)
		totalReceivedCashed := big.NewInt(0)
		totalSent := big.NewInt(0)

		settlementResponses := make(map[string]settlementResponse)

		for a, b := range settlementsSent {
			settlementResponses[a] = settlementResponse{
				Peer:               a,
				SettlementSent:     b,
				SettlementReceived: big.NewInt(0),
			}
			totalSent.Add(b, totalSent)
		}

		for a, b := range settlementsReceived {
			if _, ok := settlementResponses[a]; ok {
				t := settlementResponses[a]
				t.SettlementReceived = b
				settlementResponses[a] = t
			} else {
				settlementResponses[a] = settlementResponse{
					Peer:               a,
					SettlementSent:     big.NewInt(0),
					SettlementReceived: b,
				}
			}
			totalReceived.Add(b, totalReceived)
			if has, err := chain.SettleObject.SwapService.HasCashoutAction(context.Background(), a, token); err == nil && has {
				totalReceivedCashed.Add(b, totalReceivedCashed)
			}
		}
		settlementResponsesArray := make([]settlementResponse, len(settlementResponses))
		i := 0
		for k := range settlementResponses {
			settlementResponsesArray[i] = settlementResponses[k]
			i++
		}

		totalPaidOut, err := chain.SettleObject.VaultService.TotalPaidOut(context.Background(), token)
		if err != nil {
			return err
		}
		rsp := settlementsResponse{
			TotalSettlementReceived:  totalReceived,
			TotalSettlementSent:      totalSent,
			SettlementReceivedCashed: totalReceivedCashed,
			SettlementSentCashed:     totalPaidOut,
			Settlements:              settlementResponsesArray,
		}

		return cmds.EmitOnce(res, &rsp)
	},
	Type: &settlementsResponse{},
}
