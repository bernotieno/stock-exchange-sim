package scheduler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// WriteLogFile writes the execution log to a .log file
func (sim *Simulator) WriteLogFile(inputFilename string) error {
	// Create examples directory if it doesn't exist
	err := os.MkdirAll("examples", 0755)
	if err != nil {
		return fmt.Errorf("failed to create examples directory: %w", err)
	}

	// Generate output filename by replacing extension with .log and placing in examples/
	baseName := filepath.Base(inputFilename)
	outputFilename := filepath.Join("examples", strings.TrimSuffix(baseName, filepath.Ext(baseName)) + ".log")
	
	file, err := os.Create(outputFilename)
	if err != nil {
		return fmt.Errorf("failed to create log file %s: %w", outputFilename, err)
	}
	defer file.Close()

	// Write each execution step in cycle:process_name format
	for _, step := range sim.State.ExecutionLog {
		_, err := fmt.Fprintf(file, "%d:%s\n", step.Cycle, step.ProcessName)
		if err != nil {
			return fmt.Errorf("failed to write to log file: %w", err)
		}
	}

	// Store the output filename for later use in PrintResults
	sim.State.LogFilename = outputFilename

	return nil
}

// PrintResults prints the final stock state and cycle count to terminal
func (sim *Simulator) PrintResults() {
	fmt.Printf("Simulation completed in %d cycles\n", sim.State.CurrentCycle)
	fmt.Println("\nFinal stock levels:")

	// Sort items for consistent output
	items := make([]string, 0, len(sim.State.CurrentStock))
	for item := range sim.State.CurrentStock {
		items = append(items, item)
	}

	// Simple sort (bubble sort for small lists)
	for i := 0; i < len(items)-1; i++ {
		for j := 0; j < len(items)-i-1; j++ {
			if items[j] > items[j+1] {
				items[j], items[j+1] = items[j+1], items[j]
			}
		}
	}

	for _, item := range items {
		quantity := sim.State.CurrentStock[item]
		fmt.Printf("  %s: %d\n", item, quantity)
	}

	// Print log file location
	if sim.State.LogFilename != "" {
		fmt.Printf("\nLog saved to: %s\n", sim.State.LogFilename)
	}
}

// GetExecutionLog returns the complete execution log
func (sim *Simulator) GetExecutionLog() []ExecutionStep {
	return sim.State.ExecutionLog
}

// GetFinalStock returns the final stock state
func (sim *Simulator) GetFinalStock() map[string]int {
	return sim.State.CurrentStock
}

// GetFinalCycle returns the final cycle count
func (sim *Simulator) GetFinalCycle() int {
	return sim.State.CurrentCycle
}
