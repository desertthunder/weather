// Geocast is a command line utility/application that allows a user
// to geocode themselves or a city, fetch the weather forecast for
// the city, and display it in a table.
package main

import (
	"log"
	"os"

	"github.com/desertthunder/weather/cmd/cli"
)

// func main is the entrypoint for the CLI.
func main() {
	app := cli.Application()

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
