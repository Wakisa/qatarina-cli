package schema

type LoginResquest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateTestCaseRequest struct {
	Title           string   `json:"title"`
	Kind            string   `json:"kind"`
	ProjectID       int64    `json:"project_id"`
	Description     string   `json:"description"`
	Code            string   `json:"code"`
	FeatureOrModule string   `json:"feature_or_module"`
	IsDraft         bool     `json:"is_draft"`
	Tags            []string `json:"tags"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expires_in"`
}
