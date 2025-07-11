package opp

import (
	"testing"
)

func TestNestedMacroEscapes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "escape #0",
			input:    "##,#0",
			expected: "#0",
		},
		{
			name:     "escape #1",
			input:    "##,#1",
			expected: "#1",
		},
		{
			name:     "escape ##",
			input:    "##,##",
			expected: "##",
		},
		{
			name:     "escape in macro definition",
			input:    "##:test ##,#0 + ##,#1",
			expected: "##:test #0 + #1",
		},
		{
			name:     "multiple escapes",
			input:    "##,#0->method(##,##, ##,#1)",
			expected: "#0->method(##, #1)",
		},
		{
			name:     "escape at end of line",
			input:    "return ##,#0",
			expected: "return #0",
		},
	}
	
	p := New()
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := p.handleNestedMacroEscapes(tt.input)
			if result != tt.expected {
				t.Errorf("handleNestedMacroEscapes(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestMacroDefiningMacros(t *testing.T) {
	p := New()
	
	// Test simple macro with escaped content
	input := `##:HEADER ##:VERSION ##,#0.##,#1
HEADER`
	
	result, err := p.Process(input)
	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}
	
	// HEADER should expand to a macro definition
	expected := `##:VERSION #0.#1`
	
	if result != expected {
		t.Errorf("Process() = %q, want %q", result, expected)
	}
}

func TestNestedMacroWithLiterals(t *testing.T) {
	p := New()
	
	// Test macro that contains literal # and ##
	input := `##:CODE ##,#include <stdio.h> ##,##define DEBUG 1
CODE`
	
	result, err := p.Process(input)
	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}
	
	expected := `#include <stdio.h> ##define DEBUG 1`
	
	if result != expected {
		t.Errorf("Process() = %q, want %q", result, expected)
	}
}

func TestNestedMacroInDefinition(t *testing.T) {
	p := New()
	
	// Test macro that outputs another macro definition
	input := `##:MAKER ##:output ##,#0
MAKER`
	
	result, err := p.Process(input)
	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}
	
	// MAKER expands to the literal text "##:output #0"
	// (macro definitions from expansions are not processed)
	expected := `##:output #0`
	
	if result != expected {
		t.Errorf("Process() = %q, want %q", result, expected)
	}
}

func TestNestedEscapeEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "incomplete escape at end",
			input:    "test ##,#",
			expected: "test ##,#", // No change, incomplete sequence
		},
		{
			name:     "escape with non-digit",
			input:    "##,#x",
			expected: "#x", // Still outputs # followed by x
		},
		{
			name:     "consecutive escapes",
			input:    "##,####,##",
			expected: "####",
		},
		{
			name:     "tricky sequence",
			input:    "##,##,##,#1",
			expected: "##,#1", // First ##,## becomes ##, then ,##,#1 becomes ,#1
		},
		{
			name:     "escape in middle of text",
			input:    "before##,#0after",
			expected: "before#0after",
		},
	}
	
	p := New()
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := p.handleNestedMacroEscapes(tt.input)
			if result != tt.expected {
				t.Errorf("handleNestedMacroEscapes(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}