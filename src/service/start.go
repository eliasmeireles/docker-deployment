package service

import (
	"context"
	"docker-deployment/src/logger"
	"docker-deployment/src/utils"
	"docker-deployment/src/validation"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func Start(timeoutStr string, dockerComposeFile string, force bool) {
	timeout := utils.DefaultTimeout
	if timeoutStr != "" {
		var err error
		timeoutSeconds, err := parseTimeoutToSeconds(timeoutStr)
		if err != nil {
			fmt.Printf(utils.ColorRed+"Invalid TIMEOUT format: %s"+utils.ColorReset+"\n", err)
			timeout = time.Duration(3000) * time.Second
		} else {
			timeout = time.Duration(timeoutSeconds+60) * time.Second
		}
	}

	// Log docker-compose file content
	err := logger.LogDockerComposeContent(dockerComposeFile)
	if err != nil {
		fmt.Printf(utils.ColorRed+"Error logging docker-compose content: %s"+utils.ColorReset+"\n", err)
		os.Exit(1)
	}

	composeRun(dockerComposeFile, force, timeout, true)
}

func composeRun(dockerComposeFile string, force bool, timeout time.Duration, retry bool) {
	// Prepare docker-compose command with optional --force-recreate
	cmdArgs := []string{"-f", dockerComposeFile, "up", "-d"}
	if force {
		cmdArgs = append(cmdArgs, "--force-recreate")
	}
	fmt.Println(utils.ColorBlue + "Starting docker-compose..." + utils.ColorReset)
	cmd := exec.Command("docker-compose", cmdArgs...)
	if output, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf(utils.ColorRed+"Error running docker-compose: %s"+utils.ColorReset+"\n", string(output))
		if retry && force {
			removeOldContainer(string(output))
			composeRun(dockerComposeFile, force, timeout, false)
			return
		}
		os.Exit(1)
	}

	// Get containers
	containerMap, err := GetContainers(dockerComposeFile)
	if err != nil {
		fmt.Printf(utils.ColorRed+"Error getting containers: %s"+utils.ColorReset+"\n", err)
		os.Exit(1)
	}

	// Context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Channels for health check and logs
	healthCheckDone := make(chan error, 1)
	logsDone := make(chan struct{})

	// Run health check in a goroutine
	go func() {
		healthCheckDone <- validation.ValidateHealthCheck(ctx, timeout, containerMap)
	}()

	// Run logs retrieval in a goroutine
	go func() {
		err := logger.GetPodLogs(ctx, dockerComposeFile)
		if err != nil {
			fmt.Printf(utils.ColorRed+"Logs retrieval error: %s"+utils.ColorReset+"\n", err)
		}
		close(logsDone)
	}()

	// Wait for either health check or logs retrieval to complete
	select {
	case <-logsDone:
		// Logs retrieval completed
		// No need to do anything, health check might still be running
		cancel()          // Cancel health check
		<-healthCheckDone // Ensure health check completes

	case err := <-healthCheckDone:
		// Health check completed
		if err != nil {
			fmt.Printf(utils.ColorRed+"Health check error: %s"+utils.ColorReset+"\n", err)
			cancel()   // Cancel logs retrieval
			<-logsDone // Ensure logs retrieval completes
			os.Exit(1)
		}

	}
	fmt.Printf(utils.ColorGreen+"Deploy for %s completed successfully\n"+utils.ColorReset, dockerComposeFile)
}

func removeOldContainer(output string) {
	// Adjusted pattern with non-greedy matching
	pattern := `Error response from daemon: Conflict. The container name "/([^"]+)" is already in use by container "([0-9a-fA-F]{12,})". You have to remove \(or rename\) that container to be able to reuse that name.`
	re := regexp.MustCompile(pattern)

	if re.MatchString(output) {
		// Extract container name and ID from the error message
		matches := re.FindStringSubmatch(output)
		if len(matches) == 3 {
			containerName := matches[1]
			containerID := matches[2]

			shortId := utils.GetShortId(containerID)

			fmt.Printf(utils.ColorYellow+"Trying to remove container: [%s] with id [%s]"+utils.ColorReset+"\n",
				containerName, shortId)

			// Remove the container
			removeCmd := exec.Command("docker", "rm", "-f", containerID)
			if removeErr := removeCmd.Run(); removeErr != nil {
				fmt.Printf(utils.ColorRed+"Failed to remove container %s: %s"+utils.ColorReset+"\n", shortId, removeErr)
				os.Exit(1)
			}

			fmt.Printf(utils.ColorYellow+"Container [%s] with id [%s] removed successful"+utils.ColorReset+"\n",
				containerName, shortId)
		} else {
			fmt.Printf(utils.ColorRed+"Failed to parse error message: %s"+utils.ColorReset+"\n", output)
			os.Exit(1)
		}
	} else {
		fmt.Printf(utils.ColorRed + "Failed to remove old container." + utils.ColorReset + "\n")
		os.Exit(1)
	}
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
