package utils

import (
	"os"
)

func IsFileExist(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}

// remember after used need close file
func OpenFile(file string) (*os.File, error) {
	return os.OpenFile(file, os.O_APPEND|os.O_CREATE, 0644)
}
