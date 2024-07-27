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
			fmt.Println(ColorYellow + "Entry already exists in /etc/hosts." + ColorReset)
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
				fmt.Printf(ColorRed+"Error closing /etc/hosts: %s"+ColorReset+"\n", err)
			}
		}(f)

		if _, err := f.WriteString(entry); err != nil {
			fmt.Printf(ColorRed+"Error updating /etc/hosts: %s"+ColorReset+"\n", err)
		}

		fmt.Println(ColorGreen + "Added entry to /etc/hosts." + ColorReset)
	}

	return nil
}
