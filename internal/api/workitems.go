package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/samuelenocsson/devops-tui/internal/models"
)

// wiqlRequest represents a WIQL query request
type wiqlRequest struct {
	Query string `json:"query"`
}

// wiqlResponse represents the response from a WIQL query
type wiqlResponse struct {
	WorkItems []struct {
		ID  int    `json:"id"`
		URL string `json:"url"`
	} `json:"workItems"`
}

// workItemsResponse represents the response for batch work item fetch
type workItemsResponse struct {
	Count int               `json:"count"`
	Value []workItemAPIItem `json:"value"`
}

// workItemAPIItem represents a work item from the API
type workItemAPIItem struct {
	ID     int              `json:"id"`
	Rev    int              `json:"rev"`
	Fields workItemFields   `json:"fields"`
	URL    string           `json:"url"`
}

type workItemFields struct {
	ID           int    `json:"System.Id"`
	Title        string `json:"System.Title"`
	State        string `json:"System.State"`
	WorkItemType string `json:"System.WorkItemType"`
	AssignedTo   *struct {
		DisplayName string `json:"displayName"`
		UniqueName  string `json:"uniqueName"`
	} `json:"System.AssignedTo"`
	IterationPath string     `json:"System.IterationPath"`
	AreaPath      string     `json:"System.AreaPath"`
	Description   string     `json:"System.Description"`
	Tags          string     `json:"System.Tags"`
	Parent        int        `json:"System.Parent"`
	Priority      int        `json:"Microsoft.VSTS.Common.Priority"`
	CreatedDate   time.Time  `json:"System.CreatedDate"`
	ChangedDate   time.Time  `json:"System.ChangedDate"`
}

// escapeWIQL escapes a string value for use in WIQL queries
func escapeWIQL(s string) string {
	// Escape single quotes by doubling them
	return strings.ReplaceAll(s, "'", "''")
}

// QueryWorkItems queries work items using WIQL
func (c *Client) QueryWorkItems(sprintPath, state, assigned, areaPath string) ([]models.WorkItem, error) {
	// Build WIQL query
	query := `SELECT [System.Id], [System.Title], [System.State], [System.WorkItemType]
FROM WorkItems
WHERE [System.TeamProject] = @project`

	// Add sprint filter
	if sprintPath != "" && sprintPath != "all" {
		query += fmt.Sprintf(`
  AND [System.IterationPath] = '%s'`, escapeWIQL(sprintPath))
	}

	// Add state filter
	if state != "" && state != "all" {
		query += fmt.Sprintf(`
  AND [System.State] = '%s'`, escapeWIQL(state))
	}

	// Add assigned filter
	if assigned == "me" {
		query += `
  AND [System.AssignedTo] = @me`
	}

	// Add area filter
	if areaPath != "" && areaPath != "all" {
		// Clean up the path
		areaPath = strings.TrimPrefix(areaPath, "\\")
		areaPath = strings.TrimSuffix(areaPath, "\\")
		query += fmt.Sprintf(`
  AND [System.AreaPath] UNDER '%s'`, escapeWIQL(areaPath))
	}

	query += `
ORDER BY [System.ChangedDate] DESC`

	// Execute WIQL query
	reqBody := wiqlRequest{Query: query}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshaling WIQL request: %w", err)
	}

	resp, err := c.post("/wit/wiql", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	var wiqlResp wiqlResponse
	if err := decode(resp, &wiqlResp); err != nil {
		return nil, err
	}

	if len(wiqlResp.WorkItems) == 0 {
		return []models.WorkItem{}, nil
	}

	// Get the IDs
	ids := make([]string, 0, len(wiqlResp.WorkItems))
	for _, wi := range wiqlResp.WorkItems {
		ids = append(ids, fmt.Sprintf("%d", wi.ID))
	}

	// Fetch the full work items
	return c.GetWorkItems(ids)
}

// GetWorkItems fetches multiple work items by ID
func (c *Client) GetWorkItems(ids []string) ([]models.WorkItem, error) {
	if len(ids) == 0 {
		return []models.WorkItem{}, nil
	}

	// API has a limit of 200 items per request
	const batchSize = 200
	var allItems []models.WorkItem

	for i := 0; i < len(ids); i += batchSize {
		end := i + batchSize
		if end > len(ids) {
			end = len(ids)
		}

		batch := ids[i:end]
		fields := "System.Id,System.Title,System.State,System.WorkItemType,System.AssignedTo,System.IterationPath,System.AreaPath,System.Description,System.Tags,System.Parent,Microsoft.VSTS.Common.Priority,System.CreatedDate,System.ChangedDate"

		endpoint := fmt.Sprintf("/wit/workitems?ids=%s&fields=%s", strings.Join(batch, ","), fields)
		resp, err := c.get(endpoint)
		if err != nil {
			return nil, err
		}

		var apiResp workItemsResponse
		if err := decode(resp, &apiResp); err != nil {
			return nil, err
		}

		for _, item := range apiResp.Value {
			wi := c.convertWorkItem(item)
			allItems = append(allItems, wi)
		}
	}

	// Fetch parent titles
	c.populateParentTitles(allItems)

	return allItems, nil
}

