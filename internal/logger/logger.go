// Package logger wraps and customizes the log package.
package logger

import (
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

// struct color is an object mapping a log level to a color.
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

// Init initializes the logger with a set of default styles and colors while
// also streaming the logs to the console.
func Init() *log.Logger {
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

	return logger
}
