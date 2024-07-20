package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	defaultTimeout = 1 * time.Minute
	colorReset     = "\033[0m"
	colorGreen     = "\033[32m"
	colorRed       = "\033[31m"
	colorYellow    = "\033[33m"
	colorBlue      = "\033[34m"
)

func main() {
	// Get environment variables
	dockerComposeFile := os.Getenv("DOCKER_COMPOSE_FILE")
	timeoutStr := os.Getenv("TIMEOUT")
	dockerServerIP := os.Getenv("DOCKER_SERVER_IP")
	force := os.Getenv("FORCE")

	if dockerComposeFile == "" {
		fmt.Println(colorRed + "Error: DOCKER_COMPOSE_FILE environment variable is not set." + colorReset)
		os.Exit(1)
	}

	if dockerServerIP == "" {
		fmt.Println(colorRed + "Error: DOCKER_SERVER_IP environment variable is not set." + colorReset)
		os.Exit(1)
	}

	timeout := defaultTimeout
	if timeoutStr != "" {
		var err error
		timeoutSeconds, err := parseTimeoutToSeconds(timeoutStr)
		if err != nil {
			fmt.Printf(colorRed+"Invalid TIMEOUT format: %s"+colorReset+"\n", err)
			os.Exit(1)
		}
		timeout = time.Duration(timeoutSeconds) * time.Second
	}

	// Update /etc/hosts with docker-server
	err := updateHostsFile(dockerServerIP)
	if err != nil {
		fmt.Printf(colorRed+"Error updating /etc/hosts: %s"+colorReset+"\n", err)
		os.Exit(1)
	}

	// Log docker-compose file content
	err = logDockerComposeContent(dockerComposeFile)
	if err != nil {
		fmt.Printf(colorRed+"Error logging docker-compose content: %s"+colorReset+"\n", err)
		os.Exit(1)
	}

	// Prepare docker-compose command with optional --force-recreate
	cmdArgs := []string{"-f", dockerComposeFile, "up", "-d"}
	if force != "" {
		cmdArgs = append(cmdArgs, "--force-recreate")
	}
	fmt.Println(colorBlue + "Starting docker-compose..." + colorReset)
	cmd := exec.Command("docker-compose", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf(colorRed+"Error running docker-compose: %s"+colorReset+"\n", err)
		os.Exit(1)
	}

	// Get containers
	containerMap, err := getContainers(dockerComposeFile)
	if err != nil {
		fmt.Printf(colorRed+"Error getting containers: %s"+colorReset+"\n", err)
		os.Exit(1)
	}

	// Context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Channels for health check and logs
	healthCheckDone := make(chan error, 1)
	logsDone := make(chan error, 1)

	// Run health check in a goroutine
	go func() {
		healthCheckDone <- validateHealthCheck(ctx, containerMap)
	}()

	// Run logs retrieval in a goroutine
	go func() {
		logsDone <- getPodLogs(containerMap)
	}()

	// Wait for both goroutines to complete
	healthCheckErr := <-healthCheckDone
	logsErr := <-logsDone

	// Check for errors from health check
	if healthCheckErr != nil {
		fmt.Printf(colorRed+"Health check error: %s"+colorReset+"\n", healthCheckErr)
		os.Exit(1)
	}

	// Check for errors from logs retrieval
	if logsErr != nil {
		fmt.Printf(colorRed+"Logs retrieval error: %s"+colorReset+"\n", logsErr)
		os.Exit(1)
	}

	fmt.Println(colorGreen + "Deploy completed successfully" + colorReset)
}

