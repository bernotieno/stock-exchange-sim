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
