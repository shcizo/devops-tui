package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/samuelenocsson/devops-tui/internal/models"
	"github.com/samuelenocsson/devops-tui/internal/ui/theme"
)

// AssignModal is a modal for assigning work items to users
type AssignModal struct {
	visible       bool
	item          *models.WorkItem
	members       []models.TeamMember
	filtered      []models.TeamMember
	cursor        int
	styles        theme.Styles
	keys          theme.KeyMap
	width         int
	height        int
	filterInput   textinput.Model
	filterEnabled bool
}

// NewAssignModal creates a new assign modal
func NewAssignModal(styles theme.Styles, keys theme.KeyMap) AssignModal {
	ti := textinput.New()
	ti.Placeholder = "Type to filter..."
	ti.CharLimit = 50
	ti.Width = 30

	return AssignModal{
		styles:      styles,
		keys:        keys,
		filterInput: ti,
	}
}

// Init initializes the modal
func (m AssignModal) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m AssignModal) Update(msg tea.Msg) (AssignModal, tea.Cmd) {
	if !m.visible {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle filter input
		if m.filterEnabled {
			switch msg.String() {
			case "esc":
				if m.filterInput.Value() != "" {
					m.filterInput.SetValue("")
					m.applyFilter()
					return m, nil
				}
				m.visible = false
				return m, func() tea.Msg { return ModalClosedMsg{} }
			case "enter":
				if len(m.filtered) > 0 && m.item != nil {
					selected := m.filtered[m.cursor]
					return m, func() tea.Msg {
						return AssignRequestMsg{
							Item:      *m.item,
							UserEmail: selected.UniqueName,
							UserName:  selected.DisplayName,
						}
					}
				}
			case "up":
				if m.cursor > 0 {
					m.cursor--
				}
				return m, nil
			case "down":
				if m.cursor < len(m.filtered)-1 {
					m.cursor++
				}
				return m, nil
			default:
				var cmd tea.Cmd
				m.filterInput, cmd = m.filterInput.Update(msg)
				m.applyFilter()
				return m, cmd
			}
			return m, nil
		}

		switch {
		case key.Matches(msg, m.keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, m.keys.Down):
			if m.cursor < len(m.filtered)-1 {
				m.cursor++
			}
		case key.Matches(msg, m.keys.Select):
			if m.item != nil && m.cursor < len(m.filtered) {
				selected := m.filtered[m.cursor]
				return m, func() tea.Msg {
					return AssignRequestMsg{
						Item:      *m.item,
						UserEmail: selected.UniqueName,
						UserName:  selected.DisplayName,
					}
				}
			}
		case key.Matches(msg, m.keys.Back):
			m.visible = false
			return m, func() tea.Msg { return ModalClosedMsg{} }
		case msg.String() == "/":
			m.filterEnabled = true
			m.filterInput.Focus()
			return m, textinput.Blink
		}
	}

	return m, nil
}

func (m *AssignModal) applyFilter() {
	filter := strings.ToLower(m.filterInput.Value())
	if filter == "" {
		m.filtered = m.members
	} else {
		m.filtered = make([]models.TeamMember, 0)
		for _, member := range m.members {
			if strings.Contains(strings.ToLower(member.DisplayName), filter) ||
				strings.Contains(strings.ToLower(member.UniqueName), filter) {
				m.filtered = append(m.filtered, member)
			}
		}
	}
	// Reset cursor if out of bounds
	if m.cursor >= len(m.filtered) {
		m.cursor = 0
	}
}

// View renders the modal
func (m AssignModal) View() string {
	if !m.visible {
		return ""
	}

	// Modal dimensions
	modalWidth := 50
	visibleItems := 8
	modalHeight := visibleItems + 10

	// Build content
	var b strings.Builder

	// Title
	title := lipgloss.NewStyle().Bold(true).Render("Assign To")
	if m.item != nil {
		itemInfo := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9CA3AF")).
			Render("#" + itoa(m.item.ID) + " " + truncateStr(m.item.Title, 35))
		b.WriteString(title + "\n")
		b.WriteString(itemInfo + "\n\n")
	}

	// Current assignee
	if m.item != nil {
		currentStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#60A5FA"))
		assigned := m.item.AssignedTo
		if assigned == "" {
			assigned = "Unassigned"
		}
		b.WriteString(currentStyle.Render("Current: "+assigned) + "\n\n")
	}

	// Filter input
	if m.filterEnabled {
		b.WriteString(m.filterInput.View() + "\n\n")
	} else {
		filterHint := lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280")).Render("Press / to filter")
		b.WriteString(filterHint + "\n\n")
	}

	// Unassign option first
	unassignCursor := "  "
	if m.cursor == 0 && len(m.filtered) == 0 {
		unassignCursor = "▸ "
	}

	// Member options
	if len(m.filtered) == 0 {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280")).Render("  No members found") + "\n")
	} else {
		// Calculate scroll offset
		offset := 0
		if m.cursor >= visibleItems {
			offset = m.cursor - visibleItems + 1
		}

		end := offset + visibleItems
		if end > len(m.filtered) {
			end = len(m.filtered)
		}

		for i := offset; i < end; i++ {
			member := m.filtered[i]
			cursor := "  "
			if i == m.cursor {
				cursor = "▸ "
			}

			style := lipgloss.NewStyle()
			if i == m.cursor {
				style = style.Bold(true).Foreground(lipgloss.Color("#7C3AED"))
			}

			// Highlight if this is the current assignee
			if m.item != nil && member.DisplayName == m.item.AssignedTo {
				style = style.Foreground(lipgloss.Color("#10B981"))
			}

			// Show name, truncate if needed
			name := truncateStr(member.DisplayName, modalWidth-10)
			b.WriteString(cursor + style.Render(name) + "\n")
		}

		// Show scroll indicator
		if len(m.filtered) > visibleItems {
			scrollInfo := lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280")).
				Render("  (" + itoa(m.cursor+1) + "/" + itoa(len(m.filtered)) + ")")
			b.WriteString(scrollInfo + "\n")
		}
	}

	_ = unassignCursor // Reserved for future unassign option

	// Help text
	b.WriteString("\n")
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280"))
	if m.filterEnabled {
		b.WriteString(helpStyle.Render("Enter: confirm  Esc: clear/close"))
	} else {
		b.WriteString(helpStyle.Render("Enter: confirm  /: filter  Esc: cancel"))
	}

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
func (m *AssignModal) SetVisible(visible bool) {
	m.visible = visible
	if visible {
		m.cursor = 0
		m.filterInput.SetValue("")
		m.filterEnabled = false
		m.filterInput.Blur()
		m.applyFilter()
		// Try to set cursor to current assignee
		if m.item != nil && m.item.AssignedTo != "" {
			for i, member := range m.filtered {
				if member.DisplayName == m.item.AssignedTo {
					m.cursor = i
					break
				}
			}
		}
	}
}

// IsVisible returns whether the modal is visible
func (m *AssignModal) IsVisible() bool {
	return m.visible
}

// SetItem sets the work item to modify
func (m *AssignModal) SetItem(item *models.WorkItem) {
	m.item = item
}

// SetMembers sets the available team members
func (m *AssignModal) SetMembers(members []models.TeamMember) {
	m.members = members
	m.applyFilter()
}

// SetSize sets the modal container size
func (m *AssignModal) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// AssignRequestMsg is sent when user confirms assignment
type AssignRequestMsg struct {
	Item      models.WorkItem
	UserEmail string
	UserName  string
}

// AssignedMsg is sent when assignment is complete
type AssignedMsg struct {
	Item models.WorkItem
}
