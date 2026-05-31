package sync

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	pkgsync "github.com/franwerner/matecito-ai/internal/setup/sync"
	"github.com/franwerner/matecito-ai/internal/tui/nav"
	"github.com/franwerner/matecito-ai/internal/tui/styles"
)

type outputLineMsg struct{ line string }

type doneMsg struct {
	result pkgsync.Result
	err    error
}

// SyncModel ejecuta sync.Sync con opts.Yes=true y transmite la salida línea a
// línea hacia la vista TUI. Al completar: si el binario fue reemplazado intenta
// re-exec; si no, vuelve al menú automáticamente.
type SyncModel struct {
	opts  pkgsync.Options
	lines []string
	done  bool
	err   error
	// selfReplaced indica que el binario fue actualizado y re-exec fue intentado.
	// En Unix el proceso no debería llegar aquí (syscall.Exec reemplaza el proceso);
	// en Windows o ante un error de exec, se muestra un aviso en pantalla.
	selfReplaced bool
}

// New construye un SyncModel con las opciones de sincronización inyectadas
// desde AppModel (SelfVersion, Timeout; Stdout/Stderr se asignan internamente).
func New(opts pkgsync.Options) SyncModel {
	return SyncModel{opts: opts}
}

func (m SyncModel) Init() tea.Cmd {
	return m.startSync()
}

func (m SyncModel) Update(msg tea.Msg) (nav.ChildModel, tea.Cmd) {
	switch msg := msg.(type) {
	case outputLineMsg:
		m.lines = append(m.lines, msg.line)
		return m, waitForLine(msg)

	case doneMsg:
		m.done = true
		m.err = msg.err
		if msg.result.SelfReplaced {
			m.selfReplaced = true
			// En Unix syscall.Exec reemplaza el proceso; si llega acá fue un error.
			// En Windows ReExec() devuelve nil sin hacer exec → aviso en View.
			_ = pkgsync.ReExec()
			return m, nil
		}
		// Sin self-update: volver al menú automáticamente.
		return m, func() tea.Msg { return nav.BackMsg{} }

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "backspace", "b":
			// La goroutine queda huérfana si el sync sigue en curso
			// (mismo comportamiento que install.go — no es interrumpible sin context).
			return m, func() tea.Msg { return nav.BackMsg{} }
		}
	}
	return m, nil
}

func (m SyncModel) View() string {
	var sb strings.Builder

	sb.WriteString(styles.Title.Render("Sincronizando") + "\n\n")

	visible := m.lines
	if len(visible) > 20 {
		visible = visible[len(visible)-20:]
	}
	for _, l := range visible {
		sb.WriteString("  " + l + "\n")
	}

	if !m.done {
		sb.WriteString("\n" + styles.Dimmed.Render("  ejecutando…") + "\n")
		sb.WriteString("\n" + styles.Footer.Render("esc volver  ctrl+c salir"))
		return sb.String()
	}

	if m.selfReplaced {
		sb.WriteString("\n" + styles.Success.Render("  matecito-ai actualizado.") + "\n")
		sb.WriteString(styles.Dimmed.Render(fmt.Sprintf("  Re-ejecutá %s para usar la nueva versión.", "matecito-ai")) + "\n")
		sb.WriteString("\n" + styles.Footer.Render("esc volver  ctrl+c salir"))
		return sb.String()
	}

	if m.err != nil {
		sb.WriteString("\n" + styles.Error.Render("  error: "+m.err.Error()) + "\n")
	} else {
		sb.WriteString("\n" + styles.Dimmed.Render("  completado.") + "\n")
	}
	sb.WriteString("\n" + styles.Footer.Render("esc volver"))
	return sb.String()
}

// startSync crea el io.Pipe, lanza pkgsync.Sync en una goroutine con streaming
// de salida, y devuelve el primer mensaje del canal de líneas.
func (m SyncModel) startSync() tea.Cmd {
	return func() tea.Msg {
		pr, pw := io.Pipe()

		opts := m.opts
		opts.Yes = true
		opts.Stdout = pw
		opts.Stderr = pw

		lines := make(chan string, 64)
		done := make(chan doneMsg, 1)

		go func() {
			result := pkgsync.Sync(opts)
			pw.Close()
			done <- doneMsg{result: result}
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
}

// waitForLine devuelve un tea.Cmd que espera la siguiente línea o el fin del sync.
func waitForLine(_ outputLineMsg) tea.Cmd {
	return func() tea.Msg {
		return waitForLineFromChannels(activeSyncLines, activeSyncDone)
	}
}

// activeSyncLines y activeSyncDone persisten entre llamadas a waitForLine para
// que el canal no se pierda entre mensajes. Se inicializan en startSync.
var (
	activeSyncLines chan string
	activeSyncDone  chan doneMsg
)

func waitForLineFromChannels(lines chan string, done chan doneMsg) tea.Msg {
	activeSyncLines = lines
	activeSyncDone = done

	select {
	case line, ok := <-lines:
		if ok {
			return outputLineMsg{line: line}
		}
		msg := <-done
		return msg
	case msg := <-done:
		// drenar líneas restantes antes de retornar done
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
