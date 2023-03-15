package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

var (
	help       bool
	inputFile  string
	outputFile string
)

var extToLang = map[string]string{
	".py":    "Python",
	".java":  "Java",
	".cpp":   "C++",
	".js":    "JavaScript",
	".php":   "PHP",
	".rb":    "Ruby",
	".cs":    "C#",
	".go":    "Go",
	".swift": "Swift",
	".ts":    "TypeScript",
	".kt":    "Kotlin",
	".rs":    "Rust",
	".scala": "Scala",
	".lua":   "Lua",
	".c":     "C",
}

func init() {
	flag.BoolVar(&help, "help", false, "display help message")
	flag.StringVar(&inputFile, "input-file", "", "path to input file")
	flag.StringVar(&outputFile, "out-file", "", "path to output file")
}

func main() {
	flag.Parse()

	if help {
		flag.Usage()
		return
	}

	if inputFile == "" {
		fmt.Println("Please specify an input file with the -input-file flag.")
		return
	}

	input, err := ioutil.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		return
	}
	ext := filepath.Ext(inputFile)
	language := extToLang[ext]

	// Create channel to wait for background processing to complete
	done := make(chan bool)

	// Start background processing
	go runBackground(string(input), language, done)

	// Wait for background processing to complete
	<-done

}

func runBackground(input string, language string, done chan<- bool) {
	output := generateDraft.generateDraft(string(input), language)

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
}
