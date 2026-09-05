package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"pulsefeed/internal/ingest/binance"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	client := binance.NewClient("wss://stream.binance.com:9443/ws")

	trades, errs := client.Trades(ctx, "BTCUSDT")

	for {
		select {
		case trade, ok := <-trades:
			if !ok {
				return
			}
			fmt.Println(trade)
		case err, ok := <-errs:
			if !ok {
				continue
			}
			fmt.Println("error: ", err)
			return
		}
	}
}
