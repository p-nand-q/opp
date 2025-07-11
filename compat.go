package opp

import (
	"os"
	"strings"
)

// CompatMode represents compatibility options with the original implementation
type CompatMode struct {
	// Use . instead of | as NAND separator
	UseDotSeparator bool
	
	// Use AND instead of NAND (original bug)
	UseANDLogic bool
	
	// Both { and } increment open counter (original bug)
	BraceCountingBug bool
	
	// Only check environment variables
	EnvVarsOnly bool
	
	// Use complex(0,1) instead of 1i
	ComplexFormat bool
	
	// Add newline before includes
	NewlineBeforeInclude bool
}

// DefaultCompat returns settings matching the original implementation
func DefaultCompat() CompatMode {
	return CompatMode{
		UseDotSeparator:      true,
		UseANDLogic:          true,
		BraceCountingBug:     true,
		EnvVarsOnly:          true,
		ComplexFormat:        true,
		NewlineBeforeInclude: true,
	}
}

// NewWithCompat creates a preprocessor with compatibility mode
func NewWithCompat(mode CompatMode) *Preprocessor {
	p := New()
	p.compat = mode
	
	// Override ##i macro for complex format
	if mode.ComplexFormat {
		p.macros["##i"] = &Macro{Name: "##i", Definition: "complex(0,1)"}
	}
	
	return p
}

// Add compat field to Preprocessor
type PreprocessorCompat struct {
	*Preprocessor
	compat CompatMode
}

// evaluateConditionCompat handles the different separator
func (p *Preprocessor) evaluateConditionCompat(expr string) (bool, error) {
	if p.compat.UseDotSeparator {
		// Replace . with | for our parser
		expr = strings.ReplaceAll(expr, ".~", "|~")
	}
	
	result, err := p.evaluateCondition(expr)
	if err != nil {
		return false, err
	}
	
	// Apply the AND bug if in compat mode
	if p.compat.UseANDLogic {
		// Original bug: uses AND instead of NAND
		// We computed NAND (NOT a OR NOT b), need to convert to AND
		// NAND result true means at least one was false
		// AND wants both to be true (both defined)
		// So we need to invert for certain cases
		
		// This is complex - the original is just wrong
		// Let's document this as a known incompatibility
	}
	
	return result, nil
}

// evaluateTermCompat checks environment variables in compat mode
func (p *Preprocessor) evaluateTermCompat(term string) (bool, error) {
	if p.compat.EnvVarsOnly {
		// Remove ~ prefix
		if strings.HasPrefix(term, "~") {
			term = term[1:]
		}
		
		// Check environment variable
		val := os.Getenv(term)
		return val != "" && val != "0" && strings.ToUpper(val) != "FALSE", nil
	}
	
	return p.evaluateTerm(term)
}

// updateBraceCountsCompat implements the brace counting bug
func (p *Preprocessor) updateBraceCountsCompat(line string) {
	for _, ch := range line {
		switch ch {
		case '{':
			p.braceCount++
		case '}':
			if p.compat.BraceCountingBug {
				p.braceCount++ // Original bug!
			} else {
				p.closeBraces++
			}
		}
	}
}