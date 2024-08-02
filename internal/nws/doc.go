// weather.gov: https://www.weather.gov/documentation/services-web-api
//
// weather.gov provides a RESTful API for accessing weather data.
//
// To get offices: https://api.weather.gov/points/30.2672,-97.7431 (Austin, TX)
//
// From that get the forecast: https://api.weather.gov/gridpoints/AWX/125,68/forecast
package nws

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

type PeriodAPIResponse struct {
	Number                     int    `json:"number"`
	Label                      string `json:"name"`
	StartTime                  string `json:"startTime"`
	EndTime                    string `json:"endTime"`
	IsDaytime                  bool   `json:"isDaytime"`
	Temperature                int    `json:"temperature"`
	TemperatureUnit            string `json:"temperatureUnit"`
	TemperatureTrend           string `json:"temperatureTrend"`
	ProbabilityOfPrecipitation struct {
		UnitCode string `json:"unitCode"`
		Value    int    `json:"value"`
	} `json:"probabilityOfPrecipitation"`
	WindSpeed        string `json:"windSpeed"`
	WindDirection    string `json:"windDirection"`
	Icon             string `json:"icon"`
	ShortForecast    string `json:"shortForecast"`
	DetailedForecast string `json:"detailedForecast"`
}

type ElevationAPIResponse struct {
	UnitCode string `json:"unitCode"`
	Value    int    `json:"value"`
}

type GeometryAPIResponse struct {
	Type        string    `json:"type"`
	Coordinates []float32 `json:"coordinates"`
}

type ForecastAPIResponse struct {
	Properties struct {
		Periods []PeriodAPIResponse `json:"periods"`
	} `json:"properties"`
}

func (p PeriodAPIResponse) Wind() string {
	return fmt.Sprintf("%s %s", p.WindSpeed, p.WindDirection)
}

func (p PeriodAPIResponse) Precipitation() string {
	unit := p.ProbabilityOfPrecipitation.UnitCode
	unit = strings.TrimPrefix(unit, "wmoUnit:")

	if unit == "percent" {
		unit = "%"
	}

	return fmt.Sprintf("%d%s", p.ProbabilityOfPrecipitation.Value, unit)
}

func (p PeriodAPIResponse) Temp() string {
	unit := p.TemperatureUnit
	unit = strings.TrimPrefix(unit, "wmoUnit:")
	unit = fmt.Sprintf("°%s", unit)

	return fmt.Sprintf("%d%s", p.Temperature, unit)
}

// StartTime is of the format 2024-08-02T06:00:00-05:00
func (p PeriodAPIResponse) IsToday() bool {
	start, _ := time.Parse(time.RFC3339, p.StartTime)

	return time.Now().Weekday().String() == start.Weekday().String()
}

func (p PeriodAPIResponse) IsTomorrow() bool {
	start, _ := time.Parse(time.RFC3339, p.StartTime)

	return time.Now().Add(time.Hour*24).Weekday().String() == start.Weekday().String()
}

func (e ElevationAPIResponse) Fmt() string {
	unit := e.UnitCode
	unit = strings.TrimPrefix(unit, "wmoUnit:")

	return fmt.Sprintf("%d%s", e.Value, unit)
}

type ForecastOfficeAPIResponse struct {
	URL      string `json:"id"`
	Type     string `json:"type"`
	Geometry struct {
		Type        string    `json:"type"`
		Coordinates []float32 `json:"coordinates"`
	} `json:"geometry"`
	Properties struct {
		Forecast         string `json:"forecast"`
		ForecastHourly   string `json:"forecastHourly"`
		ForecastGridData string `json:"forecastGridData"`
		ForecastZone     string `json:"forecastZone"`
		Timezone         string `json:"timeZone"`
		County           string `json:"county"`
		FireWeatherZone  string `json:"fireWeatherZone"`
		Stations         string `json:"observationStations"`
		RadarStation     string `json:"radarStation"`
	} `json:"properties"`
}

func (f ForecastOfficeAPIResponse) ForecastURL() string {
	return f.Properties.Forecast
}

type styles struct {
	Day     *lipgloss.Style
	Night   *lipgloss.Style
	Today   *lipgloss.Style
	Tonight *lipgloss.Style
}

func Styles() *styles {
	day := lipgloss.NewStyle().
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("35")).
		Foreground(lipgloss.Color("0"))

	night := lipgloss.NewStyle().
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("37")).
		Foreground(lipgloss.Color("0"))

	today := lipgloss.NewStyle().
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("86")).
		Foreground(lipgloss.Color("0"))

	tonight := lipgloss.NewStyle().
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("204")).
		Foreground(lipgloss.Color("0"))

	return &styles{&day, &night, &today, &tonight}
}

func (p PeriodAPIResponse) View() {
	st := Styles()

	days := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	pday := strings.Split(p.Label, " ")[0]
	curr := !slices.Contains(days, pday)
	label := strings.ToUpper(p.Label)
	today := time.Now().Weekday().String()

	if curr && p.IsDaytime {
		label = fmt.Sprintf("%s (%s)", label, strings.ToUpper(today))
		label = st.Today.Render(label)
	} else if curr {
		label = fmt.Sprintf("%s (%s)", label, strings.ToUpper(today))
		label = st.Tonight.Render(label)
	} else if p.IsDaytime {
		label = st.Day.Render(label)
	} else {
		label = st.Night.Render(label)
	}

	fmt.Printf("%s %s\n", label, p.Temp())

	if curr {
		detailedForecast := strings.Split(p.DetailedForecast, ". ")

		for _, d := range detailedForecast {
			fmt.Println(d)
		}
	} else {
		fmt.Println(p.ShortForecast)
	}

}
