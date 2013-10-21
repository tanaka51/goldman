package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type options struct {
	Target     string
	FileTypes  []string
	TargetDirs []string
	Color      bool
}

func parseOptions() (*options, error) {
	options := new(options)
	var (
		targetDirs string
		fileTypes  string
		color      bool
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

	splitWithBlankSlice := func(s, sep string) []string {
		if s == "" {
			return []string{}
		}
		return strings.Split(s, sep)
	}

	options.Target = os.Args[1]
	options.TargetDirs = splitWithBlankSlice(targetDirs, ",")
	options.FileTypes = splitWithBlankSlice(fileTypes, ",")
	options.Color = color

	return options, nil
}

func main() {
	options, err := parseOptions()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// TODO: Use goroutine to directory walking

	fileResults := []fileResult{}
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		fr, err := goldman(path, options.Target, options)
		if err != nil {
			fmt.Println(err)
			return err
		}
		if fr.Exists {
			fileResults = append(fileResults, *fr)
		}

		return nil
	}

	for _, dir := range options.TargetDirs {
		filepath.Walk(dir, walkFn)
	}

	for i, fr := range fileResults {
		fr.Puts(options.Color)
		if i != len(fileResults)-1 {
			fmt.Println("")
		}
	}
}
