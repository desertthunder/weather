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
	"strings"
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
