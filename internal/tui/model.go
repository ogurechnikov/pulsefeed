package tui

import (
	"pulsefeed/internal/aggregator"
	"pulsefeed/internal/domain"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/NimbleMarkets/ntcharts/v2/linechart/streamlinechart"
)

type Model struct {
	snapshots     <-chan aggregator.Snapshot
	trades        []domain.Trade
	pricePoints   []aggregator.PricePoint
	chart         streamlinechart.Model
	lastPointTime time.Time
	width, height int
}

func NewModel(snapshots <-chan aggregator.Snapshot) Model {
	return Model{
		snapshots: snapshots,
		chart:     streamlinechart.New(80, 20),
	}
}

func (m Model) Init() tea.Cmd {
	return waitForSnapshot(m.snapshots)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case snapshotMsg:
		m.trades = msg.RecentTrades
		m.pricePoints = msg.PricePoints

		if len(msg.PricePoints) > 0 {
			latest := msg.PricePoints[len(msg.PricePoints)-1]
			if latest.Time.After(m.lastPointTime) {
				m.chart.Push(latest.Price)
				m.chart.Draw()
				m.lastPointTime = latest.Time
			}
		}
		return m, waitForSnapshot(m.snapshots)
	case snapshotsClosedMsg:
		return m, tea.Quit
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() tea.View {
	view := tea.NewView(m.chart.View())
	view.AltScreen = true
	return view
}
