package cmd

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/wakisa/qatarina-cli/internal/client"
	"github.com/wakisa/qatarina-cli/internal/schema"

	"github.com/spf13/cobra"
)

var testCaseCmd = &cobra.Command{
	Use:   "test-case",
	Short: "Test Case commands",
}

var createTestCaseCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new test case manually",
	RunE: func(cmd *cobra.Command, args []string) error {
		title, _ := cmd.Flags().GetString("title")
		description, _ := cmd.Flags().GetString("description")
		kind, _ := cmd.Flags().GetString("kind")
		code, _ := cmd.Flags().GetString("code")
		projectID, _ := cmd.Flags().GetInt64("project")
		featureOrModule, _ := cmd.Flags().GetString("feature-or-module")
		isDraft, _ := cmd.Flags().GetBool("draft")
		tags, _ := cmd.Flags().GetStringSlice("tags")

		if title == "" || kind == "" || projectID == 0 {
			return fmt.Errorf("title, kind, and project ID are required")
		}

		payload := schema.CreateTestCaseRequest{
			Title:           title,
			Kind:            kind,
			ProjectID:       projectID,
			Description:     description,
			Code:            code,
			FeatureOrModule: featureOrModule,
			IsDraft:         isDraft,
			Tags:            tags,
		}
		body, _ := json.Marshal(payload)

		resp, err := client.Default().Post("v1/test-cases", body)
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
		fmt.Println(" " + messageResponse.Message)
		return nil
	},
}

func init() {
	createTestCaseCmd.Flags().String("title", "", "Title of the test case")
	createTestCaseCmd.Flags().String("kind", "", "Kind of the test case")
	createTestCaseCmd.Flags().Int64("project", 0, "Project ID")
	createTestCaseCmd.Flags().String("description", "", "Description of test case")
	createTestCaseCmd.Flags().String("code", "", "Code identifier")
	createTestCaseCmd.Flags().String("feature-or-module", "", "Feature or module name")
	createTestCaseCmd.Flags().Bool("draft", false, "Is this a draft")
	createTestCaseCmd.Flags().StringSlice("tags", []string{}, "Comma-separated tags")

	testCaseCmd.AddCommand(createTestCaseCmd)
	rootCmd.AddCommand(testCaseCmd)
}
