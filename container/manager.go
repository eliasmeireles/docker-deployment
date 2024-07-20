package container

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"io/ioutil"

	"docker-deployment/model"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func Exists(cli *client.Client, containerName string) (bool, error) {
	containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		return false, err
	}
	for _, c := range containers {
		for _, name := range c.Names {
			if name == "/"+containerName {
				return true, nil
			}
		}
	}
	return false, nil
}

func Remove(cli *client.Client, containerName string) error {
	return cli.ContainerRemove(context.Background(), containerName, container.RemoveOptions{Force: true})
}

func Run(cli *client.Client, config model.Container) error {
	ctx := context.Background()

	var envs []string
	for k, v := range config.Env {
		envs = append(envs, fmt.Sprintf("%s=%s", k, v))
	}

	containerConfig := &container.Config{
		Image: fmt.Sprintf("%s:%s", config.Image, config.Version),
		Env:   envs,
	}

	hostConfig := &container.HostConfig{
		Binds:        config.Volumes,
		PortBindings: map[nat.Port][]nat.PortBinding{},
	}

	for _, port := range config.Port {
		hostConfig.PortBindings[nat.Port(port)] = []nat.PortBinding{
			{
				HostPort: port,
			},
		}
	}

	resp, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, config.Name)
	if err != nil {
		return err
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return err
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNextExit)
	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case status := <-statusCh:
		if status.StatusCode != 0 {
			return fmt.Errorf("container exited with status %d", status.StatusCode)
		}
	}

	logs, err := cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true, Follow: true})
	if err != nil {
		return err
	}
	defer logs.Close()

	logData, err := ioutil.ReadAll(logs)
	if err != nil {
		return err
	}

	fmt.Printf("Container Logs:\n%s\n", string(logData))
	return nil
}
