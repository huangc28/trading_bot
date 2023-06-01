package main

import (
	"context"
	"fmt"

	"github.com/huangc28/go-binance/v2"
	"github.com/huangc28/go-binance/v2/futures"
	log "github.com/sirupsen/logrus"
)

// 下 strategy 進場單，如果 strategy state 為 empty，create 一個 long order,
// 其入場價格為 market price -5usd
func wsAggTradeHandler(event *binance.WsAggTradeEvent, futureClient *futures.Client) {
	ctx := context.Background()

	ss := GetStrategyState()
	if ss.FSM.Is(Empty) {
		entryPrice := calcEntryPrice(event.Price)

		createOrderResp, err := futureClient.
			NewCreateOrderService().
			Symbol("BTCUSDT").
			Side(futures.SideTypeBuy).
			PositionSide(futures.PositionSideTypeLong).
			Type(futures.OrderTypeLimit).
			TimeInForce(futures.TimeInForceTypeGTC).
			Price(entryPrice).
			Quantity("0.01").
			NewOrderResponseType(futures.NewOrderRespTypeRESULT).
			Do(ctx)

		if err != nil {
			log.Fatalf("create order response error %v", err)
		}

		// We'll need to remember the created order so we can toggle the strategy state properly.
		// If we're being notified that the order has entered the longing position,
		// we can toggle the finite state to `longing`
		if err := ss.
			FSM.
			Event(ctx, OpenLimitOrder, createOrderResp); err != nil {
			log.Fatalf("failed to change fsm state %v", err)
		}
	}

	// 這邊我們要聽 market price，如果 market price 接近多單的 SL，馬上開空單 market price 進場。
}

func wsAggTradeErrHandler(err error) {
	fmt.Printf("trade err %v", err)
}
