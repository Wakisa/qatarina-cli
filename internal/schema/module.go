package schema

type CreateModuleRequest struct {
	ProjectID   int32  `json:"projectID"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Priority    int32  `json:"priority"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

type UpdateModuleRequest struct {
	ID          int32  `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Priority    int32  `json:"priority"`
	Type        string `josn:"name"`
	Description string `json:"description"`
}

type ModulesResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
