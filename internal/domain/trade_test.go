package domain

import (
	"errors"
	"testing"
	"time"
)

func TestNewTrade(t *testing.T) {
	validTime := time.Now()

	tests := []struct {
		name     string
		symbol   string
		price    float64
		quantity float64
		dealTime time.Time
		side     string
		wantErr  error
	}{
		{name: "valid trade", symbol: "BTCUSDT", price: 50000.0, quantity: 1.5, dealTime: validTime, side: "buy", wantErr: nil},
		{name: "empty symbol", symbol: "", price: 50000.0, quantity: 1.5, dealTime: validTime, side: "buy", wantErr: ErrEmptySymbol},
		{name: "zero price", symbol: "BTCUSDT", price: 0, quantity: 1.5, dealTime: validTime, side: "buy", wantErr: ErrInvalidPrice},
		{name: "negative price", symbol: "BTCUSDT", price: -100.0, quantity: 1.5, dealTime: validTime, side: "buy", wantErr: ErrInvalidPrice},
		{name: "zero quantity", symbol: "BTCUSDT", price: 50000.0, quantity: 0, dealTime: validTime, side: "buy", wantErr: ErrInvalidQty},
		{name: "negative quantity", symbol: "BTCUSDT", price: 50000.0, quantity: -1.5, dealTime: validTime, side: "buy", wantErr: ErrInvalidQty},
		{name: "empty side", symbol: "BTCUSDT", price: 50000.0, quantity: 1.5, dealTime: validTime, side: "", wantErr: ErrEmptySide},
		{name: "all fields invalid", symbol: "", price: 0, quantity: 0, dealTime: validTime, side: "", wantErr: ErrEmptySymbol},
		{name: "fractional values", symbol: "ETHUSDT", price: 0.00001, quantity: 0.001, dealTime: validTime, side: "sell", wantErr: nil},
		{name: "large values", symbol: "BTCUSDT", price: 999999.99, quantity: 1000.5, dealTime: validTime, side: "buy", wantErr: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTrade(
				tt.symbol,
				tt.price,
				tt.quantity,
				tt.dealTime,
				tt.side,
			)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewTrade() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr != nil {
				if got != nil {
					t.Errorf("NewTrade() got = %v, want nil", got)
				}
				return
			}

			if got == nil {
				t.Fatal("NewTrade() got nil, want non-nil")
			}

			if got.Symbol() != tt.symbol {
				t.Errorf("Symbol() = %v, want %v", got.Symbol(), tt.symbol)
			}
			if got.Price() != tt.price {
				t.Errorf("Price() = %v, want %v", got.Price(), tt.price)
			}
			if got.Quantity() != tt.quantity {
				t.Errorf("Quantity() = %v, want %v", got.Quantity(), tt.quantity)
			}
			if !got.DealTime().Equal(tt.dealTime) {
				t.Errorf("DealTime() = %v, want %v", got.DealTime(), tt.dealTime)
			}
			if got.Side() != tt.side {
				t.Errorf("Side() = %v, want %v", got.Side(), tt.side)
			}
		})
	}
}
