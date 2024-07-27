package logger

import (
	"bufio"
	"context"
	"docker-deployment/src/utils"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func GetPodLogs(ctx context.Context, dockerComposeFile string) error {
	time.Sleep(2 * time.Second)
	cmd := exec.Command("docker-compose", "-f", dockerComposeFile, "logs", "-f")
	utils.Logger(utils.ColorBlue, "Getting %s logs%s", dockerComposeFile)

	// Create a pipe to capture the command output
	stdoutPipe, err := cmd.StdoutPipe()

	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %s", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %s", err)
	}

	// Read and print logs from stdout in a goroutine
	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			utils.Logger("", scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			utils.Logger(utils.ColorRed, "Error reading logs: %s", err)
		}
	}()

	// Wait for the context to be done or the command to finish
	select {
	case <-ctx.Done():
		// Context is done, cancel logs retrieval
		if err := cmd.Process.Kill(); err != nil {
			utils.Logger(utils.ColorRed, "Failed to kill logs process: %s", err)
		}
		return ctx.Err()
	case err := <-waitCmd(cmd):
		// Command finished
		return err
	}
}

// Helper function to wait for command completion
func waitCmd(cmd *exec.Cmd) <-chan error {
	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()
	return done
}

func LogDockerComposeContent(dockerComposeFile string) error {
	utils.Logger(utils.ColorBlue, "Logging content of docker-compose file: %s", dockerComposeFile)
	file, err := os.ReadFile(dockerComposeFile)
	if err != nil {
		return err
	}
	utils.Logger(utils.ColorBlue, string(file))
	return nil
}
