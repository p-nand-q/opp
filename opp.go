// Package opp implements the Obfuscated Pre-Processor
package opp

import (
	"fmt"
	"strings"
)

// Preprocessor represents an OPP preprocessor instance
type Preprocessor struct {
	macros      map[string]*Macro
	variables   map[string]bool
	random      *RandomGenerator
	lineNumber  int
	braceCount  int
	closeBraces int
	compat      CompatMode
}

// Macro represents a macro definition
type Macro struct {
	Name       string
	Definition string
	IsOperator bool
}

// RandomGenerator provides pseudo-random numbers for ##$
type RandomGenerator struct {
	seed int
}

// New creates a new OPP preprocessor instance
func New() *Preprocessor {
	p := &Preprocessor{
		macros:      make(map[string]*Macro),
		variables:   make(map[string]bool),
		random:      &RandomGenerator{seed: 42},
		lineNumber:  1,
		braceCount:  0,
		closeBraces: 0,
	}
	
	// Initialize predefined macros
	p.initPredefinedMacros()
	
	return p
}

// Define sets a variable as defined
func (p *Preprocessor) Define(name string, value string) {
	p.variables[name] = true
	if value != "" {
		p.macros[name] = &Macro{
			Name:       name,
			Definition: value,
		}
	}
}

// Undefine removes a variable or macro
func (p *Preprocessor) Undefine(name string) {
	delete(p.variables, name)
	delete(p.macros, name)
}

// Process processes the input source code
func (p *Preprocessor) Process(input string) (string, error) {
	lines := strings.Split(input, "\n")
	output := &strings.Builder{}
	
	conditionalStack := &ConditionalStack{}
	
	for i, line := range lines {
		p.lineNumber = i + 1
		
		processedLine, err := p.processLine(line, conditionalStack)
		if err != nil {
			return "", fmt.Errorf("line %d: %w", p.lineNumber, err)
		}
		
		if processedLine != "" {
			if output.Len() > 0 {
				output.WriteString("\n")
			}
			output.WriteString(processedLine)
		}
		
		// Update brace counts after processing the line
		for _, ch := range line {
			switch ch {
			case '{':
				p.braceCount++
			case '}':
				p.closeBraces++
			}
		}
	}
	
	if !conditionalStack.IsEmpty() {
		return "", fmt.Errorf("unclosed conditional block")
	}
	
	return output.String(), nil
}

// ProcessFile processes a file
func (p *Preprocessor) ProcessFile(filename string) (string, error) {
	// TODO: Implement file reading
	return "", fmt.Errorf("file processing not implemented yet")
}

func (p *Preprocessor) initPredefinedMacros() {
	// ##i - imaginary unit (requires complex.h)
	p.macros["##i"] = &Macro{Name: "##i", Definition: "1i"}
	
	// ##_, ##$, ##{, ##} are handled dynamically in expansion
}

func (p *Preprocessor) updateBraceCounts(line string) {
	// Count braces before processing the line
	// This ensures ##{ reflects count at start of line
	// Note: This is done after processing to match the behavior
}

func (r *RandomGenerator) Next() int {
	r.seed = (r.seed*1103515245 + 12345) & 0x7fffffff
	return r.seed
}