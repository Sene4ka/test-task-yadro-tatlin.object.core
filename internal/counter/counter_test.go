package counter

import (
	"strings"
	"testing"
)

func TestCountNames(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		preserveCase bool
		want         map[string]int
	}{
		{
			name:         "basic counting",
			input:        "Алёна\nМиша\nАлёна\nДима\n",
			preserveCase: false,
			want:         map[string]int{"Алёна": 2, "Миша": 1, "Дима": 1},
		},
		{
			name:         "empty file",
			input:        "",
			preserveCase: false,
			want:         map[string]int{},
		},
		{
			name:         "case normalization",
			input:        "алёна\nАлёна\nАЛЁНА",
			preserveCase: false,
			want:         map[string]int{"Алёна": 3},
		},
		{
			name:         "preserve case",
			input:        "алёна\nАлёна",
			preserveCase: true,
			want:         map[string]int{"алёна": 1, "Алёна": 1},
		},
		{
			name:         "whitespace trimming",
			input:        "  Алёна  \n  Миша",
			preserveCase: false,
			want:         map[string]int{"Алёна": 1, "Миша": 1},
		},
		{
			name:         "empty lines ignored",
			input:        "Алёна\n\nМиша\n  \n",
			preserveCase: false,
			want:         map[string]int{"Алёна": 1, "Миша": 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CountNames(strings.NewReader(tt.input), tt.preserveCase)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != len(tt.want) {
				t.Errorf("got %d entries, want %d", len(got), len(tt.want))
			}
			for k, v := range tt.want {
				if got[k] != v {
					t.Errorf("for key %q: got %d, want %d", k, got[k], v)
				}
			}
		})
	}
}
