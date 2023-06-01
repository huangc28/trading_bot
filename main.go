package main

import (
	"context"
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"

	"github.com/huangc28/go-binance/v2"
	"github.com/huangc28/go-binance/v2/futures"
)

var (
	// APIKey    = "Os9hSbx9BLZmfoCPQcpM8ABGXV5r1DqS946JfleEx73xr8VCULkEo8aE7rKAlerp"
	// APISecret = "y7OXWWRb7da0vClEthHUUowbsDH6gJJtVc0tQvLATm9ii6o3NRbpfpqw4W37oyPo"

	// Long order account
	APIKey    = "14c8d9ffa86eb2d8148c1671297be01497363b6116f9a15c93e0eb378c43a2c2"
	APISecret = "522643237acaee837a8b7d1c9c2435cc0648ffc29bf00a89c9ddcd5fb37122ee"

	// Short order account
	ShortAPIKey   = "d6805cc1c9039f08a24053165b8585c8498117d473fdb901da51bdf52c7ecc9e"
	LongAPISecret = "164d174af6d9b4c95e94a38e95e48fda32913f0373c5612d78ab3985bd8bed43"
)

func init() {
	// valid testnet websocket
	binance.SetBaseWsTestnetURL("wss://stream.binancefuture.com/ws")
	binance.UseTestnet = true
	futures.UseTestnet = true
}

func main() {
	ctx := context.Background()
	futureClient := futures.NewClient(APIKey, APISecret)

	shortClient := futures.NewClient(
		ShortAPIKey,
		LongAPISecret,
	)

	wscArr := make([]chan struct{}, 0, 2)
	wscChan := make(chan chan struct{})

	if _, err := futureClient.
		NewChangeLeverageService().
		Leverage(2).
		Symbol("BTCUSDT").
		Do(ctx); err != nil {
		log.Fatalf("failed to change initial leverage %v", err)
	}

	// Listen to price tick to create order.
	go func(futureClient *futures.Client) {
		doneC, _, err := binance.WsAggTradeServe(
			"BTCUSDT",
			func(event *binance.WsAggTradeEvent) {
				wsAggTradeHandler(event, futureClient)
			},

			func(err error) {
				wsAggTradeErrHandler(err)
			},
		)

		if err != nil {
			log.Errorf("failed to listen to price tick %v", err)
			return
		}

		wscChan <- doneC

		<-doneC
	}(futureClient)

	// Listen to profile info to know if order has been executed.
	go func(futureClient *futures.Client, shortClient *futures.Client) {
		ctx := context.Background()

		// Create user data listenKey
		listenKey, err := futureClient.
			NewStartUserStreamService().
			Do(ctx)

		if err != nil {
			log.Errorf("failed to create listen key %v", err)
			return
		}

		log.Infof("listenKey: %v", listenKey)

		// Keep listenKey alive
		if err := futureClient.
			NewKeepaliveUserStreamService().
			ListenKey(listenKey).
			Do(ctx); err != nil {
			log.Errorf("failed to keepalive listen key %v", err)
			return
		}

		log.Infof("keepalive listenKey: %v", listenKey)

		doneC, _, err := binance.WsUserDataServe(
			listenKey,
			func(event *binance.WsUserDataEvent) {
				wsUserDataHandler(event, futureClient, shortClient)
			},

			func(err error) {
				wsUserDataErrorHandler(err, futureClient)
			},
		)

		if err != nil {
			log.Errorf("failed to listen to user data", err)
			return
		}

		wscChan <- doneC

		<-doneC
	}(
		futureClient,
		shortClient,
	)

	for wsc := range wscChan {
		wscArr = append(wscArr, wsc)

		if len(wscArr) == cap(wscArr) {
			close(wscChan)
			break
		}
	}

	log.Infof("ws context %v", len(wscArr))

	c := make(chan os.Signal, 1)

	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive interrupt signal
	<-c

	for _, wsc := range wscArr {
		close(wsc)
	}
}
