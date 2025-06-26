package parser

import (
	"bufio"
	"fmt"
	"io"
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
