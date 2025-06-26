package parser

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
