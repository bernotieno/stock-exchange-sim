package parser

import (
	"time"
)

// Process represents a manufacturing process with inputs, outputs, and duration
type Process struct {
	Name     string         // Process name
	Inputs   map[string]int // Input items and their required quantities
	Outputs  map[string]int // Output items and their produced quantities
	Duration time.Duration  // Time required to complete one cycle
}
