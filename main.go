package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
)

type City struct {
	Name string
	Lat  float64
	Long float64
}

func (c City) OfficeURL() string {
	return fmt.Sprintf("https://api.weather.gov/points/%f,%f", c.Lat, c.Long)
}

func (c City) Fmt() string {
	return fmt.Sprintf("%s (%f, %f)", c.Name, c.Lat, c.Long)
}

func Seattle() City {
	return City{
		Name: "Seattle",
		Lat:  47.6062,
		Long: -122.3321,
	}
}

func Austin() City {
	return City{
		Name: "Austin",
		Lat:  30.2672,
		Long: -97.7431,
	}
}

func Cleveland() City {
	return City{
		Name: "Cleveland",
		Lat:  41.4993,
		Long: -81.6944,
	}
}

func Boston() City {
	return City{
		Name: "Boston",
		Lat:  42.3601,
		Long: -71.0589,
	}
}

func LosAngeles() City {
	return City{
		Name: "Los Angeles",
		Lat:  34.0522,
		Long: -118.2437,
	}
}

func Pittsburgh() City {
	return City{
		Name: "Pittsburgh",
		Lat:  40.4406,
		Long: -79.9959,
	}
}

func Hartford() City {
	return City{
		Name: "Hartford",
		Lat:  41.7658,
		Long: -72.6734,
	}
}

func CityNames() []string {
	return []string{
		"Seattle",
		"Austin",
		"Cleveland",
		"Hartford",
		"Boston",
		"Los Angeles",
		"Pittsburgh",
	}
}

func Cities() map[string]City {
	return map[string]City{
		"Seattle":     Seattle(),
		"Austin":      Austin(),
		"Cleveland":   Cleveland(),
		"Hartford":    Hartford(),
		"Boston":      Boston(),
		"Los Angeles": LosAngeles(),
		"Pittsburgh":  Pittsburgh(),
	}
}

type ForecastAPIResponse struct {
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

func selectCity() City {
	var options []huh.Option[City]
	var selected City

	for _, city := range Cities() {
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
	} else {
		forecast := ForecastAPIResponse{}

		data, err := io.ReadAll(r.Body)

		if err != nil {
			fmt.Printf("Failed to read response body: %s\n", err.Error())
			return ""
		}

		json.Unmarshal(data, &forecast)

		filename := fmt.Sprintf("%s-office.json", strings.ToLower(c.Name))
		f, err := os.Create(filename)

		if err != nil {
			fmt.Printf("Failed to create file %s: %s\n", filename, err.Error())

			return forecast.Properties.Forecast
		} else {
			defer f.Close()

			// Write the response body to the file
			_, err = io.Copy(f, bytes.NewReader(data))

			if err != nil {
				fmt.Printf("Failed to write to file %s: %s\n", filename, err.Error())
			}

			return forecast.Properties.Forecast
		}
	}
}

// func fetchForecast(uri string) {}

func main() {
	selected := selectCity()

	fmt.Println(
		"You picked",
		selected.Fmt(),
	)

	forecast := getForecastURL(selected)

	fmt.Println(forecast)
}
