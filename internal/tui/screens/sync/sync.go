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

// planReadyMsg carries the plan computed during the awaitingConfirm phase.
type planReadyMsg struct {
	states  []pkgsync.ComponentState
	actions []pkgsync.SyncAction
	err     error
}

// SyncModel ejecuta sync.Sync con opts.Yes=true y transmite la salida línea a
// línea hacia la vista TUI. Al completar: si el binario fue reemplazado intenta
// re-exec; si no, vuelve al menú automáticamente.
//
// Cuando se entra manualmente (desde el menú), la pantalla primero muestra el
// plan y espera confirmación (awaitingConfirm=true). "y"/"enter" ejecuta;
// "n"/"esc" regresa al menú sin ejecutar.
type SyncModel struct {
	opts            pkgsync.Options
	lines           []string
	done            bool
	err             error
	awaitingConfirm bool
	planActions     []pkgsync.SyncAction
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
	// Mostrar el plan y esperar confirmación antes de ejecutar.
	return m.detectPlanCmd()
}

// detectPlanCmd resuelve el plan de acciones (usando PreDetected si ya están
// disponibles) y lo emite como planReadyMsg para que Update lo presente al usuario.
func (m SyncModel) detectPlanCmd() tea.Cmd {
	return func() tea.Msg {
		opts := m.opts
		var states []pkgsync.ComponentState
		if len(opts.PreDetected) > 0 {
			states = opts.PreDetected
		} else {
			var err error
			states, err = pkgsync.Detect(opts)
			if err != nil {
				return planReadyMsg{err: err}
			}
		}
		actions := pkgsync.PlanSync(states)
		active := make([]pkgsync.SyncAction, 0, len(actions))
		for _, a := range actions {
			if a.Kind != pkgsync.ActionSkip {
				active = append(active, a)
			}
		}
		return planReadyMsg{states: states, actions: active}
	}
}

func (m SyncModel) Update(msg tea.Msg) (nav.ChildModel, tea.Cmd) {
	switch msg := msg.(type) {
	case planReadyMsg:
		if msg.err != nil {
			m.done = true
			m.err = msg.err
			return m, nil
		}
		// Inyectar los estados detectados en las opts para que startSync no llame
		// a Detect de nuevo.
		m.opts.PreDetected = msg.states
		m.planActions = msg.actions
		m.awaitingConfirm = true
		return m, nil

	case tea.KeyMsg:
		if m.awaitingConfirm {
			switch msg.String() {
			case "y", "enter":
				m.awaitingConfirm = false
				return m, m.startSync()
			case "n", "esc":
				return m, func() tea.Msg { return nav.BackMsg{} }
			}
			return m, nil
		}

		switch msg.String() {
		case "esc", "backspace", "b":
			// La goroutine queda huérfana si el sync sigue en curso
			// (mismo comportamiento que install.go — no es interrumpible sin context).
			return m, func() tea.Msg { return nav.BackMsg{} }
		}

	case outputLineMsg:
		m.lines = append(m.lines, msg.line)
		return m, waitForLine(msg)

	case doneMsg:
		m.done = true
		m.err = msg.err
		if msg.result.SelfReplaced {
			m.selfReplaced = true
			// No re-ejecutar inline: syscall.Exec reemplazaría el proceso ahora,
			// antes de que bubbletea restaure la terminal, dejándola rota. Pedir
			// salir vía ReExecMsg; tui.Run hace el ReExec después de p.Run().
			return m, func() tea.Msg { return nav.ReExecMsg{} }
		}
		// Sin self-update: volver al menú automáticamente.
		return m, func() tea.Msg { return nav.BackMsg{} }
	}
	return m, nil
}

func (m SyncModel) View() string {
	var sb strings.Builder

	if m.awaitingConfirm {
		sb.WriteString(styles.Title.Render("Actualizar") + "\n\n")
		if len(m.planActions) == 0 {
			sb.WriteString(styles.Dimmed.Render("  Nada para hacer — todo está instalado y actualizado.") + "\n")
		} else {
			sb.WriteString("  Plan:\n")
			for i, a := range m.planActions {
				verb := "instalar"
				if a.Kind == pkgsync.ActionUpdate {
					verb = "actualizar"
				}
				sb.WriteString(fmt.Sprintf("    %d. %s — %s\n", i+1, a.Component, verb))
			}
			sb.WriteString("\n" + styles.Warn.Render("  ¿Ejecutar? [y/n]") + "\n")
		}
		sb.WriteString("\n" + styles.Footer.Render("y/enter confirmar  n/esc cancelar"))
		return sb.String()
	}

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
