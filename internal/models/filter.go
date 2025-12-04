package models

// FilterType represents the type of filter
type FilterType int

const (
	FilterTypeSprint FilterType = iota
	FilterTypeState
	FilterTypeAssigned
	FilterTypeArea
)

// FilterOption represents a selectable filter option
type FilterOption struct {
	Label    string
	Value    string
	Selected bool
}

// FilterGroup represents a group of filter options
type FilterGroup struct {
	Type    FilterType
	Title   string
	Options []FilterOption
	Cursor  int
	Offset  int // Scroll offset for viewing
}

// SelectedOption returns the currently selected option
func (f *FilterGroup) SelectedOption() *FilterOption {
	for i := range f.Options {
		if f.Options[i].Selected {
			return &f.Options[i]
		}
	}
	if len(f.Options) > 0 {
		return &f.Options[0]
	}
	return nil
}

// Select marks the option at the given index as selected
func (f *FilterGroup) Select(index int) {
	if index < 0 || index >= len(f.Options) {
		return
	}
	for i := range f.Options {
		f.Options[i].Selected = i == index
	}
}

// SelectCurrent marks the option at the current cursor as selected
func (f *FilterGroup) SelectCurrent() {
	f.Select(f.Cursor)
}

// MoveUp moves the cursor up
func (f *FilterGroup) MoveUp() {
	if f.Cursor > 0 {
		f.Cursor--
	}
}

// MoveDown moves the cursor down
func (f *FilterGroup) MoveDown() {
	if f.Cursor < len(f.Options)-1 {
		f.Cursor++
	}
}

// MoveToTop moves cursor to the first option
func (f *FilterGroup) MoveToTop() {
	f.Cursor = 0
}

// MoveToBottom moves cursor to the last option
func (f *FilterGroup) MoveToBottom() {
	if len(f.Options) > 0 {
		f.Cursor = len(f.Options) - 1
	}
}

// FilterState holds the complete filter state
type FilterState struct {
	Groups       []*FilterGroup
	ActiveGroup  int
	SearchQuery  string
}

// NewFilterState creates a new filter state with default groups
func NewFilterState(iterations []Iteration, areas []Area, statesByType map[string][]WorkItemStateInfo) *FilterState {
	// Build sprint options from iterations
	sprintOptions := []FilterOption{
		{Label: "All", Value: "all", Selected: false},
	}

	for _, iter := range iterations {
		selected := iter.IsCurrent()
		sprintOptions = append(sprintOptions, FilterOption{
			Label:    iter.DisplayName(),
			Value:    iter.Path,
			Selected: selected,
		})
	}

	// If no current sprint found, select "All"
	hasSelected := false
	for _, opt := range sprintOptions {
		if opt.Selected {
			hasSelected = true
			break
		}
	}
	if !hasSelected && len(sprintOptions) > 0 {
		sprintOptions[0].Selected = true
	}

	// Build area options from areas
	areaOptions := []FilterOption{
		{Label: "All", Value: "all", Selected: true},
	}
	for _, area := range areas {
		areaOptions = append(areaOptions, FilterOption{
			Label:    area.DisplayName(),
			Value:    area.Path,
			Selected: false,
		})
	}

	// Build state options from all work item types (unique states)
	stateOptions := []FilterOption{
		{Label: "All", Value: "all", Selected: true},
	}
	if len(statesByType) > 0 {
		// Collect unique states preserving order by category
		seenStates := make(map[string]bool)
		// Process in category order: Proposed, InProgress, Resolved, Completed
		categoryOrder := []string{"Proposed", "InProgress", "Resolved", "Completed", "Removed"}
		for _, category := range categoryOrder {
			for _, states := range statesByType {
				for _, state := range states {
					if state.Category == category && !seenStates[state.Name] {
						seenStates[state.Name] = true
						stateOptions = append(stateOptions, FilterOption{
							Label:    state.Name,
							Value:    state.Name,
							Selected: false,
						})
					}
				}
			}
		}
		// Also add any states that didn't match known categories
		for _, states := range statesByType {
			for _, state := range states {
				if !seenStates[state.Name] {
					seenStates[state.Name] = true
					stateOptions = append(stateOptions, FilterOption{
						Label:    state.Name,
						Value:    state.Name,
						Selected: false,
					})
				}
			}
		}
	}
	// Fallback to default states if we only have "All"
	if len(stateOptions) <= 1 {
		defaultStates := []string{"New", "Active", "Resolved", "Closed"}
		for _, state := range defaultStates {
			stateOptions = append(stateOptions, FilterOption{
				Label:    state,
				Value:    state,
				Selected: false,
			})
		}
	}

	return &FilterState{
		Groups: []*FilterGroup{
			{
				Type:    FilterTypeSprint,
				Title:   "Sprint",
				Options: sprintOptions,
				Cursor:  0,
			},
			{
				Type:    FilterTypeState,
				Title:   "State",
				Options: stateOptions,
				Cursor:  0,
			},
			{
				Type:  FilterTypeAssigned,
				Title: "Assigned",
				Options: []FilterOption{
					{Label: "All", Value: "all", Selected: false},
					{Label: "Me", Value: "me", Selected: true},
				},
				Cursor: 0,
			},
			{
				Type:    FilterTypeArea,
				Title:   "Area",
				Options: areaOptions,
				Cursor:  0,
			},
		},
		ActiveGroup: 0,
	}
}

