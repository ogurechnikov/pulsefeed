package ingest

import (
	"context"
	"pulsefeed/internal/domain"
)

type Ingester interface {
	Trades(ctx context.Context, symbol string) (<-chan domain.Trade, <-chan error)
}
