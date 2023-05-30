package tests

import (
	"context"
	"log"
	"testing"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/stretchr/testify/suite"
)

var (
	APIKey    = "14c8d9ffa86eb2d8148c1671297be01497363b6116f9a15c93e0eb378c43a2c2"
	APISecret = "522643237acaee837a8b7d1c9c2435cc0648ffc29bf00a89c9ddcd5fb37122ee"
)

type TradingBotTestSuite struct {
	suite.Suite
	futureClient *futures.Client
}

func (suite *TradingBotTestSuite) SetupTest() {
	futures.UseTestnet = true
	suite.futureClient = futures.NewClient(APIKey, APISecret)
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

func TestTradingBotTestSuite(t *testing.T) {
	suite.Run(t, new(TradingBotTestSuite))
}
