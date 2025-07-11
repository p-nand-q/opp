package opp

import (
	"testing"
)

func TestProcessStringizeCharize(t *testing.T) {
	tests := []struct {
		name       string
		definition string
		args       []string
		expected   string
	}{
		{
			name:       "simple stringize",
			definition: `#"#0`,
			args:       []string{"hello"},
			expected:   `"hello"`,
		},
		{
			name:       "simple charize",
			definition: `#'#0`,
			args:       []string{"a"},
			expected:   `'a'`,
		},
		{
			name:       "stringize with spaces",
			definition: `#"#0`,
			args:       []string{"hello world"},
			expected:   `"hello world"`,
		},
		{
			name:       "stringize with quotes",
			definition: `#"#0`,
			args:       []string{`test"quote`},
			expected:   `"test\"quote"`,
		},
		{
			name:       "charize with quotes",
			definition: `#'#0`,
			args:       []string{"it's"},
			expected:   `'it\'s'`,
		},
		{
			name:       "mixed usage",
			definition: `printf("%s: %s\n", #"#0, #"#1)`,
			args:       []string{"main", "Starting program"},
			expected:   `printf("%s: %s\n", "main", "Starting program")`,
		},
		{
			name:       "with regular substitution",
			definition: `if (!(#0)) error(#"#0)`,
			args:       []string{"x > 0"},
			expected:   `if (!(x > 0)) error("x > 0")`,
		},
		{
			name:       "multiple arguments",
			definition: `{#"#0, #'#1}`,
			args:       []string{"name", "c"},
			expected:   `{"name", 'c'}`,
		},
		{
			name:       "escape backslashes",
			definition: `#"#0`,
			args:       []string{`path\to\file`},
			expected:   `"path\\to\\file"`,
		},
		{
			name:       "empty argument",
			definition: `#"#0`,
			args:       []string{""},
			expected:   `""`,
		},
		{
			name:       "numeric argument",
			definition: `#"#0`,
			args:       []string{"123"},
			expected:   `"123"`,
		},
		{
			name:       "expression argument",
			definition: `#"#0`,
			args:       []string{"2+2"},
			expected:   `"2+2"`,
		},
		{
			name:       "no space between operator and arg",
			definition: `#"#0 and #'#1`,
			args:       []string{"foo", "bar"},
			expected:   `"foo" and 'bar'`,
		},
		{
			name:       "operator not before arg ref",
			definition: `#" #0`,
			args:       []string{"test"},
			expected:   `#" test`, // #" is literal, #0 is substituted
		},
		{
			name:       "incomplete operator",
			definition: `#"text`,
			args:       []string{},
			expected:   `#"text`,
		},
		{
			name:       "arg ref without operator",
			definition: `value: #0`,
			args:       []string{"42"},
			expected:   `value: 42`,
		},
	}
	
	p := New()
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := p.processStringizeCharize(tt.definition, tt.args)
			if result != tt.expected {
				t.Errorf("processStringizeCharize(%q, %v) = %q, want %q", 
					tt.definition, tt.args, result, tt.expected)
			}
		})
	}
}

// Test stringize/charize within actual macro processing
func TestStringizeInMacros(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "basic stringize macro",
			input: `##:STR(x) #"#0
STR(hello)`,
			expected: `"hello"`,
		},
		{
			name: "basic charize macro",
			input: `##:CHR(x) #'#0
CHR(a)`,
			expected: `'a'`,
		},
		{
			name: "stringize with spaces",
			input: `##:STR(x) #"#0
STR(hello world)`,
			expected: `"hello world"`,
		},
		{
			name: "stringize with quotes to escape",
			input: `##:STR(x) #"#0
STR(test"quote)`,
			expected: `"test\"quote"`,
		},
		{
			name: "charize with quotes to escape",
			input: `##:CHR(x) #'#0
CHR(it's)`,
			expected: `'it\'s'`,
		},
		{
			name: "multiple arguments",
			input: `##:LOG(func,msg) printf("%s: %s\n", #"#0, #"#1)
LOG(main, Starting program)`,
			expected: `printf("%s: %s\n", "main", "Starting program")`,
		},
		{
			name: "assert macro",
			input: `##:ASSERT(expr) if (!(#0)) error(#"#0)
ASSERT(x > 0)`,
			expected: `if (!(x > 0)) error("x > 0")`,
		},
		{
			name: "pair macro",
			input: `##:PAIR(a,b) {#"#0, #'#1}
PAIR(name, c)`,
			expected: `{"name", 'c'}`,
		},
		{
			name: "empty argument",
			input: `##:STR(x) #"#0
STR()`,
			expected: `""`,
		},
		{
			name: "numeric argument",
			input: `##:STR(x) #"#0
STR(123)`,
			expected: `"123"`,
		},
		{
			name: "expression argument",
			input: `##:STR(x) #"#0
STR(2+2)`,
			expected: `"2+2"`,
		},
		{
			name: "nested parentheses in argument",
			input: `##:STR(x) #"#0
STR(f(a, b))`,
			expected: `"f(a, b)"`,
		},
		{
			name: "backslashes in argument",
			input: `##:STR(x) #"#0
STR(path\to\file)`,
			expected: `"path\\to\\file"`,
		},
		{
			name: "macro without stringize",
			input: `##:NORMAL(x) [#0]
NORMAL(test)`,
			expected: `[test]`,
		},
		{
			name: "mixed stringize and normal",
			input: `##:MIX(a,b) #0 = #"#1
MIX(var, value)`,
			expected: `var = "value"`,
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