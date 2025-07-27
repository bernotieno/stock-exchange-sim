package scheduler

import (
	"fmt"
	"os"
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
			s.logCompletion(rp.Process.Name, rp.EndTime)
		} else {
			stillRunning = append(stillRunning, rp)
		}
	}

	s.RunningProcesses = stillRunning
}

// logCompletion logs when a process completes
func (s *SchedulerState) logCompletion(processName string, endTime int) {
	// Process completion is implicit in the execution log
	// The start time is already logged, completion happens at endTime
}

// isTimeOptimized checks if time is in optimization targets
func (s *SchedulerState) isTimeOptimized() bool {
	for _, target := range s.OptimizeTargets {
		if target == "time" {
			return true
		}
	}
	return false
}

// prioritizeForTime sorts processes to minimize total execution time
func (s *SchedulerState) prioritizeForTime(processes []*Process) []*Process {
	// Run shorter processes first to minimize total cycle count
	for i := 0; i < len(processes)-1; i++ {
		for j := i + 1; j < len(processes); j++ {
			if processes[i].Duration > processes[j].Duration {
				processes[i], processes[j] = processes[j], processes[i]
			}
		}
	}
	return processes
}

// WriteLogFile writes execution log to file in cycle:process_name format
func (s *SchedulerState) WriteLogFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, step := range s.ExecutionLog {
		_, err := fmt.Fprintf(file, "%d:%s\n", step.Cycle, step.Process)
		if err != nil {
			return err
		}
	}

	_, err = fmt.Fprintf(file, "No more process doable at cycle %d\n", s.CurrentCycle)
	return err
}

// PrintResults prints final stock levels and cycle count to terminal
func (s *SchedulerState) PrintResults() {
	fmt.Println("Main Processes:")
	for _, step := range s.ExecutionLog {
		fmt.Printf(" %d:%s\n", step.Cycle, step.Process)
	}
	fmt.Printf("No more process doable at cycle %d\n", s.CurrentCycle)
	fmt.Println("Stock:")
	for item, quantity := range s.Stocks {
		fmt.Printf(" %s => %d\n", item, quantity)
	}
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

// RunSimulation executes the main simulation loop
func (s *SchedulerState) RunSimulation(maxCycles int) {
	for s.CurrentCycle < maxCycles {
		s.CompleteFinishedProcesses()
		
		available := s.GetAvailableProcesses()
		if len(available) == 0 {
			break
		}
		
		if s.isTimeOptimized() {
			available = s.prioritizeForTime(available)
		}
		
		for _, process := range available {
			s.StartProcess(process)
		}
		
		s.AdvanceCycle()
	}
}

// AdvanceCycle moves to the next cycle or jumps to next process completion
func (s *SchedulerState) AdvanceCycle() {
	if len(s.RunningProcesses) == 0 {
		s.CurrentCycle++
		return
	}
	
	nextCompletion := s.RunningProcesses[0].EndTime
	for _, rp := range s.RunningProcesses[1:] {
		if rp.EndTime < nextCompletion {
			nextCompletion = rp.EndTime
		}
	}
	
	if nextCompletion > s.CurrentCycle {
		s.CurrentCycle = nextCompletion
	} else {
		s.CurrentCycle++
	}
}

// StepSimulation executes one simulation step
func (s *SchedulerState) StepSimulation() bool {
	s.CompleteFinishedProcesses()
	
	available := s.GetAvailableProcesses()
	if len(available) == 0 {
		return false
	}
	
	for _, process := range available {
		s.StartProcess(process)
	}
	
	s.AdvanceCycle()
	return true
}
