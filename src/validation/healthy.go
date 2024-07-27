package validation

import (
	"bytes"
	"context"
	"docker-deployment/src/utils"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func ValidateHealthCheck(ctx context.Context, containers map[string]string) error {
	for name, containerID := range containers {
		err := validatePodsStatus(ctx, name, containerID)
		if err != nil {
			return err
		}
	}
	return nil
}

func validatePodsStatus(ctx context.Context, name string, containerID string) error {
	shortContainerID := utils.GetShortId(containerID)

	// Check if the container has a health check defined
	healthCheckCmd := exec.Command("docker", "inspect", "--format={{.State.Health.Status}}", containerID)
	var healthCheckOut bytes.Buffer
	healthCheckCmd.Stdout = &healthCheckOut
	healthCheckCmd.Stderr = os.Stderr

	if err := healthCheckCmd.Run(); err != nil {
		fmt.Printf(utils.ColorRed+"Error checking health status for container %s (%s)\n"+utils.ColorReset, name, shortContainerID)
		fmt.Printf(utils.ColorYellow + "Note: It is always a good idea to use container health check " +
			"configuration to monitor container health properly. See more in: https://docs.docker.com/reference/dockerfile/#healthcheck" + utils.ColorReset + "\n")
		return checkPosIsRunning(ctx, name, shortContainerID, containerID)
	}

	healthCheckConfig := strings.TrimSpace(healthCheckOut.String())
	if healthCheckConfig == "" || healthCheckConfig == "<no value>" {
		// Health check not provided, check if container is running
		return checkPosIsRunning(ctx, name, shortContainerID, containerID)
	} else {
		// Health check is provided, validate health status
		return checkPosIsHealthy(ctx, name, containerID, shortContainerID)
	}
}

func checkPosIsHealthy(checkCtx context.Context, name string, containerID string, shortContainerID string) error {
	// Context with timeout for each health check operation
	ctx, cancel := context.WithTimeout(checkCtx, utils.DefaultTimeout)
	defer cancel()

	fmt.Printf(utils.ColorBlue+"Checking is healthy for container %s (%s)..."+utils.ColorReset+"\n", name, shortContainerID)
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout while waiting for container %s (%s) to become healthy", name, shortContainerID)
		default:
			cmd := exec.Command("docker", "inspect", "--format={{.State.Health.Status}}", containerID)
			var out bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				return fmt.Errorf("error inspecting container %s (%s): %s", name, shortContainerID, err)
			}

			healthStatus := strings.TrimSpace(out.String())
			fmt.Printf(utils.ColorYellow+"Container %s (%s) health status: %s"+utils.ColorReset+"\n", name, shortContainerID, healthStatus)

			switch healthStatus {
			case "healthy":
				fmt.Printf(utils.ColorGreen+"Container %s (%s) is healthy."+utils.ColorReset+"\n", name, shortContainerID)
				return nil
			case "unhealthy":
				return fmt.Errorf("container %s (%s) is unhealthy", name, shortContainerID)
			case "starting":
				// Continue the loop to keep checking
			default:
				return fmt.Errorf("unknown health status for container %s (%s): %s", name, shortContainerID, healthStatus)
			}
			time.Sleep(5 * time.Second)
		}
	}
}

func checkPosIsRunning(checkCtx context.Context, name string, shortContainerID string, containerID string) error {
	// Context with timeout for each status check operation
	ctx, cancel := context.WithTimeout(checkCtx, utils.DefaultTimeout)
	defer cancel()

	fmt.Printf(utils.ColorBlue+"Checking running status for container %s (%s)..."+utils.ColorReset+"\n", name, shortContainerID)
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout while waiting for container %s (%s) to become running", name, shortContainerID)
		default:
			cmd := exec.Command("docker", "inspect", "--format={{.State.Status}}", containerID)
			var out bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				return fmt.Errorf("error inspecting container %s (%s): %s", name, shortContainerID, err)
			}

			status := strings.TrimSpace(out.String())
			fmt.Printf(utils.ColorYellow+"Container %s (%s) status: %s"+utils.ColorReset+"\n", name, shortContainerID, status)

			switch status {
			case "running":
				fmt.Printf(utils.ColorGreen+"Container %s (%s) is running."+utils.ColorReset+"\n", name, shortContainerID)
				return nil
			case "created", "restarting":
				// Continue the loop to keep checking
			default:
				return fmt.Errorf("unknown status for container %s (%s): %s", name, shortContainerID, status)
			}
			time.Sleep(5 * time.Second)
		}
	}
}