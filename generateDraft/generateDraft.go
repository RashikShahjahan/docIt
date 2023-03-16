package generateDraft

import (
	"os"
	"regexp"
	"strings"
	"text/scanner"

	gpt35 "github.com/AlmazDelDiablo/gpt3-5-turbo-go"
)

var languageRegex = map[string]*regexp.Regexp{
	"Python":     regexp.MustCompile(`^\s*def\s+(?P<function_name>\w+)\(`),
	"Java":       regexp.MustCompile(`^\s*(public|protected|private)?\s*(static)?\s*(final)?\s*(\w+\s+)*\s*(?P<function_name>\w+)\s*\(`),
	"C++":        regexp.MustCompile(`^\s*(template\s*<.*>\s*)?\w+\s+\w+::(?P<function_name>\w+)\(`),
	"JavaScript": regexp.MustCompile(`^\s*(async\s+)?function\s+(?P<function_name>\w+)\s*\(`),
	"PHP":        regexp.MustCompile(`^\s*(public|private|protected)?\s*function\s+(?P<function_name>\w+)\(`),
	"Ruby":       regexp.MustCompile(`^\s*(def)\s+(?P<function_name>\w+)\s*(\(|\s)`),
	"C#":         regexp.MustCompile(`^\s*(public|private|protected)?\s*(static)?\s*(async)?\s*(\w+\s+)*\s*(?P<function_name>\w+)\s*\(`),
	"Go":         regexp.MustCompile(`^\s*func\s+(?P<function_name>\w+)\s*\(`),
	"Swift":      regexp.MustCompile(`^\s*(func)\s+(?P<function_name>\w+)\(`),
	"TypeScript": regexp.MustCompile(`^\s*(async\s+)?function\s+(?P<function_name>\w+)\s*\(`),
	"Kotlin":     regexp.MustCompile(`^\s*(fun)\s+(?P<function_name>\w+)\(`),
	"Rust":       regexp.MustCompile(`^\s*(fn)\s+(?P<function_name>\w+)\(`),
	"Scala":      regexp.MustCompile(`^\s*(def)\s+(?P<function_name>\w+)\(`),
	"Lua":        regexp.MustCompile(`^\s*(function)\s+(?P<function_name>\w+)\(`),
	"C":          regexp.MustCompile(`^\s*(\w+\s+)+\s*(?P<function_name>\w+)\s*\(`),
}

func GenerateDraft(code string, language string) string {
	codeArray := splitCodeByFunction(code, languageRegex[language])
	codeChunks := createBlocksOfMaxTokens(codeArray)
	libDesc := processBlocks(codeChunks, createDesc)
	results := generateFunctions(libDesc)

	return results
}

func createDesc(code string) string {
	apiKey := os.Getenv("OPENAI_API_KEY")
	c := gpt35.NewClient(apiKey)
	req := &gpt35.Request{
		Model:       gpt35.ModelGpt35Turbo,
		Temperature: 0,
		Messages: []*gpt35.Message{
			{
				Role:    gpt35.RoleSystem,
				Content: "You are a helpful assistant who represents key information in a condensed format given the content's of a code file.",
			},
			{
				Role:    gpt35.RoleUser,
				Content: code,
			},
		},
	}

	resp, err := c.GetChat(req)
	if err != nil {
		panic(err)
	}
	desc := resp.Choices[0].Message.Content
	return desc

}

func generateFunctions(desc string) string {
	apiKey := os.Getenv("OPENAI_API_KEY")
	c := gpt35.NewClient(apiKey)
	req := &gpt35.Request{
		Model:       gpt35.ModelGpt35Turbo,
		Temperature: 0,

		Messages: []*gpt35.Message{
			{
				Role:    gpt35.RoleSystem,
				Content: "You are a helpful assistant who drafts code documentation in Markdown format from content's of a function description file.",
			},
			{
				Role:    gpt35.RoleUser,
				Content: desc,
			},
		},
	}

	resp, err := c.GetChat(req)
	if err != nil {
		panic(err)
	}

	functions := resp.Choices[0].Message.Content

	return functions
}

func splitCodeByFunction(code string, functionRegex *regexp.Regexp) []string {
	functionBlocks := []string{}
	currentBlock := []string{}
	for _, line := range strings.Split(code, "\n") {
		if functionRegex.MatchString(line) {
			if len(currentBlock) > 0 {
				functionBlocks = append(functionBlocks, strings.Join(currentBlock, "\n"))
			}
			currentBlock = []string{line}
		} else {
			currentBlock = append(currentBlock, line)
		}
	}
	if len(currentBlock) > 0 {
		functionBlocks = append(functionBlocks, strings.Join(currentBlock, "\n"))
	}
	return functionBlocks
}

func createBlocksOfMaxTokens(functionBlocks []string) []string {
	maxTokens := 2048
	blocksOfMaxTokens := []string{}
	currentChunk := []string{}
	for _, block := range functionBlocks {
		block = regexp.MustCompile(`(\"\"\"|''')[\s\S]*?`).ReplaceAllString(block, `""`)
		var s scanner.Scanner
		s.Init(strings.NewReader(block))
		tokenCount := 0
		for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
			tokenCount++
		}
		if len(strings.Join(currentChunk, "\n"))+tokenCount <= maxTokens {
			currentChunk = append(currentChunk, block)
		} else {
			blocksOfMaxTokens = append(blocksOfMaxTokens, strings.Join(currentChunk, "\n"))
			currentChunk = []string{block}
		}
	}
	if len(currentChunk) > 0 {
		blocksOfMaxTokens = append(blocksOfMaxTokens, strings.Join(currentChunk, "\n"))
	}
	return blocksOfMaxTokens
}

func processBlocks(blocks []string, processFunc func(string) string) string {
	output := ""
	previousOutput := ""
	for _, block := range blocks {
		inputBlock := previousOutput + block
		outputBlock := processFunc(inputBlock)
		output += outputBlock
		previousOutput = outputBlock
	}
	return output
}
