package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/desertthunder/weather/internal/nws"
	"github.com/desertthunder/weather/internal/view"
)

type test struct {
	name string
	data [][]string
}

func TestTable(t *testing.T) {
	headers := []string{"City", "Latitude", "Longitude"}
	tests := []test{
		{
			name: "Single row",
			data: [][]string{{"Seattle", "47.6062", "-122.3321"}},
		},
		{
			name: "Multiple rows",
			data: [][]string{
				{"Seattle", "47.6062", "-122.3321"},
				{"Portland", "45.5152", "-122.6784"},
			},
		},
		{
			name: "No data",
			data: [][]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			captured := CaptureOutput(func() {
				tbl := view.Table(headers, tt.data)
				fmt.Println(tbl.Render())
			})

			for _, header := range headers {
				if !strings.Contains(captured, header) {
					t.Errorf("Expected header %s not found in output %s", header, captured)
				}
			}

			for _, row := range tt.data {
				for _, cell := range row {
					if !strings.Contains(captured, cell) {
						t.Errorf("Expected cell %s not found in output %s", cell, captured)
					}
				}
			}
		})
	}
}

func TestLines(t *testing.T) {
	t.Run("CityLine", func(t *testing.T) {
		city := nws.Seattle()

		captured := CaptureOutput(func() {
			view.CityLine(&city)
		})

		if !strings.Contains(captured, city.Name) {
			t.Errorf("Expected city name not found in output %s", captured)
		}
	})

	t.Run("ForecastLine", func(t *testing.T) {
		period := nws.PeriodAPIResponse{
			Number:                     1,
			Label:                      "Tonight",
			StartTime:                  "2024-08-02T06:00:00-05:00",
			EndTime:                    "2024-08-02T18:00:00-05:00",
			IsDaytime:                  true,
			Temperature:                98,
			TemperatureUnit:            "F",
			TemperatureTrend:           "",
			ProbabilityOfPrecipitation: nws.ProbabilityOfPrecipitation{},
			WindSpeed:                  "5 to 10 mph",
			WindDirection:              "S",
			Icon:                       "https://api.weather.gov/icons/land/day/hot?size=medium",
			ShortForecast:              "Sunny",
			DetailedForecast:           "Sunny, with a high near 98. South wind 5 to 10 mph.",
		}

		t.Run("Verbosity 0", func(t *testing.T) {
			want := period.Temp()

			buf := CaptureOutput(func() {
				view.ForecastLine(period, 0)
			})

			if !strings.Contains(buf, want) {
				t.Errorf("Expected %s not found in output %s", want, buf)
			}
		})

		t.Run("Verbosity 1", func(t *testing.T) {
			want := period.Temp()

			buf := CaptureOutput(func() {
				view.ForecastLine(period, 1)
			})

			if !strings.Contains(buf, want) {
				t.Errorf("Expected %s not found in output %s", want, buf)
			}
		})
		t.Run("Verbosity 2", func(t *testing.T) {
			want := period.ShortForecast

			buf := CaptureOutput(func() {
				view.ForecastLine(period, 2)
			})

			if !strings.Contains(buf, want) {
				t.Errorf("Expected %s not found in output %s", want, buf)
			}
		})

		t.Run("Verbosity 3", func(t *testing.T) {
			want := strings.Split(period.DetailedForecast, ". ")

			buf := CaptureOutput(func() {
				view.ForecastLine(period, 3)
			})

			if !strings.Contains(buf, want[0]) {
				t.Errorf("Expected %s not found in output %s", want[0], buf)
			}

			if !strings.Contains(buf, want[1]) {
				t.Errorf("Expected %s not found in output %s", want[1], buf)
			}
		})
	})
}
