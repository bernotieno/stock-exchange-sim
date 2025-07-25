package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/bernotieno/stock-exchange-simulator/internal/parser"
	"github.com/bernotieno/stock-exchange-simulator/internal/scheduler"
)

func main() {
	// Check command line arguments
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Stock Exchange Simulator\n")
		fmt.Fprintf(os.Stderr, "Usage: %s <config_file> <timeout>\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Arguments:\n")
		fmt.Fprintf(os.Stderr, "  config_file: Path to the configuration file\n")
		fmt.Fprintf(os.Stderr, "  timeout:     Maximum simulation time in seconds\n\n")
		fmt.Fprintf(os.Stderr, "Example:\n")
		fmt.Fprintf(os.Stderr, "  %s examples/finite.conf 60\n", os.Args[0])
		os.Exit(1)
	}

	configFile := os.Args[1]
	timeoutStr := os.Args[2]

	// Parse timeout
	timeoutSeconds, err := strconv.Atoi(timeoutStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: timeout must be a valid integer (seconds): %v\n", err)
		os.Exit(1)
	}
	if timeoutSeconds <= 0 {
		fmt.Fprintf(os.Stderr, "Error: timeout must be positive\n")
		os.Exit(1)
	}
	timeout := time.Duration(timeoutSeconds) * time.Second

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

	err = simulator.RunSimulationWithOutputAndTimeout(configFile, timeout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Simulation error: %v\n", err)
		os.Exit(1)
	}
}

