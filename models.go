package main

import (
	"github.com/huangc28/go-binance/v2"
	"github.com/huangc28/go-binance/v2/futures"
)

// LongPositionInfo Store long position with it's paired tp and sl order
// when strategy status is `longing`
type LongPositionInfo struct {
	LongPosition    *binance.WsOrderTradeUpdate
	TakeProfitOrder *futures.CreateOrderResponse
	StopLossOrder   *futures.CreateOrderResponse
}
