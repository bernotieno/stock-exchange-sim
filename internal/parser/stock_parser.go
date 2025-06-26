package parser

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// ParseStocks parses stock definitions from a reader until it encounters
// a process definition or EOF. Stock lines have the format "name:quantity".
// Empty lines and lines starting with '#' are skipped.
//
// Returns a map of stock names to quantities and any parsing error encountered.
func ParseStocks(r io.Reader) (map[string]int, error) {
	stocks := make(map[string]int)
	scanner := bufio.NewScanner(r)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Stop parsing stocks when we hit a process definition
		if isProcessLine(line) {
			break
		}

		// Parse stock line
		if err := parseStockLine(line, stocks, lineNum); err != nil {
			return nil, err
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading input: %w", err)
	}

	return stocks, nil
}

// parseStockLine parses a single stock line and adds it to the stocks map
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

// isProcessLine checks if a line appears to be a process definition
// Process lines contain at least 3 colons (name:inputs:outputs:duration)
func isProcessLine(line string) bool {
	return strings.Count(line, ":") >= 3
}
