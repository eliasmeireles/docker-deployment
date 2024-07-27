package service

import (
	"context"
	"docker-deployment/src/utils"
	"os"
	"time"
)

func DockerLogin(dockerHost string, dockerUsername string, dockerPassword string) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	err := utils.RunCommand(ctx, "docker", "login", dockerHost, "-p", dockerPassword, "-u", dockerUsername)
	if err != nil {
		utils.Logger(utils.ColorRed, "Error: Docker login failed.", err)
		os.Exit(1)
	}
}
