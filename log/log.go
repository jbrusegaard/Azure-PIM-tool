package log

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

var loggerInstance *log.Logger

func InitializeLogger() *log.Logger {
	// Initialize the logger with default settings
	// This function can be expanded to include more complex initialization logic if needed

	styles := log.DefaultStyles()
	styles.Keys["role"] = lipgloss.NewStyle().Foreground(lipgloss.Color("#f305f0")).Bold(true)
	styles.Values["role"] = lipgloss.NewStyle().Bold(true)

	styles.Levels[log.WarnLevel] = lipgloss.NewStyle().
		SetString("WARN").
		// Background(lipgloss.Color("#000000")).
		Foreground(lipgloss.Color("#F59C27"))

	styles.Levels[log.FatalLevel] = lipgloss.NewStyle().
		SetString("EXIT").
		// Background(lipgloss.Color("#000000")).
		Bold(true).
		Foreground(lipgloss.Color("#D41515"))

	logger := log.New(os.Stdout)
	logger.SetLevel(log.InfoLevel) // Set default log level
	logger.SetTimeFormat("2006-01-02 15:04:05")

	logger.SetStyles(styles)

	loggerInstance = logger

	return logger
}

func GetLogger() *log.Logger {
	if loggerInstance == nil {
		return InitializeLogger()
	}
	return loggerInstance
}
