package domain

import (
	"fmt"
	"pulsefeed/internal/pkg/errs"
	"time"
)

var (
	ErrEmptySymbol = fmt.Errorf(
		"%w: symbol is empty",
		errs.ErrValidation)
	ErrInvalidPrice = fmt.Errorf(
		"%w: price must be greater than zero",
		errs.ErrValidation)
	ErrInvalidQty = fmt.Errorf(
		"%w: quantity must be greater than zero",
		errs.ErrValidation)
	ErrUnknownSide = fmt.Errorf(
		"%w: side is unknown",
		errs.ErrValidation)
)

type Trade struct {
	symbol   string
	price    float64
	quantity float64
	dealTime time.Time
	side     Side
}

func NewTrade(
	symbol string,
	price float64,
	quantity float64,
	dealtime time.Time,
	side Side,
) (*Trade, error) {
	if symbol == "" {
		return nil, ErrEmptySymbol
	}
	if price <= 0 {
		return nil, ErrInvalidPrice
	}
	if quantity <= 0 {
		return nil, ErrInvalidQty
	}
	if side == SideUnknown {
		return nil, ErrUnknownSide
	}

	trade := &Trade{
		symbol:   symbol,
		price:    price,
		quantity: quantity,
		dealTime: dealtime,
		side:     side,
	}
	return trade, nil
}

// Геттеры состояния

func (t *Trade) Symbol() string {
	return t.symbol
}

func (t *Trade) Price() float64 {
	return t.price
}

func (t *Trade) Quantity() float64 {
	return t.quantity
}

func (t *Trade) DealTime() time.Time {
	return t.dealTime
}

func (t *Trade) Side() Side {
	return t.side
}
