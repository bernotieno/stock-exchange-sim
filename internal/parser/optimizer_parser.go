package parser

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// ParseOptimize extracts optimization targets from input.
// Accepts a line like: optimize:(item1;item2;item3)
// Returns a list of item names or an empty slice on error.
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

		// Parse optimize line
		if strings.HasPrefix(line, "optimize:") {
			targets, err := parseOptimizeLine(line, lineNum)
			if err != nil {
				// Non-fatal error, fallback to empty list
				fmt.Printf("Warning: %v, using empty optimization targets\n", err)
				return []string{}, nil
			}
			return targets, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading input: %w", err)
	}

	// No optimize line found
	return []string{}, nil
}

// parseOptimizeLine parses one optimize line into target item names.
// Expects input in the form: optimize:(item1;item2;...)
func parseOptimizeLine(line string, lineNum int) ([]string, error) {
	content := strings.TrimSpace(line[9:]) // remove "optimize:"

	if !strings.HasPrefix(content, "(") || !strings.HasSuffix(content, ")") {
		return nil, fmt.Errorf("line %d: optimize targets must be in parentheses", lineNum)
	}

	// Get text inside parentheses
	targetList := strings.TrimSpace(content[1 : len(content)-1])

	if targetList == "" {
		return []string{}, nil
	}

	parts := strings.Split(targetList, ";")
	var targets []string

	for _, part := range parts {
		target := strings.TrimSpace(part)
		if target == "" {
			continue
		}
		if !isValidItemName(target) {
			return nil, fmt.Errorf("line %d: invalid optimization target '%s'", lineNum, target)
		}
		targets = append(targets, target)
	}

	return targets, nil
}

// isValidItemName returns true if the item name has no invalid characters.
func isValidItemName(name string) bool {
	if name == "" {
		return false
	}

	// Disallow characters that interfere with parsing
	invalidChars := []string{":", ";", "(", ")", "#", "\n", "\r", "\t"}
	for _, char := range invalidChars {
		if strings.Contains(name, char) {
			return false
		}
	}

	return true
}
