package opp

import (
	"strings"
	"testing"
)

func TestConditionalCompilation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		defines  map[string]bool
		expected string
	}{
		{
			name: "undefined variable",
			input: `##~DEBUG|~DEBUG
debug line
##.
normal line`,
			defines:  map[string]bool{},
			expected: "debug line\nnormal line",
		},
		{
			name: "defined variable", 
			input: `##~DEBUG|~DEBUG
debug line
##.
normal line`,
			defines:  map[string]bool{"DEBUG": true},
			expected: "normal line",
		},
		{
			name: "check if defined using double negation",
			input: `##~(~DEBUG|~DEBUG)|~(~DEBUG|~DEBUG)
debug is defined
##.`,
			defines:  map[string]bool{"DEBUG": true},
			expected: "debug is defined",
		},
		{
			name: "if-else pattern",
			input: `##~DEBUG|~DEBUG
not debug
##@~(~DEBUG|~DEBUG)|~(~DEBUG|~DEBUG)
is debug
##.`,
			defines:  map[string]bool{"DEBUG": true},
			expected: "is debug",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()
			for k := range tt.defines {
				p.Define(k, "1")
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

func TestMacros(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "simple macro definition and use",
			input: `##:FOO bar
FOO`,
			expected: "bar",
		},
		{
			name: "predefined macro line number",
			input: `line ##_`,
			expected: "line -4", // Line 1 minus 5
		},
		{
			name: "brace counting",
			input: `{ count: ##{`,
			expected: "{ count: 1",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()
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

func TestFileEscaping(t *testing.T) {
	tests := []struct {
		escaped   string
		unescaped string
	}{
		{
			escaped:   `//server..users..opp..sample\.h`,
			unescaped: `\\server\users\opp\sample.h`,
		},
		{
			escaped:   `test\.txt`,
			unescaped: `test.txt`,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.escaped, func(t *testing.T) {
			result := unescapeFilename(tt.escaped)
			if result != tt.unescaped {
				t.Errorf("unescapeFilename(%q) = %q, want %q", tt.escaped, result, tt.unescaped)
			}
		})
	}
}

func TestNANDLogic(t *testing.T) {
	p := New()
	
	// Test basic NAND truth table
	tests := []struct {
		a, b     bool
		expected bool // ~a | ~b
	}{
		{false, false, true},  // ~false | ~false = true | true = true
		{false, true, true},   // ~false | ~true = true | false = true  
		{true, false, true},   // ~true | ~false = false | true = true
		{true, true, false},   // ~true | ~true = false | false = false
	}
	
	for _, tt := range tests {
		if tt.a {
			p.Define("A", "1")
		} else {
			p.Undefine("A")
		}
		if tt.b {
			p.Define("B", "1")
		} else {
			p.Undefine("B")
		}
		
		result, err := p.evaluateCondition("~A|~B")
		if err != nil {
			t.Fatalf("evaluateCondition() error = %v", err)
		}
		
		if result != tt.expected {
			t.Errorf("NAND(%v, %v) = %v, want %v", tt.a, tt.b, result, tt.expected)
		}
	}
}

func TestRandomNumber(t *testing.T) {
	p := New()
	
	// Test that ##$ produces consistent pseudo-random numbers
	input := "##$\n##$"
	result, err := p.Process(input)
	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}
	
	lines := strings.Split(result, "\n")
	if len(lines) != 2 {
		t.Fatalf("Expected 2 lines, got %d", len(lines))
	}
	
	// Numbers should be different (pseudo-random sequence)
	if lines[0] == lines[1] {
		t.Errorf("Random numbers should differ: %s == %s", lines[0], lines[1])
	}
}