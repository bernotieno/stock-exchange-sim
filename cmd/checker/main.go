package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/bernotieno/stock-exchange-simulator/internal/checker"
	"github.com/bernotieno/stock-exchange-simulator/internal/parser"
)

func main() {
	// Define command line flags
	var configFile string
	var logFile string

	flag.StringVar(&configFile, "config", "", "Path to the configuration file")
	flag.StringVar(&logFile, "log", "", "Path to the log file to check")
	flag.Parse()

	// Validate required flags
	if configFile == "" || logFile == "" {
		fmt.Fprintf(os.Stderr, "Stock Exchange Log Checker\n")
		fmt.Fprintf(os.Stderr, "Usage: %s -config <config_file> -log <log_file>\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "  -config string\n")
		fmt.Fprintf(os.Stderr, "        Path to the configuration file (required)\n")
		fmt.Fprintf(os.Stderr, "  -log string\n")
		fmt.Fprintf(os.Stderr, "        Path to the log file to check (required)\n\n")
		fmt.Fprintf(os.Stderr, "Example:\n")
		fmt.Fprintf(os.Stderr, "  %s -config examples/finite.conf -log examples/finite.log\n", os.Args[0])
		os.Exit(1)
	}

	// Load and parse the configuration file
	configFileHandle, err := os.Open(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening config file %s: %v\n", configFile, err)
		os.Exit(1)
	}
	defer configFileHandle.Close()

	config, err := parser.ParseConfig(configFileHandle)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing config file: %v\n", err)
		os.Exit(1)
	}

	// Create and run the checker
	checkerInstance := checker.NewChecker(config)

	err = checkerInstance.CheckLogFile(logFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Checker error: %v\n", err)
		os.Exit(1)
	}
}
