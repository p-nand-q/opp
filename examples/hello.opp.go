package main

import "fmt"

##:DEBUG_PRINT fmt.Println(##0..n)

func main() {
	##~(~DEBUG|~DEBUG)|~(~DEBUG|~DEBUG)
	DEBUG_PRINT("Running in debug mode")
	DEBUG_PRINT("Line number minus 5:", ##_)
	DEBUG_PRINT("Random number:", ##$)
	##.
	
	fmt.Println("Hello from OPP!")
	
	##~WINDOWS|~WINDOWS
	fmt.Println("Not on Windows")
	##@~(~WINDOWS|~WINDOWS)|~(~WINDOWS|~WINDOWS)
	fmt.Println("On Windows")
	##.
}