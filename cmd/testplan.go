package cmd

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"github.com/wakisa/qatarina-cli/internal/client"
	"github.com/wakisa/qatarina-cli/internal/schema"
)

var assignCasesCmd = &cobra.Command{
	Use:   "assign-cases",
	Short: "Assign test cases to a test plan",
	RunE: func(cmd *cobra.Command, args []string) error {
		projectID, _ := cmd.Flags().GetInt64("project")
		planID, _ := cmd.Flags().GetInt64("plan")
		assignmentsJSON, _ := cmd.Flags().GetString("assignments")

		if projectID == 0 || planID == 0 || assignmentsJSON == "" {
			return fmt.Errorf("project, plan and assignments are required")
		}

		var plannedTests []schema.TestCaseAssignment
		if err := json.Unmarshal([]byte(assignmentsJSON), &plannedTests); err != nil {
			return fmt.Errorf("failed to parse assignements JSON: %w", err)
		}

		payload := schema.AssignTestToPlanRequest{
			ProjectID:    projectID,
			PlanID:       planID,
			PlannedTests: plannedTests,
		}
		body, _ := json.Marshal(payload)

		path := fmt.Sprintf("v1/test-plans/%d/test-cases", planID)
		resp, err := client.Default().Post(path, body)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		if resp.StatusCode != 200 {
			return fmt.Errorf("API error: %s", string(bodyBytes))
		}

		var messageResponse schema.MessageResponse
		if err := json.Unmarshal(bodyBytes, &messageResponse); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
		fmt.Println(messageResponse.Message)
		return nil
	},
}

func init() {
	assignCasesCmd.Flags().Int64("project", 0, "Project ID")
	assignCasesCmd.Flags().Int64("plan", 0, "Test Plan ID")
	assignCasesCmd.Flags().String("assignments", "", "JSON array of assignments")

	rootCmd.AddCommand(assignCasesCmd)
}
