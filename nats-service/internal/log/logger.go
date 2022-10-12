package log

import (
	"github.com/sirupsen/logrus"
)

// Logger is main logger for the project.
var (
	Logger *logrus.Logger
)

// ConfigureLogger sets log level for Logger.
func ConfigureLogger(logLevel string) {
	Logger = logrus.New()

	switch logLevel {
	case "trace":
		Logger.SetLevel(logrus.TraceLevel)
	case "debug":
		Logger.SetLevel(logrus.DebugLevel)
	case "warning":
		Logger.SetLevel(logrus.WarnLevel)
	case "error":
		Logger.SetLevel(logrus.ErrorLevel)
	case "fatal":
		Logger.SetLevel(logrus.FatalLevel)
	case "panic":
		Logger.SetLevel(logrus.PanicLevel)
	default:
		Logger.SetLevel(logrus.InfoLevel)
	}
}
