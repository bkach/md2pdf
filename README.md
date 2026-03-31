# md2pdf

`md2pdf` is a small Go command-line tool that renders Markdown to PDF using a browser print pipeline inspired by the VS Code/Cursor Markdown PDF extension.

It follows the same broad approach:

1. Parse Markdown into HTML.
2. Apply a polished document stylesheet and code highlighting.
3. Render that HTML in headless Chrome.
4. Print the page to PDF.

## Features

- Clean command-line interface: `md2pdf input.md`
- Optional output override: `md2pdf -o custom.pdf input.md`
- Sensible default output path next to the input file
- GitHub-like document styling with print-friendly pagination
- GFM-oriented Markdown support including tables, task lists, footnotes, fenced code blocks, blockquotes, and inline HTML
- Local image support via relative paths
- Tests for CLI behavior and HTML rendering

## Requirements

- Go 1.24 or newer
- Google Chrome installed locally

The current implementation looks for Chrome in common locations, including:

- `/Applications/Google Chrome.app/Contents/MacOS/Google Chrome`
- `google-chrome`
- `chrome`
- `chromium`

## Usage

```bash
go run ./cmd/md2pdf -- testdata/comprehensive.md
go run ./cmd/md2pdf -- -o ./output/custom.pdf testdata/comprehensive.md
```

After building:

```bash
go build -o md2pdf ./cmd/md2pdf
./md2pdf testdata/comprehensive.md
./md2pdf -o testdata/comprehensive-output.pdf testdata/comprehensive.md
```

If `-o` is omitted, `input.md` becomes `input.pdf` in the same directory as the Markdown file.

## Project Layout

- `cmd/md2pdf`: CLI entrypoint
- `internal/cli`: argument parsing and output path resolution
- `internal/render`: Markdown to standalone HTML document rendering
- `internal/pdf`: Chrome discovery and HTML-to-PDF generation
- `testdata`: comprehensive Markdown fixture and local assets

## Extending

Good next extension points:

- front-matter driven metadata and page settings
- custom CSS/theme loading from disk
- header/footer templates
- page size and margin flags
- remote asset fetching policies

## Verification

The repository includes `testdata/comprehensive.md`, which exercises most commonly used Markdown features. The intended validation loop is:

```bash
go test ./...
go run ./cmd/md2pdf -- testdata/comprehensive.md
```

