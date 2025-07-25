package scheduler

import (
	"time"

	"github.com/bernotieno/stock-exchange-simulator/internal/parser"
)

// NewSchedulerState creates a new scheduler state from a config
func NewSchedulerState(config *parser.Config) *SchedulerState {
	// Initialize stock with starting values
	stock := make(map[string]int)
	for item, quantity := range config.Stocks {
		stock[item] = quantity
	}

	return &SchedulerState{
		CurrentCycle:     0,
		CurrentStock:     stock,
		RunningProcesses: make([]RunningProcess, 0),
		ExecutionLog:     make([]ExecutionStep, 0),
		Config:           config,
	}
}

// CanRunProcess checks if a process can be started given current stock levels
func (s *SchedulerState) CanRunProcess(process *parser.Process) bool {
	// Check if we have enough of each required input
	for item, required := range process.Inputs {
		if s.CurrentStock[item] < required {
			return false
		}
	}
	return true
}

// ConsumeInputs removes the required inputs from stock when starting a process
func (s *SchedulerState) ConsumeInputs(process *parser.Process) {
	for item, required := range process.Inputs {
		s.CurrentStock[item] -= required
	}
}

// ProduceOutputs adds the process outputs to stock when a process completes
func (s *SchedulerState) ProduceOutputs(process *parser.Process) {
	for item, produced := range process.Outputs {
		s.CurrentStock[item] += produced
	}
}

// StartProcess begins execution of a process
func (s *SchedulerState) StartProcess(process *parser.Process) {
	// Consume inputs
	s.ConsumeInputs(process)
	
	// Calculate end cycle (duration is in seconds, convert to cycles)
	durationCycles := int(process.Duration / time.Second)
	if durationCycles == 0 {
		durationCycles = 1 // Minimum 1 cycle
	}
	
	// Add to running processes
	runningProcess := RunningProcess{
		Process:    process,
		StartCycle: s.CurrentCycle,
		EndCycle:   s.CurrentCycle + durationCycles,
	}
	s.RunningProcesses = append(s.RunningProcesses, runningProcess)
	
	// Log the execution step
	step := ExecutionStep{
		Cycle:       s.CurrentCycle,
		ProcessName: process.Name,
	}
	s.ExecutionLog = append(s.ExecutionLog, step)
}

// CompleteFinishedProcesses checks for and completes any processes that have finished
func (s *SchedulerState) CompleteFinishedProcesses() {
	var stillRunning []RunningProcess
	
	for _, rp := range s.RunningProcesses {
		if rp.EndCycle <= s.CurrentCycle {
			// Process has completed, produce outputs
			s.ProduceOutputs(rp.Process)
		} else {
			// Process still running
			stillRunning = append(stillRunning, rp)
		}
	}
	
	s.RunningProcesses = stillRunning
}

// GetRunnableProcesses returns a list of processes that can be started this cycle
func (s *SchedulerState) GetRunnableProcesses() []*parser.Process {
	var runnable []*parser.Process
	
	for i := range s.Config.Processes {
		process := &s.Config.Processes[i]
		if s.CanRunProcess(process) {
			runnable = append(runnable, process)
		}
	}
	
	return runnable
}

// HasRunningProcesses returns true if there are any processes currently running
func (s *SchedulerState) HasRunningProcesses() bool {
	return len(s.RunningProcesses) > 0
}

// AdvanceCycle moves the simulation forward by one cycle
func (s *SchedulerState) AdvanceCycle() {
	s.CurrentCycle++
}
