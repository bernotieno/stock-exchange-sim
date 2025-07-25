package checker

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/bernotieno/stock-exchange-simulator/internal/parser"
)

// Helper function to create a test config
func createTestConfig() *parser.Config {
	configContent := `# Stocks
iron_ore:100
coal:50
wood:25

# Processes
smelting:(iron_ore:2;coal:1):(iron_ingot:1):30
crafting:(iron_ingot:3;wood:5):(sword:1):45
mining:(coal:1):(iron_ore:2):20

# Optimize
optimize:(time;sword)`

	config, err := parser.ParseConfig(strings.NewReader(configContent))
	if err != nil {
		panic("Failed to create test config: " + err.Error())
	}
	return config
}

// Helper function to create a temporary log file
func createTempLogFile(content string) (string, func()) {
	tmpFile, err := os.CreateTemp("", "test_*.log")
	if err != nil {
		panic("Failed to create temp file: " + err.Error())
	}

	_, err = tmpFile.WriteString(content)
	if err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		panic("Failed to write to temp file: " + err.Error())
	}

	tmpFile.Close()
	return tmpFile.Name(), func() { os.Remove(tmpFile.Name()) }
}

func TestChecker_ValidTrace(t *testing.T) {
	config := createTestConfig()
	checker := NewChecker(config)

	// Valid log: mining (uses 1 coal, produces 2 iron_ore), then smelting (uses 2 iron_ore, 1 coal)
	logContent := `0:mining
1:smelting
2:mining`

	logFile, cleanup := createTempLogFile(logContent)
	defer cleanup()

	err := checker.CheckLogFile(logFile)
	if err != nil {
		t.Errorf("Expected valid trace to pass, but got error: %v", err)
	}
}

func TestChecker_InsufficientStock(t *testing.T) {
	config := createTestConfig()
	checker := NewChecker(config)

	// Invalid log: try to run smelting 60 times (needs 120 iron_ore, 60 coal, but we only have 100 iron_ore, 50 coal)
	// Each smelting needs 2 iron_ore + 1 coal
	// We have 100 iron_ore, 50 coal
	// So we can run at most 50 times (limited by coal), but let's try 60
	var logLines []string
	for i := 0; i < 60; i++ {
		logLines = append(logLines, fmt.Sprintf("%d:smelting", i))
	}
	logContent := strings.Join(logLines, "\n")

	logFile, cleanup := createTempLogFile(logContent)
	defer cleanup()

	err := checker.CheckLogFile(logFile)
	if err == nil {
		t.Error("Expected insufficient stock error, but validation passed")
		return
	}

	// Check that it's a CheckerError with the right details
	checkerErr, ok := err.(CheckerError)
	if !ok {
		t.Errorf("Expected CheckerError, got %T: %v", err, err)
		return
	}

	if !strings.Contains(checkerErr.Message, "not enough") {
		t.Errorf("Expected 'not enough' in error message, got: %s", checkerErr.Message)
	}
}

func TestChecker_InvalidCycleOrder(t *testing.T) {
	config := createTestConfig()
	checker := NewChecker(config)

	// Invalid log: cycles go backwards
	logContent := `0:mining
2:smelting
1:mining`

	logFile, cleanup := createTempLogFile(logContent)
	defer cleanup()

	err := checker.CheckLogFile(logFile)
	if err == nil {
		t.Error("Expected invalid cycle order error, but validation passed")
		return
	}

	// Check that it's a CheckerError with the right details
	checkerErr, ok := err.(CheckerError)
	if !ok {
		t.Errorf("Expected CheckerError, got %T: %v", err, err)
		return
	}

	if !strings.Contains(checkerErr.Message, "less than previous cycle") {
		t.Errorf("Expected 'less than previous cycle' in error message, got: %s", checkerErr.Message)
	}

	if checkerErr.Cycle != 1 {
		t.Errorf("Expected error on cycle 1, got cycle %d", checkerErr.Cycle)
	}
}

func TestChecker_UnknownProcess(t *testing.T) {
	config := createTestConfig()
	checker := NewChecker(config)

	// Invalid log: unknown process
	logContent := `0:unknown_process`

	logFile, cleanup := createTempLogFile(logContent)
	defer cleanup()

	err := checker.CheckLogFile(logFile)
	if err == nil {
		t.Error("Expected unknown process error, but validation passed")
		return
	}

	if !strings.Contains(err.Error(), "not found in configuration") {
		t.Errorf("Expected 'not found in configuration' in error message, got: %s", err.Error())
	}
}

func TestChecker_InvalidLogFormat(t *testing.T) {
	config := createTestConfig()
	checker := NewChecker(config)

	// Invalid log: missing colon
	logContent := `0mining`

	logFile, cleanup := createTempLogFile(logContent)
	defer cleanup()

	err := checker.CheckLogFile(logFile)
	if err == nil {
		t.Error("Expected invalid format error, but validation passed")
		return
	}

	if !strings.Contains(err.Error(), "invalid format") {
		t.Errorf("Expected 'invalid format' in error message, got: %s", err.Error())
	}
}

func TestChecker_InvalidCycleNumber(t *testing.T) {
	config := createTestConfig()
	checker := NewChecker(config)

	// Invalid log: non-numeric cycle
	logContent := `abc:mining`

	logFile, cleanup := createTempLogFile(logContent)
	defer cleanup()

	err := checker.CheckLogFile(logFile)
	if err == nil {
		t.Error("Expected invalid cycle number error, but validation passed")
		return
	}

	if !strings.Contains(err.Error(), "invalid cycle number") {
		t.Errorf("Expected 'invalid cycle number' in error message, got: %s", err.Error())
	}
}

func TestChecker_EmptyLogFile(t *testing.T) {
	config := createTestConfig()
	checker := NewChecker(config)

	// Empty log file
	logContent := ``

	logFile, cleanup := createTempLogFile(logContent)
	defer cleanup()

	err := checker.CheckLogFile(logFile)
	if err != nil {
		t.Errorf("Expected empty log to be valid, but got error: %v", err)
	}
}
