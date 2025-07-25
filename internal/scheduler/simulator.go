package scheduler

import (
	"fmt"

	"github.com/bernotieno/stock-exchange-simulator/internal/parser"
)

// NewSimulator creates a new simulator with the given configuration
func NewSimulator(config *parser.Config) *Simulator {
	return &Simulator{
		State: NewSchedulerState(config),
	}
}

// RunSimulation executes the main simulation loop
func (sim *Simulator) RunSimulation() error {
	maxCycles := 10000 // Safety limit to prevent infinite loops
	
	for sim.State.CurrentCycle < maxCycles {
		// Complete any processes that finished this cycle
		sim.State.CompleteFinishedProcesses()
		
		// Get processes that can be started this cycle, optimized for time
		runnableProcesses := sim.State.GetOptimizedRunnableProcesses()
		
		// Start all runnable processes in priority order
		processesStarted := false
		for _, process := range runnableProcesses {
			// Check again in case starting previous processes consumed resources
			if sim.State.CanRunProcess(process) {
				sim.State.StartProcess(process)
				processesStarted = true
			}
		}
		
		// If no processes are running and none can be started, simulation is complete
		if !sim.State.HasRunningProcesses() && !processesStarted {
			break
		}
		
		// Advance to next cycle
		sim.State.AdvanceCycle()
	}
	
	// Complete any remaining processes
	sim.State.CompleteFinishedProcesses()
	
	return nil
}

// RunSimulationWithOutput runs the simulation and generates output
func (sim *Simulator) RunSimulationWithOutput(inputFilename string) error {
	// Run the simulation
	if err := sim.RunSimulation(); err != nil {
		return fmt.Errorf("simulation failed: %w", err)
	}
	
	// Write log file
	if err := sim.WriteLogFile(inputFilename); err != nil {
		return fmt.Errorf("failed to write log file: %w", err)
	}
	
	// Print results to terminal
	sim.PrintResults()
	
	return nil
}
