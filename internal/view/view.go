package view

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func Table(
	headers []string,
	data [][]string,
) *table.Table {
	re := lipgloss.NewRenderer(os.Stdout)
	baseStyle := lipgloss.NewStyle().Padding(0, 1)
	headerStyle := baseStyle.Foreground(lipgloss.Color("#005fd7")).Bold(true)
	oddStyle := baseStyle.Foreground(lipgloss.Color("252"))
	evenStyle := baseStyle.Foreground(lipgloss.Color("245"))

	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(re.NewStyle().Foreground(lipgloss.Color("238"))).
		Headers(headers...).
		Width(48).
		Rows(data...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == 0 {
				return headerStyle
			}

			if row%2 == 0 {
				return evenStyle
			}

			return oddStyle
		})

	return t
}
