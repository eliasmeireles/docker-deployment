package utils

import (
	"fmt"
	"os"
	"strings"
)

func GetBoolEnv(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return false
	}
	// Convert to lower case for comparison
	value = strings.ToLower(value)
	if value == "true" || value == "1" {
		return true
	} else if value == "false" || value == "0" {
		return false
	}
	_ = fmt.Errorf("%s environment variable must be 'true', 'false', '1', or '0'", key)
	return false
}
