package render

import (
	"fmt"
	"io"
	"os"

	"github.com/franwerner/matecito-ai/internal/check"
)

const (
	colorReset  = "\033[0m"
	colorBold   = "\033[1m"
	colorGreen  = "\033[32m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
)

func Section(w io.Writer, title string, results []check.Result) {
	useColor := isTTY(w)
	fmt.Fprintln(w, paint(title, colorBold, useColor))
	for _, r := range results {
		fmt.Fprintln(w, renderLine(r, useColor))
	}
	fmt.Fprintln(w)
}

func Suggestions(w io.Writer, results []check.Result) {
	useColor := isTTY(w)
	seen := map[string]bool{}
	hints := make([]string, 0, len(results))
	for _, r := range results {
		if r.FixHint == "" || seen[r.FixHint] {
			continue
		}
		seen[r.FixHint] = true
		hints = append(hints, r.FixHint)
	}
	if len(hints) == 0 {
		return
	}
	fmt.Fprintln(w, paint("Comandos sugeridos", colorBold, useColor))
	for _, h := range hints {
		fmt.Fprintln(w, "  • "+h)
	}
	fmt.Fprintln(w)
}

func Summary(w io.Writer, results []check.Result) int {
	useColor := isTTY(w)
	missing := 0
	for _, r := range results {
		if r.Required && r.Status != check.StatusOK {
			missing++
		}
	}
	if missing == 0 {
		fmt.Fprintln(w, "Estado:", paint("OK", colorGreen, useColor))
		return 0
	}
	msg := fmt.Sprintf("faltan %d chequeos críticos", missing)
	fmt.Fprintln(w, "Estado:", paint(msg, colorRed, useColor))
	return 1
}

func renderLine(r check.Result, useColor bool) string {
	mark, color := markFor(r)
	name := fmt.Sprintf("%-14s", r.Name)
	switch r.Status {
	case check.StatusOK:
		extra := r.Version
		if r.Detail != "" {
			if extra != "" {
				extra += "  "
			}
			extra += r.Detail
		}
		return fmt.Sprintf("  %s %s %s", paint(mark, color, useColor), name, extra)
	case check.StatusOutdated:
		line := fmt.Sprintf("  %s %s %s", paint(mark, color, useColor), name, r.Detail)
		if r.FixHint != "" {
			line += "\n                 → " + r.FixHint
		}
		return line
	case check.StatusMissing:
		detail := r.Detail
		if detail == "" {
			detail = "no instalado"
		}
		if !r.Required {
			detail += " (opcional)"
		}
		line := fmt.Sprintf("  %s %s %s", paint(mark, color, useColor), name, detail)
		if r.FixHint != "" {
			line += "\n                 → " + r.FixHint
		}
		return line
	}
	return ""
}

func markFor(r check.Result) (string, string) {
	switch r.Status {
	case check.StatusOK:
		return "✓", colorGreen
	case check.StatusOutdated:
		return "✗", colorRed
	case check.StatusMissing:
		if r.Required {
			return "✗", colorRed
		}
		return "⚠", colorYellow
	}
	return "?", ""
}

func paint(s, color string, useColor bool) string {
	if !useColor || color == "" {
		return s
	}
	return color + s + colorReset
}

func isTTY(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}
