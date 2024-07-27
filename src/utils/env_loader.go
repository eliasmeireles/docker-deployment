package utils

import (
	"os"
)

func EnvLoader(dockerComposeFile string, dockerServerIP string) {
	// Check required environment variables
	if dockerComposeFile == "" {
		Logger(ColorRed, "Error: DOCKER_COMPOSE_FILE environment variable is not set.")
		displayUsage(true)
	}

	if dockerServerIP == "" {
		Logger(ColorBlue, "Error: DOCKER_SERVER_IP environment variable is not set.")
		displayUsage(false)
	}
}

// displayUsage prints the usage help message
func displayUsage(required bool) {
	Logger(ColorBlue, "Usage: Set the following environment variables:")
	Logger(ColorBlue, "  DOCKER_COMPOSE_FILE - Path to the docker-compose file")
	Logger(ColorBlue, "  TIMEOUT - Timeout for the service start (optional), default is 60 seconds")
	Logger(ColorBlue, "  DOCKER_SERVER_IP - IP address of the Docker server")
	Logger(ColorBlue, "  FORCE - Force restart of containers (optional), default false")
	if required {
		os.Exit(1)
	}
}
