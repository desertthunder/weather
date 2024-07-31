package test

import (
	"testing"

	"github.com/desertthunder/weather/internal/nws"
)

func TestCities(t *testing.T) {
	cities := nws.Cities()

	for _, city := range nws.CityNames() {
		if _, ok := cities[city]; !ok {
			t.Errorf("City %s not found", city)
		}
	}
}
