package utils

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

func RunCommand(ctx context.Context, name string, arg ...string) *exec.Cmd {
	command := exec.CommandContext(ctx, name, arg...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err := command.Run()
	if err != nil {
		fmt.Println(ColorGreen + "Failed to run command." + ColorReset)
	}
	return command
}

func RunCommandCheck(ctx context.Context, name string, arg ...string) error {
	command := exec.CommandContext(ctx, name, arg...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}
