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
	BuyTrades   []domain.Trade
	SellTrades  []domain.Trade
	PricePoints []PricePoint
}

type Aggregator struct {
	buyTrades   []domain.Trade
	sellTrades  []domain.Trade
	pricePoints []PricePoint
	recentPrice float64
}

func NewAggregator() *Aggregator {
	return &Aggregator{
		buyTrades:   make([]domain.Trade, 0, 100),
		sellTrades:  make([]domain.Trade, 0, 100),
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
				if trade.Side().IsBuy() {
					a.buyTrades = append(a.buyTrades, trade)
					if len(a.buyTrades) > 100 {
						a.buyTrades = a.buyTrades[1:]
					}
				} else {
					a.sellTrades = append(a.sellTrades, trade)
					if len(a.sellTrades) > 100 {
						a.sellTrades = a.sellTrades[1:]
					}
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
		BuyTrades:   slices.Clone(a.buyTrades),
		SellTrades:  slices.Clone(a.sellTrades),
		PricePoints: slices.Clone(a.pricePoints),
	}
}
