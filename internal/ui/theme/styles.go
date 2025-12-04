package theme

import "github.com/charmbracelet/lipgloss"

// Colors
var (
	colorPrimary    = lipgloss.Color("#7C3AED") // Purple
	colorSecondary  = lipgloss.Color("#06B6D4") // Cyan
	colorSuccess    = lipgloss.Color("#10B981") // Green
	colorWarning    = lipgloss.Color("#F59E0B") // Yellow
	colorError      = lipgloss.Color("#EF4444") // Red
	colorMuted      = lipgloss.Color("#6B7280") // Gray
	colorBorder     = lipgloss.Color("#374151") // Dark gray
	colorHighlight  = lipgloss.Color("#1F2937") // Very dark gray
	colorText       = lipgloss.Color("#F9FAFB") // White
	colorTextMuted  = lipgloss.Color("#9CA3AF") // Light gray
)

// Work item type colors
var typeColors = map[string]lipgloss.Color{
	"User Story": lipgloss.Color("#3B82F6"), // Blue
	"Story":      lipgloss.Color("#3B82F6"),
	"Task":       lipgloss.Color("#F59E0B"), // Yellow
	"Bug":        lipgloss.Color("#EF4444"), // Red
	"Feature":    lipgloss.Color("#8B5CF6"), // Purple
	"Epic":       lipgloss.Color("#EC4899"), // Pink
}

// Work item state colors
var stateColors = map[string]lipgloss.Color{
	// Common states
	"New":         lipgloss.Color("#6B7280"), // Gray
	"Active":      lipgloss.Color("#3B82F6"), // Blue
	"Resolved":    lipgloss.Color("#10B981"), // Green
	"Closed":      lipgloss.Color("#6B7280"), // Gray
	// Agile states
	"To Do":       lipgloss.Color("#F97316"), // Orange
	"In Progress": lipgloss.Color("#8B5CF6"), // Purple
	"Done":        lipgloss.Color("#10B981"), // Green
	"Testing":     lipgloss.Color("#FBBF24"), // Yellow
	// Additional states
	"Removed":     lipgloss.Color("#6B7280"), // Gray
	"Approved":    lipgloss.Color("#10B981"), // Green
}

// Styles defines all UI styles
type Styles struct {
	// Base styles
	App       lipgloss.Style
	Title     lipgloss.Style
	Subtitle  lipgloss.Style
	StatusBar lipgloss.Style

	// Panel styles
	PanelActive    lipgloss.Style
	PanelInactive  lipgloss.Style
	PanelTitle     lipgloss.Style

	// Filter styles
	FilterGroup      lipgloss.Style
	FilterGroupTitle lipgloss.Style
	FilterOption     lipgloss.Style
	FilterSelected   lipgloss.Style
	FilterCursor     lipgloss.Style

	// Work items list styles
	ListHeader       lipgloss.Style
	ListItem         lipgloss.Style
	ListItemSelected lipgloss.Style
	ListItemCursor   lipgloss.Style

	// Detail view styles
	DetailTitle       lipgloss.Style
	DetailSection     lipgloss.Style
	DetailSectionTitle lipgloss.Style
	DetailLabel       lipgloss.Style
	DetailValue       lipgloss.Style
	DetailDescription lipgloss.Style
	DetailTag         lipgloss.Style

	// Help styles
	HelpKey   lipgloss.Style
	HelpDesc  lipgloss.Style
	HelpPanel lipgloss.Style

	// Type and state badges
	TypeBadge  func(string) lipgloss.Style
	StateBadge func(string) lipgloss.Style
}

// DefaultStyles returns the default styles
func DefaultStyles() Styles {
	return Styles{
		// Base styles
		App: lipgloss.NewStyle().
			Background(lipgloss.Color("#111827")),

		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(colorText).
			Padding(0, 1),

		Subtitle: lipgloss.NewStyle().
			Foreground(colorTextMuted),

		StatusBar: lipgloss.NewStyle().
			Foreground(colorTextMuted).
			Background(lipgloss.Color("#1F2937")).
			Padding(0, 1),

		// Panel styles
		PanelActive: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorPrimary).
			Padding(0, 1),

		PanelInactive: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			Padding(0, 1),

		PanelTitle: lipgloss.NewStyle().
			Bold(true).
			Foreground(colorText).
			Background(colorHighlight).
			Padding(0, 1),

		// Filter styles
		FilterGroup: lipgloss.NewStyle().
			MarginBottom(1),

		FilterGroupTitle: lipgloss.NewStyle().
			Bold(true).
			Foreground(colorText).
			MarginBottom(0),

		FilterOption: lipgloss.NewStyle().
			Foreground(colorTextMuted).
			PaddingLeft(1),

		FilterSelected: lipgloss.NewStyle().
			Foreground(colorSuccess).
			PaddingLeft(1),

		FilterCursor: lipgloss.NewStyle().
			Foreground(colorPrimary).
			Bold(true),

		// Work items list styles
		ListHeader: lipgloss.NewStyle().
			Bold(true).
			Foreground(colorTextMuted).
			BorderBottom(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(colorBorder),

		ListItem: lipgloss.NewStyle().
			Foreground(colorText),

		ListItemSelected: lipgloss.NewStyle().
			Foreground(colorText).
			Background(colorHighlight),

		ListItemCursor: lipgloss.NewStyle().
			Foreground(colorPrimary).
			Bold(true),

		// Detail view styles
		DetailTitle: lipgloss.NewStyle().
			Bold(true).
			Foreground(colorText).
			MarginBottom(1),

		DetailSection: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			Padding(0, 1).
			MarginBottom(1),

		DetailSectionTitle: lipgloss.NewStyle().
			Bold(true).
			Foreground(colorTextMuted).
			MarginBottom(0),

		DetailLabel: lipgloss.NewStyle().
			Foreground(colorTextMuted).
			Width(12),

		DetailValue: lipgloss.NewStyle().
			Foreground(colorText),

		DetailDescription: lipgloss.NewStyle().
			Foreground(colorText),

		DetailTag: lipgloss.NewStyle().
			Foreground(colorSecondary).
			Background(lipgloss.Color("#1E3A5F")).
			Padding(0, 1).
			MarginRight(1),

		// Help styles
		HelpKey: lipgloss.NewStyle().
			Foreground(colorPrimary).
			Bold(true),

		HelpDesc: lipgloss.NewStyle().
			Foreground(colorTextMuted),

		HelpPanel: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			Padding(1, 2),

		// Type and state badge generators
		TypeBadge: func(t string) lipgloss.Style {
			color := typeColors[t]
			if color == "" {
				color = colorMuted
			}
			return lipgloss.NewStyle().
				Foreground(color).
				Bold(true)
		},

		StateBadge: func(s string) lipgloss.Style {
			color := stateColors[s]
			if color == "" {
				color = colorMuted
			}
			return lipgloss.NewStyle().
				Foreground(color)
		},
	}
}
