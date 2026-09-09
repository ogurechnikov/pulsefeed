package aggregator

import (
	"context"
	"pulsefeed/internal/domain"
	"slices"
	"time"
)

type PricePoint struct {
	Time  time.Time
	Price float64
}

type Snapshot struct {
	RecentTrades []domain.Trade
	PricePoints  []PricePoint
}

type Aggregator struct {
	trades      []domain.Trade
	pricePoints []PricePoint
	recentPrice float64
}

func NewAggregator() *Aggregator {
	return &Aggregator{
		trades:      make([]domain.Trade, 0, 20),
		pricePoints: make([]PricePoint, 0, 60),
	}
}

func (a *Aggregator) Run(
	ctx context.Context,
	trades <-chan domain.Trade,
) <-chan Snapshot {
	snapshots := make(chan Snapshot)

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case trade := <-trades:
				a.trades = append(a.trades, trade)
				if len(a.trades) > 20 {
					a.trades = a.trades[1:]
				}
				a.recentPrice = trade.Price()
				snapshots <- a.snapshot()
			case <-ticker.C:
				if a.recentPrice == 0 {
					continue
				}
				a.pricePoints = append(
					a.pricePoints,
					PricePoint{
						Time:  time.Now(),
						Price: a.recentPrice,
					},
				)
				if len(a.pricePoints) > 60 {
					a.pricePoints = a.pricePoints[1:]
				}
				snapshots <- a.snapshot()
			case <-ctx.Done():
				return
			}
		}
	}()

	return snapshots
}

func (a *Aggregator) snapshot() Snapshot {
	return Snapshot{
		RecentTrades: slices.Clone(a.trades),
		PricePoints:  slices.Clone(a.pricePoints),
	}
}
