package main

import "os"

func main() {
	// NOTE: an expectation of a Diagnostic is specified by a string literal
	// containing a regular expression that must match the diagnostic message
	os.Exit(0) // want `direct call of os.Exit\(\) in main\(\) function of main package is not allowed`
}
