package scheduler

import (
	"os"
	"testing"
	"time"

	"github.com/bernotieno/stock-exchange-simulator/internal/parser"
)

func TestNewSchedulerState(t *testing.T) {
	config := &parser.Config{
		Stocks: map[string]int{"wood": 5, "metal": 3},
		Processes: []parser.Process{
			{Name: "cut", Inputs: map[string]int{"wood": 1}, Outputs: map[string]int{"plank": 1}, Duration: 10 * time.Second},
		},
		Optimize: []string{"time"},
	}

	state := NewSchedulerState(config)

	if state.CurrentCycle != 0 {
		t.Errorf("Expected CurrentCycle 0, got %d", state.CurrentCycle)
	}
	if state.Stocks["wood"] != 5 {
		t.Errorf("Expected wood stock 5, got %d", state.Stocks["wood"])
	}
	if len(state.Processes) != 1 {
		t.Errorf("Expected 1 process, got %d", len(state.Processes))
	}
}

func TestCanRunProcess(t *testing.T) {
	state := &SchedulerState{
		Stocks: map[string]int{"wood": 2, "metal": 1},
	}

	process := &Process{
		Inputs: map[string]int{"wood": 1, "metal": 1},
	}

	if !state.CanRunProcess(process) {
		t.Error("Expected process to be runnable")
	}

	process.Inputs["wood"] = 3
	if state.CanRunProcess(process) {
		t.Error("Expected process to not be runnable")
	}
}

func TestConsumeInputs(t *testing.T) {
	state := &SchedulerState{
		Stocks: map[string]int{"wood": 5, "metal": 3},
	}

	process := &Process{
		Inputs: map[string]int{"wood": 2, "metal": 1},
	}

	state.ConsumeInputs(process)

	if state.Stocks["wood"] != 3 {
		t.Errorf("Expected wood stock 3, got %d", state.Stocks["wood"])
	}
	if state.Stocks["metal"] != 2 {
		t.Errorf("Expected metal stock 2, got %d", state.Stocks["metal"])
	}
}

func TestProduceOutputs(t *testing.T) {
	state := &SchedulerState{
		Stocks: map[string]int{"plank": 0},
	}

	process := &Process{
		Outputs: map[string]int{"plank": 2, "sawdust": 1},
	}

	state.ProduceOutputs(process)

	if state.Stocks["plank"] != 2 {
		t.Errorf("Expected plank stock 2, got %d", state.Stocks["plank"])
	}
	if state.Stocks["sawdust"] != 1 {
		t.Errorf("Expected sawdust stock 1, got %d", state.Stocks["sawdust"])
	}
}

func TestStartProcess(t *testing.T) {
	state := &SchedulerState{
		CurrentCycle: 0,
		Stocks:       map[string]int{"wood": 2},
		RunningProcesses: []RunningProcess{},
		ExecutionLog:     []ExecutionStep{},
	}

	process := &Process{
		Name:     "cut",
		Inputs:   map[string]int{"wood": 1},
		Duration: 10 * time.Second,
	}

	state.StartProcess(process)

	if len(state.RunningProcesses) != 1 {
		t.Errorf("Expected 1 running process, got %d", len(state.RunningProcesses))
	}
	if state.Stocks["wood"] != 1 {
		t.Errorf("Expected wood stock 1, got %d", state.Stocks["wood"])
	}
	if len(state.ExecutionLog) != 1 {
		t.Errorf("Expected 1 execution log entry, got %d", len(state.ExecutionLog))
	}
}

func TestCompleteFinishedProcesses(t *testing.T) {
	process := &Process{
		Name:    "cut",
		Outputs: map[string]int{"plank": 1},
	}

	state := &SchedulerState{
		CurrentCycle: 10,
		Stocks:       map[string]int{"plank": 0},
		RunningProcesses: []RunningProcess{
			{Process: process, StartTime: 0, EndTime: 5},
			{Process: process, StartTime: 5, EndTime: 15},
		},
	}

	state.CompleteFinishedProcesses()

	if len(state.RunningProcesses) != 1 {
		t.Errorf("Expected 1 running process, got %d", len(state.RunningProcesses))
	}
	if state.Stocks["plank"] != 1 {
		t.Errorf("Expected plank stock 1, got %d", state.Stocks["plank"])
	}
}

func TestIsTimeOptimized(t *testing.T) {
	state := &SchedulerState{
		OptimizeTargets: []string{"time", "wood"},
	}

	if !state.isTimeOptimized() {
		t.Error("Expected time optimization to be true")
	}

	state.OptimizeTargets = []string{"wood"}
	if state.isTimeOptimized() {
		t.Error("Expected time optimization to be false")
	}
}

