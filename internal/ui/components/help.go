package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
	"github.com/samuelenocsson/devops-tui/internal/ui/theme"
)

// HelpPanel displays keyboard shortcuts
type HelpPanel struct {
	keys    theme.KeyMap
	styles  theme.Styles
	visible bool
	width   int
	height  int
}

// NewHelpPanel creates a new help panel
func NewHelpPanel(keys theme.KeyMap, styles theme.Styles) HelpPanel {
	return HelpPanel{
		keys:   keys,
		styles: styles,
	}
}

// View renders the help panel
func (h HelpPanel) View() string {
	if !h.visible {
		return ""
	}

	var b strings.Builder

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#F9FAFB")).
		MarginBottom(1).
		Render("Keyboard Shortcuts")

	b.WriteString(title)
	b.WriteString("\n\n")

	// Group shortcuts by category
	sections := []struct {
		title    string
		bindings []key.Binding
	}{
		{
			title: "Navigation",
			bindings: []key.Binding{
				h.keys.Up,
				h.keys.Down,
				h.keys.Top,
				h.keys.Bottom,
				h.keys.NextPanel,
				h.keys.PrevPanel,
			},
		},
		{
			title: "Actions",
			bindings: []key.Binding{
				h.keys.Select,
				h.keys.Open,
				h.keys.View,
				h.keys.Search,
				h.keys.Refresh,
			},
		},
		{
			title: "General",
			bindings: []key.Binding{
				h.keys.Help,
				h.keys.Back,
				h.keys.Quit,
			},
		},
	}

	for i, section := range sections {
		// Section title
		sectionTitle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7C3AED")).
			Render(section.title)
		b.WriteString(sectionTitle)
		b.WriteString("\n")

		// Key bindings
		for _, binding := range section.bindings {
			keyStyle := h.styles.HelpKey.Width(12)
			descStyle := h.styles.HelpDesc

			help := binding.Help()
			line := keyStyle.Render(help.Key) + descStyle.Render(help.Desc)
			b.WriteString(line)
			b.WriteString("\n")
		}

		// Add spacing between sections
		if i < len(sections)-1 {
			b.WriteString("\n")
		}
	}

	// Footer
	b.WriteString("\n")
	footer := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280")).
		Italic(true).
		Render("Press ? or Esc to close")
	b.WriteString(footer)

	content := b.String()

	// Center the help panel
	helpWidth := 40
	helpHeight := 25

	panel := h.styles.HelpPanel.
		Width(helpWidth).
		Height(helpHeight).
		Background(lipgloss.Color("#1F2937")).
		Render(content)

	// Create overlay positioning
	return lipgloss.Place(
		h.width,
		h.height,
		lipgloss.Center,
		lipgloss.Center,
		panel,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("#000000")),
	)
}

// SetVisible sets whether the help panel is visible
func (h *HelpPanel) SetVisible(visible bool) {
	h.visible = visible
}

// IsVisible returns whether the help panel is visible
func (h *HelpPanel) IsVisible() bool {
	return h.visible
}

// Toggle toggles the visibility of the help panel
func (h *HelpPanel) Toggle() {
	h.visible = !h.visible
}

// SetSize sets the size of the help panel
func (h *HelpPanel) SetSize(width, height int) {
	h.width = width
	h.height = height
}

// ShortHelp returns a short help string for the status bar
func ShortHelp(keys theme.KeyMap, styles theme.Styles) string {
	bindings := keys.ShortHelp()
	var parts []string

	for _, b := range bindings {
		help := b.Help()
		key := styles.HelpKey.Render(help.Key)
		desc := styles.HelpDesc.Render(help.Desc)
		parts = append(parts, key+" "+desc)
	}

	return strings.Join(parts, "  ")
}
