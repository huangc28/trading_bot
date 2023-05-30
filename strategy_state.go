package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/looplab/fsm"
)

type StrategyState struct {
	FSM *fsm.FSM

	// The long order of the current strategy deployment.
	LongOrder *futures.CreateOrderResponse

	// The long position of the current strategy deployment.
	// LongPosition
}

var (
	_strategyState *StrategyState
	once           sync.Once
)

// List of strategy states.
var (
	Empty string = "empty"

	Done string = "done"

	// WaitingToEnterLong order has been created, waiting to be executed
	// to enter long position.
	WaitingToEnterLong = "WaitingToEnterLong"

	// Long order has been executed.
	Longing = "Longing"

	// Short order is at the other account. Thus, we'll check the short order
	// on the other account is bought in before toggling from `WaitingToEnterShort`
	// to `Shorting`
	WaitingToEnterShort = "WaitingToEnterShort"

	Shorting = "Shorting"

	ShortLoss = "ShortLoss"
)

// List of strategy events.
var (
	OpenLimitOrder = "OpenLimitOrder"

	EnterLong = "EnterLong"

	HitLongTakeProfit = "HitLongTakeProfit"

	HitLongStopLoss = "HitLongStopLoss"

	EnterShort = "EnterShort"

	HitShortStopLoss = "HitShortStopLoss"

	HitShortTakeProfit = "HitShortTakeProfit"
)

func GetStrategyState() *StrategyState {
	once.Do(func() {
		_strategyState = &StrategyState{}

		_strategyState.FSM = fsm.NewFSM(
			Empty,
			fsm.Events{
				{
					Name: OpenLimitOrder,
					Src:  []string{Empty, Done},
					Dst:  WaitingToEnterLong,
				},
				{
					Name: EnterLong,
					Src:  []string{WaitingToEnterLong},
					Dst:  Longing,
				},
				{
					Name: HitLongTakeProfit,
					Src:  []string{Longing},
					Dst:  Done,
				},
				{
					Name: HitLongStopLoss,
					Src:  []string{Longing},
					Dst:  WaitingToEnterShort,
				},
				{
					Name: EnterShort,
					Src:  []string{WaitingToEnterShort},
					Dst:  Shorting,
				},
				{
					Name: HitShortStopLoss,
					Src:  []string{Shorting},
					Dst:  ShortLoss,
				},
				{
					Name: HitShortTakeProfit,
					Src:  []string{Shorting},
					Dst:  Done,
				},
			},
			fsm.Callbacks{
				fmt.Sprintf("enter_%s", WaitingToEnterLong): func(ctx context.Context, evt *fsm.Event) {
					log.Printf("debug enter %v", WaitingToEnterLong)
					log.Printf("debug enter 2 %v", evt.Args)
					if len(evt.Args) > 0 {
						longOrder := evt.Args[0].(*futures.CreateOrderResponse)
						_strategyState.LongOrder = longOrder
					} else {
						// WTF... handle this error
					}
				},
			},
		)
	})

	return _strategyState
}
