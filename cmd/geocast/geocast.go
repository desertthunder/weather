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
