// Submodule cli provides the definitions for the cli commands
// via the "geocast" application.
package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/desertthunder/weather/cmd/actions"
	"github.com/desertthunder/weather/internal/nominatim"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

type conf struct {
	v *viper.Viper
}

func Config() *conf {
	v := viper.New()
	v.SetConfigFile(".env")

	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s\n", err)
	}

	return &conf{v: v}
}

func (c *conf) Get(key string) string {
	return c.v.GetString(key)
}

func logging() *log.Logger {
	styles := log.DefaultStyles()
	logger := log.New(os.Stdout)
	styles.Levels[log.DebugLevel] = lipgloss.NewStyle().
		SetString("DEBUG").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("#005fd7")).
		Foreground(lipgloss.Color("0"))

	logger.SetStyles(styles)

	return logger
}

var logger = logging()

func GeocastCommand() *cli.Command {
	return &cli.Command{
		Name:    "geocode",
		Aliases: []string{"code", "gc"},
		Usage:   "Geocode a location.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "city",
				Aliases: []string{"c", "n", "cn"},
				Usage:   "City name.",
			},
		},
		Action: func(ctx *cli.Context) error {
			city := ctx.String("city")

			fmt.Printf("Fetching weather forecast for %s...\n", city)

			return nil
		},
	}
}

// func ForecastCommand defines a pointer to the forecast command.
//
// Usage: geocast f[orecast] [-city]
func ForecastCommand(config *conf) *cli.Command {
	return &cli.Command{
		Name: "forecast",
		Aliases: []string{
			"f",
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
			osm := nominatim.Init(
				config.Get("NOMINATIM_USER_AGENT"),
			)

			logger.Debug(fmt.Sprintf("Initialized client with base url: %s", osm.GetBaseURL()))
			c := ctx.String("city")
			if c != "" {
				osm.SetParams(nominatim.Params{
					Q: c,
				})

				logger.Debug(fmt.Sprintf("Set params to: %s", c))

				osm.Search()

				return nil
			}

			fmt.Println("Please provide a city name to fetch the forecast for.")

			return nil
		},
	}
}

func commands() []*cli.Command {
	config := Config()
	return []*cli.Command{
		ForecastCommand(config),
	}
}

// func application acts as a constant and is the entry
// point for the application.
func application() *cli.App {
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

			actions.RootAction(ip, nil)

			return nil
		},
	}
}
