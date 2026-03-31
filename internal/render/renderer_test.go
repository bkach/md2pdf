package render

import (
	"strings"
	"testing"
)

func TestMarkdownToHTMLProducesStandaloneDocument(t *testing.T) {
	t.Parallel()

	markdown := []byte("# Title\n\n- [x] done\n\nA footnote.[^1]\n\n[^1]: note")
	html, err := MarkdownToHTML(markdown, "/tmp/docs/example.md")
	if err != nil {
		t.Fatalf("MarkdownToHTML returned error: %v", err)
	}

	assertContains(t, html, "<!doctype html>")
	assertContains(t, html, `<base href="file:///tmp/docs/">`)
	assertContains(t, html, "<h1 id=\"title\">Title</h1>")
	assertContains(t, html, "checkbox")
	assertContains(t, html, "footnotes")
}

func TestMarkdownToHTMLPreservesRelativeImageReferences(t *testing.T) {
	t.Parallel()

	html, err := MarkdownToHTML([]byte("![Alt](assets/diagram.svg)"), "/tmp/docs/example.md")
	if err != nil {
		t.Fatalf("MarkdownToHTML returned error: %v", err)
	}

	assertContains(t, html, `src="assets/diagram.svg"`)
}

func assertContains(t *testing.T, body, want string) {
	t.Helper()
	if !strings.Contains(body, want) {
		t.Fatalf("expected output to contain %q", want)
	}
}
