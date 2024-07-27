package logger

import (
	"bytes"
	"docker-deployment/src/utils"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func GetPodLogs(containers map[string]string) error {
	for name, containerID := range containers {
		shortContainerID := utils.GetShortId(containerID)
		cmd := exec.Command("docker", "logs", containerID)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf(utils.ColorRed+"Error retrieving logs for container %s (%s): %s"+utils.ColorReset, name, shortContainerID, err)
		}
		// Prepend each line of the log with the container name and short ID
		logLines := strings.Split(out.String(), "\n")
		for _, line := range logLines {
			fmt.Printf("[%s (%s)] %s\n", name, shortContainerID, line)
		}
	}
	return nil
}

func LogDockerComposeContent(dockerComposeFile string) error {
	fmt.Printf(utils.ColorBlue+"Logging content of docker-compose file: %s"+utils.ColorReset+"\n", dockerComposeFile)
	file, err := os.ReadFile(dockerComposeFile)
	if err != nil {
		return err
	}
	fmt.Println(string(file))
	return nil
}
