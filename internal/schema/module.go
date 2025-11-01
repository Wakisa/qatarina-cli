package schema

type ModuleRequest struct {
	ProjectID   int32  `json:"projectID"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Priority    int32  `json:"priority"`
	Type        string `json:"type"`
	Description string `json:"description"`
}
