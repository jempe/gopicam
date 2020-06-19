package utils

import (
	"os"
)

// Exists checks if file exists
//
func Exists(file string) bool {
	if _, err := os.Stat(file); err == nil {
		return true
	}
	return false
}

// IsDir checks if path is a directory
//
func IsDirectory(file string) bool {
	if stat, err := os.Stat(file); err == nil && stat.IsDir() {
		return true
	}
	return false
}
