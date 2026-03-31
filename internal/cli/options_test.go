package cli

import (
	"path/filepath"
	"testing"
)

func TestResolveOutputPathDefaultsNextToInput(t *testing.T) {
	t.Parallel()

	input := filepath.Join(string(filepath.Separator), "tmp", "docs", "test.md")
	got, err := ResolveOutputPath(input, "")
	if err != nil {
		t.Fatalf("ResolveOutputPath returned error: %v", err)
	}

	want := filepath.Join(string(filepath.Separator), "tmp", "docs", "test.pdf")
	if got != want {
		t.Fatalf("ResolveOutputPath = %q, want %q", got, want)
	}
}

func TestResolveOutputPathHonorsOverride(t *testing.T) {
	t.Parallel()

	got, err := ResolveOutputPath("/tmp/docs/test.md", "/tmp/out/custom.pdf")
	if err != nil {
		t.Fatalf("ResolveOutputPath returned error: %v", err)
	}

	want, _ := filepath.Abs("/tmp/out/custom.pdf")
	if got != want {
		t.Fatalf("ResolveOutputPath = %q, want %q", got, want)
	}
}
