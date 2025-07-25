package checker

import (
	"fmt"

	"github.com/bernotieno/stock-exchange-simulator/internal/parser"
)

// Checker validates execution logs against configuration
type Checker struct {
	config *parser.Config
}

// NewChecker creates a new checker with the given configuration
func NewChecker(config *parser.Config) *Checker {
	return &Checker{
		config: config,
	}
}

// CheckLogFile validates a log file against the configuration
func (c *Checker) CheckLogFile(logFilePath string) error {
	// Parse the log file
	logEntries, err := c.parseLogFile(logFilePath)
	if err != nil {
		return fmt.Errorf("failed to parse log file: %w", err)
	}

	// Validate the log entries
	err = c.validateLogEntries(logEntries)
	if err != nil {
		return err
	}

	// If we get here, validation passed
	c.printSuccessMessage(logEntries)
	return nil
}