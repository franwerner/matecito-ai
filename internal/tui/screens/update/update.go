package update

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/franwerner/matecito-ai/internal/setup/install"
	"github.com/franwerner/matecito-ai/internal/tui/nav"
	"github.com/franwerner/matecito-ai/internal/tui/styles"
)

type outputLineMsg struct{ line string }
type doneMsg struct{ err error }

// UpdateModel ejecuta la misma secuencia de actualización que internal/cli/update.go:
// InstallSelf, InstallEngram, UpdateEngramPlugin, InstallCodegraph.
// Usa io.Pipe con opts.Yes=true — idéntico al patrón de InstallModel (S3.2/R3.5).
type UpdateModel struct {
	lines []string
	done  bool
	err   error
}

func New() UpdateModel { return UpdateModel{} }

func (m UpdateModel) Init() tea.Cmd {
	return startUpdate
}

func (m UpdateModel) Update(msg tea.Msg) (nav.ChildModel, tea.Cmd) {
	switch msg := msg.(type) {
	case outputLineMsg:
		m.lines = append(m.lines, msg.line)
		return m, waitForLine(msg)
	case doneMsg:
		m.done = true
		m.err = msg.err
	case tea.KeyMsg:
		// esc/back vuelve al menú en cualquier momento; si el update sigue en
		// curso, su goroutine queda huérfano (no es interrumpible sin context).
		switch msg.String() {
		case "esc", "backspace", "b":
			return m, func() tea.Msg { return nav.BackMsg{} }
		}
	}
	return m, nil
}

func (m UpdateModel) View() string {
	var sb strings.Builder

	sb.WriteString(styles.Title.Render("Update") + "\n\n")

	visible := m.lines
	if len(visible) > 20 {
		visible = visible[len(visible)-20:]
	}
	for _, l := range visible {
		sb.WriteString("  " + l + "\n")
	}

	if !m.done {
		sb.WriteString("\n" + styles.Dimmed.Render("  actualizando…") + "\n")
		sb.WriteString("\n" + styles.Footer.Render("esc volver  ctrl+c salir"))
		return sb.String()
	}

	if m.err != nil {
		sb.WriteString("\n" + styles.Dimmed.Render("  error: "+m.err.Error()) + "\n")
	} else {
		sb.WriteString("\n" + styles.Dimmed.Render("  completado.") + "\n")
	}
	sb.WriteString("\n" + styles.Footer.Render("esc volver"))
	return sb.String()
}

var (
	activeLines chan string
	activeDone  chan doneMsg
)

// runUpdate ejecuta la misma secuencia que cli/update.go: cuatro tareas en orden.
func runUpdate(opts install.Options) error {
	tasks := []struct {
		name string
		fn   func(install.Options) error
	}{
		{"matecito-ai", install.InstallSelf},
		{"engram (binario)", install.InstallEngram},
		{"engram (plugin)", install.UpdateEngramPlugin},
		{"codegraph", install.InstallCodegraph},
	}

	var failed []string
	for _, t := range tasks {
		fmt.Fprintf(opts.Stdout, "Actualizando %s…\n", t.name)
		if err := t.fn(opts); err != nil {
			fmt.Fprintf(opts.Stderr, "  ✗ %s: %v\n", t.name, err)
			failed = append(failed, t.name)
		}
	}

	fmt.Fprintln(opts.Stdout, "context7 corre con npx @latest en cada sesión — nada que actualizar.")

	if len(failed) > 0 {
		return fmt.Errorf("falló la actualización de: %s", strings.Join(failed, ", "))
	}
	fmt.Fprintln(opts.Stdout, "Listo.")
	return nil
}

func startUpdate() tea.Msg {
	pr, pw := io.Pipe()

	opts := install.Options{
		Yes:    true,
		Stdout: pw,
		Stderr: pw,
	}

	lines := make(chan string, 64)
	done := make(chan doneMsg, 1)

	go func() {
		err := runUpdate(opts)
		pw.Close()
		done <- doneMsg{err: err}
	}()

	go func() {
		sc := bufio.NewScanner(pr)
		for sc.Scan() {
			lines <- sc.Text()
		}
		close(lines)
	}()

	return waitForLineFromChannels(lines, done)
}

func waitForLine(_ outputLineMsg) tea.Cmd {
	return func() tea.Msg {
		return waitForLineFromChannels(activeLines, activeDone)
	}
}

func waitForLineFromChannels(lines chan string, done chan doneMsg) tea.Msg {
	activeLines = lines
	activeDone = done

	select {
	case line, ok := <-lines:
		if ok {
			return outputLineMsg{line: line}
		}
		msg := <-done
		return msg
	case msg := <-done:
		for {
			select {
			case _, ok := <-lines:
				if !ok {
					return msg
				}
			default:
				return msg
			}
		}
	}
}
