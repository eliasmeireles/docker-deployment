package service

import (
	"bytes"
	"docker-deployment/src/utils"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func GetContainers(dockerComposeFile string) (map[string]string, error) {
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
		shortContainerID := utils.GetShortId(containerID)

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
		utils.Logger(utils.ColorGreen, "Container %s (%s) started.", name, shortContainerID)
	}

	return containerMap, nil
}
