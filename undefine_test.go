package opp

import (
	"testing"
)

func TestUndefineDirective(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "undefine macro",
			input: `##:TEST hello world
TEST
##-TEST
TEST`,
			expected: `hello world
TEST`,
		},
		{
			name: "undefine variable",
			input: `##:DEBUG 1
##~(~DEBUG|~DEBUG)|~(~DEBUG|~DEBUG)
DEBUG is defined
##.
##-DEBUG
##~(~DEBUG|~DEBUG)|~(~DEBUG|~DEBUG)
DEBUG is defined
##@~DEBUG|~DEBUG
DEBUG is undefined
##.`,
			expected: `1 is defined
DEBUG is undefined`,
		},
		{
			name: "undefine multiple macros",
			input: `##:A aaa
##:B bbb
##:C ccc
A B C
##-A
##-B
A B C`,
			expected: `aaa bbb ccc
A B ccc`,
		},
		{
			name: "undefine non-existent macro",
			input: `##-NONEXISTENT
##:TEST test
TEST`,
			expected: `test`,
		},
		{
			name: "undefine and redefine",
			input: `##:GREETING Hello
GREETING
##-GREETING
GREETING
##:GREETING Goodbye
GREETING`,
			expected: `Hello
GREETING
Goodbye`,
		},
		{
			name: "cannot undefine predefined macro",
			input: `##i
##-##i
##i`,
			expected: `1i
1i`, // Predefined macros cannot be undefined
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

func TestUndefineInConditionals(t *testing.T) {
	p := New()
	
	input := `##:DEBUG 1
##~DEBUG|~DEBUG
unreachable
##@~(~DEBUG|~DEBUG)|~(~DEBUG|~DEBUG)
##-DEBUG
inside conditional
##.
##~DEBUG|~DEBUG
DEBUG now undefined
##.`
	
	result, err := p.Process(input)
	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}
	
	expected := `inside conditional
DEBUG now undefined`
	
	if result != expected {
		t.Errorf("Process() = %q, want %q", result, expected)
	}
}

func TestUndefineEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "undefine with leading space",
			input: `##:MACRO test
##- MACRO
MACRO`,
			expected: `MACRO`, // After undefine, MACRO is no longer expanded
		},
		{
			name: "undefine empty name",
			input: `##-
##:TEST test
TEST`,
			expected: `test`,
		},
		{
			name: "undefine with special characters",
			input: `##:test! value
test!
##-test!
test!`,
			expected: `value
test!`,
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