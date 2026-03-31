package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/boriskachscovsky/md2pdf/internal/cli"
	"github.com/boriskachscovsky/md2pdf/internal/pdf"
	"github.com/boriskachscovsky/md2pdf/internal/render"
)

func main() {
	options, err := cli.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	markdown, err := os.ReadFile(options.InputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: reading %s: %v\n", options.InputPath, err)
		os.Exit(1)
	}

	htmlDoc, err := render.MarkdownToHTML(markdown, options.InputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: rendering markdown: %v\n", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	renderer, err := pdf.NewRenderer()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if err := renderer.Render(ctx, htmlDoc, options.OutputPath); err != nil {
		fmt.Fprintf(os.Stderr, "error: creating PDF: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("wrote %s\n", options.OutputPath)
}
