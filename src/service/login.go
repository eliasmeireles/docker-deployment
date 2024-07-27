package service

import (
	"context"
	"docker-deployment/src/utils"
	"fmt"
	"os"
	"time"
)

func DockerLogin(dockerHost string, dockerUsername string, dockerPassword string) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	err := utils.RunCommand(ctx, "docker", "login", dockerHost, "-p", dockerPassword, "-u", dockerUsername)
	if err != nil {
		fmt.Println(utils.ColorRed+"Error: Docker login failed."+utils.ColorReset, err)
		os.Exit(1)
	}
}
