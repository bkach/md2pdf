package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Options struct {
	InputPath  string
	OutputPath string
}

func Parse(args []string) (Options, error) {
	fs := flag.NewFlagSet("md2pdf", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var output string
	fs.StringVar(&output, "o", "", "output PDF path")

	if err := fs.Parse(args); err != nil {
		return Options{}, err
	}

	if fs.NArg() != 1 {
		return Options{}, errors.New("usage: md2pdf [-o output.pdf] input.md")
	}

	inputPath, err := filepath.Abs(fs.Arg(0))
	if err != nil {
		return Options{}, fmt.Errorf("resolve input path: %w", err)
	}

	info, err := os.Stat(inputPath)
	if err != nil {
		return Options{}, fmt.Errorf("stat input path: %w", err)
	}
	if info.IsDir() {
		return Options{}, errors.New("input path must be a markdown file, not a directory")
	}

	outputPath, err := ResolveOutputPath(inputPath, output)
	if err != nil {
		return Options{}, err
	}

	return Options{
		InputPath:  inputPath,
		OutputPath: outputPath,
	}, nil
}

func ResolveOutputPath(inputPath, outputFlag string) (string, error) {
	if outputFlag != "" {
		return filepath.Abs(outputFlag)
	}

	dir := filepath.Dir(inputPath)
	base := filepath.Base(inputPath)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	if name == "" {
		return "", errors.New("input file must have a valid name")
	}

	return filepath.Join(dir, name+".pdf"), nil
}
