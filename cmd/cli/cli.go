// Submodule cli provides the definitions for the cli commands
// via the "geocast" application.
package main

import (
	"fmt"

	"github.com/desertthunder/weather/cmd/actions"
	"github.com/urfave/cli/v2"
)

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

func commands() []*cli.Command {
	return []*cli.Command{
		{
			Name:  "hello",
			Usage: "Prints 'Hello world!'",
			Action: func(ctx *cli.Context) error {
				fmt.Println("Hello world!")

				return nil
			},
		},
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
