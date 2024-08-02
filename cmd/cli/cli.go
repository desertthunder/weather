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
			logger := config.log

			logger.Warn("Not implemented.")

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
		Flags: []cli.Flag{
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
			&cli.BoolFlag{
				Name: "verbose",
				Aliases: []string{
					"v",
					"extended",
					"e",
				},
				Usage: "Include forecast data beyond the next day.",
			},
			&cli.BoolFlag{
				Name: "interactive",
				Aliases: []string{
					"i",
				},
				Usage: "Interactive mode",
			},
		},
		Category: "Core",
		Action: func(ctx *cli.Context) error {
			logger := config.log
			logger.Debug("Forecast command invoked.")

			city := ctx.String("city")
			ip := ctx.String("ip")
			pt := ctx.StringSlice("pt")

			var n *nominatim.Nominatim

			if len(pt) > 0 || city != "" {
				n = nominatim.Client()
			}

			// Default to point, otherwise proceed to city.
			if len(pt) > 0 {
				lat, _ := strconv.ParseFloat(pt[0], 64)
				lng, _ := strconv.ParseFloat(pt[1], 64)

				logger.Debug(fmt.Sprintf("Set params to lat: %f, lon: %f", lat, lng))

				city, err := n.GeocodeByPoint(lat, lng)

				if err != nil {
					logger.Error(err.Error())
				} else {
					logger.Debug(fmt.Sprintf("Found: %s", city.Fmt()))
				}

				return err
			}

			// Geocode the city.
			if city != "" {
				logger.Debug(fmt.Sprintf("Set params to city: %s", city))

				city, err := n.GeocodeByCity(city)

				if err != nil {
					logger.Error(err.Error())
				} else {
					logger.Debug(fmt.Sprintf("Found: %s", city.Fmt()))
				}

				return err
			}

			// If no city or IP is provided, use the current user's IP.
			// Initialize the IPInfo client.
			ipc := ipinfo.NewIPInfoClient(config.Get("IPINFO_TOKEN"))

			// Instantiate the IPInfo client.
			// Get the city information via the IP.
			var loc ipinfo.IPInfoResponse
			var err error
			if ip != "" {
				logger.Debug(fmt.Sprintf("Set params to ip: %s", ip))

				loc, err = ipc.Geolocate(&ip)

			} else {
				logger.Debug("No IP address provided, will attempt to use device IP.")

				loc, err = ipc.Geolocate(nil)
			}

			if err != nil {
				logger.Error(err.Error())
				return err
			} else {
				city := loc.BuildCity()

				wea := nws.NewWeatherClient()
				wea.SetLogger(logger)
				logger.Debug(fmt.Sprintf("Found: %s", city.Fmt()))

				forecast, err := wea.GetWeather(city)

				if err != nil {
					logger.Error(err.Error())
					return err
				}

				if ctx.Bool("interactive") {
					interactive(city)

					return nil
				}

				for _, period := range forecast.Properties.Periods {
					period.View()
					if ctx.Bool("extended") {
						time.Sleep(time.Millisecond * 500)
					} else {
						break
					}
				}
			}

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

func commands() []*cli.Command {
	config := Config()
	return []*cli.Command{
		ForecastCommand(config),
		GeocodeCommand(config),
		InteractiveCommand(),
	}
}

// func application acts as a constant and is
// the entry point for the application.
func Application() *cli.App {
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
		Version:  "0.1.0",
		Compiled: time.Now(),
		Commands: commands(),
		Action: func(ctx *cli.Context) error {
			return nil
		},
	}
}
