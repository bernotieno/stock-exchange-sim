package checker

import (
	"fmt"
)

// validateLogEntries validates the log entries by simulating execution with proper timing
func (c *Checker) validateLogEntries(logEntries []LogEntry) error {
	// Check that cycles are in non-decreasing order
	for i := 1; i < len(logEntries); i++ {
		if logEntries[i].Cycle < logEntries[i-1].Cycle {
			return CheckerError{
				LineNumber:  logEntries[i].LineNumber,
				Cycle:       logEntries[i].Cycle,
				ProcessName: logEntries[i].ProcessName,
				Message:     fmt.Sprintf("cycle %d is less than previous cycle %d", logEntries[i].Cycle, logEntries[i-1].Cycle),
			}
		}
	}

	// Initialize shadow stock map with starting stocks
	shadowStock := make(map[string]int)
	for item, quantity := range c.config.Stocks {
		shadowStock[item] = quantity
	}

	// Track running processes
	var runningProcesses []RunningProcessChecker
	currentCycle := 0

	// Simulate execution step by step
	for _, entry := range logEntries {
		// Advance to the cycle of this entry
		for currentCycle < entry.Cycle {
			currentCycle++
			// Complete any processes that finish at this cycle
			runningProcesses = c.completeFinishedProcesses(runningProcesses, currentCycle, shadowStock)
		}

		process := c.config.GetProcessByName(entry.ProcessName)
		if process == nil {
			// This should have been caught in parsing, but double-check
			return CheckerError{
				LineNumber:  entry.LineNumber,
				Cycle:       entry.Cycle,
				ProcessName: entry.ProcessName,
				Message:     "process not found in configuration",
			}
		}

		// Check if we have enough stock for inputs
		for item, required := range process.Inputs {
			available := shadowStock[item]
			if available < required {
				return CheckerError{
					LineNumber:  entry.LineNumber,
					Cycle:       entry.Cycle,
					ProcessName: entry.ProcessName,
					Message:     fmt.Sprintf("not enough %s (need %d, have %d)", item, required, available),
				}
			}
		}

		// Consume inputs (subtract when process starts)
		for item, required := range process.Inputs {
			shadowStock[item] -= required
		}

		// Calculate when this process will complete
		durationCycles := int(process.Duration.Seconds())
		if durationCycles == 0 {
			durationCycles = 1 // Minimum 1 cycle
		}

		// Add to running processes
		runningProcess := RunningProcessChecker{
			Process:    process,
			StartCycle: entry.Cycle,
			EndCycle:   entry.Cycle + durationCycles,
		}
		runningProcesses = append(runningProcesses, runningProcess)
	}

	return nil
}

// completeFinishedProcesses completes any processes that have finished and returns the updated list
func (c *Checker) completeFinishedProcesses(runningProcesses []RunningProcessChecker, currentCycle int, shadowStock map[string]int) []RunningProcessChecker {
	var stillRunning []RunningProcessChecker

	for _, rp := range runningProcesses {
		if rp.EndCycle <= currentCycle {
			// Process has completed, produce outputs
			for item, produced := range rp.Process.Outputs {
				shadowStock[item] += produced
			}
		} else {
			// Process still running
			stillRunning = append(stillRunning, rp)
		}
	}

	return stillRunning
}
