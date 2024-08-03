// Submodule conf defines a wrapper around viper to manage settings.
package cli

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/desertthunder/weather/internal/logger"
	"github.com/spf13/viper"
)

// struct conf is a wrapper around the viper package, used for configuration
// management.
type conf struct {
	v   *viper.Viper
	log *log.Logger
}

// func Config is the conf constructor.
//
// It reads the configuration from the .env file and returns a pointer to an
// instance of conf.
func Config() *conf {
	v := viper.New()
	v.SetConfigFile(".env")

	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s\n", err)
	}

	c := &conf{v: v}

	c.log = logger.Init()

	return c
}

// Get is the accessor method for any configuration values/environment variables.
func (c *conf) Get(key string) string {
	return c.v.GetString(key)
}
