package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/samuelenocsson/devops-tui/internal/models"
	"github.com/samuelenocsson/devops-tui/internal/ui/theme"
)

var defaultStates = []string{"New", "Active", "Resolved", "Closed"}

// StateModal is a modal for changing work item state
type StateModal struct {
	visible      bool
	item         *models.WorkItem
	states       []string
	statesByType map[string][]models.WorkItemStateInfo
	cursor       int
	styles       theme.Styles
	keys         theme.KeyMap
	width        int
	height       int
}

// NewStateModal creates a new state modal
func NewStateModal(styles theme.Styles, keys theme.KeyMap) StateModal {
	return StateModal{
		states: defaultStates,
		styles: styles,
		keys:   keys,
	}
}

// SetStatesByType sets the available states per work item type
func (m *StateModal) SetStatesByType(statesByType map[string][]models.WorkItemStateInfo) {
	m.statesByType = statesByType
}

// Init initializes the modal
func (m StateModal) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m StateModal) Update(msg tea.Msg) (StateModal, tea.Cmd) {
	if !m.visible {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, m.keys.Down):
			if m.cursor < len(m.states)-1 {
				m.cursor++
			}
		case key.Matches(msg, m.keys.Select):
			if m.item != nil {
				newState := m.states[m.cursor]
				return m, func() tea.Msg {
					return StateChangeRequestMsg{
						Item:     *m.item,
						NewState: newState,
					}
				}
			}
		case key.Matches(msg, m.keys.Back):
			m.visible = false
			return m, func() tea.Msg { return ModalClosedMsg{} }
		}
	}

	return m, nil
}

// View renders the modal
func (m StateModal) View() string {
	if !m.visible {
		return ""
	}

	// Modal dimensions
	modalWidth := 40
	modalHeight := len(m.states) + 6

	// Build content
	var b strings.Builder

	// Title
	title := "Change State"
	if m.item != nil {
		title = lipgloss.NewStyle().Bold(true).Render("Change State")
		itemInfo := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9CA3AF")).
			Render("#" + itoa(m.item.ID) + " " + truncateStr(m.item.Title, 25))
		b.WriteString(title + "\n")
		b.WriteString(itemInfo + "\n\n")
	}

	// Current state indicator
	if m.item != nil {
		currentStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#60A5FA"))
		b.WriteString(currentStyle.Render("Current: "+string(m.item.State)) + "\n\n")
	}

	// State options
	for i, state := range m.states {
		cursor := "  "
		if i == m.cursor {
			cursor = "▸ "
		}

		style := lipgloss.NewStyle()
		if i == m.cursor {
			style = style.Bold(true).Foreground(lipgloss.Color("#7C3AED"))
		}

		// Highlight if this is the current state
		if m.item != nil && state == string(m.item.State) {
			style = style.Foreground(lipgloss.Color("#10B981"))
		}

		b.WriteString(cursor + style.Render(state) + "\n")
	}

	// Help text
	b.WriteString("\n")
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280"))
	b.WriteString(helpStyle.Render("Enter: confirm  Esc: cancel"))

	// Modal style
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7C3AED")).
		Padding(1, 2).
		Width(modalWidth).
		Height(modalHeight).
		Background(lipgloss.Color("#1F2937"))

	modal := modalStyle.Render(b.String())

	// Center the modal
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, modal)
}

// SetVisible sets the visibility
func (m *StateModal) SetVisible(visible bool) {
	m.visible = visible
	if visible {
		m.cursor = 0
		// Try to set cursor to current state
		if m.item != nil {
			for i, state := range m.states {
				if state == string(m.item.State) {
					m.cursor = i
					break
				}
			}
		}
	}
}

// IsVisible returns whether the modal is visible
func (m *StateModal) IsVisible() bool {
	return m.visible
}

// SetItem sets the work item to modify
func (m *StateModal) SetItem(item *models.WorkItem) {
	m.item = item
	// Update states based on work item type
	if item != nil && m.statesByType != nil {
		if states, ok := m.statesByType[string(item.Type)]; ok && len(states) > 0 {
			m.states = make([]string, len(states))
			for i, s := range states {
				m.states[i] = s.Name
			}
			return
		}
	}
	// Fallback to defaults
	m.states = defaultStates
}

// SetSize sets the modal container size
func (m *StateModal) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// StateChangeRequestMsg is sent when user confirms state change
type StateChangeRequestMsg struct {
	Item     models.WorkItem
	NewState string
}

// StateChangedMsg is sent when state change is complete
type StateChangedMsg struct {
	Item models.WorkItem
}

// ModalClosedMsg is sent when a modal is closed
type ModalClosedMsg struct{}

// Helper function
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}
