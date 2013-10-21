package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type fileResult struct {
	filename    string
	lineResults []lineResult
	Exists      bool
}

func newFileResult(filename string) *fileResult {
	return &fileResult{filename: filename, Exists: false}
}

func (fr *fileResult) Puts(color bool) {
	if color {
		fmt.Println("\x1b[1m" + fr.filename + "\x1b[0m")
		for _, lr := range fr.lineResults {
			fmt.Printf("%d:%s\n", lr.lineNumber, strings.Replace(lr.line, lr.target, "\x1b[43m"+lr.target+"\x1b[0m", -1))
		}
		return
	}

	fmt.Println(fr.filename)
	for _, lr := range fr.lineResults {
		fmt.Printf("%d: %s\n", lr.lineNumber, lr.target)
	}
}

type lineResult struct {
	lineNumber int
	line       string
	target     string
}

// grep per file
func goldman(filename, target string, options *options) (*fileResult, error) {
	fileResult := newFileResult(filename)

	f, err := os.Open(filename)
	if err != nil {
		return fileResult, err
	}
	defer f.Close()

	lineResults := []lineResult{}

	scanner := bufio.NewScanner(f)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber += 1
		line := string(scanner.Text())
		if strings.Contains(line, target) {
			lineResults = append(lineResults, lineResult{lineNumber, line, target})
		}
	}

	if err := scanner.Err(); err != nil {
		return fileResult, err
	}

	if len(lineResults) == 0 {
		return fileResult, nil
	}

	fileResult.lineResults = lineResults
	fileResult.Exists = true
	return fileResult, nil
}
