// Submodule flags provides the definitions for the flags used by the "geocast"
// application.
package cli

import "github.com/urfave/cli/v2"

func cityFlag() cli.Flag {
	return &cli.StringFlag{
		Name:    "city",
		Aliases: []string{"c", "n", "cn"},
		Usage:   "The city name to fetch the forecast for.",
	}
}

func ipFlag() cli.Flag {
	return &cli.StringFlag{
		Name:  "ip",
		Usage: "The IP address to fetch the forecast for.",
	}
}

func pointFlag() cli.Flag {
	return &cli.StringSliceFlag{
		Name:    "pt",
		Aliases: []string{"p"},
		Usage:   "The point to fetch the forecast for (lat,lon).",
	}
}

func verbosityFlag() cli.Flag {
	return &cli.IntFlag{
		Name: "verbosity",
		Aliases: []string{
			"v",
		},
		Usage: `Verbosity level. 0 - Default, 1 - Temperature,
		2 - Short Forecast, 3 - Detailed Forecast.`,
		Value: 0,
	}
}

func extendedFlag() cli.Flag {
	return &cli.BoolFlag{
		Name: "extended",
		Aliases: []string{
			"e",
		},
		Usage: "Include forecast data beyond the next day.",
	}
}

func interactiveFlag() cli.Flag {
	return &cli.BoolFlag{
		Name: "interactive",
		Aliases: []string{
			"i",
		},
		Usage: "Interactive mode (not *quite* implemented)",
	}
}

func flags() []cli.Flag {
	return []cli.Flag{
		cityFlag(),
		ipFlag(),
		pointFlag(),
		verbosityFlag(),
		extendedFlag(),
		interactiveFlag(),
	}
}
