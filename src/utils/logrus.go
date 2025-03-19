package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

type CustomFormatter struct {
	logrus.TextFormatter
}

var Log *logrus.Logger

func init() {
	Log = logrus.New()

	// Set logger to use the custom text formatter
	Log.SetFormatter(&CustomFormatter{
		TextFormatter: logrus.TextFormatter{
			TimestampFormat: "15:04:05.000",
			FullTimestamp:   true,
			ForceColors:     true,
		},
	})

	Log.SetOutput(os.Stdout)
}

// LogRequestBody logs the request body with user information
func LogRequestBody(path string, method string, body interface{}, ip string) {
	Log.WithFields(logrus.Fields{
		"path":   path,
		"method": method,
		"body":   body,
		"ip":     ip,
	}).Info("Request Body")
}
