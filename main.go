package main

import (
	"fmt"
	"github.com/RashikShahjahan/docIt/generateDraft"
	"io/ioutil"
	"os"
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

func main() {
	args := os.Args
	for i := range args {
		input, err := ioutil.ReadFile(args[i])
		if err != nil {
			fmt.Printf("Error reading input file: %v\n", err)
			return
		}
		ext := filepath.Ext(args[i])

		if language, ok := extToLang[ext]; ok {
			done := make(chan bool)

			go singleFile(string(input), language, args[i][:len(args[i])-len(ext)]+".md", done)

			<-done
		}

	}
}

func singleFile(input string, language string, outputFile string, done chan<- bool) {
	output := generateDraft.GenerateDraft(string(input), language)
	err := ioutil.WriteFile(outputFile, []byte(output), 0644)
	if err != nil {
		fmt.Printf("Error writing output file: %v\n", err)
		return
	}
	done <- true
}
