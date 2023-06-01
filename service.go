package main

import (
	"context"

	"github.com/huangc28/go-binance/v2/futures"
)

type createHedgeShortOrderParams struct {
	ShortClient *futures.Client
	LongSLPrice string
	Quantity    string
}

// When long order is executed, we need to create another short order
// to hedge the long position
func createHedgeShortOrder(ctx context.Context, params createHedgeShortOrderParams) error {
	// shortOrder, shortOrderErr := params.
	// 	ShortClient.
	// 	NewCreateOrderService().
	// 	Symbol("BTCUSDT").
	// 	Side(futures.SideTypeSell).
	// 	PositionSide(futures.PositionSideTypeShort).
	// 	Type(futures.OrderTypeLimit).
	// 	TimeInForce(futures.TimeInForceTypeGTC).
	// 	Price(params.LongSLPrice).
	// 	Quantity(params.Quantity).
	// 	NewOrderResponseType(futures.NewOrderRespTypeRESULT).
	// 	Do(ctx)

	// if shortOrderErr != nil {
	// 	return shortOrderErr
	// }

	// shortSLPrice, err := calcShortOrderStopLoss(
	// 	params.LongSLPrice,
	// 	35*BasisPoint,
	// )

	// if err != nil {
	// 	log.Fatalf("failed to calculate short stop loss price %v", err)
	// }

	// if shortSLOrder, err := params.
	// 	ShortClient.
	// 	NewCreateOrderService().
	// 	Symbol("BTCUSDT").
	// 	Side(futures.SideTypeBuy).
	// 	PositionSide(futures.PositionSideTypeLong).
	// 	Type(futures.OrderTypeStop).
	// 	TimeInForce(futures.TimeInForceTypeGTC).
	// 	StopPrice(shortSLPrice).
	// 	Price(shortSLPrice).
	// 	Quantity(params.Quantity).
	// 	NewOrderResponseType(futures.NewOrderRespTypeRESULT).
	// 	Do(ctx); err != nil {

	// 	log.Fatalf("failed to create short stop loss order %v", err)
	// }

	// shortTPPrice, err := calcShortOrderTakeProfit(
	// 	params.LongSLPrice,
	// 	30*BasisPoint,
	// )

	// if err != nil {
	// 	log.Fatalf("failed to calculate short take profit price %v", err)
	// }

	return nil
}
