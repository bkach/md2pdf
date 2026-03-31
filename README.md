# md2pdf

Small Go CLI for rendering Markdown to PDF with headless Chrome.

## Requirements

- Go 1.24+
- Google Chrome installed locally

## Usage

```bash
go run ./cmd/md2pdf -- input.md
go run ./cmd/md2pdf -- -o output.pdf input.md
```

Or build it:

```bash
go build -o md2pdf ./cmd/md2pdf
./md2pdf input.md
./md2pdf -o output.pdf input.md
```

If `-o` is not provided, the output PDF is written next to the input file using the same base name.

Examples:

```bash
./md2pdf testdata/comprehensive.md
./md2pdf -o output/custom-result.pdf testdata/comprehensive.md
```

## Notes

- Supports common GitHub-flavored Markdown features including tables, task lists, fenced code blocks, and footnotes
- Supports local relative images
- Uses Chrome to print styled HTML to PDF

## Development

```bash
go test ./...
```
