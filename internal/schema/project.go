package schema

type NewProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	WebsiteURL  string `json:"website_url"`
	Version     string `json:"version"`
	GitHubURL   string `json:"github_url"`
}

type ProjectListResponse struct {
	Projects []ProjectResponse `json:"projects"`
}

type ProjectResponse struct {
	ID          int32  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
	IsActive    bool   `json:"is_active"`
	IsPublic    bool   `json:"is_public"`
	WebsiteUrl  string `json:"website_url"`
	GithubUrl   string `json:"github_url"`
	TrelloUrl   string `json:"trello_url"`
	JiraUrl     string `json:"jira_url"`
	MondayUrl   string `json:"monday_url"`
	OwnerUserID int32  `json:"owner_user_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
