package parser

import (
	"reflect"
	"strings"
	"testing"
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
