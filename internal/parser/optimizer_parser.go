package parser

import (
	"strings"
)

// isValidItemName checks if an item name is valid (basic validation)
func isValidItemName(name string) bool {
	if name == "" {
		return false
	}

	// Check for invalid characters (basic check)
	invalidChars := []string{":", ";", "(", ")", "#", "\n", "\r", "\t"}
	for _, char := range invalidChars {
		if strings.Contains(name, char) {
			return false
		}
	}

	return true
}
