package state

import (
	"github.com/looplab/fsm"

	"order-service/constant"
)

var status = "status"
var mapProcessFlow = map[string]fsm.Events{
	status: {
		{
			Name: constant.Pending.String(),
			Src:  []string{constant.Inital.String()},
			Dst:  constant.Pending.String(),
		},
		{
			Name: constant.PendingPayment.String(),
			Src:  []string{constant.Pending.String()},
			Dst:  constant.PendingPayment.String(),
		},
		{
			Name: constant.PaymentSuccess.String(),
			Src:  []string{constant.PendingPayment.String()},
			Dst:  constant.PaymentSuccess.String(),
		},
		{
			Name: constant.Completed.String(),
			Src:  []string{constant.PaymentSuccess.String()},
			Dst:  constant.Completed.String(),
		},
		{
			Name: constant.Cancelled.String(),
			Src: []string{
				constant.Pending.String(),
				constant.PendingPayment.String(),
			},
			Dst: constant.Cancelled.String(),
		},
	},
}

type StatusState struct {
	FSM *fsm.FSM
}

func NewStatusState(statusFlow constant.OrderStatus) *StatusState {
	return &StatusState{
		FSM: fsm.NewFSM(
			statusFlow.String(),
			mapProcessFlow[status],
			fsm.Callbacks{},
		),
	}
}