// populateParentTitles fetches titles for all parent work items
func (c *Client) populateParentTitles(items []models.WorkItem) {
	// Collect unique parent IDs
	parentIDs := make(map[int]bool)
	for _, item := range items {
		if item.ParentID > 0 {
			parentIDs[item.ParentID] = true
		}
	}

	if len(parentIDs) == 0 {
		return
	}

	// Convert to string slice
	ids := make([]string, 0, len(parentIDs))
	for id := range parentIDs {
		ids = append(ids, fmt.Sprintf("%d", id))
	}

	// Fetch parent work items (only need ID and Title)
	endpoint := fmt.Sprintf("/wit/workitems?ids=%s&fields=System.Id,System.Title", strings.Join(ids, ","))
	resp, err := c.get(endpoint)
	if err != nil {
		return // Silently fail - parent titles are optional
	}

	var apiResp workItemsResponse
	if err := decode(resp, &apiResp); err != nil {
		return
	}

	// Build ID -> Title map
	titleMap := make(map[int]string)
	for _, item := range apiResp.Value {
		titleMap[item.ID] = item.Fields.Title
	}

	// Update items with parent titles
	for i := range items {
		if items[i].ParentID > 0 {
			if title, ok := titleMap[items[i].ParentID]; ok {
				items[i].ParentTitle = title
			}
		}
	}
}

// GetWorkItem fetches a single work item by ID
func (c *Client) GetWorkItem(id int) (*models.WorkItem, error) {
	fields := "System.Id,System.Title,System.State,System.WorkItemType,System.AssignedTo,System.IterationPath,System.AreaPath,System.Description,System.Tags,System.Parent,Microsoft.VSTS.Common.Priority,System.CreatedDate,System.ChangedDate"

	endpoint := fmt.Sprintf("/wit/workitems/%d?fields=%s", id, fields)
	resp, err := c.get(endpoint)
	if err != nil {
		return nil, err
	}

	var item workItemAPIItem
	if err := decode(resp, &item); err != nil {
		return nil, err
	}

	wi := c.convertWorkItem(item)

	// Fetch parent title if parent exists
	if wi.ParentID > 0 {
		parentEndpoint := fmt.Sprintf("/wit/workitems/%d?fields=System.Title", wi.ParentID)
		parentResp, err := c.get(parentEndpoint)
		if err == nil {
			var parentItem workItemAPIItem
			if decode(parentResp, &parentItem) == nil {
				wi.ParentTitle = parentItem.Fields.Title
			}
		}
	}

	return &wi, nil
}

// convertWorkItem converts an API work item to our model
func (c *Client) convertWorkItem(item workItemAPIItem) models.WorkItem {
	wi := models.WorkItem{
		ID:            item.ID,
		Rev:           item.Rev,
		Title:         item.Fields.Title,
		State:         models.WorkItemState(item.Fields.State),
		Type:          models.WorkItemType(item.Fields.WorkItemType),
		IterationPath: item.Fields.IterationPath,
		AreaPath:      item.Fields.AreaPath,
		Description:   stripHTML(item.Fields.Description),
		ParentID:      item.Fields.Parent,
		Priority:      item.Fields.Priority,
		CreatedDate:   item.Fields.CreatedDate,
		ChangedDate:   item.Fields.ChangedDate,
		URL:           item.URL,
		WebURL:        c.WorkItemWebURL(item.ID),
	}

	if item.Fields.AssignedTo != nil {
		wi.AssignedTo = item.Fields.AssignedTo.DisplayName
	}

	// Parse tags
	if item.Fields.Tags != "" {
		tags := strings.Split(item.Fields.Tags, ";")
		for _, tag := range tags {
			tag = strings.TrimSpace(tag)
			if tag != "" {
				wi.Tags = append(wi.Tags, tag)
			}
		}
	}

	return wi
}

// UpdateWorkItemState updates a work item's state
func (c *Client) UpdateWorkItemState(id int, newState string) error {
	// Azure DevOps uses JSON Patch format
	patchDoc := []map[string]interface{}{
		{
			"op":    "add",
			"path":  "/fields/System.State",
			"value": newState,
		},
	}

	bodyBytes, err := json.Marshal(patchDoc)
	if err != nil {
		return fmt.Errorf("marshaling patch document: %w", err)
	}

	endpoint := fmt.Sprintf("/wit/workitems/%d", id)
	resp, err := c.patch(endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}
	resp.Body.Close()

	return nil
}

// AssignWorkItem assigns a work item to a user
// Pass empty string to unassign
func (c *Client) AssignWorkItem(id int, userEmail string) error {
	// Azure DevOps uses JSON Patch format
	patchDoc := []map[string]interface{}{
		{
			"op":    "add",
			"path":  "/fields/System.AssignedTo",
			"value": userEmail,
		},
	}

	bodyBytes, err := json.Marshal(patchDoc)
	if err != nil {
		return fmt.Errorf("marshaling patch document: %w", err)
	}

	endpoint := fmt.Sprintf("/wit/workitems/%d", id)
	resp, err := c.patch(endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}
	resp.Body.Close()

	return nil
}

// stripHTML removes HTML tags from a string
func stripHTML(s string) string {
	// Simple HTML tag removal
	re := regexp.MustCompile(`<[^>]*>`)
	s = re.ReplaceAllString(s, "")

	// Replace common HTML entities
	s = strings.ReplaceAll(s, "&nbsp;", " ")
	s = strings.ReplaceAll(s, "&amp;", "&")
	s = strings.ReplaceAll(s, "&lt;", "<")
	s = strings.ReplaceAll(s, "&gt;", ">")
	s = strings.ReplaceAll(s, "&quot;", "\"")
	s = strings.ReplaceAll(s, "&#39;", "'")

	// Trim whitespace
	s = strings.TrimSpace(s)

	return s
}
