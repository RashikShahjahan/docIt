package docIt

import (
	"flag"
	"fmt"
	"io/ioutil"
)

var (
	help       bool
	inputFile  string
	outputFile string
)

func init() {
	flag.BoolVar(&help, "help", false, "display help message")
	flag.StringVar(&inputFile, "input-file", "", "path to input file")
	flag.StringVar(&outputFile, "out-file", "", "path to output file")
}

func docIt() {
	flag.Parse()

	if help {
		flag.Usage()
		return
	}

	// Check if an input file was specified
	if inputFile == "" {
		fmt.Println("Please specify an input file with the -input-file flag.")
		return
	}

	// Read the input file
	input, err := ioutil.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		return
	}

	// Start background processing
	done := make(chan bool)
	go func() {
		// Process the input file
		output := processFile(input)

		// Write the output file
		if outputFile == "" {
			fmt.Println(output)
		} else {
			err := ioutil.WriteFile(outputFile, []byte(output), 0644)
			if err != nil {
				fmt.Printf("Error writing output file: %v\n", err)
				return
			}
		}
		done <- true
	}()

	// Wait for background processing to complete
	<-done
}

func processFile(input []byte) string {
	// TODO: Implement processing logic here
	// For example, split the code by functions and create blocks of maximum token size.
	return string(input)
}
