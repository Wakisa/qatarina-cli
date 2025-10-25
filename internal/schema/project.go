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
	WebsiteURL  string `json:"website_url"`
	GithubURL   string `json:"github_url"`
	TrelloURL   string `json:"trello_url"`
	JiraURL     string `json:"jira_url"`
	MondayURL   string `json:"monday_url"`
	OwnerUserID int32  `json:"owner_user_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type ModuleResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
