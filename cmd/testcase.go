package cmd

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/wakisa/qatarina-cli/common"
	"github.com/wakisa/qatarina-cli/internal/client"
	"github.com/wakisa/qatarina-cli/internal/schema"

	"github.com/spf13/cobra"
)

var testCaseCmd = &cobra.Command{
	Use:   "test-case",
	Short: "Test Case commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		viewID, _ := cmd.Flags().GetString("view")
		deleteID, _ := cmd.Flags().GetString("delete")

		switch {
		case viewID != "":
			return runViewTestCasesByID(viewID)
		case deleteID != "":
			return runDeleteTestCase(deleteID)
		default:
			return cmd.Help()

		}
	},
}

var createTestCaseCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a new test case manually",
	Example: `qatarina-cli test-case create --title "Login flow" --kind "functional" --project 2 --code "TC-001"`,
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

var listTestCasesCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Example: "qatarina-cli testcase list --project 2",
	Short:   "list all test cases for a project",
	RunE: func(cmd *cobra.Command, args []string) error {
		projectID, _ := cmd.Flags().GetInt64("project")
		return runViewTestCases(projectID)
	},
}

func runViewTestCases(projectID int64) error {
	testCases, err := fetchTestCases(projectID)
	if err != nil {
		return err
	}

	if len(testCases) == 0 {
		fmt.Println("No test cases found.")
		return nil
	}

	fmt.Printf("Test Cases for Project %d:\n", projectID)
	for _, tc := range testCases {
		fmt.Printf("• %s\n Code: %s\n Kind: %s\n ID: %s\n\n", tc.Title, tc.Code, tc.Kind, tc.ID)
	}
	return nil
}

func runViewTestCasesByID(id string) error {
	path := fmt.Sprintf("v1/test-cases/%s", id)
	resp, err := client.Default().Get(path)
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

	var wrapper struct {
		TestCase schema.TestCaseResponse `json:"test_case"`
	}
	if err := json.Unmarshal(bodyBytes, &wrapper); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	tc := wrapper.TestCase

	fmt.Printf("Test Case Details:\n")
	fmt.Printf("• ID: %s\n", tc.ID)
	fmt.Printf("• Project ID: %d\n", tc.ProjectID)
	fmt.Printf("• Title: %s\n", tc.Title)
	fmt.Printf("• Code: %s\n", tc.Code)
	fmt.Printf("• Kind: %s\n", tc.Kind)
	fmt.Printf("• Description: %s\n", tc.Description)
	fmt.Printf("• Feature/Module: %s\n", tc.FeatureOrModule)
	fmt.Printf("• Tags: %v\n", tc.Tags)
	fmt.Printf("• Draft: %v\n", tc.IsDraft)
	fmt.Printf("• Created By: %d\n", tc.CreatedByID)
	fmt.Printf("• Created At: %s\n", common.DefaultIfEmpty(tc.CreatedAt))
	fmt.Printf("• Updated At: %s\n", common.DefaultIfEmpty(tc.UpdatedAt))

	return nil
}

func runDeleteTestCase(id string) error {
	path := fmt.Sprintf("v1/test-cases/%s", id)
	resp, err := client.Default().Delete(path)
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

	var message schema.MessageResponse
	if err := json.Unmarshal(bodyBytes, &message); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Println(message.Message)
	return nil
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

	testCaseCmd.Flags().String("view", "", "View test case by ID")
	testCaseCmd.Flags().String("delete", "", "Delete test case by ID")

	listTestCasesCmd.Flags().Int64("project", 0, "Project ID")

	testCaseCmd.AddCommand(createTestCaseCmd)
	testCaseCmd.AddCommand(listTestCasesCmd)

	rootCmd.AddCommand(testCaseCmd)

}
