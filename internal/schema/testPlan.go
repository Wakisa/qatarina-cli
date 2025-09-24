package schema

type TestCaseAssignment struct {
	TestCaseID string  `json:"test_case_id"`
	UserIDs    []int64 `json:"user_ids"`
}

type AssignTestToPlanRequest struct {
	ProjectID    int64                `json:"project_id"`
	PlanID       int64                `json:"test_plan_id"`
	PlannedTests []TestCaseAssignment `json:"planned_tests"`
}
