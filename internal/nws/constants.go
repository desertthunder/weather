// Submodule constants for the nws package.
//
// Contains cities used as constants for the NWS API.
package nws

import "github.com/charmbracelet/lipgloss"

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
