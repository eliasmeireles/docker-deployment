package utils

import (
	"fmt"
	"os"
)

func EnvLoader(dockerComposeFile string, dockerServerIP string) {
	// Check required environment variables
	if dockerComposeFile == "" {
		fmt.Println(ColorRed + "Error: DOCKER_COMPOSE_FILE environment variable is not set." + ColorReset)
		displayUsage(true)
	}

	if dockerServerIP == "" {
		fmt.Println(ColorRed + "Error: DOCKER_SERVER_IP environment variable is not set." + ColorReset)
		displayUsage(false)
	}
}

// displayUsage prints the usage help message
func displayUsage(required bool) {
	fmt.Println(ColorRed + "Usage: Set the following environment variables:" + ColorReset)
	fmt.Println("  DOCKER_COMPOSE_FILE - Path to the docker-compose file")
	fmt.Println("  TIMEOUT - Timeout for the service start (optional), default is 60 seconds")
	fmt.Println("  DOCKER_SERVER_IP - IP address of the Docker server")
	fmt.Println("  FORCE - Force restart of containers (optional), default false")
	if required {
		os.Exit(1)
	}
}
