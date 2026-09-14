# PulseFeed

A real-time terminal dashboard for crypto market data, built with Go and Bubble Tea.

PulseFeed connects to Binance's public WebSocket trade stream, aggregates
incoming trades, and renders a live price chart alongside a running feed of
recent buy/sell orders — all inside your terminal.

## Features

- Live price chart (updated every 5 seconds) with automatic Y-axis scaling
- Separate BUY/SELL trade feeds, color-coded (green/red)
- Graceful shutdown on Ctrl+C or connection loss
- Responsive layout that adapts to terminal size (down to a minimum
  standardized height)

## Architecture

Binance WebSocket
↓
internal/ingest/binance — WS client, DTO, mapping to domain.Trade
↓
internal/domain — Trade, Side value object (validated, tested)
↓
internal/aggregator — ring buffers, price sampling, Snapshot channel
↓
internal/tui — Bubble Tea dashboard (chart + trade feeds)


Each layer only knows about the one below it — the domain layer has no
knowledge of Binance or WebSockets, and the TUI layer has no business
logic of its own; it simply renders whatever `Snapshot` the aggregator
sends.

## Running

```bash
go run ./cmd/pulsefeed
```

Press `q` or `Ctrl+C` to quit.

## Development approach

The domain layer (`Trade`, `Side`) was built test-first (TDD): tests were
written before the implementation for each unit of behavior, following a
red-green-refactor cycle.

## Known limitations

- No automatic reconnection if the WebSocket connection drops — the
  program currently logs the error and exits. Retry-with-backoff is a
  planned follow-up.
- `internal/ingest/binance` (the WebSocket client) and
  `internal/aggregator` (ring buffers and timers) are not yet covered by
  unit tests — this was a conscious time trade-off during initial
  development. Domain logic (`Trade`, `Side`) is fully tested.
- Currently supports a single trading pair (BTCUSDT), hardcoded.

## Tech stack

- [Bubble Tea v2](https://charm.land/bubbletea) — TUI framework
- [Lip Gloss v2](https://charm.land/lipgloss) — styling and layout
- [ntcharts](https://github.com/NimbleMarkets/ntcharts) — terminal charts
- [coder/websocket](https://github.com/coder/websocket) — WebSocket client

## License

MIT