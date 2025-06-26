package parser

import "fmt"

// Config represents a complete parsed manufacturing configuration.
// It contains all the information needed by a production scheduler.
type Config struct {
	// Stocks maps item names to their initial quantities in inventory
	Stocks map[string]int `json:"stocks"`

	// Processes contains all defined manufacturing processes
	Processes []Process `json:"processes"`

	// Optimize lists the items that should be optimized for maximum production
	Optimize []string `json:"optimize"`
}

// Validate performs consistency checks on the parsed configuration.
// It ensures that optimization targets reference valid items and that
// processes have valid input/output relationships.
func (c *Config) Validate() error {
	// Collect all known items from stocks and processes
	knownItems := make(map[string]bool)

	// Add stock items
	for item := range c.Stocks {
		knownItems[item] = true
	}

	// Add items from process inputs and outputs
	for _, process := range c.Processes {
		for item := range process.Inputs {
			knownItems[item] = true
		}
		for item := range process.Outputs {
			knownItems[item] = true
		}
	}

	// Validate optimization targets reference known items
	for _, target := range c.Optimize {
		if !knownItems[target] {
			return fmt.Errorf("optimization target '%s' is not defined in stocks or processes", target)
		}
	}

	// Validate that processes have at least one input or output
	for _, process := range c.Processes {
		if len(process.Inputs) == 0 && len(process.Outputs) == 0 {
			return fmt.Errorf("process '%s' must have at least one input or output", process.Name)
		}
	}

	return nil
}

// GetProcessByName returns a process by its name, or nil if not found.
func (c *Config) GetProcessByName(name string) *Process {
	for i := range c.Processes {
		if c.Processes[i].Name == name {
			return &c.Processes[i]
		}
	}
	return nil
}

// GetStockQuantity returns the initial stock quantity for an item.
// Returns 0 if the item is not in the initial stock.
func (c *Config) GetStockQuantity(item string) int {
	return c.Stocks[item] // Returns 0 for missing keys
}

// GetAllItems returns a slice of all unique item names mentioned in the configuration.
func (c *Config) GetAllItems() []string {
	items := make(map[string]bool)

	// Add stock items
	for item := range c.Stocks {
		items[item] = true
	}

	// Add items from processes
	for _, process := range c.Processes {
		for item := range process.Inputs {
			items[item] = true
		}
		for item := range process.Outputs {
			items[item] = true
		}
	}

	// Convert to slice
	result := make([]string, 0, len(items))
	for item := range items {
		result = append(result, item)
	}

	return result
}
