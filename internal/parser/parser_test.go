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
