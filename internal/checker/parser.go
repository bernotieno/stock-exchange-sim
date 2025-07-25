package checker

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// parseLogFile reads and parses a log file into LogEntry structs
func (c *Checker) parseLogFile(logFilePath string) ([]LogEntry, error) {
	file, err := os.Open(logFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	var logEntries []LogEntry
	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines
		if line == "" {
			continue
		}

		// Parse line in format "cycle:process_name"
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("line %d: invalid format, expected 'cycle:process_name', got '%s'", lineNumber, line)
		}

		// Parse cycle number
		cycleStr := strings.TrimSpace(parts[0])
		cycle, err := strconv.Atoi(cycleStr)
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid cycle number '%s': %w", lineNumber, cycleStr, err)
		}

		if cycle < 0 {
			return nil, fmt.Errorf("line %d: cycle number must be non-negative, got %d", lineNumber, cycle)
		}

		// Parse process name
		processName := strings.TrimSpace(parts[1])
		if processName == "" {
			return nil, fmt.Errorf("line %d: process name cannot be empty", lineNumber)
		}

		// Check that the process exists in config
		process := c.config.GetProcessByName(processName)
		if process == nil {
			return nil, fmt.Errorf("line %d: process '%s' not found in configuration", lineNumber, processName)
		}

		logEntries = append(logEntries, LogEntry{
			Cycle:       cycle,
			ProcessName: processName,
			LineNumber:  lineNumber,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading log file: %w", err)
	}

	return logEntries, nil
}
