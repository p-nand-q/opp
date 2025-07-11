package opp

import (
	"fmt"
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
	
	// TODO: Actually read and process the file
	return "", fmt.Errorf("file inclusion not yet implemented: %s", filename)
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