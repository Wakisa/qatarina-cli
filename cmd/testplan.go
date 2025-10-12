package cmd

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"github.com/wakisa/qatarina-cli/internal/client"
	"github.com/wakisa/qatarina-cli/internal/schema"
	"github.com/wakisa/qatarina-cli/internal/tui"
)

var assignCasesCmd = &cobra.Command{
	Use:   "assign-cases",
	Short: "Assign test cases to a test plan",
	RunE: func(cmd *cobra.Command, args []string) error {
		projectID, _ := cmd.Flags().GetInt64("project")
		planID, _ := cmd.Flags().GetInt64("plan")

		if projectID == 0 || planID == 0 {
			return fmt.Errorf("project and plan are required")
		}

		testCases, err := fetchTestCases(projectID)
		if err != nil {
			return err
		}

		model, err := tui.RunAssignUI(projectID, planID, testCases)
		if err != nil {
			return err
		}

		assignments := model.CollectedAssignments()
		if len(assignments) == 0 {
			fmt.Println("No test cases selected.")
			return nil
		}

		payload := schema.AssignTestToPlanRequest{
			ProjectID:    projectID,
			PlanID:       planID,
			PlannedTests: assignments,
		}

		body, _ := json.Marshal(payload)
		path := fmt.Sprintf("v1/test-plans/%d/test-cases", planID)
		resp, err := client.Default().Post(path, body)
		if err != nil {
			return fmt.Errorf("API error: %w", err)
		}
		defer resp.Body.Close()

		bodyBytes, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != 200 {
			return fmt.Errorf("API error: %s", string(bodyBytes))
		}

		var message schema.MessageResponse
		if err := json.Unmarshal(bodyBytes, &message); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
		fmt.Println(message.Message)
		return nil
	},
}

func fetchTestCases(projectID int64) ([]schema.TestCaseResponse, error) {
	path := fmt.Sprintf("v1/projects/%d/test-cases", projectID)
	resp, err := client.Default().Get(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API error: %s", string(bodyBytes))
	}

	var wrapper struct {
		TestCases []schema.TestCaseResponse `json:"test_cases"`
	}
	if err := json.Unmarshal(bodyBytes, &wrapper); err != nil {
		return nil, err
	}
	return wrapper.TestCases, nil
}

func init() {
	assignCasesCmd.Flags().Int64("project", 0, "Project ID")
	assignCasesCmd.Flags().Int64("plan", 0, "Test Plan ID")

	rootCmd.AddCommand(assignCasesCmd)
}
