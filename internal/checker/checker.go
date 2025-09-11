package checker

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/bernotieno/stock-exchange-simulator/internal/parser"
)

// LogEntry represents a single execution event from the log file
type LogEntry struct {
	Cycle   int
	Process string
}

// ParseLogFile loads and parses a log file into execution events
func ParseLogFile(filename string, config *parser.Config) ([]LogEntry, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries []LogEntry
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "No more process") {
			continue
		}

		entry, err := parseLogLine(line, config)
		if err != nil {
			return nil, err
		}

		entries = append(entries, entry)
	}

	return entries, scanner.Err()
}

// parseLogLine parses a single log line in format "cycle:process"
func parseLogLine(line string, config *parser.Config) (LogEntry, error) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return LogEntry{}, fmt.Errorf("invalid log format: %s", line)
	}

	cycle, err := strconv.Atoi(parts[0])
	if err != nil {
		return LogEntry{}, fmt.Errorf("invalid cycle number '%s': %v", parts[0], err)
	}

	processName := parts[1]
	if config.GetProcessByName(processName) == nil {
		return LogEntry{}, fmt.Errorf("process '%s' not found in config", processName)
	}

	return LogEntry{Cycle: cycle, Process: processName}, nil
}
