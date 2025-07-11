package opp

import (
	"fmt"
	"strconv"
	"strings"
)

func (p *Preprocessor) defineMacro(definition string) error {
	// Find first space to separate name from body
	spaceIdx := strings.Index(definition, " ")
	if spaceIdx == -1 {
		// Macro with no body
		name := strings.TrimSpace(definition)
		// Check if it's a function-like macro with no body
		if parenIdx := strings.Index(name, "("); parenIdx >= 0 {
			// Extract just the name part
			name = name[:parenIdx]
		}
		p.macros[name] = &Macro{Name: name, Definition: ""}
		return nil
	}
	
	nameAndArgs := definition[:spaceIdx]
	body := definition[spaceIdx+1:]
	
	// Check if it's a function-like macro
	name := nameAndArgs
	isFunctionLike := false
	if parenIdx := strings.Index(nameAndArgs, "("); parenIdx >= 0 {
		// Function-like macro
		name = nameAndArgs[:parenIdx]
		isFunctionLike = true
		// We don't need to parse parameter names since we use positional args (#0, #1, etc)
	}
	
	// Process ##,# escapes in the macro body
	body = p.handleNestedMacroEscapes(body)
	
	p.macros[name] = &Macro{
		Name:           name,
		Definition:     body,
		IsFunctionLike: isFunctionLike,
	}
	
	return nil
}

// handleNestedMacroEscapes processes ##,# escape sequences
// ##,#0 → #0 (literal)
// ##,## → ## (literal)
func (p *Preprocessor) handleNestedMacroEscapes(line string) string {
	result := ""
	i := 0
	
	for i < len(line) {
		// Check if we have ##,# sequence
		if i+3 < len(line) && line[i:i+4] == "##,#" {
			// Check what follows the ##,#
			if i+4 < len(line) {
				nextChar := line[i+4]
				if nextChar == '#' {
					// ##,## → ##
					result += "##"
					i += 5 // Skip ##,##
				} else if nextChar >= '0' && nextChar <= '9' {
					// ##,#0 → #0
					result += "#" + string(nextChar)
					i += 5 // Skip ##,#N
				} else {
					// ##,#x where x is not # or digit → #x
					result += "#" + string(nextChar)
					i += 5
				}
			} else {
				// ##,# at end of line - don't process as escape
				result += line[i:]
				break
			}
		} else {
			// Regular character
			result += string(line[i])
			i++
		}
	}
	
	return result
}

// processStringizeCharize handles #" and #' operators in macro definitions
// Returns the processed string with stringize/charize applied
func (p *Preprocessor) processStringizeCharize(definition string, args []string) string {
	result := ""
	i := 0
	
	for i < len(definition) {
		// Check for stringize operator #"
		if i+3 < len(definition) && definition[i:i+2] == "#\"" && definition[i+2] == '#' {
			// Look for the digit after #"#
			if i+3 < len(definition) && definition[i+3] >= '0' && definition[i+3] <= '9' {
				argNum := int(definition[i+3] - '0')
				if argNum < len(args) {
					// Escape quotes and backslashes in the argument
					escaped := strings.ReplaceAll(args[argNum], "\\", "\\\\")
					escaped = strings.ReplaceAll(escaped, "\"", "\\\"")
					result += "\"" + escaped + "\""
				} else {
					// Argument index out of bounds - empty string
					result += "\"\""
				}
				i += 4 // Skip #"#N
				continue
			}
		}
		
		// Check for charize operator #'
		if i+3 < len(definition) && definition[i:i+2] == "#'" && definition[i+2] == '#' {
			// Look for the digit after #'#
			if i+3 < len(definition) && definition[i+3] >= '0' && definition[i+3] <= '9' {
				argNum := int(definition[i+3] - '0')
				if argNum < len(args) {
					// Escape quotes and backslashes in the argument
					escaped := strings.ReplaceAll(args[argNum], "\\", "\\\\")
					escaped = strings.ReplaceAll(escaped, "'", "\\'")
					result += "'" + escaped + "'"
				} else {
					// Argument index out of bounds - empty string
					result += "''"
				}
				i += 4 // Skip #'#N
				continue
			}
		}
		
		// Check for regular argument substitution #N
		if i+1 < len(definition) && definition[i] == '#' && definition[i+1] >= '0' && definition[i+1] <= '9' {
			argNum := int(definition[i+1] - '0')
			if argNum < len(args) {
				result += args[argNum]
			}
			// If arg doesn't exist, #N remains as is
			i += 2 // Skip #N
			continue
		}
		
		// Regular character
		result += string(definition[i])
		i++
	}
	
	return result
}

