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

func renderTrades(trades []domain.Trade) string {
	var lines []string
	for _, t := range trades {
		lines = append(lines, formatTradeLine(t))
	}
	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}
