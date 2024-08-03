package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/desertthunder/weather/internal/logger"
	"github.com/desertthunder/weather/internal/nws"
)

func TestConstants(t *testing.T) {
	cities := nws.Cities()

	for _, c := range nws.CityNames() {
		city, ok := cities[c]

		if !ok {
			t.Errorf("City %s not found", c)
		}

		t.Run(fmt.Sprintf("Computed fields for %s", c), func(t *testing.T) {
			t.Run("OfficeURL", func(t *testing.T) {
				url := city.OfficeURL()
				expected := fmt.Sprintf("https://api.weather.gov/points/%f,%f", city.Lat, city.Long)

				if url != expected {
					t.Errorf("Expected %s, got %s", expected, url)
				}
			})

			t.Run("Fmt", func(t *testing.T) {
				city_str := city.Fmt()
				expected := fmt.Sprintf("%s (%f, %f)", city.Name, city.Lat, city.Long)

				if city_str != expected {
					t.Errorf("Expected %s, got %s", expected, city_str)
				}
			})
		})
	}
}

func TestWeatherClient(t *testing.T) {

	t.Run("Base URL", func(t *testing.T) {
		want := "https://new_url.com"

		client := nws.NewWeatherClient()
		client.SetURL(want)
		got := client.BaseURL()

		if got != want {
			t.Errorf("Expected base URL to be %s, got %s", want, got)
		}
	})

	t.Run("GetWeather", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"properties": {"periods": [{"number": 1, "name": "Tonight", "startTime": "2024-08-02T06:00:00-05:00", "endTime": "2024-08-02T18:00:00-05:00", "isDaytime": true, "temperature": 98, "temperatureUnit": "F", "temperatureTrend": "", "probabilityOfPrecipitation": {"unitCode": "wmoUnit:percent", "value": null}, "windSpeed": "5 to 10 mph", "windDirection": "S", "icon": "https://api.weather.gov/icons/land/day/hot?size=medium", "shortForecast": "Sunny", "detailedForecast": "Sunny, with a high near 98. South wind 5 to 10 mph."}]}}`))
		}))

		defer server.Close()

		client := nws.NewWeatherClient()

		client.SetURL(server.URL)
		client.SetLogger(logger.Init())

		city := nws.Seattle()

		forecast, err := client.GetWeather(city)

		if err != nil {
			t.Fatalf("Expected no error, got %s", err.Error())
		}

		if len(forecast.Properties.Periods) < 1 {
			t.Errorf("Expected at least one period, got %d", len(forecast.Properties.Periods))
		}
	})
}
