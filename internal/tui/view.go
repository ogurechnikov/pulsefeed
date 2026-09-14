package tui

import (
	"fmt"
	"pulsefeed/internal/domain"

	"charm.land/lipgloss/v2"
)

func formatTradeLine(t domain.Trade) string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Green)
	label := "BUY"
	if t.Side().IsSell() {
		style = lipgloss.NewStyle().Foreground(lipgloss.Red)
		label = "SELL"
	}
	return style.Render(fmt.Sprintf("%s %.2f", label, t.Price()))
}

func renderTradeBlock(title string, trades []domain.Trade, limit int) string {
	if limit < 0 {
		limit = 0
	}

	if len(trades) > limit {
		trades = trades[len(trades)-limit:]
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Green)
	if title == "SELL" {
		titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Red)
	}

	lines := []string{titleStyle.Render(title)}

	for _, t := range trades {
		lines = append(lines, formatTradeLine(t))
	}

	for len(lines) < limit+1 {
		lines = append(lines, "")
	}
	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}
