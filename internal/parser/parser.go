package parser

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ParseStockDefinitions reads a configuration file and extracts stock definitions
// that appear before any process lines. Stock definitions are in the format
// "stock_name:quantity" where quantity must be a valid integer.
//
// The function stops parsing stock definitions when it encounters the first
// line starting with "process:" and returns all valid stock definitions
// found up to that point.
//
// Parameters:
//   - filename: path to the configuration file to parse
//
// Returns:
//   - map[string]int: stock names mapped to their quantities 
//   - error: any error encountered during file operations or parsing
//
// The function skips:
//   - Empty lines (after trimming whitespace)
//   - Comment lines (starting with '#')
//
// Error conditions:
//   - File cannot be opened or read
//   - Stock definition line missing colon separator
//   - Stock quantity is not a valid integer
func ParseStockDefinitions(filename string) (map[string]int, error) {
	// Open the configuration file for reading
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %q: %w", filename, err)
	}
	defer file.Close()

	// Initialize the stock map to store parsed definitions
	stocks := make(map[string]int)
	
	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	lineNum := 0

	// Process each line in the file
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		
		// Skip empty lines - no processing needed
		if line == "" {
			continue
		}
		
		// Skip comment lines (starting with '#')
		if strings.HasPrefix(line, "#") {
			continue
		}
		
		// Stop parsing stock definitions when we encounter a process line
		// Process lines mark the beginning of a different configuration section
		if strings.HasPrefix(line, "process:") {
			break
		}
		
		// Parse stock definition lines in format "stock_name:quantity"
		if err := parseStockLine(line, lineNum, stocks); err != nil {
			return nil, err
		}
	}
	
	// Check for any scanner errors that occurred during file reading
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}
	
	return stocks, nil
}

// parseStockLine processes a single line expected to contain a stock definition
// and adds the parsed stock to the provided stocks map.
//
// Parameters:
//   - line: the line to parse (should be in format "stock_name:quantity")
//   - lineNum: current line number for error reporting
//   - stocks: map to store the parsed stock definition
//
// Returns:
//   - error: parsing error with line number context, or nil if successful
func parseStockLine(line string, lineNum int, stocks map[string]int) error {
	// Split the line on the colon separator
	// We expect exactly two parts: stock name and quantity
	parts := strings.Split(line, ":")
	if len(parts) != 2 {
		return fmt.Errorf("line %d: malformed stock definition %q (expected format: stock_name:quantity)", 
			lineNum, line)
	}
	
	// Extract and clean the stock name and quantity string
	stockName := strings.TrimSpace(parts[0])
	quantityStr := strings.TrimSpace(parts[1])
	
	// Validate that stock name is not empty after trimming
	if stockName == "" {
		return fmt.Errorf("line %d: empty stock name in %q", lineNum, line)
	}
	
	// Convert quantity string to integer
	// This will fail if the quantity is not a valid integer
	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		return fmt.Errorf("line %d: invalid quantity %q for stock %q (must be an integer)", 
			lineNum, quantityStr, stockName)
	}
	
	// Store the parsed stock definition in the map
	// Note: if a stock name appears multiple times, the last value wins
	stocks[stockName] = quantity
	
	return nil
}