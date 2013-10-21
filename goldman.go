package main

import (
	"fmt"
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"flag"
)

type options struct {
	Target string
	FileTypes []string
	TargetDirs []string
	Color bool
}

// grep per file
func goldman(filename, target string, options *options) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
  defer f.Close()

	type pageResult struct {
		lineNumber int
		line string
	}
	results := []pageResult{}

	scanner := bufio.NewScanner(f)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber += 1
		line := string(scanner.Text())
		if strings.Contains(line, target) {
			results = append(results, pageResult{lineNumber, line})
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if len(results) != 0 {
		if options.Color {
			fmt.Println("\x1b[1b" + filename + "\x1b[0m")
		} else {
			fmt.Println(filename)
		}
		for _, result := range results {
			if options.Color {
				result.line = strings.Replace(result.line, target, "\x1b[43m" + target + "\x1b[0m", -1)
			}
			fmt.Printf("%d: %s\n", result.lineNumber, result.line)
		}
		fmt.Println("")
	}

	return nil
}

func SplitExt(s, sep string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, sep)
}

func parseOptions() (*options, error) {
	options := new(options)
	var (
		targetDirs string
		fileTypes string
		color bool
	)

	f := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	// target directories
	f.StringVar(&targetDirs, "target-dirs", ".", "specify target dirs. you can use ',' for specify multi value")
	f.StringVar(&targetDirs, "d", ".", "alias for target-dirs")
	// color
	f.BoolVar(&color, "color", true, "enable color (default true)")
	f.BoolVar(&color, "c", true, "alias for color")

	f.Parse(os.Args[1:])
	for 0 < f.NArg() {
		f.Parse(f.Args()[1:])
	}

	options.Target = os.Args[1]
	options.TargetDirs = SplitExt(targetDirs, ",")
	options.FileTypes = SplitExt(fileTypes, ",")
	options.Color = color

	return options, nil
}

func main() {

	// TODO: Parse command line options
	options, err := parseOptions()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// TODO: Use goroutine to directory walking

	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		err = goldman(path, options.Target, options)
		if err != nil {
			fmt.Println(err)
			return err
		}

		return nil
	}

	for _, dir := range options.TargetDirs {
		filepath.Walk(dir, walkFn)
	}
}
