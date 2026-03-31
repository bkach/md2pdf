package render

import (
	"bytes"
	"fmt"
	"html/template"
	"net/url"
	"path/filepath"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

type documentData struct {
	BaseURL template.URL
	Body    template.HTML
	CSS     template.CSS
	Title   string
}

func MarkdownToHTML(markdown []byte, inputPath string) (string, error) {
	var body bytes.Buffer

	engine := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.DefinitionList,
			extension.Footnote,
			extension.Linkify,
			extension.Typographer,
			highlighting.NewHighlighting(
				highlighting.WithStyle("github"),
				highlighting.WithGuessLanguage(true),
				highlighting.WithFormatOptions(chromahtml.WithClasses(false)),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)

	if err := engine.Convert(markdown, &body); err != nil {
		return "", fmt.Errorf("convert markdown: %w", err)
	}

	baseURL, err := directoryBaseURL(inputPath)
	if err != nil {
		return "", err
	}

	data := documentData{
		BaseURL: template.URL(baseURL),
		Body:    template.HTML(body.String()),
		CSS:     template.CSS(documentCSS),
		Title:   filepath.Base(inputPath),
	}

	var page bytes.Buffer
	if err := documentTemplate.Execute(&page, data); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}

	return page.String(), nil
}

func directoryBaseURL(inputPath string) (string, error) {
	abs, err := filepath.Abs(filepath.Dir(inputPath))
	if err != nil {
		return "", fmt.Errorf("resolve base dir: %w", err)
	}

	u := &url.URL{Scheme: "file", Path: filepath.ToSlash(abs) + "/"}
	return u.String(), nil
}

var documentTemplate = template.Must(template.New("page").Parse(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{{ .Title }}</title>
  <base href="{{ .BaseURL }}">
  <style>{{ .CSS }}</style>
</head>
<body>
  <main class="page-shell">
    <article class="markdown-body">
{{ .Body }}
    </article>
  </main>
</body>
</html>
`))

const documentCSS = `
@page {
  size: A4;
  margin: 18mm 16mm 20mm 16mm;
}

:root {
  color-scheme: light;
  --bg: #ffffff;
  --fg: #1f2328;
  --muted: #59636e;
  --subtle: #f6f8fa;
  --border: #d0d7de;
  --accent: #0969da;
  --quote: #eaeef2;
  --code-bg: #f6f8fa;
  --pre-bg: #0d1117;
}

* {
  box-sizing: border-box;
}

html {
  background: var(--bg);
}

body {
  margin: 0;
  color: var(--fg);
  background: var(--bg);
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif;
  font-size: 12px;
  line-height: 1.6;
}

.page-shell {
  width: 100%;
}

.markdown-body > :first-child {
  margin-top: 0;
}

.markdown-body > :last-child {
  margin-bottom: 0;
}

.markdown-body h1,
.markdown-body h2,
.markdown-body h3,
.markdown-body h4,
.markdown-body h5,
.markdown-body h6 {
  margin: 1.4em 0 0.6em;
  line-height: 1.25;
  font-weight: 600;
  break-after: avoid-page;
}

.markdown-body h1,
.markdown-body h2 {
  padding-bottom: 0.3em;
  border-bottom: 1px solid var(--border);
}

.markdown-body h1 { font-size: 2em; }
.markdown-body h2 { font-size: 1.5em; }
.markdown-body h3 { font-size: 1.25em; }
.markdown-body h4 { font-size: 1em; }
.markdown-body h5 { font-size: 0.875em; }
.markdown-body h6 { font-size: 0.85em; color: var(--muted); }

.markdown-body p,
.markdown-body ul,
.markdown-body ol,
.markdown-body dl,
.markdown-body table,
.markdown-body pre,
.markdown-body blockquote {
  margin: 0 0 1em;
}

.markdown-body ul,
.markdown-body ol {
  padding-left: 1.6em;
}

.markdown-body li + li {
  margin-top: 0.25em;
}

.markdown-body li > p {
  margin-bottom: 0.5em;
}

.markdown-body a {
  color: var(--accent);
  text-decoration: none;
}

.markdown-body a:hover {
  text-decoration: underline;
}

.markdown-body hr {
  border: 0;
  border-top: 1px solid var(--border);
  margin: 1.8em 0;
}

.markdown-body blockquote {
  padding: 0.1em 1em;
  color: var(--muted);
  border-left: 0.25em solid var(--quote);
}

.markdown-body code {
  padding: 0.16em 0.35em;
  border-radius: 6px;
  background: var(--code-bg);
  font-family: "SFMono-Regular", SFMono-Regular, ui-monospace, Menlo, monospace;
  font-size: 0.9em;
}

.markdown-body pre {
  overflow: auto;
  padding: 1em;
  border-radius: 10px;
  background: var(--pre-bg);
  color: #e6edf3;
  break-inside: avoid-page;
}

.markdown-body pre code {
  padding: 0;
  background: transparent;
  color: inherit;
}

.markdown-body table {
  width: 100%;
  border-collapse: collapse;
  break-inside: avoid-page;
}

.markdown-body th,
.markdown-body td {
  padding: 0.55em 0.75em;
  border: 1px solid var(--border);
  text-align: left;
  vertical-align: top;
}

.markdown-body thead th {
  background: var(--subtle);
}

.markdown-body img {
  max-width: 100%;
  border-radius: 8px;
}

.markdown-body input[type="checkbox"] {
  margin-right: 0.5em;
}

.markdown-body .footnotes {
  margin-top: 2em;
  color: var(--muted);
  font-size: 0.92em;
}

.markdown-body .footnotes hr {
  margin-bottom: 1em;
}

.markdown-body .highlight pre,
.markdown-body pre.chroma {
  background: var(--pre-bg) !important;
}

@media print {
  a {
    color: inherit;
  }
}
`
