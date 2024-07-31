// Submodule types for the nws package.
//
// Contains types and getter implementations for NWS API
// fetching and parsing.
package nws

import "fmt"

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
