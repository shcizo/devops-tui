package models

import "strings"

// Area represents an Azure DevOps area
type Area struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
}

// DisplayName returns a shortened display name
func (a Area) DisplayName() string {
	parts := strings.Split(a.Path, "\\")
	if len(parts) > 1 {
		// Return the last part for brevity
		return parts[len(parts)-1]
	}
	return a.Name
}

// FullPath returns the full path
func (a Area) FullPath() string {
	return a.Path
}
