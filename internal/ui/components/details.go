package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/samuelenocsson/devops-tui/internal/models"
	"github.com/samuelenocsson/devops-tui/internal/ui/theme"
)

// DetailsPanel shows details for a selected work item
type DetailsPanel struct {
	item               *models.WorkItem
	styles             theme.Styles
	width              int
	height             int
	renderedDesc       string
	renderedDescWidth  int
}

// NewDetailsPanel creates a new details panel
func NewDetailsPanel(styles theme.Styles) DetailsPanel {
	return DetailsPanel{
		styles: styles,
	}
}

// View renders the details panel
func (d DetailsPanel) View() string {
	if d.item == nil {
		content := d.styles.Subtitle.Render("Select a work item to view details")
		return d.styles.PanelInactive.
			Width(d.width).
			Height(d.height).
			Render(content)
	}

	var b strings.Builder

	// Title
	title := fmt.Sprintf("#%d %s", d.item.ID, d.item.Title)
	if len(title) > d.width-6 {
		title = title[:d.width-9] + "..."
	}
	b.WriteString(d.styles.DetailTitle.Render(title))
	b.WriteString("\n\n")

	// Metadata section with aligned columns
	typeStyle := d.styles.TypeBadge(string(d.item.Type))
	stateStyle := d.styles.StateBadge(string(d.item.State))

	labelWidth := 10
	valueWidth := 12

	// Type and State row
	b.WriteString(d.styles.DetailLabel.Width(labelWidth).Render("Type:"))
	b.WriteString(typeStyle.Width(valueWidth).Render(d.item.ShortType()))
	b.WriteString(d.styles.DetailLabel.Width(labelWidth).Render("State:"))
	b.WriteString(stateStyle.Render(string(d.item.State)))
	b.WriteString("\n")

	// Assigned and Sprint row
	assignedTo := d.item.AssignedTo
	if assignedTo == "" {
		assignedTo = "Unassigned"
	}
	b.WriteString(d.styles.DetailLabel.Width(labelWidth).Render("Assigned:"))
	b.WriteString(d.styles.DetailValue.Width(valueWidth).Render(truncate(assignedTo, valueWidth)))
	b.WriteString(d.styles.DetailLabel.Width(labelWidth).Render("Sprint:"))
	b.WriteString(d.styles.DetailValue.Render(d.item.SprintName()))
	b.WriteString("\n")

	// Area row
	b.WriteString(d.styles.DetailLabel.Width(labelWidth).Render("Area:"))
	b.WriteString(d.styles.DetailValue.Render(d.item.AreaName()))
	b.WriteString("\n")

	// Parent (if exists)
	if d.item.ParentID > 0 {
		b.WriteString("\n")
		parentLabel := fmt.Sprintf("Parent: #%d", d.item.ParentID)
		if d.item.ParentTitle != "" {
			parentLabel += " " + d.item.ParentTitle
		}
		b.WriteString(d.styles.Subtitle.Render(parentLabel))
		b.WriteString("\n")
	}

	// Description section
	if d.item.Description != "" {
		b.WriteString("\n")
		b.WriteString(d.styles.DetailSectionTitle.Render("─── Description ───"))
		b.WriteString("\n")

		maxWidth := d.width - 8
		if maxWidth < 20 {
			maxWidth = 20
		}

		// Use cached rendered description if available and width matches
		desc := d.renderedDesc
		if desc == "" || d.renderedDescWidth != maxWidth {
			desc = d.renderMarkdown(d.item.Description, maxWidth)
		}

		// Limit description height
		lines := strings.Split(desc, "\n")
		maxLines := d.height - 12
		if maxLines < 3 {
			maxLines = 3
		}
		if len(lines) > maxLines {
			lines = lines[:maxLines]
			lines = append(lines, "...")
		}

		b.WriteString(strings.Join(lines, "\n"))
		b.WriteString("\n")
	}

	// Tags section
	if len(d.item.Tags) > 0 {
		b.WriteString("\n")
		b.WriteString(d.styles.DetailSectionTitle.Render("─── Tags ───"))
		b.WriteString("\n")

		var tagStrings []string
		for _, tag := range d.item.Tags {
			tagStrings = append(tagStrings, d.styles.DetailTag.Render(tag))
		}
		b.WriteString(strings.Join(tagStrings, " "))
	}

	content := b.String()
	return d.styles.PanelInactive.
		Width(d.width).
		Height(d.height).
		BorderForeground(lipgloss.Color("#374151")).
		Render(content)
}

// SetItem sets the work item to display
func (d *DetailsPanel) SetItem(item *models.WorkItem) {
	d.item = item
	d.renderedDesc = ""
	d.renderedDescWidth = 0
}

// renderMarkdown renders markdown content and caches the result
func (d *DetailsPanel) renderMarkdown(content string, width int) string {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithStylePath("dark"),
		glamour.WithWordWrap(width),
	)

	var result string
	if err == nil {
		rendered, renderErr := renderer.Render(content)
		if renderErr == nil {
			result = strings.TrimSpace(rendered)
		} else {
			result = wordWrap(content, width)
		}
	} else {
		result = wordWrap(content, width)
	}

	d.renderedDesc = result
	d.renderedDescWidth = width
	return result
}

// SetSize sets the size of the details panel
func (d *DetailsPanel) SetSize(width, height int) {
	d.width = width
	d.height = height
}

// Helper functions

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func wordWrap(text string, width int) string {
	if width <= 0 {
		return text
	}

	var result strings.Builder
	words := strings.Fields(text)
	currentLineLength := 0

	for i, word := range words {
		if currentLineLength+len(word)+1 > width {
			result.WriteString("\n")
			currentLineLength = 0
		} else if i > 0 {
			result.WriteString(" ")
			currentLineLength++
		}
		result.WriteString(word)
		currentLineLength += len(word)
	}

	return result.String()
}
