package main

import (
	"docker-deployment/src/container"
	"docker-deployment/src/docker"
	"docker-deployment/src/model"
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

func readConfig(filename string) (model.Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return model.Config{}, err
	}

	var config model.Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return model.Config{}, err
	}

	if config.DockerConfig.Host == "" {
		return model.Config{}, fmt.Errorf("docker host is required")
	}

	if config.ContainerConfig.Name == "" {
		return model.Config{}, fmt.Errorf("container name is required")
	}
	if config.ContainerConfig.Image == "" {
		return model.Config{}, fmt.Errorf("container image is required")
	}
	if config.ContainerConfig.Version == "" {
		return model.Config{}, fmt.Errorf("container version is required")
	}

	return config, nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <config.yml>", os.Args[0])
	}
	configFile := os.Args[1]

	config, err := readConfig(configFile)
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	cli, err := docker.Setup(config.DockerConfig)
	if err != nil {
		log.Fatalf("Error setting up Docker client: %v", err)
	}

	exists, err := container.Exists(cli, config.ContainerConfig.Name)
	if err != nil {
		log.Fatalf("Error checking container existence: %v", err)
	}

	if exists && config.ContainerConfig.Force {
		err := container.Remove(cli, config.ContainerConfig.Name)
		if err != nil {
			log.Fatalf("Error removing existing container: %v", err)
		}
	}

	err = container.Run(cli, config.ContainerConfig)
	if err != nil {
		log.Fatalf("Error running container: %v", err)
	}
}
