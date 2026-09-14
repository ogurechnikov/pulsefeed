package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"pulsefeed/internal/aggregator"
	"pulsefeed/internal/ingest/binance"
	"pulsefeed/internal/tui"

	tea "charm.land/bubbletea/v2"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	client := binance.NewClient("wss://stream.binance.com:9443/ws")

	trades, errs := client.Trades(ctx, "BTCUSDT")

	agg := aggregator.NewAggregator()
	snapshots := agg.Run(ctx, trades)

	go func() {
		select {
		case err, ok := <-errs:
			if ok {
				fmt.Println("connection error:", err)
				stop()
			}
		case <-ctx.Done():
		}
	}()

	model := tui.NewModel(snapshots)
	tea.NewProgram(model).Run()
}
