package opp

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileInclusion(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "opp-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create test files
	headerFile := filepath.Join(tempDir, "header.h")
	headerContent := `#define HEADER_INCLUDED
const char* version = "1.0";`
	
	if err := os.WriteFile(headerFile, []byte(headerContent), 0644); err != nil {
		t.Fatalf("Failed to write header file: %v", err)
	}
	
	mainFile := filepath.Join(tempDir, "main.c")
	mainContent := `##<header\.h.
int main() {
    return 0;
}`
	
	if err := os.WriteFile(mainFile, []byte(mainContent), 0644); err != nil {
		t.Fatalf("Failed to write main file: %v", err)
	}
	
	// Test basic inclusion
	p := New()
	p.currentFile = mainFile
	result, err := p.Process(mainContent)
	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}
	
	expected := `#define HEADER_INCLUDED
const char* version = "1.0";
int main() {
    return 0;
}`
	
	if result != expected {
		t.Errorf("Process() = %q, want %q", result, expected)
	}
}

func TestFileInclusionEscaping(t *testing.T) {
	tests := []struct {
		name     string
		escaped  string
		expected string
	}{
		{
			name:     "simple dot escape",
			escaped:  "file\\.h",
			expected: "file.h",
		},
		{
			name:     "backslash to parent dir",
			escaped:  "parent..child",
			expected: "parent\\child",
		},
		{
			name:     "UNC path",
			escaped:  "//server..share..file\\.h",
			expected: "\\\\server\\share\\file.h",
		},
		{
			name:     "complex path",
			escaped:  "C:..Users..opp..test\\.h",
			expected: "C:\\Users\\opp\\test.h",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := unescapeFilename(tt.escaped)
			if result != tt.expected {
				t.Errorf("unescapeFilename(%q) = %q, want %q", tt.escaped, result, tt.expected)
			}
		})
	}
}

func TestFileInclusionErrors(t *testing.T) {
	p := New()
	
	tests := []struct {
		name  string
		input string
		error string
	}{
		{
			name:  "missing dot terminator",
			input: "##<file.h",
			error: "invalid include syntax",
		},
		{
			name:  "missing angle bracket",
			input: "##file.h.",
			error: "invalid include syntax",
		},
		{
			name:  "nonexistent file",
			input: "##<nonexistent\\.file.>",
			error: "cannot read file",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := p.Process(tt.input)
			if err == nil {
				t.Errorf("Expected error containing %q, got nil", tt.error)
			} else if !contains(err.Error(), tt.error) {
				t.Errorf("Expected error containing %q, got %q", tt.error, err.Error())
			}
		})
	}
}

func TestNestedInclusion(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "opp-nested-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create nested include files
	deepFile := filepath.Join(tempDir, "deep.h")
	deepContent := `#define DEEP_INCLUDED`
	if err := os.WriteFile(deepFile, []byte(deepContent), 0644); err != nil {
		t.Fatalf("Failed to write deep file: %v", err)
	}
	
	middleFile := filepath.Join(tempDir, "middle.h")
	middleContent := `#define MIDDLE_INCLUDED
##<deep\.h.`
	if err := os.WriteFile(middleFile, []byte(middleContent), 0644); err != nil {
		t.Fatalf("Failed to write middle file: %v", err)
	}
	
	// Test nested inclusion
	p := New()
	p.currentFile = filepath.Join(tempDir, "main.c")
	input := `##<middle\.h.
int main() {}`
	
	result, err := p.Process(input)
	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}
	
	expected := `#define MIDDLE_INCLUDED
#define DEEP_INCLUDED
int main() {}`
	
	if result != expected {
		t.Errorf("Process() = %q, want %q", result, expected)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && contains(s[1:], substr) || len(substr) > 0 && contains(s, substr[1:]))
}