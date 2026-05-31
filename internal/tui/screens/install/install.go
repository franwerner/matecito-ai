package install

import (
	"bufio"
	"io"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/franwerner/matecito-ai/internal/setup/install"
	"github.com/franwerner/matecito-ai/internal/tui/nav"
	"github.com/franwerner/matecito-ai/internal/tui/styles"
)

type outputLineMsg struct{ line string }
type doneMsg struct{ err error }

// InstallModel ejecuta el flujo de install existente con streaming de salida línea a
// línea hacia la vista TUI. opts.Yes=true es obligatorio porque bubbletea posee stdin;
// sin ese flag, el confirm() del install bloquea indefinidamente (S3.2).
type InstallModel struct {
	lines []string
	done  bool
	err   error
}

func New() InstallModel { return InstallModel{} }

func (m InstallModel) Init() tea.Cmd {
	return startInstall
}

func (m InstallModel) Update(msg tea.Msg) (nav.ChildModel, tea.Cmd) {
	switch msg := msg.(type) {
	case outputLineMsg:
		m.lines = append(m.lines, msg.line)
		return m, waitForLine(msg)
	case doneMsg:
		m.done = true
		m.err = msg.err
	case tea.KeyMsg:
		// esc/back vuelve al menú en cualquier momento; si el install sigue en
		// curso, su goroutine queda huérfano (no es interrumpible sin context).
		switch msg.String() {
		case "esc", "backspace", "b":
			return m, func() tea.Msg { return nav.BackMsg{} }
		}
	}
	return m, nil
}

func (m InstallModel) View() string {
	var sb strings.Builder

	sb.WriteString(styles.Title.Render("Install") + "\n\n")

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

	if m.err != nil {
		sb.WriteString("\n" + styles.Dimmed.Render("  error: "+m.err.Error()) + "\n")
	} else {
		sb.WriteString("\n" + styles.Dimmed.Render("  completado.") + "\n")
	}
	sb.WriteString("\n" + styles.Footer.Render("esc volver"))
	return sb.String()
}

// startInstall crea el io.Pipe, lanza install.Run en un goroutine, lanza el
// lector de líneas en otro goroutine, y devuelve el primer mensaje del canal.
func startInstall() tea.Msg {
	pr, pw := io.Pipe()

	opts := install.Options{
		Yes:    true,
		Stdout: pw,
		Stderr: pw,
	}

	lines := make(chan string, 64)
	done := make(chan doneMsg, 1)

	// goroutine: ejecuta install.Run y cierra el pipe al terminar
	go func() {
		err := install.Run(opts)
		pw.Close()
		done <- doneMsg{err: err}
	}()

	// goroutine: lee líneas del pipe y las envía al canal
	go func() {
		sc := bufio.NewScanner(pr)
		for sc.Scan() {
			lines <- sc.Text()
		}
		close(lines)
	}()

	return waitForLineFromChannels(lines, done)
}

// waitForLine bloquea hasta recibir la siguiente línea o el mensaje de fin.
// Recibe el outputLineMsg anterior sólo para acceder al canal mediante el
// patrón de closure; en la práctica usamos variables de nivel de paquete.
func waitForLine(_ outputLineMsg) tea.Cmd {
	return func() tea.Msg {
		return waitForLineFromChannels(activeLines, activeDone)
	}
}

// activeLines y activeDone persisten entre llamadas a waitForLine para que el
// canal no se pierda entre mensajes. Se inicializan en startInstall.
var (
	activeLines chan string
	activeDone  chan doneMsg
)

func waitForLineFromChannels(lines chan string, done chan doneMsg) tea.Msg {
	// guardar los canales activos para el próximo waitForLine
	activeLines = lines
	activeDone = done

	select {
	case line, ok := <-lines:
		if ok {
			return outputLineMsg{line: line}
		}
		// canal cerrado: esperar el doneMsg
		msg := <-done
		return msg
	case msg := <-done:
		// drenar líneas restantes antes de retornar done
		for {
			select {
			case line, ok := <-lines:
				if ok {
					_ = line
				}
			default:
				return msg
			}
		}
	}
}
