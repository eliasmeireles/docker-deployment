package service

import (
	"docker-deployment/src/utils"
	"os/exec"
)

func Prune(arg ...string) error {
	// Prepare docker-compose command with optional --force-recreate

	cmdArgs := []string{"system", "prune", "-f"}
	cmdArgs = append(cmdArgs, arg...)

	utils.Logger(utils.ColorYellow, "Running docker system prune -f...")
	cmd := exec.Command("docker", cmdArgs...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		utils.Logger(utils.ColorRed, "docker system prune -f : %s", string(output))
		return err
	}

	utils.Logger(utils.ColorYellow, "docker system prune -f completed successful")
	utils.Logger("", string(output))

	return nil
}
