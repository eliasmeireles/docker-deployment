package utils

import (
	"os"
)

func EnvLoader(dockerComposeFile string) {
	// Check required environment variables
	if dockerComposeFile == "" {
		Logger(ColorRed, "Error: DOCKER_COMPOSE_FILE environment variable is not set.")
		displayUsage(true)
	}
}

// displayUsage prints the usage help message
func displayUsage(required bool) {
	Logger(ColorBlue, "Usage: Set the following environment variables:")
	Logger(ColorBlue, "  DOCKER_COMPOSE_FILE - Path to the docker-compose file")
	Logger(ColorBlue, "  TIMEOUT - Timeout for the service start (optional), default is 5 minutes")
	Logger(ColorBlue, "  FORCE - Force restart of containers (optional), default false")
	if required {
		os.Exit(1)
	}
}
