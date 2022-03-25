package appconfig

import (
	"os"

	log "github.com/sirupsen/logrus"
)

var SearchServerAddress = GetEnvOrDefault("GH_SEARCH_SERVER_ADDRESS", "0.0.0.0")
var SearchServerPort = GetEnvOrDefault("GH_SEARCH_SERVER_PORT", "9090")
var ResultsPerPage = GetEnvOrDefault("GH_RESULTS_PER_PAGE", "30")

func SetupLogging() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func GetEnvOrDefault(env, defaultValue string) string {
	value := os.Getenv(env)
	if value != "" {
		return value
	}
	return defaultValue
}
