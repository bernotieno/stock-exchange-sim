package parser

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestParseStocks(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    map[string]int
		wantErr bool
	}{
		{
			name:  "valid stocks",
			input: "iron_ore:100\ncoal:50\nwood:25",
			want:  map[string]int{"iron_ore": 100, "coal": 50, "wood": 25},
		},
		{
			name:  "with comments and empty lines",
			input: "# Initial resources\niron_ore:100\n\n# More resources\ncoal:50",
			want:  map[string]int{"iron_ore": 100, "coal": 50},
		},
		{
			name:  "stop at process line",
			input: "iron_ore:100\ncoal:50\nsmelting:(iron_ore:2):(iron_ingot:1):30",
			want:  map[string]int{"iron_ore": 100, "coal": 50},
		},
		{
			name:    "invalid format missing colon",
			input:   "iron_ore100",
			wantErr: true,
		},
		{
			name:    "invalid format too many colons",
			input:   "iron_ore:100:extra",
			wantErr: true,
		},
		{
			name:    "empty name",
			input:   ":100",
			wantErr: true,
		},
		{
			name:    "invalid quantity",
			input:   "iron_ore:abc",
			wantErr: true,
		},
		{
			name:    "negative quantity",
			input:   "iron_ore:-10",
			wantErr: true,
		},
		{
			name:  "zero quantity allowed",
			input: "iron_ore:0",
			want:  map[string]int{"iron_ore": 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			got, err := ParseStocks(r)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseStocks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseStocks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseProcesses(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []Process
		wantErr bool
	}{
		{
			name:  "valid process",
			input: "smelting:(iron_ore:2):(iron_ingot:1):30",
			want: []Process{
				{
					Name:     "smelting",
					Inputs:   map[string]int{"iron_ore": 2},
					Outputs:  map[string]int{"iron_ingot": 1},
					Duration: 30 * time.Second,
				},
			},
		},
		{
			name:  "process with duration units",
			input: "smelting:(iron_ore:2):(iron_ingot:1):30s",
			want: []Process{
				{
					Name:     "smelting",
					Inputs:   map[string]int{"iron_ore": 2},
					Outputs:  map[string]int{"iron_ingot": 1},
					Duration: 30 * time.Second,
				},
			},
		},
		{
			name:  "process with multiple inputs and outputs",
			input: "crafting:(iron_ingot:2;wood:1):(sword:1;shield:1):120",
			want: []Process{
				{
					Name:     "crafting",
					Inputs:   map[string]int{"iron_ingot": 2, "wood": 1},
					Outputs:  map[string]int{"sword": 1, "shield": 1},
					Duration: 120 * time.Second,
				},
			},
		},
		{
			name:  "process with empty inputs",
			input: "gathering:():(wood:1):10",
			want: []Process{
				{
					Name:     "gathering",
					Inputs:   map[string]int{},
					Outputs:  map[string]int{"wood": 1},
					Duration: 10 * time.Second,
				},
			},
		},
		{
			name:  "process with empty outputs",
			input: "consuming:(wood:1):():5",
			want: []Process{
				{
					Name:     "consuming",
					Inputs:   map[string]int{"wood": 1},
					Outputs:  map[string]int{},
					Duration: 5 * time.Second,
				},
			},
		},
		{
			name:  "multiple processes",
			input: "smelting:(iron_ore:2):(iron_ingot:1):30\ncrafting:(iron_ingot:1):(sword:1):60",
			want: []Process{
				{
					Name:     "smelting",
					Inputs:   map[string]int{"iron_ore": 2},
					Outputs:  map[string]int{"iron_ingot": 1},
					Duration: 30 * time.Second,
				},
				{
					Name:     "crafting",
					Inputs:   map[string]int{"iron_ingot": 1},
					Outputs:  map[string]int{"sword": 1},
					Duration: 60 * time.Second,
				},
			},
		},
		{
			name:    "invalid format insufficient colons",
			input:   "smelting:(iron_ore:2):(iron_ingot:1)",
			wantErr: true,
		},
		{
			name:    "invalid format too many colons",
			input:   "smelting:(iron_ore:2):(iron_ingot:1):30:extra",
			wantErr: true,
		},
		{
			name:    "empty process name",
			input:   ":(iron_ore:2):(iron_ingot:1):30",
			wantErr: true,
		},
		{
			name:    "invalid inputs format",
			input:   "smelting:iron_ore:2:(iron_ingot:1):30",
			wantErr: true,
		},
		{
			name:    "invalid outputs format",
			input:   "smelting:(iron_ore:2):iron_ingot:1:30",
			wantErr: true,
		},
		{
			name:    "invalid duration",
			input:   "smelting:(iron_ore:2):(iron_ingot:1):abc",
			wantErr: true,
		},
		{
			name:    "zero duration",
			input:   "smelting:(iron_ore:2):(iron_ingot:1):0",
			wantErr: true,
		},
		{
			name:    "negative duration",
			input:   "smelting:(iron_ore:2):(iron_ingot:1):-30",
			wantErr: true,
		},
		{
			name:    "invalid item quantity",
			input:   "smelting:(iron_ore:abc):(iron_ingot:1):30",
			wantErr: true,
		},
		{
			name:    "zero item quantity",
			input:   "smelting:(iron_ore:0):(iron_ingot:1):30",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			got, err := ParseProcesses(r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseProcesses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseProcesses() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseOptimize(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []string
		wantErr bool
	}{
		{
			name:  "valid single target",
			input: "optimize:(iron_ingot)",
			want:  []string{"iron_ingot"},
		},
		{
			name:  "valid multiple targets",
			input: "optimize:(iron_ingot;sword;shield)",
			want:  []string{"iron_ingot", "sword", "shield"},
		},
		{
			name:  "empty targets",
			input: "optimize:()",
			want:  []string{},
		},
		{
			name:  "no optimize line",
			input: "iron_ore:100\nsmelting:(iron_ore:2):(iron_ingot:1):30",
			want:  []string{},
		},
		{
			name:  "optimize with comments",
			input: "# Optimization targets\noptimize:(iron_ingot;sword)",
			want:  []string{"iron_ingot", "sword"},
		},
		{
			name:  "malformed line fallback",
			input: "optimize:iron_ingot;sword",
			want:  []string{}, // Should fallback to empty slice
		},
		{
			name:  "targets with spaces",
			input: "optimize:(iron ingot; steel plate)",
			want:  []string{"iron ingot", "steel plate"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			got, err := ParseOptimize(r)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseOptimize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseOptimize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseConfig(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantConfig  *Config
		wantErr     bool
		errContains string
	}{
		{
			name: "empty config",
			input: `
# Stocks
wood:0
nails:0
`,
			wantConfig: &Config{
				Stocks:    map[string]int{"wood": 0, "nails": 0},
				Processes: []Process{},
				Optimize:  []string{},
			},
			wantErr: false,
		},
		{
			name: "invalid stock line - missing quantity",
			input: `
# Stocks
wood:
`,
			wantErr:     true,
			errContains: "invalid quantity",
		},
		{
			name: "invalid stock line - negative quantity",
			input: `
# Stocks
wood:-5
`,
			wantErr:     true,
			errContains: "stock quantity cannot be negative",
		},
		{
			name: "invalid config - unknown optimize target",
			input: `
# Stocks
wood:100
nails:200

# Processes
process build chair {
  (wood:5) -> (chair:1)
}

# Optimize
optimize:(table)
`,
			wantErr:     true,
			errContains: "optimization target 'table' is not defined",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			got, err := ParseConfig(reader)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("ParseConfig() error = %v, should contain %v", err, tt.errContains)
				}
				return
			}

			if !configEqual(got, tt.wantConfig) {
				t.Errorf("ParseConfig() = %+v, want %+v", got, tt.wantConfig)
			}
		})
	}
}

func configEqual(a, b *Config) bool {
	// Compare stocks
	if len(a.Stocks) != len(b.Stocks) {
		return false
	}
	for k, v := range a.Stocks {
		if b.Stocks[k] != v {
			return false
		}
	}

	// Compare processes
	if len(a.Processes) != len(b.Processes) {
		return false
	}
	for i := range a.Processes {
		if !processEqual(&a.Processes[i], &b.Processes[i]) {
			return false
		}
	}

	// Compare optimize
	if len(a.Optimize) != len(b.Optimize) {
		return false
	}
	for i := range a.Optimize {
		if a.Optimize[i] != b.Optimize[i] {
			return false
		}
	}

	return true
}

func processEqual(a, b *Process) bool {
	if a.Name != b.Name {
		return false
	}

	if len(a.Inputs) != len(b.Inputs) {
		return false
	}
	for k, v := range a.Inputs {
		if b.Inputs[k] != v {
			return false
		}
	}

	if len(a.Outputs) != len(b.Outputs) {
		return false
	}
	for k, v := range a.Outputs {
		if b.Outputs[k] != v {
			return false
		}
	}

	return true
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		wantErr     bool
		errContains string
	}{
		{
			name: "valid config",
			config: &Config{
				Stocks: map[string]int{"wood": 100},
				Processes: []Process{
					{
						Name:    "build",
						Inputs:  map[string]int{"wood": 5},
						Outputs: map[string]int{"chair": 1},
					},
				},
				Optimize: []string{"chair"},
			},
			wantErr: false,
		},
		{
			name: "invalid - optimize target not in stocks or processes",
			config: &Config{
				Stocks:    map[string]int{"wood": 100},
				Processes: []Process{},
				Optimize:  []string{"chair"},
			},
			wantErr:     true,
			errContains: "optimization target 'chair' is not defined",
		},
		{
			name: "invalid - process with no inputs or outputs",
			config: &Config{
				Stocks: map[string]int{},
				Processes: []Process{
					{
						Name: "empty",
					},
				},
				Optimize: []string{},
			},
			wantErr:     true,
			errContains: "process 'empty' must have at least one input or output",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("Validate() error = %v, should contain %v", err, tt.errContains)
			}
		})
	}
}

func TestConfig_GetProcessByName(t *testing.T) {
	config := &Config{
		Processes: []Process{
			{Name: "process1"},
			{Name: "process2"},
		},
	}

	tests := []struct {
		name     string
		query    string
		expected *Process
	}{
		{"existing process", "process1", &config.Processes[0]},
		{"non-existent process", "process3", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := config.GetProcessByName(tt.query)
			if (got == nil) != (tt.expected == nil) {
				t.Errorf("GetProcessByName() = %v, want %v", got, tt.expected)
			}
			if got != nil && tt.expected != nil && got.Name != tt.expected.Name {
				t.Errorf("GetProcessByName() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestConfig_GetStockQuantity(t *testing.T) {
	config := &Config{
		Stocks: map[string]int{
			"wood":  100,
			"nails": 200,
		},
	}

	tests := []struct {
		name     string
		item     string
		expected int
	}{
		{"existing item", "wood", 100},
		{"existing item", "nails", 200},
		{"non-existent item", "metal", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := config.GetStockQuantity(tt.item)
			if got != tt.expected {
				t.Errorf("GetStockQuantity() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestConfig_GetAllItems(t *testing.T) {
	config := &Config{
		Stocks: map[string]int{
			"wood":  100,
			"nails": 200,
		},
		Processes: []Process{
			{
				Name: "build chair",
				Inputs: map[string]int{
					"wood":  5,
					"nails": 10,
				},
				Outputs: map[string]int{
					"chair": 1,
				},
			},
			{
				Name: "build table",
				Inputs: map[string]int{
					"wood": 20,
				},
				Outputs: map[string]int{
					"table": 1,
				},
			},
		},
	}

	got := config.GetAllItems()
	expectedItems := map[string]bool{
		"wood":  true,
		"nails": true,
		"chair": true,
		"table": true,
	}

	if len(got) != len(expectedItems) {
		t.Errorf("GetAllItems() returned %d items, expected %d", len(got), len(expectedItems))
	}

	for _, item := range got {
		if !expectedItems[item] {
			t.Errorf("GetAllItems() returned unexpected item: %s", item)
		}
	}
}