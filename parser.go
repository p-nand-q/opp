package opp

import (
	"fmt"
	"strings"
)

// ConditionalStack manages nested conditional compilation
type ConditionalStack struct {
	conditions []bool
	inElse     []bool
}

func (c *ConditionalStack) Push(condition bool) {
	c.conditions = append(c.conditions, condition)
	c.inElse = append(c.inElse, false)
}

func (c *ConditionalStack) Pop() error {
	if len(c.conditions) == 0 {
		return fmt.Errorf("no matching conditional to close")
	}
	c.conditions = c.conditions[:len(c.conditions)-1]
	c.inElse = c.inElse[:len(c.inElse)-1]
	return nil
}

func (c *ConditionalStack) ToggleElse() error {
	if len(c.conditions) == 0 {
		return fmt.Errorf("##@ without matching conditional")
	}
	idx := len(c.inElse) - 1
	c.inElse[idx] = !c.inElse[idx]
	// Flip the condition when entering else
	c.conditions[idx] = !c.conditions[idx]
	return nil
}

func (c *ConditionalStack) ShouldProcess() bool {
	for _, cond := range c.conditions {
		if !cond {
			return false
		}
	}
	return true
}

func (c *ConditionalStack) IsEmpty() bool {
	return len(c.conditions) == 0
}

func (p *Preprocessor) processLine(line string, stack *ConditionalStack) (string, error) {
	trimmed := strings.TrimSpace(line)
	
	// Check for OPP directives
	if strings.HasPrefix(trimmed, "##") {
		return p.processDirective(trimmed, stack)
	}
	
	// If we're in a false conditional block, skip the line
	if !stack.ShouldProcess() {
		return "", nil
	}
	
	// Process macros in the line
	return p.expandMacros(line)
}

func (p *Preprocessor) processDirective(line string, stack *ConditionalStack) (string, error) {
	// Remove ## prefix
	directive := line[2:]
	
	switch {
	case directive == ".":
		// End conditional block
		return "", stack.Pop()
		
	case strings.HasPrefix(directive, "@"):
		// Else-like behavior - the rest is a new condition
		err := stack.ToggleElse()
		if err != nil {
			return "", err
		}
		// If there's a condition after @, evaluate it
		if len(directive) > 1 {
			condition, err := p.evaluateCondition(directive[1:])
			if err != nil {
				return "", err
			}
			// Replace the toggled condition with the new one
			if len(stack.conditions) > 0 {
				stack.conditions[len(stack.conditions)-1] = condition
			}
		}
		return "", nil
		
	case strings.HasPrefix(directive, "~"):
		// Conditional compilation
		var condition bool
		condition, err := p.evaluateCondition(directive)
		if err != nil {
			return "", err
		}
		stack.Push(condition)
		return "", nil
		
	case strings.HasPrefix(directive, "<"):
		// Include file
		if !stack.ShouldProcess() {
			return "", nil
		}
		return p.processInclude(directive)
		
	case strings.HasPrefix(directive, ":"):
		// Define macro
		if !stack.ShouldProcess() {
			return "", nil
		}
		return "", p.defineMacro(directive[1:])
		
	case strings.HasPrefix(directive, "-"):
		// Undefine macro
		if !stack.ShouldProcess() {
			return "", nil
		}
		p.Undefine(directive[1:])
		return "", nil
		
	default:
		// Check for predefined macros
		if !stack.ShouldProcess() {
			return "", nil
		}
		return p.expandPredefinedMacro(line)
	}
}

func (p *Preprocessor) evaluateCondition(expr string) (bool, error) {
	// Handle parentheses at the expression level
	expr = strings.TrimSpace(expr)
	
	// Find the rightmost | that's not inside parentheses
	parenDepth := 0
	splitPos := -1
	
	for i := len(expr) - 1; i >= 0; i-- {
		switch expr[i] {
		case ')':
			parenDepth++
		case '(':
			parenDepth--
		case '|':
			if parenDepth == 0 {
				splitPos = i
				break
			}
		}
	}
	
	if splitPos == -1 {
		// No | found, this should be a single term
		return p.evaluateTerm(expr)
	}
	
	// Split at the |
	left := expr[:splitPos]
	right := expr[splitPos+1:]
	
	// Evaluate each part
	leftVal, err := p.evaluateTerm(strings.TrimSpace(left))
	if err != nil {
		return false, err
	}
	
	rightVal, err := p.evaluateTerm(strings.TrimSpace(right))
	if err != nil {
		return false, err
	}
	
	// NAND logic: ~a | ~b
	return !leftVal || !rightVal, nil
}

func (p *Preprocessor) evaluateTerm(term string) (bool, error) {
	// Remove ~ prefix
	if !strings.HasPrefix(term, "~") {
		return false, fmt.Errorf("term must start with ~: %s", term)
	}
	term = term[1:]
	
	// Handle parentheses
	if strings.HasPrefix(term, "(") && strings.HasSuffix(term, ")") {
		inner := term[1 : len(term)-1]
		return p.evaluateCondition(inner)
	}
	
	// Check if variable is defined
	_, defined := p.variables[term]
	return defined, nil
}