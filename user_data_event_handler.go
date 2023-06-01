package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/huangc28/go-binance/v2"
	"github.com/huangc28/go-binance/v2/futures"
)

func calcTakeProfitPrice(filledPrice string, tpPercentage float64) (string, error) {
	fpf, err := strconv.ParseFloat(filledPrice, 64)
	if err != nil {
		return "", err
	}
	tp := fpf + (fpf * tpPercentage)
	return fmt.Sprintf("%.1f", tp), nil
}

func calcStopLossPrice(filledPrice string, slPercentage float64) (string, error) {
	fpf, err := strconv.ParseFloat(filledPrice, 64)
	if err != nil {
		return "", err
	}

	sl := fpf * (1 - slPercentage)
	return fmt.Sprintf("%.1f", sl), nil
}

func wsUserDataHandler(event *binance.WsUserDataEvent, futuresClient, shortClient *futures.Client) {
	ctx := context.Background()

	ss := GetStrategyState()

	// If order has been executed and the excuted order is the order
	// that we created and strategy status state machine is `waiting_to_enter_long`.
	if event.Event == binance.UserDataEventTypeOrderTradeUpdate {
		// Makesure the following:
		//
		//   1. our long order has been executed.
		//   2. the executed order is indeed our order
		//   3. order status is `FILLED`
		//   4. filled quantity equals to order quantity
		//
		// Once the above are assured, we can start create take profit order
		// and stop loss order.
		//
		// start creating
		//   - take profit order
		//   - stop loss order
		//
		// @TODO
		//  deal with cenario when failed to create
		//    - stop loss order
		//    - take profit order
		//  write data to file log
		if event.OrderTradeUpdate.OriginalOrderType == "LIMIT" &&
			event.OrderTradeUpdate.OrderStatus == "FILLED" &&
			ss.LongOrder.ClientOrderID == event.OrderTradeUpdate.ClientOrderId {

			// 測試先用 0.25% take profit
			takeProfitPrice, err := calcTakeProfitPrice(
				event.OrderTradeUpdate.AveragePrice,
				25*BasisPoint,
			)

			if err != nil {
				log.Fatalf("failed to calc take profit price %v", err)
			}

			// 測試先用 0.3% stop loss
			slPrice, err := calcStopLossPrice(
				event.OrderTradeUpdate.AveragePrice,
				30*BasisPoint,
			)

			if err != nil {
				log.Fatalf("failed to calc stop loss price %v", err)
			}

			log.Printf("debug 1 %v", slPrice)

			// If error happens, notify admin to either stop the bot or reset the bot.
			tpOrder, tpOrderErr := futuresClient.
				NewCreateOrderService().
				Symbol("BTCUSDT").
				Side(futures.SideTypeSell).
				PositionSide(futures.PositionSideTypeLong).
				Type(futures.OrderTypeTakeProfitMarket).
				TimeInForce(futures.TimeInForceTypeGTC).
				StopPrice(takeProfitPrice).
				Quantity(event.OrderTradeUpdate.FilledVolume).
				NewOrderResponseType(futures.NewOrderRespTypeRESULT).
				Do(ctx)

			if tpOrderErr != nil {
				log.Fatalf("failed to create take profit order %v", tpOrderErr)
			}

			slOrder, slOrderErr := futuresClient.
				NewCreateOrderService().
				Symbol("BTCUSDT").
				Side(futures.SideTypeSell).
				PositionSide(futures.PositionSideTypeLong).
				Type(futures.OrderTypeStopMarket).
				TimeInForce(futures.TimeInForceTypeGTC).
				StopPrice(slPrice).
				Quantity(event.OrderTradeUpdate.FilledVolume).
				NewOrderResponseType(futures.NewOrderRespTypeRESULT).
				Do(ctx)

			if slOrderErr != nil {
				log.Fatalf("failed to create stop loss order %v", slOrderErr)
			}

			// both tp & sl have been created, toggle state machine to `longing`
			// and store all tp, sl and long position.
			ss.FSM.Event(
				ctx,
				EnterLong,

				// Long position.
				event.OrderTradeUpdate,

				// Long position take profit order
				tpOrder,

				// Long position stop loss order
				slOrder,
			)
		}

		// Long position hits stop loss
		if event.OrderTradeUpdate.OriginalOrderType == "TAKE_PROFIT_MARKET" &&
			event.OrderTradeUpdate.OrderStatus == "FILLED" &&
			ss.LongPositionInfo.LongPosition.ClientOrderId == event.OrderTradeUpdate.ClientOrderId {
			// Cancel all orders
			if err := futuresClient.
				NewCancelAllOpenOrdersService().
				Symbol("BTCUSDT").
				Do(ctx); err != nil {
				log.Fatalf("failed to cancel all open order %v", err)
			}

			// Write result to log
			ss.FSM.Event(ctx, HitLongTakeProfit)
		}
	}
}

func wsUserDataErrorHandler(err error, futureClient *futures.Client) {
	log.Printf("debug wsUserDataHandler err %v", err)
}