func parseTimeoutToSeconds(timeoutStr string) (int64, error) {
	timeoutStr = strings.TrimSpace(timeoutStr)
	if strings.HasSuffix(timeoutStr, "s") {
		timeoutStr = strings.TrimSuffix(timeoutStr, "s")
	}

	timeoutSeconds, err := strconv.ParseInt(timeoutStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return timeoutSeconds, nil
}

func logDockerComposeContent(dockerComposeFile string) error {
	fmt.Printf(colorBlue+"Logging content of docker-compose file: %s"+colorReset+"\n", dockerComposeFile)
	file, err := os.ReadFile(dockerComposeFile)
	if err != nil {
		return err
	}
	fmt.Println(string(file))
	return nil
}

func updateHostsFile(ip string) error {
	const hostsFilePath = "/etc/hosts"
	entry := fmt.Sprintf("%s docker-server\n", ip)

	file, err := os.ReadFile(hostsFilePath)
	if err != nil {
		return fmt.Errorf("error reading /etc/hosts: %w", err)
	}

	// Check if the entry already exists
	if strings.Contains(string(file), entry) {
		fmt.Println(colorYellow + "Entry already exists in /etc/hosts." + colorReset)
		return nil
	}

	// Append the entry to the hosts file
	f, err := os.OpenFile(hostsFilePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening /etc/hosts for writing: %w", err)
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Printf(colorRed+"Error closing /etc/hosts: %s"+colorReset+"\n", err)
		}
	}(f)

	if _, err := f.WriteString(entry); err != nil {
		return fmt.Errorf("error writing to /etc/hosts: %w", err)
	}

	fmt.Println(colorGreen + "Added entry to /etc/hosts." + colorReset)
	return nil
}

func getContainers(dockerComposeFile string) (map[string]string, error) {
	cmd := exec.Command("docker-compose", "-f", dockerComposeFile, "ps", "-q")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	containerIDs := strings.TrimSpace(out.String())
	if containerIDs == "" {
		return nil, fmt.Errorf("no containers found")
	}

	idList := strings.Split(containerIDs, "\n")
	containerMap := make(map[string]string)

	for _, containerID := range idList {
		shortContainerID := getShortId(containerID)

		cmdName := exec.Command("docker", "inspect", "--format={{.Name}}", containerID)
		var nameOut bytes.Buffer
		cmdName.Stdout = &nameOut
		cmdName.Stderr = os.Stderr
		if err := cmdName.Run(); err != nil {
			return nil, err
		}

		name := strings.TrimSpace(nameOut.String())
		// Strip leading '/' from container name
		name = strings.TrimPrefix(name, "/")
		containerMap[name] = containerID
		fmt.Printf(colorGreen+"Container %s (%s) started."+colorReset+"\n", name, shortContainerID)
	}

	return containerMap, nil
}

func getShortId(containerID string) string {
	shortContainerID := containerID
	if len(containerID) > 10 {
		shortContainerID = containerID[:10]
	}
	return shortContainerID
}

func validateHealthCheck(ctx context.Context, containers map[string]string) error {
	for name, containerID := range containers {
		return validatePodsStatus(ctx, name, containerID)
	}
	return nil
}

func validatePodsStatus(ctx context.Context, name string, containerID string) error {

	shortContainerID := getShortId(containerID)

	// Check if the container has a health check defined

	time.Sleep(2 * time.Second)

	healthCheckCmd := exec.Command("docker", "inspect", "--format={{.Config.Healthcheck}}", containerID)
	var healthCheckOut bytes.Buffer
	healthCheckCmd.Stdout = &healthCheckOut
	healthCheckCmd.Stderr = os.Stderr

	healthCheckConfig := strings.TrimSpace(healthCheckOut.String())
	if healthCheckConfig == "map[]" || healthCheckConfig == "" {
		// Health check not provided, check if container is running
		return checkPosIsRunning(ctx, name, shortContainerID, containerID)
	} else {
		// Health check is provided, validate health status
		return checkPosIsHealthy(ctx, name, containerID, shortContainerID)
	}
}

