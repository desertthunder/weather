// Submodule commands provides the definitions for the commands via the
// "geocast" application.
package cli

import (
	"fmt"
	"strconv"
	"time"

	"github.com/desertthunder/weather/internal/ipinfo"
	"github.com/desertthunder/weather/internal/nominatim"
	"github.com/desertthunder/weather/internal/nws"
	"github.com/desertthunder/weather/internal/view"
	"github.com/urfave/cli/v2"
)

func geocode(i *ipinfo.IPInfoClient, n *nominatim.Nominatim, ctx *cli.Context) *nws.City {
	pt := ctx.StringSlice("pt")
	c := ctx.String("city")
	ip := ctx.String("ip")

	var city *nws.City
	var err error
	var ipc ipinfo.IPInfoResponse

	if len(pt) > 0 {
		lat, _ := strconv.ParseFloat(pt[0], 64)
		lng, _ := strconv.ParseFloat(pt[1], 64)

		city, err = n.GeocodeByPoint(lat, lng)
	} else if c != "" {
		city, err = n.GeocodeByCity(c)
	}

	if err != nil {
		i.Log.Error(err.Error())

		return nil
	} else if city != nil {
		return city
	}

	if ip == "" {
		i.Log.Debug("No IP address provided, will attempt to use device IP.")

		ipc, err = i.Geolocate(nil)

	} else {
		i.Log.Debug(fmt.Sprintf("Set params to ip: %s", ip))

		ipc, err = i.Geolocate(&ip)
	}

	if err != nil {
		i.Log.Error(err.Error())

		return nil
	}

	cityV := ipc.BuildCity()

	return &cityV
}

// func forecast defines the shared functionality for the forecast command.
func forecast(city *nws.City, w *nws.WeatherClient, ctx *cli.Context) {
	forecast, err := w.GetWeather(*city)

	if err != nil {
		w.Log.Error(err.Error())

		return
	}

	v := ctx.Int("verbosity")

	w.Log.Debug(fmt.Sprintf("Verbosity level: %d", v))

	extended := ctx.Bool("extended")

	w.Log.Debug(fmt.Sprintf("Extended: %t", extended))

	for _, period := range forecast.Properties.Periods {
		view.ForecastLine(period, v)

		if extended {
			time.Sleep(time.Millisecond * 500)
		} else {
			break
		}
	}
}

// func ForecastCommand defines a pointer to the forecast command.
//
// Usage: geocast g[eocode] [-city]
//
// If no city is provided, the current IP address is used.
func GeocodeCommand(config *conf) *cli.Command {
	return &cli.Command{
		Name: "geocode",
		Aliases: []string{
			"g",
			"gc",
		},
		Category: "Core",
		Usage:    "Geocode a city or IP address, or reverse geocode a latitude and longitude.",
		Action: func(ctx *cli.Context) error {
			i := ipinfo.NewIPInfoClient(config.Get("IPINFO_TOKEN"))
			n := nominatim.Client()
			i.SetLogger(config.log)

			city := geocode(i, n, ctx)

			if city == nil {
				return nil
			}

			view.CityLine(city)

			return nil
		},
	}
}

// ForecastCommand defines a pointer to the forecast command.
func ForecastCommand(config *conf) *cli.Command {
	return &cli.Command{
		Name: "forecast",
		Aliases: []string{
			"f",
		},
		Usage:     "Fetch the weather forecast.",
		UsageText: "geocast f[orecast] [--c]ity [--i]p [--p]t",
		Args:      true,
		Flags:     flags(),
		Action: func(ctx *cli.Context) error {
			i := ipinfo.NewIPInfoClient(config.Get("IPINFO_TOKEN"))
			n := nominatim.Client()
			w := nws.NewWeatherClient()
			i.SetLogger(config.log)
			w.SetLogger(config.log)

			DefaultAction(i, n, w, ctx)

			return nil
		},
	}
}

// InteractiveCommand defines a pointer to the charm/bubble
// table-based interactive mode. It starts a bubbletea application
// that displays the weather forecast for a selected city.
func InteractiveCommand(config *conf) *cli.Command {
	return &cli.Command{
		Name:  "interactive",
		Usage: "Interactive mode",
		Aliases: []string{
			"i",
		},
		Category: "Core",
		Action: func(ctx *cli.Context) error {
			config.log.Debug("Interactive mode invoked.")

			Interactive()

			return nil
		},
	}
}
