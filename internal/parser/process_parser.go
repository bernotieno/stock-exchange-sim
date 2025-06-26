package parser

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

// Process represents a manufacturing process with inputs, outputs, and duration
type Process struct {
	Name     string         // Process name
	Inputs   map[string]int // Input items and their required quantities
	Outputs  map[string]int // Output items and their produced quantities
	Duration time.Duration  // Time required to complete one cycle
}

// parseDuration parses a duration string, supporting both plain seconds and Go duration format
func parseDuration(durationStr string, lineNum int) (time.Duration, error) {
	durationStr = strings.TrimSpace(durationStr)
	if durationStr == "" {
		return 0, fmt.Errorf("line %d: duration cannot be empty", lineNum)
	}

	// Try parsing as Go duration first (e.g., "30s", "2m", "1h30m")
	if duration, err := time.ParseDuration(durationStr); err == nil {
		if duration <= 0 {
			return 0, fmt.Errorf("line %d: duration must be positive", lineNum)
		}
		return duration, nil
	}

	// Try parsing as plain seconds
	if seconds, err := strconv.ParseInt(durationStr, 10, 64); err == nil {
		if seconds <= 0 {
			return 0, fmt.Errorf("line %d: duration must be positive", lineNum)
		}
		return time.Duration(seconds) * time.Second, nil
	}

	return 0, fmt.Errorf("line %d: invalid duration format '%s'", lineNum, durationStr)
}

func parseProcessLine(line string, lineNum int) (Process, error) {
	// Find the process name (everything before the first colon)
	firstColon := strings.Index(line, ":")
	if firstColon == -1 {
		return Process{}, fmt.Errorf("line %d: invalid process format, expected 'name:inputs:outputs:duration'", lineNum)
	}

	name := strings.TrimSpace(line[:firstColon])
	if name == "" {
		return Process{}, fmt.Errorf("line %d: process name cannot be empty", lineNum)
	}

	remainder := line[firstColon+1:]

	// Find the inputs section (from first '(' to matching ')')
	if !strings.HasPrefix(remainder, "(") {
		return Process{}, fmt.Errorf("line %d: inputs section must start with '('", lineNum)
	}

	parenCount := 0
	inputsEnd := -1
	for i, char := range remainder {
		if char == '(' {
			parenCount++
		} else if char == ')' {
			parenCount--
			if parenCount == 0 {
				inputsEnd = i
				break
			}
		}
	}

	if inputsEnd == -1 {
		return Process{}, fmt.Errorf("line %d: inputs section missing closing ')'", lineNum)
	}

	inputsSection := remainder[:inputsEnd+1]
	remainder = remainder[inputsEnd+1:]

	// Next should be a colon
	if !strings.HasPrefix(remainder, ":") {
		return Process{}, fmt.Errorf("line %d: expected ':' after inputs section", lineNum)
	}
	remainder = remainder[1:]

	// Find the outputs section (from next '(' to matching ')')
	if !strings.HasPrefix(remainder, "(") {
		return Process{}, fmt.Errorf("line %d: outputs section must start with '('", lineNum)
	}

	parenCount = 0
	outputsEnd := -1
	for i, char := range remainder {
		if char == '(' {
			parenCount++
		} else if char == ')' {
			parenCount--
			if parenCount == 0 {
				outputsEnd = i
				break
			}
		}
	}

	if outputsEnd == -1 {
		return Process{}, fmt.Errorf("line %d: outputs section missing closing ')'", lineNum)
	}

	outputsSection := remainder[:outputsEnd+1]
	remainder = remainder[outputsEnd+1:]

	// Next should be a colon followed by duration
	if !strings.HasPrefix(remainder, ":") {
		return Process{}, fmt.Errorf("line %d: expected ':' after outputs section", lineNum)
	}
	durationStr := remainder[1:]

	// Check for extra colons in duration (indicating too many parts)
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
