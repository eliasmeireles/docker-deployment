package main

import (
	"docker-deployment/src/service"
	"docker-deployment/src/utils"
	"os"
)

func main() {
	dockerComposeFile := os.Getenv("DOCKER_COMPOSE_FILE")
	timeoutStr := os.Getenv("TIMEOUT")
	dockerServerIP := os.Getenv("DOCKER_SERVER_IP")
	force := utils.GetBoolEnv("FORCE", false)

	utils.EnvLoader(dockerComposeFile, dockerServerIP)

	_ = utils.UpdateHostsFile(dockerServerIP)

	//service.DockerLogin(dockerHost, dockerUsername, dockerPassword)

	service.Start(timeoutStr, dockerComposeFile, force)
}
