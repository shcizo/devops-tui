package models

// TeamMember represents a team member
type TeamMember struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	UniqueName  string `json:"uniqueName"`
}
