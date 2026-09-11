package tui

import (
	"pulsefeed/internal/aggregator"
	"pulsefeed/internal/domain"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/NimbleMarkets/ntcharts/v2/linechart/streamlinechart"
)

type Model struct {
	snapshots        <-chan aggregator.Snapshot
	trades           []domain.Trade
	pricePoints      []aggregator.PricePoint
	chart            streamlinechart.Model
	lastPointTime    time.Time
	width, height    int
	chartWidth       int
	tradeWidth       int
	contentHeight    int
	rangeInitialized bool
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
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.chartWidth = m.width * 80 / 100
		m.tradeWidth = m.width - m.chartWidth
		m.contentHeight = m.height - 4

		m.chart = streamlinechart.New(m.chartWidth, m.contentHeight)
		if len(m.pricePoints) > 0 {
			min, max := m.pricePoints[0].Price, m.pricePoints[0].Price
			for _, p := range m.pricePoints {
				if p.Price < min {
					min = p.Price
				}
				if p.Price > max {
					max = p.Price
				}
			}
			m.chart.SetYRange(min*0.999, max*1.001)
			m.chart.SetViewYRange(min*0.999, max*1.001)
			m.rangeInitialized = true
		} else {
			m.rangeInitialized = false
		}
		for _, p := range m.pricePoints {
			m.chart.Push(p.Price)
		}
		m.chart.Draw()

	case snapshotMsg:
		m.trades = msg.RecentTrades
		m.pricePoints = msg.PricePoints

		if len(msg.PricePoints) > 0 {
			latest := msg.PricePoints[len(msg.PricePoints)-1]
			if latest.Time.After(m.lastPointTime) {
				m.lastPointTime = latest.Time
				if !m.rangeInitialized {
					m.chart.SetYRange(latest.Price*0.999, latest.Price*1.001)
					m.chart.SetViewYRange(latest.Price*0.999, latest.Price*1.001)
					m.rangeInitialized = true
				}
				m.chart.Push(latest.Price)
				m.chart.Draw()
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

	title := lipgloss.NewStyle().Bold(true).Render("PulseFeed")

	chartPanel := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Width(m.chartWidth).
		Height(m.contentHeight).
		Render(m.chart.View())

	tradesPanel := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Width(m.tradeWidth).
		Height(m.contentHeight).
		Render(renderTrades(m.trades))

	middle := lipgloss.JoinHorizontal(lipgloss.Top, chartPanel, tradesPanel)

	footer := lipgloss.NewStyle().Faint(true).Render("q: quit")

	content := lipgloss.JoinVertical(lipgloss.Left, title, middle, footer)

	view := tea.NewView(content)
	view.AltScreen = true

	return view
}
