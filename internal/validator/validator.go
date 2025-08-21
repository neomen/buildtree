package validator

import (
	"path/filepath"
	"strings"
	"unicode/utf8"
)

// IsValidPath checks if a path is safe and valid
func IsValidPath(name string) bool {
	// Removing the ending slash for verification
	cleanName := strings.TrimSuffix(name, "/")
	cleanName = strings.TrimSpace(cleanName)

	if cleanName == "" {
		return false
	}

	// Check for dangerous path patterns
	if cleanName == "." || cleanName == ".." {
		return false
	}

	// Check for path components with ".." or "."
	parts := strings.Split(cleanName, "/")
	for _, part := range parts {
		if part == ".." || part == "." {
			return false
		}
	}

	// Check for prohibited characters
	forbiddenChars := []string{"//", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range forbiddenChars {
		if strings.Contains(cleanName, char) {
			return false
		}
	}

	// Checking for reserved Windows names
	forbiddenNames := []string{"CON", "PRN", "AUX", "NUL", "COM1", "COM2", "COM3", "COM4", "COM5",
		"COM6", "COM7", "COM8", "COM9", "LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9"}

	base := filepath.Base(cleanName)
	for _, forbidden := range forbiddenNames {
		if strings.EqualFold(base, forbidden) {
			return false
		}
	}

	// Checking the length of the name (maximum 255 characters)
	if utf8.RuneCountInString(cleanName) > 255 {
		return false
	}

	// Check that the name does not consist only of spaces
	if strings.TrimSpace(cleanName) == "" {
		return false
	}

	return true
}
