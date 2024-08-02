// Submodule cli provides the definitions for the cli commands
// via the "geocast" application.
package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"

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

type color struct {
	level log.Level
	color string
}

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
		Usage: "Fetch the weather forecast.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "city",
				Aliases: []string{"c", "n", "cn"},
				Usage:   "The city name to fetch the forecast for.",
			},
		},
		Action: func(ctx *cli.Context) error {
			logger := config.log
			osm := nominatim.Init()
			ua := config.Get("NOMINATIM_USER_AGENT")

			if ua != "" {
				osm.SetUserAgent(ua)
			}

			logger.Debug(fmt.Sprintf("Initialized client at: %s", osm.GetBaseURL()))

			c := ctx.String("city")

			if c == "" {
				fmt.Println("Please provide a city name to fetch the forecast for.")

				return nil
			}

			osm.SetParams(nominatim.Params{
				Q: c,
			})

			logger.Debug(fmt.Sprintf("Set params to q: %s", c))

			results := osm.Search()

			if len(results) == 0 {

				fmt.Println("No results found for the provided city name.")

				return nil
			}

			result := results[0]

			city := nws.BuildCity(result.DisplayName, result.Lat, result.Lon)

			logger.Debug(fmt.Sprintf("Found: %s", city.Fmt()))

			return nil
		},
	}
}

// ForecastCommand defines a pointer to the forecast command.
//
// Usage: geocast f[orecast] [-city] [-ip] [-pt]
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
				Name:    "ip",
				Aliases: []string{"i"},
				Usage:   "The IP address to fetch the forecast for.",
			},
			&cli.StringSliceFlag{
				Name:    "pt",
				Aliases: []string{"p"},
				Usage:   "The point to fetch the forecast for (lat,lon).",
			},
		},
		Category: "Forecast",
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

			if err == nil {
				city := loc.BuildCity()

				logger.Debug(fmt.Sprintf("Found: %s", city.Fmt()))
			} else {
				logger.Error(err.Error())
			}

			return err
		},
	}
}

func commands() []*cli.Command {
	config := Config()
	return []*cli.Command{
		ForecastCommand(config),
		GeocodeCommand(config),
	}
}

// func application acts as a constant and is
// the entry point for the application.
func Application() *cli.App {
	return &cli.App{
		Name:  "geocast",
		Usage: "Hello world example.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "ip",
				Usage:   "IP address to geocode",
				Aliases: []string{"i"},
			},
		},
		Commands: commands(),
		Action: func(ctx *cli.Context) error {
			ip := ctx.String("ip")

			if ip != "" {
				fmt.Printf("Geocoding IP address: %s\n", ip)
			}

			RootAction(ip, nil)

			return nil
		},
	}
}
