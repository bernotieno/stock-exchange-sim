package parser

import (
	"fmt"
	"io"
	"strings"
)

// Config holds parsed data from a manufacturing config file.
// It includes starting stocks, process definitions, and optimization targets.
type Config struct {
	Stocks    map[string]int `json:"stocks"`    // Initial item quantities
	Processes []Process      `json:"processes"` // All defined processes
	Optimize  []string       `json:"optimize"`  // Items to maximize in output
}

// ParseConfig reads and parses a config file from the given reader.
// It extracts stocks, processes, and optimize sections in any order,
// validates the results, and returns a complete Config object.
func ParseConfig(r io.Reader) (*Config, error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	config := &Config{
		Stocks:    make(map[string]int),
		Processes: []Process{},
		Optimize:  []string{},
	}

	sections := splitIntoSections(string(content))

	if stocksSection, exists := sections["stocks"]; exists {
		config.Stocks, err = ParseStocks(strings.NewReader(stocksSection))
		if err != nil {
			return nil, fmt.Errorf("failed to parse stocks: %w", err)
		}
	}

	if processesSection, exists := sections["processes"]; exists {
		config.Processes, err = ParseProcesses(strings.NewReader(processesSection))
		if err != nil {
			return nil, fmt.Errorf("failed to parse processes: %w", err)
		}
	}

	if optimizeSection, exists := sections["optimize"]; exists {
		config.Optimize, err = ParseOptimize(strings.NewReader(optimizeSection))
		if err != nil {
			return nil, fmt.Errorf("failed to parse optimization targets: %w", err)
		}
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}

// splitIntoSections separates the raw config content into named sections.
// Recognizes headers like "# Stocks", "# Processes", and "# Optimize".
func splitIntoSections(content string) map[string]string {
	sections := make(map[string]string)
	currentSection := ""
	var sectionContent strings.Builder

	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "# Stocks") {
			if currentSection != "" {
				sections[currentSection] = sectionContent.String()
				sectionContent.Reset()
			}
			currentSection = "stocks"
		} else if strings.HasPrefix(trimmed, "# Processes") {
			if currentSection != "" {
				sections[currentSection] = sectionContent.String()
				sectionContent.Reset()
			}
			currentSection = "processes"
		} else if strings.HasPrefix(trimmed, "# Optimize") {
			if currentSection != "" {
				sections[currentSection] = sectionContent.String()
				sectionContent.Reset()
			}
			currentSection = "optimize"
		} else if currentSection != "" {
			sectionContent.WriteString(line)
			sectionContent.WriteString("\n")
		}
	}

	if currentSection != "" {
		sections[currentSection] = sectionContent.String()
	}

	return sections
}

// Validate checks the integrity of the parsed config.
// Ensures all optimize targets exist and each process has inputs or outputs.
func (c *Config) Validate() error {
	knownItems := make(map[string]bool)

	for item := range c.Stocks {
		knownItems[item] = true
	}
	for _, process := range c.Processes {
		for item := range process.Inputs {
			knownItems[item] = true
		}
		for item := range process.Outputs {
			knownItems[item] = true
		}
	}

	for _, target := range c.Optimize {
		if !knownItems[target] {
			return fmt.Errorf("optimization target '%s' is not defined in stocks or processes", target)
		}
	}

	for _, process := range c.Processes {
		if len(process.Inputs) == 0 && len(process.Outputs) == 0 {
			return fmt.Errorf("process '%s' must have at least one input or output", process.Name)
		}
	}

	return nil
}

// GetProcessByName finds and returns a process by name.
// Returns nil if the process doesn't exist.
func (c *Config) GetProcessByName(name string) *Process {
	for i := range c.Processes {
		if c.Processes[i].Name == name {
			return &c.Processes[i]
		}
	}
	return nil
}

// GetStockQuantity returns the starting quantity of a given stock item.
// If the item isn't defined, returns 0.
func (c *Config) GetStockQuantity(item string) int {
	return c.Stocks[item]
}

// GetAllItems returns a unique list of all items used in stocks and processes.
func (c *Config) GetAllItems() []string {
	items := make(map[string]bool)

	for item := range c.Stocks {
		items[item] = true
	}
	for _, process := range c.Processes {
		for item := range process.Inputs {
			items[item] = true
		}
		for item := range process.Outputs {
			items[item] = true
		}
	}

	result := make([]string, 0, len(items))
	for item := range items {
		result = append(result, item)
	}
	return result
}
