package utils

import (
	"time"
)

const (
	DefaultTimeout = 3000 * time.Minute
)

func CurrentTimeFormatted() string {
	now := time.Now()
	return now.Format("2006/01/02 15:04:05")
}
