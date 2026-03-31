# Zed Integration

`md2pdf` works well from a Zed task.

## Global task

Create `~/.config/zed/tasks.json` and add:

```json
[
  {
    "label": "Export $ZED_FILENAME to PDF",
    "command": "case \"$ZED_FILE\" in (*.md|*.markdown) ;; (*) echo \"Current file is not Markdown: $ZED_FILE\"; exit 1 ;; esac; pdf_path=\"$ZED_DIRNAME/$ZED_STEM.pdf\"; \"$HOME/src/md2pdf/md2pdf\" \"$ZED_FILE\" && open -a Preview \"$pdf_path\" && osascript -e 'tell application \"Preview\" to activate'",
    "reveal": "always",
    "use_new_terminal": false
  }
]
```

If your binary lives somewhere else, change `"$HOME/src/md2pdf/md2pdf"` accordingly.

## Usage

1. Open a Markdown file in Zed.
2. Open the command palette with `Cmd+Shift+P`.
3. Run `task: spawn`.
4. Choose `Export <current file> to PDF`.

The task:

- exports the active Markdown file
- writes the PDF next to it
- opens the PDF in Preview
- brings Preview to the front on macOS

## Notes

- The task label includes `$ZED_FILENAME` so the picker shows which file Zed thinks is active.
- Zed task variables are documented here: https://zed.dev/docs/tasks

