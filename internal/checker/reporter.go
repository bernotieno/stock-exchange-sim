package checker

import (
	"fmt"
	"sort"
)

// printSuccessMessage prints a success message with final cycle and stock information
func (c *Checker) printSuccessMessage(logEntries []LogEntry) {
	fmt.Println("✓ Log validation successful!")

	if len(logEntries) == 0 {
		fmt.Println("No processes were executed.")
		return
	}

	// Find the final cycle
	finalCycle := 0
	for _, entry := range logEntries {
		if entry.Cycle > finalCycle {
			finalCycle = entry.Cycle
		}
	}

	fmt.Printf("Final cycle: %d\n", finalCycle)
	fmt.Printf("Total processes executed: %d\n", len(logEntries))

	// Calculate final stock state by simulating the execution
	shadowStock := make(map[string]int)
	for item, quantity := range c.config.Stocks {
		shadowStock[item] = quantity
	}

	// Subtract inputs for each process execution
	for _, entry := range logEntries {
		process := c.config.GetProcessByName(entry.ProcessName)
		for item, required := range process.Inputs {
			shadowStock[item] -= required
		}
	}

	// Print final stock levels
	fmt.Println("\nFinal stock levels after input consumption:")

	// Get all items and sort them for consistent output
	items := make([]string, 0, len(shadowStock))
	for item := range shadowStock {
		items = append(items, item)
	}
	sort.Strings(items)

	for _, item := range items {
		quantity := shadowStock[item]
		fmt.Printf("  %s: %d\n", item, quantity)
	}
}
