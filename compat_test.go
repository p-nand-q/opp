package opp

import (
	"os"
	"testing"
)

func TestCompatMode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		compat   bool
		expected string
		setup    func()
		cleanup  func()
	}{
		{
			name: "dot separator in compat mode",
			input: `##~DEBUG.~DEBUG
debug line
##.`,
			compat: true,
			setup: func() {
				os.Setenv("DEBUG", "")
			},
			expected: "debug line",
			cleanup: func() {
				os.Unsetenv("DEBUG")
			},
		},
		{
			name: "complex format",
			input: `##i`,
			compat: true,
			expected: "complex(0,1)",
		},
		{
			name: "environment variable check",
			input: `##~(~TEST.~TEST).~(~TEST.~TEST)
test defined
##.`,
			compat: true,
			setup: func() {
				os.Setenv("TEST", "1")
			},
			expected: "test defined",
			cleanup: func() {
				os.Unsetenv("TEST")
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			if tt.cleanup != nil {
				defer tt.cleanup()
			}
			
			var p *Preprocessor
			if tt.compat {
				p = NewWithCompat(DefaultCompat())
			} else {
				p = New()
			}
			
			result, err := p.Process(tt.input)
			if err != nil {
				t.Fatalf("Process() error = %v", err)
			}
			
			if result != tt.expected {
				t.Errorf("Process() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestOriginalBugs(t *testing.T) {
	t.Run("brace counting bug", func(t *testing.T) {
		// In original, both { and } increment open counter
		p := NewWithCompat(CompatMode{BraceCountingBug: true})
		input := `{ } ##{ ##}`
		result, err := p.Process(input)
		if err != nil {
			t.Fatalf("Process() error = %v", err)
		}
		
		// With bug: { increments to 1, } increments to 2
		// So ##{ should show 2, ##} should show 0 (no close braces counted)
		expected := "{ } 2 0"
		if result != expected {
			t.Errorf("Brace bug: got %q, want %q", result, expected)
		}
	})
}

func TestANDLogicBug(t *testing.T) {
	// The original implementation has AND instead of NAND
	// This is fundamentally broken and we document it as incompatible
	t.Skip("Original AND logic bug is documented as incompatible")
}