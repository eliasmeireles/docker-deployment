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
	fmt.Printf("%sGetting %s logs%s\n", utils.ColorBlue, dockerComposeFile, utils.ColorReset)

	fmt.Printf("Starting pipe logs\n")

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
		fmt.Printf("Starting scaner logs\n")
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading logs: %s\n", err)
		}
	}()

	// Wait for the context to be done or the command to finish
	select {
	case <-ctx.Done():
		// Context is done, cancel logs retrieval
		if err := cmd.Process.Kill(); err != nil {
			fmt.Printf("Failed to kill logs process: %s\n", err)
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
	fmt.Printf(utils.ColorBlue+"Logging content of docker-compose file: %s"+utils.ColorReset+"\n", dockerComposeFile)
	file, err := os.ReadFile(dockerComposeFile)
	if err != nil {
		return err
	}
	fmt.Println(utils.ColorBlue + string(file) + utils.ColorReset)
	return nil
}
