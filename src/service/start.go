package service

import (
	"context"
	"docker-deployment/src/logger"
	"docker-deployment/src/utils"
	"docker-deployment/src/validation"
	"fmt"
	"github.com/google/uuid"
	"io"
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
			utils.Logger(utils.ColorRed, "Invalid TIMEOUT format: %s", err)
			timeout = utils.DefaultTimeout
		} else {
			timeout = time.Duration(timeoutSeconds+60) * time.Second
		}
	}

	// Log docker-compose file content
	err := logger.LogDockerComposeContent(dockerComposeFile)
	if err != nil {
		utils.Logger(utils.ColorRed, "Error logging docker-compose content: %s", err)
		os.Exit(1)
	}

	composeRun(dockerComposeFile, force, timeout, true)
}

func composeRun(dockerComposeFile string, force bool, timeout time.Duration, retry bool) {
	_ = Prune()

	// Generate a UUID and create the path with it
	tempSource := uuid.New().String()
	tempPath := fmt.Sprintf("/_temp/%s/docker-compose.yaml", tempSource)

	// Create the destination directory if it does not exist
	err := os.MkdirAll(fmt.Sprintf("/_temp/%s", tempSource), os.ModePerm)
	if err != nil {
		utils.Logger(utils.ColorRed, "Error creating directory: %s", err)
		os.Exit(1)
	}

	// Copy the docker-compose file to the destination path
	err = copyFile(dockerComposeFile, tempPath)
	if err != nil {
		utils.Logger(utils.ColorRed, "Error copying docker-compose file: %s", err)
		os.Exit(1)
	}

	err = Pull(tempPath)
	if err != nil {
		os.Exit(1)
	}

	// Prepare docker-compose command with optional --force-recreate
	cmdArgs := []string{"-f", tempPath, "up", "-d"}
	if force {
		cmdArgs = append(cmdArgs, "--force-recreate")
	}
	utils.Logger(utils.ColorBlue, "Starting docker-compose...")
	cmd := exec.Command("docker-compose", cmdArgs...)
	if output, err := cmd.CombinedOutput(); err != nil {
		utils.Logger(utils.ColorRed, "Error running docker-compose: %s", string(output))
		if retry && force {
			removeOldContainer(string(output))
			composeRun(tempPath, force, timeout, false)
			return
		}
		os.Exit(1)
	}

	// Get containers
	containerMap, err := GetContainers(tempPath)
	if err != nil {
		utils.Logger(utils.ColorRed, "Error getting containers: %s", err)
	}

	// Context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Channels for health check and logs
	healthCheckDone := make(chan error, 1)
	logsDone := make(chan struct{})

	// Run health check in a goroutine
	go func() {
		healthCheckDone <- validation.ValidateHealthCheck(ctx, timeout, containerMap, dockerComposeFile)
	}()

	// Run logs retrieval in a goroutine
	go func() {
		err := logger.GetPodLogs(ctx, tempPath)
		if err != nil {
			utils.Logger(utils.ColorRed, "Logs retrieval error: %s", err)
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
		_ = Prune()
		// Health check completed
		if err != nil {
			utils.Logger(utils.ColorRed, "Health check error: %s", err)
			cancel()
			<-logsDone
			os.Exit(1)
		}
	}
}

// copyFile copies a file from src to dst
func copyFile(src string, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("source file is not a regular file")
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func(source *os.File) {
		err := source.Close()
		if err != nil {

		}
	}(source)

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func(destination *os.File) {
		err := destination.Close()
		if err != nil {

		}
	}(destination)

	_, err = io.Copy(destination, source)
	if err != nil {
		return err
	}

	return nil
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

			utils.Logger(utils.ColorYellow, "Trying to remove container: [%s] with id [%s]",
				containerName, shortId)

			// Remove the container
			removeCmd := exec.Command("docker", "rm", "-f", containerID)
			if removeErr := removeCmd.Run(); removeErr != nil {
				utils.Logger(utils.ColorRed, "Failed to remove container %s: %s", shortId, removeErr)
				os.Exit(1)
			}

			utils.Logger(utils.ColorYellow, "Container [%s] with id [%s] removed successful",
				containerName, shortId)
		} else {
			utils.Logger(utils.ColorRed, "Failed to parse error message: %s", output)
			os.Exit(1)
		}
	} else {
		utils.Logger(utils.ColorRed, "Failed to remove old container.")
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
