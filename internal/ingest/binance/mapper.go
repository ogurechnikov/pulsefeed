package binance

import (
	"fmt"
	"strconv"
	"time"

	"pulsefeed/internal/domain"
	"pulsefeed/internal/pkg/errs"
)

var (
	ErrInvalidPriceFormat = fmt.Errorf(
		"%w: price is not a valid number",
		errs.ErrValidation,
	)
	ErrInvalidQuantityFormat = fmt.Errorf(
		"%w: quantity is not a valid number",
		errs.ErrValidation,
	)
)

func mapTradeEvent(e binanceTradeEvent) (*domain.Trade, error) {
	symbol := e.Symbol

	price, err := strconv.ParseFloat(e.Price, 64)
	if err != nil {
		return nil, ErrInvalidPriceFormat
	}

	quantity, err := strconv.ParseFloat(e.Quantity, 64)
	if err != nil {
		return nil, ErrInvalidQuantityFormat
	}

	dealTime := time.UnixMilli(e.TradeTime)

	side := mapSide(e.IsBuyerMaker)

	return domain.NewTrade(
		symbol,
		price,
		quantity,
		dealTime,
		side,
	)
}

func mapSide(isBuyerMaker bool) domain.Side {
	if isBuyerMaker {
		return domain.SideSell
	}
	return domain.SideBuy
}
