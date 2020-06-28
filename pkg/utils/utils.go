package utils

import (
	"encoding/base64"
	"fmt"
	"net/http"
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

// Base64 encode
func Base64Encode(content []byte) string {
	mimeType := http.DetectContentType(content)

	encodedContent := base64.StdEncoding.EncodeToString(content)

	return fmt.Sprintf("data:%s;base64,%s", mimeType, encodedContent)
}
