package container

import (
	"context"
	"docker-deployment/src/model"
	"fmt"
	"github.com/docker/go-connections/nat"
	"io"
	"os"
	"strings"
	"time"

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

// GetContainerLogs retrieves the logs for a given container ID
func GetContainerLogs(cli *client.Client, containerID string, follow bool) error {
	ctx := context.Background()

	// Fetch logs
	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     follow,
		Tail:       "all", // Retrieve all logs. Modify if you need to limit log entries.
	}

	logs, err := cli.ContainerLogs(ctx, containerID, options)
	if err != nil {
		return fmt.Errorf("error fetching container logs: %v", err)
	}
	defer func(logs io.ReadCloser) {
		if err := logs.Close(); err != nil {
			fmt.Printf("Error closing logs: %v\n", err)
		}
	}(logs)

	logData, err := io.ReadAll(logs)
	if err != nil {
		return fmt.Errorf("error reading container logs: %v", err)
	}
	fmt.Printf(string(logData))
	return nil
}

func healthCheckLoader(config model.Container, containerConfig *container.Config) {
	// Add health check configuration if provided
	if config.HealthCheck != nil {
		containerConfig.Healthcheck = &container.HealthConfig{
			Test:        config.HealthCheck.Test,
			Interval:    time.Duration(config.HealthCheck.Interval) * time.Second,
			Timeout:     time.Duration(config.HealthCheck.Timeout) * time.Second,
			Retries:     config.HealthCheck.Retries,
			StartPeriod: time.Duration(config.HealthCheck.StartPeriod) * time.Second,
		}
	}
}

func bindPortsLoader(config model.Container, hostConfig *container.HostConfig) error {
	for _, port := range config.Ports {
		portParts := strings.Split(port, ":")
		if len(portParts) != 2 && len(portParts) != 3 {
			return fmt.Errorf("invalid port specification: %s", port)
		}

		var (
			hostIP        string
			hostPort      string
			containerPort string
		)

		if len(portParts) == 3 {
			hostIP = portParts[0]
			hostPort = portParts[1]
			containerPort = portParts[2]
		} else {
			hostPort = portParts[0]
			containerPort = portParts[1]
		}

		// Define the port with protocol
		portWithProtocol := fmt.Sprintf("%s/tcp", containerPort)

		// Create port binding
		portBinding := nat.PortBinding{
			HostIP:   hostIP,
			HostPort: hostPort,
		}

		// Define the port bindings
		hostConfig.PortBindings[nat.Port(portWithProtocol)] = []nat.PortBinding{
			portBinding,
		}
	}
	return nil
}

func CheckContainerStatus(cli *client.Client, containerID string, healthCheck *model.HealthCheck, timeout int) error {
	ctx := context.Background()
	start := time.Now()

	// Default timeout if not specified
	if timeout == 0 {
		timeout = 60 // 1 minute default timeout
	}

	for {
		// Check the container status
		containerInfo, err := cli.ContainerInspect(ctx, containerID)
		if err != nil {
			return fmt.Errorf("error inspecting container: %v", err)
		}

		// If health check is provided
		if healthCheck != nil {
			if containerInfo.State.Health != nil {
				switch containerInfo.State.Health.Status {
				case "healthy":
					fmt.Println("Container is healthy. Deployment completed successfully")
					return nil
				case "unhealthy":
					return fmt.Errorf("container is unhealthy")
				}
			}

			// Check for timeout
			if time.Since(start) > time.Duration(timeout)*time.Second {
				return fmt.Errorf("timeout exceeded while waiting for container to be healthy")
			}

			// Wait before checking again
			time.Sleep(time.Duration(healthCheck.Interval) * time.Second)
		} else {
			// No health check provided, just wait for container to be running
			if containerInfo.State.Running {
				fmt.Println("Container is running. Deployment completed successfully")
				return nil
			}

			// Check for timeout
			if time.Since(start) > time.Duration(timeout)*time.Second {
				return fmt.Errorf("timeout exceeded while waiting for container to be running")
			}

			// Wait before checking again
			fmt.Println("Container is not running yet, waiting...")
			time.Sleep(5 * time.Second) // Check every 5 seconds if no health check is provided
		}
	}
}

func Run(cli *client.Client, config model.Container) error {
	ctx := context.Background()

	// Convert environment variables from map to slice of strings
	var envs []string
	for k, v := range config.Env {
		envs = append(envs, fmt.Sprintf("%s=%s", k, v))
	}

	// Define container configuration
	containerConfig := &container.Config{
		Image: fmt.Sprintf("%s:%s", config.Image, config.Version),
		Env:   envs,
	}

	// Define host configuration
	hostConfig := &container.HostConfig{
		Binds:        config.Volumes,
		PortBindings: map[nat.Port][]nat.PortBinding{},
	}

	// Convert port strings to Docker's format
	if err := bindPortsLoader(config, hostConfig); err != nil {
		return fmt.Errorf("error binding ports: %v", err)
	}

	healthCheckLoader(config, containerConfig)

	// Create the container
	resp, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, config.Name)
	if err != nil {
		return fmt.Errorf("error creating container: %v", err)
	}

	// Start the container
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return fmt.Errorf("error starting container: %v", err)
	}

	fmt.Printf("Container %s started\n", config.Name)

	// Channels to synchronize goroutines
	statusChan := make(chan error)

	// Start goroutine to check container status
	go func() {
		statusErr := CheckContainerStatus(cli, resp.ID, config.HealthCheck, config.Startup.Timeout)
		statusChan <- statusErr
	}()

	// Start goroutine to get container logs
	go func() {
		err := GetContainerLogs(cli, resp.ID, true) // Set `false` to not follow logs
		if err != nil {
			statusChan <- fmt.Errorf("error getting container logs: %v", err)
			return
		}
	}()

	// Wait for either status check or logs retrieval
	var statusErr error
	var logs string
	select {
	case err := <-statusChan:
		if err != nil {
			statusErr = err
		} else {
			statusErr = nil
		}
	}

	// Print logs if available
	if logs != "" {
		fmt.Printf("Container Logs:\n%s\n", logs)
	}

	// Handle the status result
	if statusErr == nil {
		_ = make(chan os.Signal)
		return nil
	} else {
		_ = make(chan os.Signal)
		return statusErr
	}
}
