package binance

import (
	"errors"
	"pulsefeed/internal/domain"
	"testing"
)

func TestMapTradeEvent(t *testing.T) {
	tests := []struct {
		name    string
		event   binanceTradeEvent
		wantErr error
	}{
		{
			name: "valid",
			event: binanceTradeEvent{
				Symbol:       "BTCUSDT",
				Price:        "80922.83000000",
				Quantity:     "0.11152000",
				TradeTime:    1788500943378,
				IsBuyerMaker: true,
			},
			wantErr: nil,
		},
		{
			name: "empty symbol",
			event: binanceTradeEvent{
				Symbol:       "",
				Price:        "80922.83000000",
				Quantity:     "0.11152000",
				TradeTime:    1788500943378,
				IsBuyerMaker: true,
			},
			wantErr: domain.ErrEmptySymbol,
		},
		{
			name: "zero price",
			event: binanceTradeEvent{
				Symbol:       "BTCUSD",
				Price:        "0",
				Quantity:     "0.11152000",
				TradeTime:    1788500943378,
				IsBuyerMaker: true,
			},
			wantErr: domain.ErrInvalidPrice,
		},
		{
			name: "negative price",
			event: binanceTradeEvent{
				Symbol:       "BTCUSD",
				Price:        "-80922.83000000",
				Quantity:     "0.11152000",
				TradeTime:    1788500943378,
				IsBuyerMaker: true,
			},
			wantErr: domain.ErrInvalidPrice,
		},
		{
			name: "zero quantity",
			event: binanceTradeEvent{
				Symbol:       "BTCUSD",
				Price:        "80922.83000000",
				Quantity:     "0",
				TradeTime:    1788500943378,
				IsBuyerMaker: true,
			},
			wantErr: domain.ErrInvalidQty,
		},
		{
			name: "negative quantity",
			event: binanceTradeEvent{
				Symbol:       "BTCUSD",
				Price:        "80922.83000000",
				Quantity:     "-1.5",
				TradeTime:    1788500943378,
				IsBuyerMaker: true,
			},
			wantErr: domain.ErrInvalidQty,
		},
		{
			name: "invalid price format",
			event: binanceTradeEvent{
				Symbol:       "BTCUSD",
				Price:        "not-a-number",
				Quantity:     "0.11152000",
				TradeTime:    1788500943378,
				IsBuyerMaker: true,
			},
			wantErr: ErrInvalidPriceFormat,
		},
		{
			name: "invalid quantity format",
			event: binanceTradeEvent{
				Symbol:       "BTCUSD",
				Price:        "80922.83000000",
				Quantity:     "not-a-number",
				TradeTime:    1788500943378,
				IsBuyerMaker: true,
			},
			wantErr: ErrInvalidQuantityFormat,
		},
		{
			name: "fractional values",
			event: binanceTradeEvent{
				Symbol:       "ETHUSDT",
				Price:        "0.00001",
				Quantity:     "0.001",
				TradeTime:    1788500943378,
				IsBuyerMaker: true,
			},
			wantErr: nil,
		},
		{
			name: "large values",
			event: binanceTradeEvent{
				Symbol:       "BTCUSDT",
				Price:        "999999.99",
				Quantity:     "1000.5",
				TradeTime:    1788500943378,
				IsBuyerMaker: true,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mapTradeEvent(tt.event)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("mapTradeEvent() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr != nil {
				if got != nil {
					t.Errorf("mapTradeEvent() got = %v, want nil", got)
				}
				return
			}

			if got == nil {
				t.Fatal("mapTradeEvent() got nil, want non-nil")
			}
		})
	}
}
