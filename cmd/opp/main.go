package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/p-nand-q/opp"
)

func main() {
	var (
		output  = flag.String("o", "", "Output file (default: stdout)")
		defines flagList
	)
	
	flag.Var(&defines, "D", "Define a variable (can be used multiple times)")
	flag.Parse()
	
	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <input-file>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	
	inputFile := flag.Arg(0)
	
	// Read input file
	input, err := ioutil.ReadFile(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}
	
	// Create preprocessor
	preprocessor := opp.New()
	
	// Apply command-line defines
	for _, def := range defines {
		parts := strings.SplitN(def, "=", 2)
		if len(parts) == 2 {
			preprocessor.Define(parts[0], parts[1])
		} else {
			preprocessor.Define(parts[0], "1")
		}
	}
	
	// Process the input
	result, err := preprocessor.Process(string(input))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Preprocessing error: %v\n", err)
		os.Exit(1)
	}
	
	// Write output
	if *output != "" {
		err = ioutil.WriteFile(*output, []byte(result), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Print(result)
	}
}

// flagList allows multiple -D flags
type flagList []string

func (f *flagList) String() string {
	return strings.Join(*f, ", ")
}

func (f *flagList) Set(value string) error {
	*f = append(*f, value)
	return nil
}