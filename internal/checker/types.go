package checker

import (
	"fmt"

	"github.com/bernotieno/stock-exchange-simulator/internal/parser"
)

// LogEntry represents a single entry in the execution log
type LogEntry struct {
	Cycle       int    // The cycle number when this step occurred
	ProcessName string // Name of the process that started
	LineNumber  int    // Line number in the log file for error reporting
}

// CheckerError represents an error found during checking
type CheckerError struct {
	LineNumber  int
	Cycle       int
	ProcessName string
	Message     string
}

func (e CheckerError) Error() string {
	return fmt.Sprintf("Line %d, Cycle %d, Process %s: %s", e.LineNumber, e.Cycle, e.ProcessName, e.Message)
}

// RunningProcessChecker represents a process that is currently executing in the checker
type RunningProcessChecker struct {
	Process    *parser.Process // Reference to the process definition
	StartCycle int             // Cycle when this process started
	EndCycle   int             // Cycle when this process will complete
}
