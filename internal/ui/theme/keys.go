package theme

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines all keyboard shortcuts
type KeyMap struct {
	// Navigation
	Up        key.Binding
	Down      key.Binding
	Left      key.Binding
	Right     key.Binding
	Top       key.Binding
	Bottom    key.Binding
	NextPanel key.Binding
	PrevPanel key.Binding

	// Actions
	Select       key.Binding
	Open         key.Binding
	View         key.Binding
	Search       key.Binding
	Refresh      key.Binding
	Help         key.Binding
	Back         key.Binding
	Quit         key.Binding
	ChangeState  key.Binding
	CreateBranch key.Binding
}

// DefaultKeyMap returns the default key bindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "right"),
		),
		Top: key.NewBinding(
			key.WithKeys("g"),
			key.WithHelp("g", "top"),
		),
		Bottom: key.NewBinding(
			key.WithKeys("G"),
			key.WithHelp("G", "bottom"),
		),
		NextPanel: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("Tab", "next panel"),
		),
		PrevPanel: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("Shift+Tab", "prev panel"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter", " "),
			key.WithHelp("Enter/Space", "select"),
		),
		Open: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("Enter", "open in browser"),
		),
		View: key.NewBinding(
			key.WithKeys("v"),
			key.WithHelp("v", "view details"),
		),
		Search: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "search"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("ctrl+r"),
			key.WithHelp("Ctrl+r", "refresh"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("Esc", "back"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		ChangeState: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "change state"),
		),
		CreateBranch: key.NewBinding(
			key.WithKeys("b"),
			key.WithHelp("b", "create branch"),
		),
	}
}

// ShortHelp returns a short help text for the status bar
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.NextPanel,
		k.Up, k.Down,
		k.Top, k.Bottom,
		k.Open,
		k.Help,
	}
}

// FullHelp returns all key bindings for the help panel
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Top, k.Bottom},
		{k.NextPanel, k.PrevPanel},
		{k.Select, k.Open, k.View},
		{k.ChangeState, k.CreateBranch},
		{k.Search, k.Refresh},
		{k.Help, k.Back, k.Quit},
	}
}
