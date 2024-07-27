package utils

import (
	"fmt"
	"os"
	"strings"
)

func UpdateHostsFile(dockerServerIP string) error {
	if dockerServerIP != "" {
		const hostsFilePath = "/etc/hosts"
		entry := fmt.Sprintf("%s docker-server\n", dockerServerIP)

		file, err := os.ReadFile(hostsFilePath)
		if err != nil {
			return fmt.Errorf("error reading /etc/hosts: %w", err)
		}

		// Check if the entry already exists
		if strings.Contains(string(file), entry) {
			Logger(ColorYellow, "Entry already exists in /etc/hosts.")
			return nil
		}

		// Append the entry to the hosts file
		f, err := os.OpenFile(hostsFilePath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("error opening /etc/hosts for writing: %w", err)
		}

		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				Logger(ColorRed, "Error closing /etc/hosts: %s", err)
			}
		}(f)

		if _, err := f.WriteString(entry); err != nil {
			Logger(ColorRed, "Error updating /etc/hosts: %s", err)
		}

		Logger(ColorGreen, "Added entry to /etc/hosts.")
	}

	return nil
}
