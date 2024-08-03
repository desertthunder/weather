// Submodule cli provides the definitions for the cli commands via the
// "geocast" application.
package cli

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/desertthunder/weather/internal/ipinfo"
	"github.com/desertthunder/weather/internal/nominatim"
	"github.com/desertthunder/weather/internal/nws"
	"github.com/desertthunder/weather/internal/view"
	"github.com/urfave/cli/v2"
)

// DefaultAction is the function called when no arguments are provided or when
// the "me" argument is provided.
//
// The default action is to first geocode the current device's IP address and
// then fetch the weather forecast for the city.
func DefaultAction(i *ipinfo.IPInfoClient, n *nominatim.Nominatim, nwsc *nws.WeatherClient, ctx *cli.Context) {
	city := geocode(i, n, ctx)

	if city == nil {
		err := errors.New("no results found for the provided city name")

		i.Log.Error(err.Error())

		return
	}

	view.CityLine(city)

	forecast(city, nwsc, ctx)
}

// func Application acts as a constant and is the entry point for the application.
func Application() *cli.App {
	config := Config()

	logger := config.log

	return &cli.App{
		Name:     "geocast",
		HelpName: "geocast (Geo[coding] + [Fore]cast)",
		Usage:    "Location aware weather forecasts for the command line.",
		UsageText: `geocast f[orecast] [--c]ity [--ip] [--p]t [--i]nteractive
geocast g[eocode] [--c]ity [--ip] [--p]t
geocast i[nteractive]`,
		Description: `Geocast is a command line utility that provides location aware weather forecasts.
It can be used to fetch the weather forecast for a specific city, latitude and
longitude, or the current device's IP address.`,
		Compiled: time.Now(),
		// Global flags for the application , i.e. the flags that apply to all commands.
		//
		// City, IP, and Point flags
		Flags: flags(),
		Commands: []*cli.Command{
			ForecastCommand(config),
			GeocodeCommand(config),
			InteractiveCommand(config),
		},
		Action: func(ctx *cli.Context) error {
			arg := ctx.Args().First()
			flags := ctx.FlagNames()

			logger.Debug(fmt.Sprintf("Flags: %s", flags))
			logger.Debug(fmt.Sprintf("Arg: %s", arg))

			ipc := ipinfo.NewIPInfoClient(config.Get("IPINFO_TOKEN"))
			n := nominatim.Client()
			nwsc := nws.NewWeatherClient()

			nwsc.SetLogger(logger)
			ipc.SetLogger(logger)

			app := ctx.Bool("interactive")

			if strings.ToLower(arg) == "me" || ctx.Args().Len() == 0 {
				logger.Debug("Default command invoked.")

				if app {
					city := geocode(ipc, n, ctx)

					if city == nil {
						err := errors.New("no results found for the provided city name")

						logger.Error(err.Error())

						return err
					} else {
						view.CityLine(city)

						interactive(*city)

						return nil
					}
				}

				DefaultAction(ipc, n, nwsc, ctx)

				return nil
			}

			// Geocode the city.
			city, err := n.GeocodeByCity(arg)

			if err != nil {
				logger.Error(err.Error())

				return err
			} else if city != nil {
				err = errors.New("no results found for the provided city name")

				logger.Error(err.Error())

				return err
			}

			view.CityLine(city)

			forecast(city, nwsc, ctx)

			return nil
		},
	}
}
