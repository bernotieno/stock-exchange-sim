package scheduler

import (
	"time"

	"github.com/bernotieno/stock-exchange-simulator/internal/parser"
)

// Process represents a manufacturing process with its requirements and outputs
type Process struct {
	Name     string
	Inputs   map[string]int
	Outputs  map[string]int
	Duration time.Duration
}

// RunningProcess tracks a process currently being executed
type RunningProcess struct {
	Process   *Process
	StartTime int
	EndTime   int
}

// ExecutionStep represents a single step in the execution schedule
type ExecutionStep struct {
	Cycle   int
	Process string
}

// SchedulerState maintains the current state of the simulation
type SchedulerState struct {
	CurrentCycle    int
	Stocks          map[string]int
	RunningProcesses []RunningProcess
	ExecutionLog    []ExecutionStep
	Processes       []Process
	OptimizeTargets []string
}

// NewSchedulerState creates a new scheduler state from config
func NewSchedulerState(config *parser.Config) *SchedulerState {
	stocks := make(map[string]int)
	for k, v := range config.Stocks {
		stocks[k] = v
	}

	processes := make([]Process, len(config.Processes))
	for i, p := range config.Processes {
		processes[i] = Process{
			Name:     p.Name,
			Inputs:   p.Inputs,
			Outputs:  p.Outputs,
			Duration: p.Duration,
		}
	}

	return &SchedulerState{
		CurrentCycle:     0,
		Stocks:           stocks,
		RunningProcesses: []RunningProcess{},
		ExecutionLog:     []ExecutionStep{},
		Processes:        processes,
		OptimizeTargets:  config.Optimize,
	}
}

// CanRunProcess checks if a process can be started with current stock levels
func (s *SchedulerState) CanRunProcess(process *Process) bool {
	for item, needed := range process.Inputs {
		if s.Stocks[item] < needed {
			return false
		}
	}
	return true
}

// ConsumeInputs removes required inputs from stock when starting a process
func (s *SchedulerState) ConsumeInputs(process *Process) {
	for item, needed := range process.Inputs {
		s.Stocks[item] -= needed
	}
}

// ProduceOutputs adds outputs to stock when a process completes
func (s *SchedulerState) ProduceOutputs(process *Process) {
	for item, produced := range process.Outputs {
		s.Stocks[item] += produced
	}
}

// StartProcess begins execution of a process
func (s *SchedulerState) StartProcess(process *Process) {
	if !s.CanRunProcess(process) {
		return
	}

	s.ConsumeInputs(process)
	endTime := s.CurrentCycle + int(process.Duration.Seconds())

	s.RunningProcesses = append(s.RunningProcesses, RunningProcess{
		Process:   process,
		StartTime: s.CurrentCycle,
		EndTime:   endTime,
	})

	s.ExecutionLog = append(s.ExecutionLog, ExecutionStep{
		Cycle:   s.CurrentCycle,
		Process: process.Name,
	})
}

// CompleteFinishedProcesses handles processes that have finished at current cycle
func (s *SchedulerState) CompleteFinishedProcesses() {
	var stillRunning []RunningProcess

	for _, rp := range s.RunningProcesses {
		if rp.EndTime <= s.CurrentCycle {
			s.ProduceOutputs(rp.Process)
		} else {
			stillRunning = append(stillRunning, rp)
		}
	}

	s.RunningProcesses = stillRunning
}

// GetAvailableProcesses returns processes that can be started now
func (s *SchedulerState) GetAvailableProcesses() []*Process {
	var available []*Process
	for i := range s.Processes {
		if s.CanRunProcess(&s.Processes[i]) {
			available = append(available, &s.Processes[i])
		}
	}
	return available
}

// HasRunnableProcesses checks if any process can be executed
func (s *SchedulerState) HasRunnableProcesses() bool {
	return len(s.GetAvailableProcesses()) > 0
}
