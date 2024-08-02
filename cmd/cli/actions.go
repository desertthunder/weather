package cli

import (
	"fmt"

	"github.com/desertthunder/weather/internal/ipinfo"
	"github.com/desertthunder/weather/internal/nws"
	"github.com/desertthunder/weather/internal/view"
	"github.com/spf13/viper"
)

// The root action for the CLI.
//
// This function is executed when there are no arguments
// passed to the application. It simply geocodes the current
// user and displays it.
func RootAction(ipaddr string, ipc *ipinfo.IPInfoClient) {
	viper.SetConfigFile(".env")

	if ipc != nil {
		ipc.SetToken(viper.GetString("IPINFO_TOKEN"))
	} else {
		ipc = ipinfo.NewIPInfoClient(viper.GetString("IPINFO_TOKEN"))
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s\n", err)
	}

	loc, err := ipc.Geolocate(&ipaddr)

	if err != nil {
		fmt.Printf("Error: %s\n", err)

		return
	}

	lat, lon := loc.Point()

	c := nws.City{
		Name: loc.City,
		Lat:  lat,
		Long: lon,
	}

	// Render a table!
	headers := []string{"City", "Latitude", "Longitude"}
	data := [][]string{
		{c.Name, fmt.Sprintf("%f", c.Lat), fmt.Sprintf("%f", c.Long)},
	}

	t := view.Table(headers, data)

	fmt.Println(t)
}
