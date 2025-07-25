package main

import (
	"fmt"
	"os"

	"github.com/bernotieno/stock-exchange-simulator/internal/parser"
	"github.com/bernotieno/stock-exchange-simulator/internal/scheduler"
)

func main() {
	// Check command line arguments
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <config_file>\n", os.Args[0])
		os.Exit(1)
	}

	configFile := os.Args[1]

	// Open and parse the configuration file
	file, err := os.Open(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening config file %s: %v\n", configFile, err)
		os.Exit(1)
	}
	defer file.Close()

	config, err := parser.ParseConfig(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing config file: %v\n", err)
		os.Exit(1)
	}



	// Create and run the simulator
	simulator := scheduler.NewSimulator(config)

	err = simulator.RunSimulationWithOutput(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Simulation error: %v\n", err)
		os.Exit(1)
	}
}

