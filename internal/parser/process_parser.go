package parser

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

// Process defines a single manufacturing step with inputs, outputs, and duration.
type Process struct {
	Name     string         // Name of the process
	Inputs   map[string]int // Required input items
	Outputs  map[string]int // Produced output items
	Duration time.Duration  // Time per production cycle
}

// ParseProcesses reads and parses all valid process lines from the reader.
// Ignores empty lines and comments. Returns a slice of Process structs.
func ParseProcesses(r io.Reader) ([]Process, error) {
	var processes []Process
	scanner := bufio.NewScanner(r)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue // Skip comments and blanks
		}

		if !isProcessLine(line) {
			continue // Ignore unrelated lines
		}

		process, err := parseProcessLine(line, lineNum)
		if err != nil {
			return nil, err
		}

		processes = append(processes, process)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading input: %w", err)
	}

	return processes, nil
}

// parseItemQuantities parses (item:qty;...) sections into a map.
// Validates structure and ensures quantities are positive integers.
func parseItemQuantities(section string, lineNum int, sectionName string) (map[string]int, error) {
	section = strings.TrimSpace(section)

	if !strings.HasPrefix(section, "(") || !strings.HasSuffix(section, ")") {
		return nil, fmt.Errorf("line %d: %s section must be enclosed in parentheses", lineNum, sectionName)
	}

	content := strings.TrimSpace(section[1 : len(section)-1])
	items := make(map[string]int)

	if content == "" {
		return items, nil // Empty section
	}

	pairs := strings.Split(content, ";")
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}

		itemParts := strings.SplitN(pair, ":", 2)
		if len(itemParts) != 2 {
			return nil, fmt.Errorf("line %d: invalid %s item format '%s', expected 'item:quantity'", lineNum, sectionName, pair)
		}

		itemName := strings.TrimSpace(itemParts[0])
		if itemName == "" {
			return nil, fmt.Errorf("line %d: item name cannot be empty in %s", lineNum, sectionName)
		}

		quantityStr := strings.TrimSpace(itemParts[1])
		quantity, err := strconv.Atoi(quantityStr)
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid quantity '%s' for item '%s' in %s: %w", lineNum, quantityStr, itemName, sectionName, err)
		}

		if quantity <= 0 {
			return nil, fmt.Errorf("line %d: quantity for item '%s' in %s must be positive", lineNum, itemName, sectionName)
		}

		items[itemName] = quantity
	}

	return items, nil
}

// parseDuration parses duration values in Go format (e.g., "30s", "1m") or plain seconds.
func parseDuration(durationStr string, lineNum int) (time.Duration, error) {
	durationStr = strings.TrimSpace(durationStr)
	if durationStr == "" {
		return 0, fmt.Errorf("line %d: duration cannot be empty", lineNum)
	}

	// Try standard Go duration format
	if duration, err := time.ParseDuration(durationStr); err == nil {
		if duration <= 0 {
			return 0, fmt.Errorf("line %d: duration must be positive", lineNum)
		}
		return duration, nil
	}

	// Try plain seconds
	if seconds, err := strconv.ParseInt(durationStr, 10, 64); err == nil {
		if seconds <= 0 {
			return 0, fmt.Errorf("line %d: duration must be positive", lineNum)
		}
		return time.Duration(seconds) * time.Second, nil
	}

	return 0, fmt.Errorf("line %d: invalid duration format '%s'", lineNum, durationStr)
}

// parseProcessLine parses a single line defining a process.
// Expected format: name:(inputs):(outputs):duration
func parseProcessLine(line string, lineNum int) (Process, error) {
	firstColon := strings.Index(line, ":")
	if firstColon == -1 {
		return Process{}, fmt.Errorf("line %d: invalid process format, expected 'name:inputs:outputs:duration'", lineNum)
	}

	name := strings.TrimSpace(line[:firstColon])
	if name == "" {
		return Process{}, fmt.Errorf("line %d: process name cannot be empty", lineNum)
	}

	remainder := line[firstColon+1:]

	// Parse inputs section
	if !strings.HasPrefix(remainder, "(") {
		return Process{}, fmt.Errorf("line %d: inputs section must start with '('", lineNum)
	}
	inputsEnd := findClosingParen(remainder)
	if inputsEnd == -1 {
		return Process{}, fmt.Errorf("line %d: inputs section missing closing ')'", lineNum)
	}
	inputsSection := remainder[:inputsEnd+1]
	remainder = remainder[inputsEnd+1:]

	if !strings.HasPrefix(remainder, ":") {
		return Process{}, fmt.Errorf("line %d: expected ':' after inputs section", lineNum)
	}
	remainder = remainder[1:]

	// Parse outputs section
	if !strings.HasPrefix(remainder, "(") {
		return Process{}, fmt.Errorf("line %d: outputs section must start with '('", lineNum)
	}
	outputsEnd := findClosingParen(remainder)
	if outputsEnd == -1 {
		return Process{}, fmt.Errorf("line %d: outputs section missing closing ')'", lineNum)
	}
	outputsSection := remainder[:outputsEnd+1]
	remainder = remainder[outputsEnd+1:]

	if !strings.HasPrefix(remainder, ":") {
		return Process{}, fmt.Errorf("line %d: expected ':' after outputs section", lineNum)
	}
	durationStr := remainder[1:]

	// Validate parts
	if strings.Contains(durationStr, ":") {
		return Process{}, fmt.Errorf("line %d: invalid process format, too many colons", lineNum)
	}

	inputs, err := parseItemQuantities(inputsSection, lineNum, "inputs")
	if err != nil {
		return Process{}, err
	}

	outputs, err := parseItemQuantities(outputsSection, lineNum, "outputs")
	if err != nil {
		return Process{}, err
	}

	duration, err := parseDuration(durationStr, lineNum)
	if err != nil {
		return Process{}, err
	}

	return Process{
		Name:     name,
		Inputs:   inputs,
		Outputs:  outputs,
		Duration: duration,
	}, nil
}

// findClosingParen finds the index of the matching ')' for the first '('.
// Returns -1 if the parentheses are unbalanced.
func findClosingParen(s string) int {
	count := 0
	for i, ch := range s {
		if ch == '(' {
			count++
		} else if ch == ')' {
			count--
			if count == 0 {
				return i
			}
		}
	}
	return -1
}
