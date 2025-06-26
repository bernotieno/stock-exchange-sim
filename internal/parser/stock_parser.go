package parser

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// ParseStocks reads stock definitions from the input until a process line is encountered.
// Accepts lines like: name:quantity. Ignores empty lines and comments.
// Returns a map of item names to their quantities.
func ParseStocks(r io.Reader) (map[string]int, error) {
	stocks := make(map[string]int)
	scanner := bufio.NewScanner(r)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue // Skip comments and blanks
		}

		if isProcessLine(line) {
			break // Stop parsing at the first process definition
		}

		if err := parseStockLine(line, stocks, lineNum); err != nil {
			return nil, err
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading input: %w", err)
	}

	return stocks, nil
}

// parseStockLine parses a single "name:quantity" line and updates the stocks map.
func parseStockLine(line string, stocks map[string]int, lineNum int) error {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("line %d: invalid stock format, expected 'name:quantity'", lineNum)
	}

	name := strings.TrimSpace(parts[0])
	if name == "" {
		return fmt.Errorf("line %d: stock name cannot be empty", lineNum)
	}

	quantityStr := strings.TrimSpace(parts[1])
	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		return fmt.Errorf("line %d: invalid quantity '%s': %w", lineNum, quantityStr, err)
	}

	if quantity < 0 {
		return fmt.Errorf("line %d: stock quantity cannot be negative", lineNum)
	}

	stocks[name] = quantity
	return nil
}

// isProcessLine returns true if the line resembles a process definition.
// Assumes a process line contains at least 3 colons (name:inputs:outputs:duration).
func isProcessLine(line string) bool {
	return strings.Count(line, ":") >= 3
}
