package view

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/desertthunder/weather/internal/nws"
)

func Table(headers []string, data [][]string) *table.Table {
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

type styles struct {
	Today         *lipgloss.Style
	Tonight       *lipgloss.Style
	Overnight     *lipgloss.Style
	Tomorrow      *lipgloss.Style
	TomorrowNight *lipgloss.Style
	Day           *lipgloss.Style
	Night         *lipgloss.Style
	City          *lipgloss.Style
}

func Styles() *styles {
	today := lipgloss.NewStyle().
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("42")). // Green
		Foreground(lipgloss.Color("0"))

	tonight := lipgloss.NewStyle().
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("92")). // High Intensity Green
		Foreground(lipgloss.Color("0"))

	overnight := lipgloss.NewStyle().
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("32")). // Bold Green
		Foreground(lipgloss.Color("0"))

	tomorrow := lipgloss.NewStyle().
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("86")). // Red
		Foreground(lipgloss.Color("0"))

	tomorrowNight := lipgloss.NewStyle().
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("204")). // Bold Red
		Foreground(lipgloss.Color("0"))
	day := lipgloss.NewStyle().
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("35")). // Blue
		Foreground(lipgloss.Color("0"))

	night := lipgloss.NewStyle().
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("37")). // Bold Blue
		Foreground(lipgloss.Color("0"))

	city := lipgloss.NewStyle().
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("86")). // Red
		Foreground(lipgloss.Color("0"))
	return &styles{&today, &tonight, &overnight, &tomorrow, &tomorrowNight, &day, &night, &city}
}

func ForecastLine(p nws.PeriodAPIResponse, v int) {
	// We have essentially three day categories: today, tomorrow, and
	// after tomorrow. We can use the start & end times to determine
	styles := Styles()
	var style *lipgloss.Style

	isToday := p.IsToday()
	isTomorrow := p.IsTomorrow()

	switch {
	case isToday && p.IsDaytime:
		style = styles.Today
	case isToday:
		style = styles.Tonight
	case isTomorrow && p.IsDaytime:
		style = styles.Tomorrow
	case isTomorrow:
		style = styles.TomorrowNight
	case p.IsDaytime:
		style = styles.Day
	default:
		style = styles.Night
	}

	tag := strings.ToUpper(p.Label)
	tag = style.Render(tag)

	switch v {
	case 0, 1:
		fmt.Printf("%s %s\n", tag, p.Temp())
	case 2:
		fmt.Printf("%s %s %s\n", tag, p.Temp(), p.ShortForecast)
	case 3:
		fmt.Printf("%s %s\n", tag, p.Temp())
		values := strings.Split(p.DetailedForecast, ". ")

		for _, v := range values {
			fmt.Println(v)
		}
	default:
		fmt.Printf("%s %s\n", tag, p.Temp())
	}
}

func CityLine(c *nws.City) {
	tag := Styles().City.Render("CITY")

	fmt.Printf("%s %s\n", tag, c.Fmt())
}
