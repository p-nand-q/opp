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
		p.macros[name] = &Macro{Name: name, Definition: ""}
		return nil
	}
	
	name := definition[:spaceIdx]
	body := definition[spaceIdx+1:]
	
	p.macros[name] = &Macro{
		Name:       name,
		Definition: body,
	}
	
	return nil
}

func (p *Preprocessor) expandMacros(line string) (string, error) {
	result := line
	
	// Expand user-defined macros
	for name, macro := range p.macros {
		if strings.Contains(result, name) {
			// Simple replacement for now - TODO: handle function-like macros
			result = strings.ReplaceAll(result, name, macro.Definition)
		}
	}
	
	// Expand predefined dynamic macros
	result = p.expandDynamicMacros(result)
	
	return result, nil
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
	
	// ##{ - number of { braces seen so far (up to this point)
	// Count braces before the ##{ token
	if idx := strings.Index(result, "##{"); idx >= 0 {
		bracesBeforeToken := 0
		for i := 0; i < idx; i++ {
			if result[i] == '{' {
				bracesBeforeToken++
			}
		}
		result = strings.ReplaceAll(result, "##{", strconv.Itoa(p.braceCount+bracesBeforeToken))
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