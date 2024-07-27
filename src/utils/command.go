package utils

import (
	"context"
	"os"
	"os/exec"
)

func RunCommand(ctx context.Context, name string, arg ...string) *exec.Cmd {
	command := exec.CommandContext(ctx, name, arg...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err := command.Run()
	if err != nil {
		Logger(ColorGreen, "Failed to run command.")
	}
	return command
}

func RunCommandCheck(ctx context.Context, name string, arg ...string) error {
	command := exec.CommandContext(ctx, name, arg...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}
