package utils

import "fmt"

func Logger(color string, message any, args ...any) {
	var formattedString string

	// Check if message is a string
	if msgStr, ok := message.(string); ok {
		formattedString = fmt.Sprintf(msgStr, args...)
	} else {
		// If message is not a string, use %v to convert it to string
		formattedString = fmt.Sprintf("%v", message)
	}

	fmt.Printf(ColorReset+"[%s%s%s] - %s%s\n", ColorYellow, CurrentTimeFormatted(), ColorReset, color, formattedString)
}
