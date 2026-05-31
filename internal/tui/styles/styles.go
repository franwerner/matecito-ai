package styles

import "github.com/charmbracelet/lipgloss"

// Paleta de marca matecito-ai (docs/visual/tokens.css). Se usan como ACENTOS
// de foreground; el texto de cuerpo se deja en el color por defecto de la
// terminal para mantener legibilidad tanto en temas claros como oscuros.
var (
	calabash = lipgloss.Color("#B26A3C") // primario de marca (calabash)
	yerba    = lipgloss.Color("#6B8E4E") // verde yerba (OK / éxito)
	accent   = lipgloss.Color("#8B7BFF") // acento IA (violeta) · sufijo "-ai"
	inkMute  = lipgloss.Color("#6F5A4A") // secundario
	mateBord = lipgloss.Color("#E8DCC8") // bordes / divisores
	cheek    = lipgloss.Color("#F2A6B0") // warnings (cachete)
	danger   = lipgloss.Color("#C0392B") // error (la paleta de marca no define uno)
)

const mateGlyph = "🧉"

var (
	// Title es el encabezado principal de cada pantalla (calabash).
	Title = lipgloss.NewStyle().Bold(true).Foreground(calabash)

	// Selected resalta el ítem activo (acento violeta de marca).
	Selected = lipgloss.NewStyle().Bold(true).Foreground(accent)

	// Dimmed: texto secundario. Faint degrada bien en cualquier terminal.
	Dimmed = lipgloss.NewStyle().Faint(true)

	// Footer: leyenda de atajos al pie.
	Footer = lipgloss.NewStyle().Faint(true).Foreground(inkMute)

	// Accent aplica el violeta de marca (sufijo "-ai", valores de scope).
	Accent = lipgloss.NewStyle().Foreground(accent)

	// Success colorea estados OK (verde yerba).
	Success = lipgloss.NewStyle().Foreground(yerba)

	// Warn usa el rosa cachete para advertencias no críticas.
	Warn = lipgloss.NewStyle().Foreground(cheek)

	// Error resalta fallos (rojo; la marca no tiene color de error).
	Error = lipgloss.NewStyle().Bold(true).Foreground(danger)

	// Border colorea bordes/divisores con el tono de marca.
	Border = lipgloss.NewStyle().Foreground(mateBord)

	// wordmarkBase deja "matecito" en el fg por defecto (legible light/dark).
	wordmarkBase = lipgloss.NewStyle().Bold(true)
)

// Wordmark renderiza el logotipo "🧉 matecito-ai" con el sufijo "-ai" en violeta,
// respetando la regla de marca (el sufijo "-ai" siempre en el acento).
func Wordmark() string {
	return mateGlyph + " " + wordmarkBase.Render("matecito") + Accent.Bold(true).Render("-ai")
}
