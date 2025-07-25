package scheduler

import (
	"github.com/bernotieno/stock-exchange-simulator/internal/parser"
)

// GetOptimizedRunnableProcesses returns processes prioritized for time optimization
func (s *SchedulerState) GetOptimizedRunnableProcesses() []*parser.Process {
	runnable := s.GetRunnableProcesses()
	
	// Check if we have a time optimization goal
	hasTimeOptimization := false
	for _, target := range s.Config.Optimize {
		if target == "time" {
			hasTimeOptimization = true
			break
		}
	}
	
	if !hasTimeOptimization {
		return runnable // No time optimization, return as-is
	}
	
	// Sort processes to prioritize time optimization
	// Priority order:
	// 1. Processes that produce items needed by other processes (enablers)
	// 2. Processes with shorter duration (faster completion)
	// 3. Processes that produce optimization targets
	
	return s.sortProcessesForTimeOptimization(runnable)
}

// sortProcessesForTimeOptimization sorts processes to minimize total execution time
func (s *SchedulerState) sortProcessesForTimeOptimization(processes []*parser.Process) []*parser.Process {
	if len(processes) <= 1 {
		return processes
	}
	
	// Create a copy to avoid modifying the original slice
	sorted := make([]*parser.Process, len(processes))
	copy(sorted, processes)
	
	// Simple bubble sort with custom comparison
	for i := 0; i < len(sorted)-1; i++ {
		for j := 0; j < len(sorted)-i-1; j++ {
			if s.shouldPrioritizeProcess(sorted[j+1], sorted[j]) {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			}
		}
	}
	
	return sorted
}

// shouldPrioritizeProcess returns true if process a should be prioritized over process b
func (s *SchedulerState) shouldPrioritizeProcess(a, b *parser.Process) bool {
	// Priority 1: Processes that enable other processes (produce items needed as inputs)
	aEnablesOthers := s.processEnablesOthers(a)
	bEnablesOthers := s.processEnablesOthers(b)
	
	if aEnablesOthers && !bEnablesOthers {
		return true
	}
	if !aEnablesOthers && bEnablesOthers {
		return false
	}
	
	// Priority 2: Shorter duration processes first
	if a.Duration != b.Duration {
		return a.Duration < b.Duration
	}
	
	// Priority 3: Processes that produce optimization targets
	aProducesTarget := s.processProducesOptimizationTarget(a)
	bProducesTarget := s.processProducesOptimizationTarget(b)
	
	if aProducesTarget && !bProducesTarget {
		return true
	}
	if !aProducesTarget && bProducesTarget {
		return false
	}
	
	// Default: maintain original order
	return false
}

// processEnablesOthers checks if a process produces items that other processes need
func (s *SchedulerState) processEnablesOthers(process *parser.Process) bool {
	for outputItem := range process.Outputs {
		// Check if any other process needs this item as input
		for i := range s.Config.Processes {
			otherProcess := &s.Config.Processes[i]
			if otherProcess.Name != process.Name {
				if _, needed := otherProcess.Inputs[outputItem]; needed {
					return true
				}
			}
		}
	}
	return false
}

// processProducesOptimizationTarget checks if a process produces any optimization targets
func (s *SchedulerState) processProducesOptimizationTarget(process *parser.Process) bool {
	for _, target := range s.Config.Optimize {
		if target == "time" {
			continue // Skip time optimization target
		}
		if _, produces := process.Outputs[target]; produces {
			return true
		}
	}
	return false
}
