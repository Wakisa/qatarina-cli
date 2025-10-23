package schema

type ExcelTestCase struct {
	Title           string   `json:"title"`
	Kind            string   `json:"kind"`
	Description     string   `json:"description"`
	Code            string   `json:"code"`
	FeatureOrModule string   `json:"feature_or_module"`
	IsDraft         bool     `json:"is_draft"`
	Tags            []string `json:"tags"`
}

type BulkCreateTestCaseRequest struct {
	ProjectID int64           `json:"project_id"`
	TestCases []ExcelTestCase `json:"test_cases"`
}
