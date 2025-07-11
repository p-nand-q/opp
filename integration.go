package opp

import (
	"io"
	"io/ioutil"
)

// ProcessReader processes input from an io.Reader
func (p *Preprocessor) ProcessReader(r io.Reader) (string, error) {
	input, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}
	return p.Process(string(input))
}

// ProcessWriter processes input and writes to an io.Writer
func (p *Preprocessor) ProcessWriter(r io.Reader, w io.Writer) error {
	result, err := p.ProcessReader(r)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(result))
	return err
}

// DefaultPreprocessor returns a preprocessor with standard configuration
func DefaultPreprocessor() *Preprocessor {
	return New()
}

// WithDefines creates a preprocessor with predefined variables
func WithDefines(defines map[string]string) *Preprocessor {
	p := New()
	for k, v := range defines {
		p.Define(k, v)
	}
	return p
}