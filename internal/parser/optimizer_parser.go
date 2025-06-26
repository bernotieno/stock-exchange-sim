package parser

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// ParseOptimize parses optimization targets from a reader.
// The optimize line has the format: optimize:(item1;item2;item3)
// If no optimize line is found or it's malformed, returns an empty slice.
//
// Returns a slice of optimization target item names.
func ParseOptimize(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Look for optimize line
		if strings.HasPrefix(line, "optimize:") {
			targets, err := parseOptimizeLine(line, lineNum)
			if err != nil {
				// Log warning but don't fail - return empty slice as fallback
				fmt.Printf("Warning: %v, using empty optimization targets\n", err)
				return []string{}, nil
			}
			return targets, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading input: %w", err)
	}

	// No optimize line found - return empty slice as fallback
	return []string{}, nil
}

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
