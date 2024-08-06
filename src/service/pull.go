package service

import (
	"docker-deployment/src/utils"
	"os/exec"
)

func Pull(dockerComposeFile string) error {
	// Prepare docker-compose command with optional --force-recreate
	cmdArgs := []string{"-f", dockerComposeFile, "pull"}

	utils.Logger(utils.ColorBlue, "Running docker-compose -f %s pull...", dockerComposeFile)
	cmd := exec.Command("docker-compose", cmdArgs...)
	if output, err := cmd.CombinedOutput(); err != nil {
		utils.Logger(utils.ColorRed, "docker-compose -f %s pull failed: %s", dockerComposeFile, string(output))
		return err
	}
	utils.Logger(utils.ColorBlue, "docker-compose -f %s pull completed successful", dockerComposeFile)

	return nil
}