// parseMacroCall attempts to parse a function-like macro call
// Returns the arguments if it's a function call, nil otherwise
func parseMacroCall(text string, pos int, macroName string) []string {
	// Check if there's a '(' immediately after the macro name
	if pos+len(macroName) >= len(text) || text[pos+len(macroName)] != '(' {
		return nil
	}
	
	// Find the matching closing parenthesis
	start := pos + len(macroName) + 1
	parenCount := 1
	current := start
	args := []string{}
	argStart := start
	
	for current < len(text) && parenCount > 0 {
		switch text[current] {
		case '(':
			parenCount++
		case ')':
			parenCount--
			if parenCount == 0 {
				// Last argument
				arg := strings.TrimSpace(text[argStart:current])
				if arg != "" || len(args) > 0 {
					// Add non-empty arg or empty arg if there were previous args
					args = append(args, arg)
				}
			}
		case ',':
			if parenCount == 1 {
				// Argument separator at top level
				args = append(args, strings.TrimSpace(text[argStart:current]))
				argStart = current + 1
			}
		}
		current++
	}
	
	if parenCount != 0 {
		// Unmatched parentheses
		return nil
	}
	
	return args
}

func (p *Preprocessor) expandMacros(line string) (string, error) {
	result := line
	changed := true
	
	// Keep expanding until no more changes (handles nested macros)
	for changed {
		changed = false
		newResult := ""
		i := 0
		
		for i < len(result) {
			found := false
			
			// Try each macro
			for name, macro := range p.macros {
				if i+len(name) <= len(result) && result[i:i+len(name)] == name {
					if macro.IsFunctionLike {
						// Function-like macro - only expand if followed by (
						args := parseMacroCall(result, i, name)
						if args != nil {
							// Function-like macro invocation
							endPos := i + len(name)
							parenCount := 1
							endPos++ // Skip opening (
							for endPos < len(result) && parenCount > 0 {
								if result[endPos] == '(' {
									parenCount++
								} else if result[endPos] == ')' {
									parenCount--
								}
								endPos++
							}
							
							// Process the macro definition with arguments
							expanded := p.processStringizeCharize(macro.Definition, args)
							newResult += expanded
							i = endPos
							found = true
							changed = true
							break
						}
					} else {
						// Object-like macro - check word boundary
						if !isAlphaNum(getCharAt(result, i+len(name))) {
							// Simple macro replacement (not followed by alnum)
							newResult += macro.Definition
							i += len(name)
							found = true
							changed = true
							break
						}
					}
				}
			}
			
			if !found {
				newResult += string(result[i])
				i++
			}
		}
		
		result = newResult
	}
	
	// Expand predefined dynamic macros
	result = p.expandDynamicMacros(result)
	
	return result, nil
}

// Helper functions
func getCharAt(s string, i int) byte {
	if i < len(s) {
		return s[i]
	}
	return 0
}

func isAlphaNum(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_'
}

func (p *Preprocessor) expandDynamicMacros(line string) string {
	result := line
	
	// ##_ - current line number minus 5
	if strings.Contains(result, "##_") {
		result = strings.ReplaceAll(result, "##_", strconv.Itoa(p.lineNumber-5))
	}
	
	// ##$ - pseudo-random number
	if strings.Contains(result, "##$") {
		result = strings.ReplaceAll(result, "##$", strconv.Itoa(p.random.Next()))
	}
	
	// ##{ - number of { braces seen so far (including in current line up to the token)
	for strings.Contains(result, "##{") {
		idx := strings.Index(result, "##{")
		if idx >= 0 {
			// Count braces in the current line up to this point
			bracesInLine := 0
			for i := 0; i < idx; i++ {
				if result[i] == '{' {
					bracesInLine++
				}
			}
			// Replace just this occurrence
			result = result[:idx] + strconv.Itoa(p.braceCount+bracesInLine) + result[idx+3:]
		}
	}
	
	// ##} - number of } braces modulo 5
	if strings.Contains(result, "##}") {
		result = strings.ReplaceAll(result, "##}", strconv.Itoa(p.closeBraces%5))
	}
	
	return result
}

func (p *Preprocessor) expandPredefinedMacro(line string) (string, error) {
	// Handle standalone predefined macros
	switch line {
	case "##i":
		return "1i", nil
	case "##_":
		return strconv.Itoa(p.lineNumber - 5), nil
	case "##$":
		return strconv.Itoa(p.random.Next()), nil
	case "##{":
		return strconv.Itoa(p.braceCount), nil
	case "##}":
		return strconv.Itoa(p.closeBraces % 5), nil
	default:
		return "", fmt.Errorf("unknown directive: %s", line)
	}
}

// expandFunctionMacro handles macro calls with arguments
func (p *Preprocessor) expandFunctionMacro(name string, args []string, definition string) string {
	result := definition
	
	// Replace #n with corresponding arguments
	for i, arg := range args {
		placeholder := fmt.Sprintf("#%d", i)
		result = strings.ReplaceAll(result, placeholder, arg)
	}
	
	// Handle varargs ##n..n
	// TODO: Implement varargs support
	
	return result
}