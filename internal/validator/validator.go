package validator

import (
	"path/filepath"
	"strings"
)

// IsValidPath checks if a path is safe and valid
func IsValidPath(name string) bool {
	cleanName := strings.TrimSuffix(name, "/")
	if cleanName == "" {
		return false
	}

	if cleanName == "." || cleanName == ".." {
		return false
	}

	// Forbidden characters
	forbiddenChars := []string{"\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range forbiddenChars {
		if strings.Contains(cleanName, char) {
			return false
		}
	}

	// Forbidden names (Windows)
	forbiddenNames := []string{"CON", "PRN", "AUX", "NUL", "COM1", "COM2", "COM3", "COM4", "COM5",
		"COM6", "COM7", "COM8", "COM9", "LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9"}

	base := filepath.Base(cleanName)
	for _, forbidden := range forbiddenNames {
		if strings.EqualFold(base, forbidden) {
			return false
		}
	}

	return true
}