func TestPrioritizeForTime(t *testing.T) {
	processes := []*Process{
		{Name: "long", Duration: 30 * time.Second},
		{Name: "short", Duration: 10 * time.Second},
		{Name: "medium", Duration: 20 * time.Second},
	}

	state := &SchedulerState{}
	sorted := state.prioritizeForTime(processes)

	if sorted[0].Name != "short" {
		t.Errorf("Expected first process to be 'short', got '%s'", sorted[0].Name)
	}
	if sorted[2].Name != "long" {
		t.Errorf("Expected last process to be 'long', got '%s'", sorted[2].Name)
	}
}

func TestWriteLogFile(t *testing.T) {
	state := &SchedulerState{
		CurrentCycle: 50,
		ExecutionLog: []ExecutionStep{
			{Cycle: 0, Process: "cut"},
			{Cycle: 10, Process: "assemble"},
		},
	}

	filename := "test.log"
	defer os.Remove(filename)

	err := state.WriteLogFile(filename)
	if err != nil {
		t.Errorf("WriteLogFile failed: %v", err)
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		t.Errorf("Failed to read log file: %v", err)
	}

	expected := "0:cut\n10:assemble\nNo more process doable at cycle 50\n"
	if string(content) != expected {
		t.Errorf("Expected log content:\n%s\nGot:\n%s", expected, string(content))
	}
}

func TestGetAvailableProcesses(t *testing.T) {
	state := &SchedulerState{
		Stocks: map[string]int{"wood": 1, "metal": 0},
		Processes: []Process{
			{Name: "cut", Inputs: map[string]int{"wood": 1}},
			{Name: "forge", Inputs: map[string]int{"metal": 1}},
		},
	}

	available := state.GetAvailableProcesses()

	if len(available) != 1 {
		t.Errorf("Expected 1 available process, got %d", len(available))
	}
	if available[0].Name != "cut" {
		t.Errorf("Expected available process 'cut', got '%s'", available[0].Name)
	}
}

func TestHasRunnableProcesses(t *testing.T) {
	state := &SchedulerState{
		Stocks: map[string]int{"wood": 1},
		Processes: []Process{
			{Name: "cut", Inputs: map[string]int{"wood": 1}},
		},
	}

	if !state.HasRunnableProcesses() {
		t.Error("Expected runnable processes to be true")
	}

	state.Stocks["wood"] = 0
	if state.HasRunnableProcesses() {
		t.Error("Expected runnable processes to be false")
	}
}

func TestAdvanceCycle(t *testing.T) {
	// Test with no running processes
	state := &SchedulerState{
		CurrentCycle:     5,
		RunningProcesses: []RunningProcess{},
	}

	state.AdvanceCycle()
	if state.CurrentCycle != 6 {
		t.Errorf("Expected cycle 6, got %d", state.CurrentCycle)
	}

	// Test with running processes
	state.CurrentCycle = 5
	state.RunningProcesses = []RunningProcess{
		{EndTime: 10},
		{EndTime: 8},
	}

	state.AdvanceCycle()
	if state.CurrentCycle != 8 {
		t.Errorf("Expected cycle 8, got %d", state.CurrentCycle)
	}
}

func TestStepSimulation(t *testing.T) {
	state := &SchedulerState{
		CurrentCycle: 0,
		Stocks:       map[string]int{"wood": 1},
		Processes: []Process{
			{Name: "cut", Inputs: map[string]int{"wood": 1}, Duration: 10 * time.Second},
		},
		RunningProcesses: []RunningProcess{},
		ExecutionLog:     []ExecutionStep{},
	}

	result := state.StepSimulation()
	if !result {
		t.Error("Expected step simulation to return true")
	}

	// No more processes can run
	state.Stocks["wood"] = 0
	result = state.StepSimulation()
	if result {
		t.Error("Expected step simulation to return false")
	}
}

func TestRunSimulation(t *testing.T) {
	state := &SchedulerState{
		CurrentCycle: 0,
		Stocks:       map[string]int{"wood": 2},
		Processes: []Process{
			{Name: "cut", Inputs: map[string]int{"wood": 1}, Outputs: map[string]int{"plank": 1}, Duration: 5 * time.Second},
		},
		RunningProcesses: []RunningProcess{},
		ExecutionLog:     []ExecutionStep{},
		OptimizeTargets:  []string{"time"},
	}

	state.RunSimulation(100)

	if len(state.ExecutionLog) != 2 {
		t.Errorf("Expected 2 execution log entries, got %d", len(state.ExecutionLog))
	}
	if state.Stocks["plank"] != 2 {
		t.Errorf("Expected plank stock 2, got %d", state.Stocks["plank"])
	}
}