package cmd

import (
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/wakisa/qatarina-cli/internal/client"
	"github.com/wakisa/qatarina-cli/internal/schema"
	"github.com/wakisa/qatarina-cli/internal/tui"

	"github.com/spf13/cobra"
)

var testCaseCmd = &cobra.Command{
	Use:   "test-case",
	Short: "Test Case commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var createTestCaseCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new test case interactively",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runCreateTestCase()
	},
}

func runCreateTestCase() error {
	// Launch Bubble Tea TUI
	m := tui.NewCreateModel()
	prog := tea.NewProgram(m)
	final, err := prog.Run()
	if err != nil {
		return err
	}
	cm, ok := final.(*tui.CreateModel)
	if !ok {
		return fmt.Errorf("unexpected model type: %T", final)
	}
	a := cm.Answers()

	// Validate answers
	if len(a) != 8 {
		return fmt.Errorf("incomplete answers: expected 8 fields, got %d", len(a))
	}
	for i, val := range a {
		if strings.TrimSpace(val) == "" {
			fieldNames := []string{
				"Title", "Kind", "Project ID", "Description", "Code",
				"Feature/Module", "Is Draft", "Tags",
			}
			return fmt.Errorf("missing value for field %s", fieldNames[i])
		}
	}

	// Parse and transform
	projectID, err := strconv.ParseInt(a[2], 10, 64)
	if err != nil || projectID <= 0 {
		return fmt.Errorf("invalid project ID: %v", a[2])
	}
	isDraft, err := strconv.ParseBool(a[6])
	if err != nil {
		return fmt.Errorf("invalid value for Is Draft: %v", a[6])
	}
	tags := []string{}
	for _, t := range strings.Split(a[7], ",") {
		if trimmed := strings.TrimSpace(t); trimmed != "" {
			tags = append(tags, trimmed)
		}
	}

	// Construct payload
	payload := schema.CreateTestCaseRequest{
		Title:           a[0],
		Kind:            a[1],
		ProjectID:       projectID,
		Description:     a[3],
		Code:            a[4],
		FeatureOrModule: a[5],
		IsDraft:         isDraft,
		Tags:            tags,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Submit to API
	resp, err := client.Default().Post("v1/test-cases", body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("API error: %s", string(bodyBytes))
	}

	var msg schema.MessageResponse
	if err := json.Unmarshal(bodyBytes, &msg); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Println(msg.Message)
	return nil
}

var listTestCasesCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "list all test cases for a project",
	Example: "qatarina-cli test-case list --project 2",
	RunE: func(cmd *cobra.Command, args []string) error {
		projectID, err := cmd.Flags().GetInt64("project")
		if err != nil {
			return fmt.Errorf("project ID is invalid: %w", err)
		}
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

var viewTestCaseCmd = &cobra.Command{
	Use:     "view [test-case-id]",
	Aliases: []string{"show", "get"},
	Short:   "View a test case by ID",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runViewTestCasesByID(args[0])
	},
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
		return fmt.Errorf("failed to read response body: %w", err)
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
	fmt.Printf("• Project: %d\n", tc.ProjectID)
	fmt.Printf("• Title: %s\n", tc.Title)
	fmt.Printf("• Code: %s\n", tc.Code)
	fmt.Printf("• Kind: %s\n", tc.Kind)
	fmt.Printf("• Description: %s\n", tc.Description)
	fmt.Printf("• Feature/Module: %s\n", tc.FeatureOrModule)
	fmt.Printf("• Tags: %v\n", tc.Tags)
	fmt.Printf("• Draft: %v\n", tc.IsDraft)
	fmt.Printf("• Created By: %d\n", tc.CreatedByID)
	fmt.Printf("• Created At: %s\n", cmp.Or(strings.TrimSpace(tc.CreatedAt), "N/A"))
	fmt.Printf("• Updated At: %s\n", cmp.Or(strings.TrimSpace(tc.UpdatedAt), "N/A"))

	return nil
}

var deleteTestCaseCmd = &cobra.Command{
	Use:     "delete [test-case-id]",
	Aliases: []string{"rm"},
	Short:   "Delete a test case by ID",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDeleteTestCase(args[0])
	},
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
		return fmt.Errorf("failed to read response body: %w", err)
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

	listTestCasesCmd.Flags().Int64("project", 0, "Project ID")

	testCaseCmd.AddCommand(createTestCaseCmd)
	testCaseCmd.AddCommand(listTestCasesCmd)
	testCaseCmd.AddCommand(viewTestCaseCmd)
	testCaseCmd.AddCommand(deleteTestCaseCmd)

	rootCmd.AddCommand(testCaseCmd)

}
