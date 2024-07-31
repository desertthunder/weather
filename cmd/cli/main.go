package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/desertthunder/weather/internal/nws"
)

type City = nws.City

func selectCity() City {
	var options []huh.Option[City]
	var selected City

	for _, city := range nws.Cities() {
		options = append(options, huh.NewOption(city.Name, city))
	}

	f := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[City]().
				Title("Choose a city to fetch the weather for").
				Options(options...).
				Value(&selected),
		),
	)

	if err := f.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	return selected
}

func getForecastURL(c City) string {
	uri := c.OfficeURL()
	r, err := http.Get(uri)

	if err != nil {
		fmt.Printf("Request to %s failed with error: %s\n", uri, err.Error())

		return ""
	}

	office := nws.ForecastOfficeAPIResponse{}

	data, err := io.ReadAll(r.Body)

	if err != nil {
		fmt.Printf("Failed to read response body: %s\n", err.Error())
		return ""
	}

	json.Unmarshal(data, &office)

	return office.Properties.Forecast
}

// fetchForecast fetches the forecast data for the locale
// retrieved from the office endpoint and displays it in
// a table
func fetchForecast(uri string) model {
	r, err := http.Get(uri)

	if err != nil {
		fmt.Printf("Request to %s failed with error: %s\n", uri, err.Error())

		return model{}
	}

	defer r.Body.Close()

	forecast := nws.ForecastAPIResponse{}

	data, err := io.ReadAll(r.Body)

	if err != nil {
		fmt.Printf("Failed to read response body: %s\n", err.Error())

		return model{}
	}

	json.Unmarshal(data, &forecast)

	columns := []table.Column{
		{Title: "ID", Width: 3},
		{Title: "Label", Width: 15},
		{Title: "T", Width: 5},
		{Title: "P", Width: 5},
		{Title: "Wind", Width: 15},
		{Title: "Forecast", Width: 25},
	}

	forecasts := forecast.Properties.Periods
	rows := []table.Row{}

	for _, period := range forecasts {
		rows = append(rows, []string{
			fmt.Sprintf("%d", period.Number),
			period.Label,
			period.Temp(),
			period.Precipitation(),
			period.Wind(),
			period.ShortForecast,
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return model{table: t, forecasts: forecasts}
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table     table.Model
	forecasts []nws.PeriodAPIResponse
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	cleared := false

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			selected := m.table.SelectedRow()
			selected_id := selected[0]

			for _, period := range m.forecasts {
				fmt_num := fmt.Sprintf("%d", period.Number)
				time_period := strings.ToLower(period.Label)

				info := tea.Printf("You selected: %s's weather forecast (id: %s).", time_period, fmt_num)
				selected_forecast := tea.Printf("%s\n", period.DetailedForecast)

				if fmt_num == selected_id {
					if cleared {
						return m, tea.Sequence(info, selected_forecast)
					} else {
						return m, tea.Sequence(tea.ClearScreen, info, selected_forecast)
					}
				}
			}

			return m, tea.Batch(
				tea.Printf("The temperature is %s °F (%s)", selected[1], selected[2]),
			)
		}
	}

	m.table, cmd = m.table.Update(msg)

	return m, cmd
}

func (m model) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

func Root() {
	selected := selectCity()

	fmt.Println(
		"You picked",
		selected.Fmt(),
	)

	forecast := getForecastURL(selected)

	fmt.Println(forecast)

	m := fetchForecast(forecast)

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

// func main is the entrypoint for the CLI.
func main() {
	app := application()

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
