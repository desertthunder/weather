// Submodule types for the nws package.
//
// Contains types and getter implementations for NWS API
// fetching and parsing.
package nws

import (
	"fmt"
	"strconv"
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

// SetLatLong is a setter method for the city's latitude and longitude.
//
// It takes in two strings representing the latitude and longitude
// and sets the corresponding fields in the city struct.
func (c *City) setLatLong(lat, long string) {
	c.Lat, _ = strconv.ParseFloat(lat, 64)
	c.Long, _ = strconv.ParseFloat(long, 64)
}

func BuildCity(name, lat, long string) City {
	city := City{Name: name}
	city.setLatLong(lat, long)

	return city
}
