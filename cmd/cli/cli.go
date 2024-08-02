// Submodule cli provides the definitions for the cli commands
// via the "geocast" application.
package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/desertthunder/weather/internal/ipinfo"
	"github.com/desertthunder/weather/internal/nominatim"
	"github.com/desertthunder/weather/internal/nws"
	"github.com/desertthunder/weather/internal/view"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

// struct conf is a wrapper around the viper package,
// used for configuration management.
type conf struct {
	v   *viper.Viper
	log *log.Logger
}

// struct color is an object mapping a log level
// to a color.
type color struct {
	level log.Level
	color string
}

// Convert log level to the string representation.
//
// ex. log.DebugLevel -> "DEBUG"
func (c color) String() string {
	return strings.ToUpper(c.level.String())
}

func colors() []color {
	return []color{
		{log.DebugLevel, "63"},
		{log.InfoLevel, "86"},
		{log.WarnLevel, "192"},
		{log.ErrorLevel, "204"},
		{log.FatalLevel, "134"},
	}
}

// func Config is the conf constructor.
//
// It reads the configuration from the .env file and returns
// a pointer to an instance of conf.
func Config() *conf {
	v := viper.New()
	v.SetConfigFile(".env")

	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s\n", err)
	}

	c := &conf{v: v}

	c.initLogger()

	return c
}

// Get is the accessor method for any
// configuration values/environment variables.
func (c *conf) Get(key string) string {
	return c.v.GetString(key)
}

// initLogger initializes the logger with a set of
// default styles and colors while also streaming
// the logs to the console.
//
// TODO: Add support for file logging.
func (c *conf) initLogger() {
	styles := log.DefaultStyles()
	logger := log.New(os.Stdout)

	for _, item := range colors() {
		styles.Levels[item.level] = lipgloss.NewStyle().
			SetString(item.String()).
			Padding(0, 1, 0, 1).
			Background(lipgloss.Color(item.color)).
			Foreground(lipgloss.Color("0"))
	}

	logger.SetStyles(styles)

	if os.Getenv("DEBUG") != "" {
		logger.SetLevel(log.DebugLevel)
	}

	c.log = logger
}

func flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "city",
			Aliases: []string{"c", "n", "cn"},
			Usage:   "The city name to fetch the forecast for.",
		},
		&cli.StringFlag{
			Name:  "ip",
			Usage: "The IP address to fetch the forecast for.",
		},
		&cli.StringSliceFlag{
			Name:    "pt",
			Aliases: []string{"p"},
			Usage:   "The point to fetch the forecast for (lat,lon).",
		},
		&cli.IntFlag{
			Name: "verbosity",
			Aliases: []string{
				"v",
			},
			Usage: "Verbosity level. 0 - Default, 1 - Temperature, 2 - Short Forecast, 3 - Detailed Forecast",
			Value: 0,
		},
		&cli.BoolFlag{
			Name: "extended",
			Aliases: []string{
				"e",
			},
			Usage: "Include forecast data beyond the next day.",
		},
		&cli.BoolFlag{
			Name: "interactive",
			Aliases: []string{
				"i",
			},
			Usage: "Interactive mode (not *quite* implemented)",
		},
	}
}

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
func InteractiveCommand() *cli.Command {
	return &cli.Command{
		Name:  "interactive",
		Usage: "Interactive mode",
		Aliases: []string{
			"i",
		},
		Category: "Core",
		Action: func(ctx *cli.Context) error {
			Interactive()
			return nil
		},
	}
}

// DefaultAction is the function called when no arguments are provided or when
// the "me" argument is provided.
//
// The default action is to first geocode the current device's IP address and
// then fetch the weather forecast for the city.
func DefaultAction(i *ipinfo.IPInfoClient, n *nominatim.Nominatim, nwsc *nws.WeatherClient, ctx *cli.Context) {
	city := geocode(i, n, ctx)

	if city == nil {
		return
	}

	forecast(city, nwsc, ctx)
}

// func application acts as a constant and is
// the entry point for the application.
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
			InteractiveCommand(),
		},
		Action: func(ctx *cli.Context) error {
			arg := ctx.Args().First()
			flags := ctx.FlagNames()

			logger.Debug(fmt.Sprintf("Flags: %s", flags))

			logger.Debug(fmt.Sprintf("Arg: %s", arg))

			// arg can be "me" or a city name.

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
						return nil
					} else {
						view.CityLine(city)

						interactive(*city)

						return nil
					}
				} else {
					DefaultAction(ipc, n, nwsc, ctx)
				}

				return nil
			}
			// Geocode the city.
			city, err := n.GeocodeByCity(arg)

			if err != nil {
				logger.Error(err.Error())
				return err
			}

			if city != nil {
				view.CityLine(city)
			}

			weather, err := nwsc.GetWeather(*city)

			if err != nil {
				logger.Error(err.Error())
				return err
			}

			for _, period := range weather.Properties.Periods {
				view.ForecastLine(period, ctx.Int("verbosity"))

				if ctx.Bool("extended") {
					time.Sleep(time.Millisecond * 500)
				} else {
					break
				}
			}
			return nil
		},
	}
}
