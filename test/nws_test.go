package test

import (
	"fmt"
	"testing"

	"github.com/desertthunder/weather/internal/nws"
)

func TestConstants(t *testing.T) {
	cities := nws.Cities()

	for _, city := range nws.CityNames() {
		if _, ok := cities[city]; !ok {
			t.Errorf("City %s not found", city)
		}
	}
}

func TestCityComputedFields(t *testing.T) {
	city := nws.Seattle()

	t.Run("Generate office URL", func(t *testing.T) {
		url := city.OfficeURL()
		expected := fmt.Sprintf("https://api.weather.gov/points/%f,%f", city.Lat, city.Long)

		if url != expected {
			t.Errorf("Expected %s, got %s", expected, url)
		}
	})

	t.Run("Generate city string", func(t *testing.T) {
		city_str := city.Fmt()
		expected := fmt.Sprintf("%s (%f, %f)", city.Name, city.Lat, city.Long)

		if city_str != expected {
			t.Errorf("Expected %s, got %s", expected, city_str)
		}
	})
}