func checkPosIsHealthy(checkCtx context.Context, name string, containerID string, shortContainerID string) error {
	// Context with timeout for each health check operation
	ctx, cancel := context.WithTimeout(checkCtx, defaultTimeout)
	defer cancel()

	fmt.Printf(colorBlue+"Checking is healthy for container %s (%s)..."+colorReset+"\n", name, shortContainerID)
	select {
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			fmt.Printf(colorRed+"Health check for container %s (%s) timed out."+colorReset+"\n", name, shortContainerID)
		} else {
			fmt.Printf(colorRed+"Health check for container %s (%s) failed: %s"+colorReset+"\n", name, shortContainerID, ctx.Err())
		}
		return ctx.Err()
	default:
		cmd := exec.Command("docker", "inspect", "--format={{.State.Health.Status}}", containerID)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf(colorRed+"Error checking health for container %s (%s): %s"+colorReset, name, shortContainerID, err)
		}

		healthStatus := strings.TrimSpace(out.String())
		if healthStatus == "healthy" {
			fmt.Printf(colorGreen+"Container %s (%s) is healthy."+colorReset+"\n", name, shortContainerID)
			return nil
		}

		if healthStatus == "unhealthy" || healthStatus == "starting" {
			fmt.Printf(colorRed+"Container %s (%s) is unhealthy."+colorReset+"\n", name, shortContainerID)
			time.Sleep(10 * time.Second)
			return checkPosIsHealthy(ctx, name, containerID, shortContainerID)
		}

		fmt.Printf(colorRed+"Container %s (%s) has unknown health status: %s."+colorReset+"\n", name, shortContainerID, healthStatus)
		return fmt.Errorf(colorRed+"Container %s (%s) has unknown health status: %s."+colorReset, name, shortContainerID, healthStatus)
	}
}

func checkPosIsRunning(ctx context.Context, name string, shortContainerID string, containerID string) error {
	// Context with timeout for each health check operation
	checkCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	fmt.Printf(colorBlue+"Checking is running for container %s (%s)..."+colorReset+"\n", name, shortContainerID)
	for {
		select {
		case <-checkCtx.Done():
			if errors.Is(checkCtx.Err(), context.DeadlineExceeded) {
				fmt.Printf(colorRed+"Health check for container %s (%s) timed out."+colorReset+"\n", name, shortContainerID)
			} else {
				fmt.Printf(colorRed+"Health check for container %s (%s) failed: %s"+colorReset+"\n", name, shortContainerID, checkCtx.Err())
			}
			return checkCtx.Err()
		default:
			cmd := exec.Command("docker", "inspect", "--format={{.State.Running}}", containerID)
			var out bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf(colorRed+"Error checking running status for container %s (%s): %s"+colorReset, name, shortContainerID, err)
			}

			runningStatus := strings.TrimSpace(out.String())
			if runningStatus == "true" {
				fmt.Printf(colorGreen+"Container %s (%s) is running."+colorReset+"\n", name, shortContainerID)
				return nil
			}

			if runningStatus == "starting" {
				fmt.Printf(colorYellow+"Container %s (%s) is starting, retrying."+colorReset+"\n", name, shortContainerID)
				time.Sleep(10 * time.Second)
				return checkPosIsRunning(checkCtx, name, shortContainerID, containerID)
			}

			fmt.Printf(colorRed+"Container %s (%s) is not running. Current status %s."+colorReset+"\n", name, shortContainerID, runningStatus)
			return fmt.Errorf(colorRed+"Container %s (%s) is not running. Current status %s."+colorReset, name, shortContainerID, runningStatus)
		}
	}
}

func getPodLogs(containers map[string]string) error {
	for name, containerID := range containers {
		shortContainerID := getShortId(containerID)
		cmd := exec.Command("docker", "logs", containerID)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf(colorRed+"Error retrieving logs for container %s (%s): %s"+colorReset, name, shortContainerID, err)
		}
		// Prepend each line of the log with the container name and short ID
		logLines := strings.Split(out.String(), "\n")
		for _, line := range logLines {
			fmt.Printf("[%s (%s)] %s\n", name, shortContainerID, line)
		}
	}
	return nil
}
