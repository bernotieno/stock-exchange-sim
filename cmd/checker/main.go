package main

import (
	"fmt"
	"os"

	"github.com/bernotieno/stock-exchange-simulator/internal/checker"
	"github.com/bernotieno/stock-exchange-simulator/internal/parser"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <config_file> <log_file>\n", os.Args[0])
		os.Exit(1)
	}

	configFile := os.Args[1]
	logFile := os.Args[2]

	config, err := loadConfig(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	entries, err := checker.ParseLogFile(logFile, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing log file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Parsed %d log entries successfully\n", len(entries))
}

func loadConfig(filename string) (*parser.Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return parser.ParseConfig(file)
}
