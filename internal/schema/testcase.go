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
	ProjectID       int64    `json:"project_id"`
	CreatedByID     int64    `json:"created_by"`
	Kind            string   `json:"kind"`
	Code            string   `json:"code"`
	FeatureOrModule string   `json:"feature_or_module"`
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	IsDraft         bool     `json:"is_draft"`
	Tags            []string `json:"tags"`
	CreatedAt       string   `json:"created_at"`
	UpdatedAt       string   `json:"updated_at"`
}

type UpdateTestCaseRequest struct {
	ID              string   `json:"id"`
	Kind            string   `json:"kind"`
	Code            string   `json:"code"`
	FeatureOrModule string   `json:"feature_or_module"`
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	IsDraft         string   `json:"is_draft"`
	Tags            []string `json:"tags"`
}
