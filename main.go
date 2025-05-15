package main

import (
	"docker-deployment/src/service"
	"docker-deployment/src/utils"
	"os"
)

func main() {
	dockerComposeFile := os.Getenv("DOCKER_COMPOSE_FILE")
	timeoutStr := os.Getenv("TIMEOUT")
	force := utils.GetBoolEnv("FORCE", false)

	utils.EnvLoader(dockerComposeFile)

	service.Start(timeoutStr, dockerComposeFile, force)
}
