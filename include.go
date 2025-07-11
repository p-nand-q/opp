package opp

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func (p *Preprocessor) processInclude(directive string) (string, error) {
	// Format: ##<<filename>.
	if !strings.HasPrefix(directive, "<") || !strings.HasSuffix(directive, ".") {
		return "", fmt.Errorf("invalid include syntax: ##%s", directive)
	}
	
	// Extract filename
	filename := directive[1 : len(directive)-1]
	
	// Unescape the bizarre OPP escape sequences
	filename = unescapeFilename(filename)
	
	// Read the file
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		// Try relative to current file if we have that context
		if p.currentFile != "" {
			relPath := filepath.Join(filepath.Dir(p.currentFile), filename)
			content, err = ioutil.ReadFile(relPath)
		}
		if err != nil {
			return "", fmt.Errorf("cannot read file %s: %w", filename, err)
		}
	}
	
	// Determine the full path for the included file
	fullPath := filename
	if !filepath.IsAbs(filename) && p.currentFile != "" {
		fullPath = filepath.Join(filepath.Dir(p.currentFile), filename)
	}
	
	// Create a new preprocessor for the included file to avoid state pollution
	// But share the macros and variables
	includeProcessor := &Preprocessor{
		macros:      p.macros,
		variables:   p.variables,
		random:      p.random,
		lineNumber:  1,
		braceCount:  p.braceCount,
		closeBraces: p.closeBraces,
		currentFile: fullPath,
	}
	
	// Process the included file
	result, err := includeProcessor.Process(string(content))
	if err != nil {
		return "", fmt.Errorf("error processing included file %s: %w", filename, err)
	}
	
	// Update our brace counts from the included file
	p.braceCount = includeProcessor.braceCount
	p.closeBraces = includeProcessor.closeBraces
	
	return result, nil
}

// unescapeFilename reverses OPP's creative path escaping
func unescapeFilename(escaped string) string {
	result := escaped
	
	// Process in reverse order of the escape rules to avoid conflicts
	
	// Step 1: Replace // with temporary marker
	result = strings.ReplaceAll(result, "//", "\x00DBLSLASH\x00")
	
	// Step 2: Replace \\ with temporary marker  
	result = strings.ReplaceAll(result, "\\\\", "\x00DBLBACK\x00")
	
	// Step 3: Replace .. with temporary marker
	result = strings.ReplaceAll(result, "..", "\x00DOTDOT\x00")
	
	// Step 4: Replace \. with temporary marker
	result = strings.ReplaceAll(result, "\\.", "\x00BSDOT\x00")
	
	// Now apply the transformations
	result = strings.ReplaceAll(result, "\x00DBLSLASH\x00", "\\\\")  // // -> \\
	result = strings.ReplaceAll(result, "\x00DBLBACK\x00", "..")    // \\ -> ..
	result = strings.ReplaceAll(result, "\x00DOTDOT\x00", "\\")     // .. -> \
	result = strings.ReplaceAll(result, "\x00BSDOT\x00", ".")       // \. -> .
	
	return result
}

// escapeFilename applies OPP's creative path escaping (for testing)
func escapeFilename(filename string) string {
	result := filename
	
	// Apply in forward order
	result = strings.ReplaceAll(result, ".", "\\.")   // . -> \.
	result = strings.ReplaceAll(result, "\\", "..")   // \ -> ..
	result = strings.ReplaceAll(result, "..", "\\\\") // .. -> \\
	result = strings.ReplaceAll(result, "\\\\", "//") // \\ -> //
	
	return result
}