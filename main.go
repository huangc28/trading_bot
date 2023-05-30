package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
)

var (
	// APIKey    = "Os9hSbx9BLZmfoCPQcpM8ABGXV5r1DqS946JfleEx73xr8VCULkEo8aE7rKAlerp"
	// APISecret = "y7OXWWRb7da0vClEthHUUowbsDH6gJJtVc0tQvLATm9ii6o3NRbpfpqw4W37oyPo"

	APIKey    = "14c8d9ffa86eb2d8148c1671297be01497363b6116f9a15c93e0eb378c43a2c2"
	APISecret = "522643237acaee837a8b7d1c9c2435cc0648ffc29bf00a89c9ddcd5fb37122ee"
)

func main() {
	futures.UseTestnet = true
	ctx := context.Background()
	futureClient := futures.NewClient(APIKey, APISecret)

	wscArr := []chan struct{}{}
	wscChan := make(chan chan struct{})

	if _, err := futureClient.
		NewChangeLeverageService().
		Leverage(2).
		Symbol("BTCUSDT").
		Do(ctx); err != nil {
		log.Fatalf("failed to change initial leverage %v", err)
	}

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
			fmt.Println(err)
			return
		}

		wscChan <- doneC

		<-doneC
	}(
		futureClient,
	)

	log.Printf("debug 1")

	for wsContext := range wscChan {
		wscArr = append(wscArr, wsContext)
	}

	log.Printf("debug 2 %v", wscArr)

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