// ActiveFilterGroup returns the currently active filter group
func (f *FilterState) ActiveFilterGroup() *FilterGroup {
	if f.ActiveGroup >= 0 && f.ActiveGroup < len(f.Groups) {
		return f.Groups[f.ActiveGroup]
	}
	return nil
}

// NextGroup moves to the next filter group
func (f *FilterState) NextGroup() {
	if f.ActiveGroup < len(f.Groups)-1 {
		f.ActiveGroup++
	}
}

// PrevGroup moves to the previous filter group
func (f *FilterState) PrevGroup() {
	if f.ActiveGroup > 0 {
		f.ActiveGroup--
	}
}

// GetSelectedSprint returns the selected sprint path
func (f *FilterState) GetSelectedSprint() string {
	for _, g := range f.Groups {
		if g.Type == FilterTypeSprint {
			if opt := g.SelectedOption(); opt != nil {
				return opt.Value
			}
		}
	}
	return "all"
}

// GetSelectedState returns the selected state filter
func (f *FilterState) GetSelectedState() string {
	for _, g := range f.Groups {
		if g.Type == FilterTypeState {
			if opt := g.SelectedOption(); opt != nil {
				return opt.Value
			}
		}
	}
	return "all"
}

// GetSelectedAssigned returns the selected assigned filter
func (f *FilterState) GetSelectedAssigned() string {
	for _, g := range f.Groups {
		if g.Type == FilterTypeAssigned {
			if opt := g.SelectedOption(); opt != nil {
				return opt.Value
			}
		}
	}
	return "all"
}

// GetSelectedArea returns the selected area path
func (f *FilterState) GetSelectedArea() string {
	for _, g := range f.Groups {
		if g.Type == FilterTypeArea {
			if opt := g.SelectedOption(); opt != nil {
				return opt.Value
			}
		}
	}
	return "all"
}

// ApplySavedSelections applies saved filter selections
func (f *FilterState) ApplySavedSelections(sprint, state, assigned, area string) {
	for _, g := range f.Groups {
		var targetValue string
		switch g.Type {
		case FilterTypeSprint:
			targetValue = sprint
		case FilterTypeState:
			targetValue = state
		case FilterTypeAssigned:
			targetValue = assigned
		case FilterTypeArea:
			targetValue = area
		}

		if targetValue == "" {
			continue
		}

		// Find and select the matching option
		found := false
		for i, opt := range g.Options {
			if opt.Value == targetValue {
				g.Select(i)
				found = true
				break
			}
		}

		// If not found and it's sprint with "current", keep the current selection
		if !found && g.Type == FilterTypeSprint && targetValue == "current" {
			// Already handled in NewFilterState
		}
	}
}
