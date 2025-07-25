package scheduler

import (
	"github.com/bernotieno/stock-exchange-simulator/internal/parser"
)

// RunningProcess represents a process that is currently executing
type RunningProcess struct {
	Process   *parser.Process // Reference to the process definition
	StartCycle int            // Cycle when this process started
	EndCycle   int            // Cycle when this process will complete
}

// ExecutionStep represents a single step in the simulation log
type ExecutionStep struct {
	Cycle       int    // The cycle number when this step occurred
	ProcessName string // Name of the process that started
}

// SchedulerState holds the current state of the simulation
type SchedulerState struct {
	CurrentCycle     int                        // Current simulation cycle
	CurrentStock     map[string]int             // Current stock levels
	RunningProcesses []RunningProcess           // Processes currently running
	ExecutionLog     []ExecutionStep            // Log of all execution steps
	Config           *parser.Config             // Configuration data
}

// Simulator represents the main simulation engine
type Simulator struct {
	State *SchedulerState
}
