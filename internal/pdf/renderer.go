package pdf

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type Renderer struct {
	chromePath string
}

func NewRenderer() (*Renderer, error) {
	chromePath, err := findChrome()
	if err != nil {
		return nil, err
	}

	return &Renderer{chromePath: chromePath}, nil
}

func (r *Renderer) Render(ctx context.Context, html, outputPath string) error {
	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	htmlPath, htmlURL, err := writeTempHTML(html)
	if err != nil {
		return err
	}
	defer os.Remove(htmlPath)

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(ctx, append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(r.chromePath),
		chromedp.Headless,
		chromedp.DisableGPU,
		chromedp.Flag("allow-file-access-from-files", true),
		chromedp.Flag("hide-scrollbars", true),
		chromedp.Flag("mute-audio", true),
	)...)
	defer cancelAlloc()

	taskCtx, cancelTask := chromedp.NewContext(allocCtx)
	defer cancelTask()

	var pdfData []byte
	if err := chromedp.Run(taskCtx,
		chromedp.Navigate(htmlURL),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.ActionFunc(waitForAssets),
		chromedp.Sleep(250*time.Millisecond),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdfData, _, err = page.PrintToPDF().
				WithPrintBackground(true).
				WithPreferCSSPageSize(true).
				Do(ctx)
			return err
		}),
	); err != nil {
		return fmt.Errorf("render with chrome: %w", err)
	}

	if err := os.WriteFile(outputPath, pdfData, 0o644); err != nil {
		return fmt.Errorf("write pdf: %w", err)
	}

	return nil
}

func writeTempHTML(html string) (string, string, error) {
	file, err := os.CreateTemp("", "md2pdf-*.html")
	if err != nil {
		return "", "", fmt.Errorf("create temp html: %w", err)
	}
	defer file.Close()

	if _, err := file.WriteString(html); err != nil {
		return "", "", fmt.Errorf("write temp html: %w", err)
	}

	absPath, err := filepath.Abs(file.Name())
	if err != nil {
		return "", "", fmt.Errorf("resolve temp html path: %w", err)
	}

	return absPath, (&url.URL{Scheme: "file", Path: filepath.ToSlash(absPath)}).String(), nil
}

func waitForAssets(ctx context.Context) error {
	script := `(async () => {
		const images = Array.from(document.images || []);
		await Promise.all(images.map((img) => {
			if (img.complete) return Promise.resolve();
			return new Promise((resolve) => {
				img.addEventListener('load', resolve, { once: true });
				img.addEventListener('error', resolve, { once: true });
			});
		}));
		if (document.fonts && document.fonts.ready) {
			await document.fonts.ready;
		}
	})()`
	return chromedp.Evaluate(script, nil).Do(ctx)
}

func findChrome() (string, error) {
	candidates := []string{
		"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		"/Applications/Google Chrome Canary.app/Contents/MacOS/Google Chrome Canary",
		"/Applications/Chromium.app/Contents/MacOS/Chromium",
		"/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge",
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}

	lookups := []string{"google-chrome", "chrome", "chromium"}
	for _, name := range lookups {
		if path, err := execLookPath(name); err == nil {
			return path, nil
		}
	}

	return "", errors.New("could not find a Chrome/Chromium executable")
}
