package binance

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"pulsefeed/internal/domain"
	"strings"

	"github.com/coder/websocket"
)

type Client struct {
	baseURL string
}

func NewClient(baseURL string) *Client {
	return &Client{baseURL: baseURL}
}

func (c *Client) Trades(ctx context.Context, symbol string,
) (<-chan domain.Trade, <-chan error) {
	trades := make(chan domain.Trade)
	errs := make(chan error)
	url := fmt.Sprintf("%s/%s@trade", c.baseURL, strings.ToLower(symbol))

	go func() {
		conn, _, err := websocket.Dial(ctx, url, nil)
		if err != nil {
			errs <- err
			return
		}
		defer conn.CloseNow()
		defer close(errs)
		defer close(trades)

		for {
			_, data, err := conn.Read(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) ||
					errors.Is(err, context.DeadlineExceeded) {
					return
				}
				errs <- err
				return
			}

			var event binanceTradeEvent
			if err := json.Unmarshal(data, &event); err != nil {
				errs <- err
				return
			}

			trade, err := mapTradeEvent(event)
			if err != nil {
				errs <- err
				return
			}

			trades <- *trade
		}
	}()

	return trades, errs
}
