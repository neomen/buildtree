package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// IsTreeSymbol checks if a rune is a tree diagram symbol
func IsTreeSymbol(r rune) bool {
	switch r {
	case ' ', '│', '├', '└', '─', '|', '-', '+', '\\', '/', '>', ':', '\'':
		return true
	default:
		return false
	}
}

// ParseSize converts a string representation of size to bytes
// Examples: "100", "100b", "100kb", "1mb", "1g"
func ParseSize(sizeStr string) (int64, error) {
	sizeStr = strings.TrimSpace(sizeStr)
	if len(sizeStr) == 0 {
		return 0, fmt.Errorf("empty size string")
	}

	// Extract numeric part and unit
	var numStr string
	var unitStr string
	for i, c := range sizeStr {
		if c >= '0' && c <= '9' {
			numStr += string(c)
		} else {
			unitStr = strings.TrimSpace(sizeStr[i:])
			break
		}
	}

	if numStr == "" {
		return 0, fmt.Errorf("no numeric part in size string: %s", sizeStr)
	}

	num, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, fmt.Errorf("invalid numeric part: %s", numStr)
	}

	// Default to KB if no unit specified
	if unitStr == "" {
		unitStr = "kb"
	}

	// Convert to bytes based on unit
	switch strings.ToLower(unitStr) {
	case "b", "bytes":
		return int64(num), nil
	case "k", "kb", "kib", "kilobytes":
		return int64(num) * 1024, nil
	case "m", "mb", "mib", "megabytes":
		return int64(num) * 1024 * 1024, nil
	case "g", "gb", "gib", "gigabytes":
		return int64(num) * 1024 * 1024 * 1024, nil
	default:
		return 0, fmt.Errorf("unknown size unit: %s", unitStr)
	}
}
