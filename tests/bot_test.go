package tests

import (
	"context"
	"log"
	"testing"

	"github.com/huangc28/go-binance/v2/futures"
	"github.com/stretchr/testify/suite"
)

var (
	APIKey    = "14c8d9ffa86eb2d8148c1671297be01497363b6116f9a15c93e0eb378c43a2c2"
	APISecret = "522643237acaee837a8b7d1c9c2435cc0648ffc29bf00a89c9ddcd5fb37122ee"

	// Short order account
	ShortAPIKey    = "164d174af6d9b4c95e94a38e95e48fda32913f0373c5612d78ab3985bd8bed43"
	ShortAPISecret = "d6805cc1c9039f08a24053165b8585c8498117d473fdb901da51bdf52c7ecc9e"
)

type TradingBotTestSuite struct {
	suite.Suite
	futureClient *futures.Client
	shortClient  *futures.Client
	ctx          context.Context
}

func (suite *TradingBotTestSuite) SetupTest() {
	futures.UseTestnet = true
	suite.futureClient = futures.NewClient(APIKey, APISecret)
	suite.shortClient = futures.NewClient(ShortAPIKey, ShortAPISecret)
	suite.ctx = context.Background()
}

func (suite *TradingBotTestSuite) TestCheckPositionSide() {
	ctx := context.Background()
	posMode, err := suite.
		futureClient.
		NewGetPositionModeService().
		Do(ctx)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("debug posMode %v", posMode.DualSidePosition)
}

func (suite *TradingBotTestSuite) TestCheckShortClientPositionSide() {
	ctx := context.Background()
	posMode, err := suite.
		shortClient.
		NewGetPositionModeService().
		Do(ctx)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("debug short posMode %v", posMode.DualSidePosition)
}

func (suite *TradingBotTestSuite) TestCreateTakeProfit() {
	if _, err := suite.
		futureClient.
		NewCreateOrderService().
		Symbol("BTCUSDT").
		Side(futures.SideTypeSell).
		PositionSide(futures.PositionSideTypeLong).
		Type(futures.OrderTypeTakeProfitMarket).
		TimeInForce(futures.TimeInForceTypeGTC).
		StopPrice("27000").
		Quantity("0.01").
		NewOrderResponseType(futures.NewOrderRespTypeRESULT).
		Do(suite.ctx); err != nil {
		suite.T().Fatal(err)
	}
}

func (suite *TradingBotTestSuite) TestCreateStopLossOrder() {
	if _, err := suite.
		futureClient.
		NewCreateOrderService().
		Symbol("BTCUSDT").
		Side(futures.SideTypeSell).
		PositionSide(futures.PositionSideTypeLong).
		Type(futures.OrderTypeStopMarket).
		TimeInForce(futures.TimeInForceTypeGTC).
		StopPrice("27000").
		Quantity("0.01").
		NewOrderResponseType(futures.NewOrderRespTypeRESULT).
		Do(suite.ctx); err != nil {
		suite.T().Fatal(err)
	}
}

func (suite *TradingBotTestSuite) TestCreateShortOrder() {
	if _, err := suite.
		shortClient.
		NewCreateOrderService().
		Symbol("BTCUSDT").
		Side(futures.SideTypeSell).
		PositionSide(futures.PositionSideTypeShort).
		Type(futures.OrderTypeLimit).
		TimeInForce(futures.TimeInForceTypeGTC).
		Price("26842").
		Quantity("0.01").
		NewOrderResponseType(futures.NewOrderRespTypeRESULT).
		Do(suite.ctx); err != nil {
		suite.T().Fatal(err)
	}
}

func (suite *TradingBotTestSuite) TestCreateShortStopLossOrder() {
	if _, err := suite.
		shortClient.
		NewCreateOrderService().
		Symbol("BTCUSDT").
		Side(futures.SideTypeBuy).
		PositionSide(futures.PositionSideTypeLong).
		Type(futures.OrderTypeStop).
		TimeInForce(futures.TimeInForceTypeGTC).
		StopPrice("30100").
		Price("30100").
		Quantity("0.01").
		NewOrderResponseType(futures.NewOrderRespTypeRESULT).
		Do(suite.ctx); err != nil {
		suite.T().Fatal(err)
	}
}

func (suite *TradingBotTestSuite) TestCreateShortTakeProfitOrder() {
	if _, err := suite.
		shortClient.
		NewCreateOrderService().
		Symbol("BTCUSDT").
		Side(futures.SideTypeBuy).
		PositionSide(futures.PositionSideTypeLong).
		Type(futures.OrderTypeTakeProfitMarket).
		TimeInForce(futures.TimeInForceTypeGTC).
		StopPrice("27000").
		Quantity("0.01").
		NewOrderResponseType(futures.NewOrderRespTypeRESULT).
		Do(suite.ctx); err != nil {
		suite.T().Fatal(err)
	}
}

func (suite *TradingBotTestSuite) TestCancelAllOrder() {
	if err := suite.
		futureClient.
		NewCancelAllOpenOrdersService().
		Symbol("BTCUSDT").
		Do(suite.ctx); err != nil {
		suite.T().Fatal(err)
	}
}

func TestTradingBotTestSuite(t *testing.T) {
	suite.Run(t, new(TradingBotTestSuite))
}
