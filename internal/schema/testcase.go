package schema

type LoginRequest struct {
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

type TestCaseResponse struct {
	ID              string   `json:"id"`
	Title           string   `json:"title"`
	Code            string   `json:"code"`
	Kind            string   `json:"kind"`
	FeatureOrModule string   `json:"feature_or_module"`
	Tags            []string `json:"tags"`
	IsDraft         bool     `json:"is_draft"`
	CreatedByID     int64    `json:"created_by_id"`
	ProjectID       int64    `json:"project_id"`
}
