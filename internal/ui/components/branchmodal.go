package components

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/samuelenocsson/devops-tui/internal/models"
	"github.com/samuelenocsson/devops-tui/internal/ui/theme"
)

// BranchModal is a modal for creating a branch linked to a work item
type BranchModal struct {
	visible   bool
	item      *models.WorkItem
	textInput textinput.Model
	styles    theme.Styles
	keys      theme.KeyMap
	width     int
	height    int
	err       error
}

// NewBranchModal creates a new branch modal
func NewBranchModal(styles theme.Styles, keys theme.KeyMap) BranchModal {
	ti := textinput.New()
	ti.Placeholder = "feature/123-task-name"
	ti.CharLimit = 100
	ti.Width = 35

	return BranchModal{
		textInput: ti,
		styles:    styles,
		keys:      keys,
	}
}

// Init initializes the modal
func (m BranchModal) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m BranchModal) Update(msg tea.Msg) (BranchModal, tea.Cmd) {
	if !m.visible {
		return m, nil
	}

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Back):
			m.visible = false
			m.err = nil
			return m, func() tea.Msg { return ModalClosedMsg{} }
		case msg.Type == tea.KeyEnter:
			branchName := strings.TrimSpace(m.textInput.Value())
			if branchName == "" {
				m.err = fmt.Errorf("branch name cannot be empty")
				return m, nil
			}
			if !isValidBranchName(branchName) {
				m.err = fmt.Errorf("invalid branch name")
				return m, nil
			}
			m.err = nil
			return m, func() tea.Msg {
				return BranchCreateRequestMsg{
					Item:       *m.item,
					BranchName: branchName,
				}
			}
		default:
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// View renders the modal
func (m BranchModal) View() string {
	if !m.visible {
		return ""
	}

	// Modal dimensions
	modalWidth := 50
	modalHeight := 10

	// Build content
	var b strings.Builder

	// Title
	title := lipgloss.NewStyle().Bold(true).Render("Create Branch")
	b.WriteString(title + "\n")

	// Item info
	if m.item != nil {
		itemInfo := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9CA3AF")).
			Render("#" + itoa(m.item.ID) + " " + truncateStr(m.item.Title, 35))
		b.WriteString(itemInfo + "\n\n")
	}

	// Label
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#D1D5DB"))
	b.WriteString(labelStyle.Render("Branch name:") + "\n")

	// Text input
	b.WriteString(m.textInput.View() + "\n")

	// Error message
	if m.err != nil {
		errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#EF4444"))
		b.WriteString(errStyle.Render(m.err.Error()) + "\n")
	} else {
		b.WriteString("\n")
	}

	// Help text
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280"))
	b.WriteString(helpStyle.Render("Enter: create  Esc: cancel"))

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
func (m *BranchModal) SetVisible(visible bool) {
	m.visible = visible
	m.err = nil
	if visible {
		m.textInput.Focus()
		// Generate suggested branch name from work item
		if m.item != nil {
			suggested := generateBranchName(m.item)
			m.textInput.SetValue(suggested)
			m.textInput.CursorEnd()
		}
	} else {
		m.textInput.Blur()
	}
}

// IsVisible returns whether the modal is visible
func (m *BranchModal) IsVisible() bool {
	return m.visible
}

// SetItem sets the work item to link
func (m *BranchModal) SetItem(item *models.WorkItem) {
	m.item = item
}

// SetSize sets the modal container size
func (m *BranchModal) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// generateBranchName generates a suggested branch name from work item
func generateBranchName(item *models.WorkItem) string {
	// Get work item type prefix
	prefix := "feature"
	switch item.Type {
	case models.WorkItemTypeBug:
		prefix = "bugfix"
	case models.WorkItemTypeTask:
		prefix = "task"
	case models.WorkItemTypeStory:
		prefix = "feature"
	case models.WorkItemTypeEpic:
		prefix = "epic"
	}

	// Sanitize title for branch name
	title := strings.ToLower(item.Title)
	// Replace spaces and special chars with hyphens
	re := regexp.MustCompile(`[^a-z0-9]+`)
	title = re.ReplaceAllString(title, "-")
	// Remove leading/trailing hyphens
	title = strings.Trim(title, "-")
	// Truncate to reasonable length
	if len(title) > 40 {
		title = title[:40]
		// Don't cut off in middle of a word if possible
		if lastHyphen := strings.LastIndex(title, "-"); lastHyphen > 20 {
			title = title[:lastHyphen]
		}
	}

	return fmt.Sprintf("%s/%d-%s", prefix, item.ID, title)
}

// isValidBranchName validates branch name
func isValidBranchName(name string) bool {
	// Basic validation for git branch names
	if name == "" {
		return false
	}
	// Check for invalid characters
	invalid := regexp.MustCompile(`[\s~^:?*\[\]\\]`)
	if invalid.MatchString(name) {
		return false
	}
	// Can't start or end with /
	if strings.HasPrefix(name, "/") || strings.HasSuffix(name, "/") {
		return false
	}
	// Can't have consecutive dots
	if strings.Contains(name, "..") {
		return false
	}
	return true
}

// BranchCreateRequestMsg is sent when user confirms branch creation
type BranchCreateRequestMsg struct {
	Item       models.WorkItem
	BranchName string
}

// BranchCreatedMsg is sent when branch creation is complete
type BranchCreatedMsg struct {
	BranchName string
}

// BranchCreateErrorMsg is sent when branch creation fails
type BranchCreateErrorMsg struct {
	Err error
}
