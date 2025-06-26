package parser

import (
	"fmt"
	"strings"
)

// parseOptimizeLine parses a single optimize line
func parseOptimizeLine(line string, lineNum int) ([]string, error) {
	// Remove "optimize:" prefix
	content := strings.TrimSpace(line[9:]) // len("optimize:") = 9

	// Check for parentheses
	if !strings.HasPrefix(content, "(") || !strings.HasSuffix(content, ")") {
		return nil, fmt.Errorf("line %d: optimize targets must be enclosed in parentheses", lineNum)
	}

	// Extract content inside parentheses
	targetList := strings.TrimSpace(content[1 : len(content)-1])

	// Handle empty target list
	if targetList == "" {
		return []string{}, nil
	}

	// Split by semicolon and clean up
	parts := strings.Split(targetList, ";")
	var targets []string

	for _, part := range parts {
		target := strings.TrimSpace(part)
		if target == "" {
			continue // Skip empty parts
		}

		// Validate target name (basic check for valid identifier)
		if !isValidItemName(target) {
			return nil, fmt.Errorf("line %d: invalid optimization target name '%s'", lineNum, target)
		}

		targets = append(targets, target)
	}

	return targets, nil
}

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
