package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/samuelenocsson/devops-tui/internal/models"
	"github.com/samuelenocsson/devops-tui/internal/ui/theme"
)

// DetailView is the fullscreen detail view component
type DetailView struct {
	item         *models.WorkItem
	styles       theme.Styles
	keys         theme.KeyMap
	width        int
	height       int
	scrollOffset int
	maxScroll    int
}

// NewDetailView creates a new detail view
func NewDetailView(styles theme.Styles, keys theme.KeyMap) DetailView {
	return DetailView{
		styles: styles,
		keys:   keys,
	}
}

// Init initializes the detail view
func (d DetailView) Init() tea.Cmd {
	return nil
}

// Update handles messages for the detail view
func (d DetailView) Update(msg tea.Msg) (DetailView, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, d.keys.Back):
			return d, func() tea.Msg { return CloseDetailViewMsg{} }
		case key.Matches(msg, d.keys.Quit) && msg.String() == "q":
			return d, func() tea.Msg { return CloseDetailViewMsg{} }
		case key.Matches(msg, d.keys.Open):
			if d.item != nil {
				return d, func() tea.Msg { return OpenWorkItemMsg{Item: *d.item} }
			}
		case key.Matches(msg, d.keys.Up):
			if d.scrollOffset > 0 {
				d.scrollOffset--
			}
		case key.Matches(msg, d.keys.Down):
			if d.scrollOffset < d.maxScroll {
				d.scrollOffset++
			}
		}
	}

	return d, nil
}

// View renders the detail view
func (d DetailView) View() string {
	if d.item == nil {
		return ""
	}

	var sections []string

	// Title bar
	title := fmt.Sprintf("#%d %s", d.item.ID, d.item.Title)
	titleBar := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#F9FAFB")).
		Background(lipgloss.Color("#7C3AED")).
		Padding(0, 1).
		Width(d.width - 2).
		Render(title)
	sections = append(sections, titleBar)

	// Metadata section
	metadataContent := d.renderMetadata()
	metadataSection := d.styles.DetailSection.
		Width(d.width - 6).
		Render("METADATA\n" + metadataContent)
	sections = append(sections, metadataSection)

	// Parent section (if exists)
	if d.item.ParentID > 0 {
		parentContent := fmt.Sprintf("#%d", d.item.ParentID)
		if d.item.ParentTitle != "" {
			parentContent += " " + d.item.ParentTitle
		}
		parentSection := d.styles.DetailSection.
			Width(d.width - 6).
			Render("PARENT\n" + parentContent)
		sections = append(sections, parentSection)
	}

	// Description section
	if d.item.Description != "" {
		desc := wordWrap(d.item.Description, d.width-10)
		descSection := d.styles.DetailSection.
			Width(d.width - 6).
			Render("DESCRIPTION\n" + desc)
		sections = append(sections, descSection)
	}

	// Tags section
	if len(d.item.Tags) > 0 {
		var tagStrings []string
		for _, tag := range d.item.Tags {
			tagStrings = append(tagStrings, d.styles.DetailTag.Render(tag))
		}
		tagsSection := d.styles.DetailSection.
			Width(d.width - 6).
			Render("TAGS\n" + strings.Join(tagStrings, " "))
		sections = append(sections, tagsSection)
	}

	// Join all sections
	content := strings.Join(sections, "\n\n")

	// Calculate scrolling
	contentLines := strings.Split(content, "\n")
	viewableHeight := d.height - 4
	d.maxScroll = len(contentLines) - viewableHeight
	if d.maxScroll < 0 {
		d.maxScroll = 0
	}

	// Apply scrolling
	if d.scrollOffset > 0 && d.scrollOffset < len(contentLines) {
		contentLines = contentLines[d.scrollOffset:]
	}
	if len(contentLines) > viewableHeight {
		contentLines = contentLines[:viewableHeight]
	}

	scrolledContent := strings.Join(contentLines, "\n")

	// Status bar
	statusBar := d.renderStatusBar()

	// Build final view
	mainContent := d.styles.PanelActive.
		Width(d.width).
		Height(d.height - 2).
		Render(scrolledContent)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		mainContent,
		statusBar,
	)
}

func (d *DetailView) renderMetadata() string {
	typeStyle := d.styles.TypeBadge(string(d.item.Type))
	stateStyle := d.styles.StateBadge(string(d.item.State))

	labelWidth := 12
	valueWidth := 20

	label := func(s string) string {
		return d.styles.DetailLabel.Width(labelWidth).Render(s)
	}
	value := func(s string) string {
		return d.styles.DetailValue.Width(valueWidth).Render(s)
	}

	rows := []string{
		label("Type:") + typeStyle.Render(d.item.ShortType()) + "     " + label("ID:") + value(fmt.Sprintf("#%d", d.item.ID)),
		label("State:") + stateStyle.Render(string(d.item.State)) + "     " + label("Created:") + value(d.item.CreatedDate.Format("2006-01-02")),
	}

	assignedTo := d.item.AssignedTo
	if assignedTo == "" {
		assignedTo = "Unassigned"
	}
	rows = append(rows, label("Assigned:")+value(assignedTo)+"  "+label("Updated:")+value(d.item.ChangedDate.Format("2006-01-02")))
	rows = append(rows, label("Sprint:")+value(d.item.SprintName())+"  "+label("Priority:")+value(fmt.Sprintf("%d", d.item.Priority)))
	rows = append(rows, label("Area:")+value(d.item.AreaName()))

	return strings.Join(rows, "\n")
}

func (d *DetailView) renderStatusBar() string {
	help := "Esc Back  Enter Open in browser  j/k Scroll"
	return d.styles.StatusBar.
		Width(d.width).
		Render(help)
}

// SetItem sets the work item to display
func (d *DetailView) SetItem(item *models.WorkItem) {
	d.item = item
	d.scrollOffset = 0
}

// SetSize sets the size of the detail view
func (d *DetailView) SetSize(width, height int) {
	d.width = width
	d.height = height
}

// CloseDetailViewMsg is sent when the detail view should be closed
type CloseDetailViewMsg struct{}
