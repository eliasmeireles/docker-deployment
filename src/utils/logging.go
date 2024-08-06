package utils

import (
	"fmt"
	"strings"
)

func Logger(color string, message any, args ...any) {
	var formattedString string

	// Check if message is a string
	if msgStr, ok := message.(string); ok {
		formattedString = fmt.Sprintf(msgStr, args...)
	} else {
		// If message is not a string, use %v to convert it to string
		formattedString = fmt.Sprintf("%v", message)
	}

	// Split the formatted string by newlines
	lines := strings.Split(formattedString, "\n")

	// Print each line with the timestamp and color
	for _, line := range lines {
		if line != "" { // Avoid printing empty lines
			fmt.Printf(
				"%s[%s%s%s] - %s%s\n",
				ColorReset,
				ColorYellow,
				CurrentTimeFormatted(),
				ColorReset,
				color,
				line,
			)
		}
	}
}
