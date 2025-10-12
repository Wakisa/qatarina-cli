package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wakisa/qatarina-cli/internal/client"
	"github.com/wakisa/qatarina-cli/internal/schema"
	"github.com/wakisa/qatarina-cli/tui"
	"github.com/xuri/excelize/v2"
)

var importExcelCmd = &cobra.Command{
	Use:   "import-excel",
	Short: "Import test cases from an Excel file",
	RunE: func(cmd *cobra.Command, args []string) error {
		projectID, err := cmd.Flags().GetInt64("project")
		if err != nil || projectID == 0 {
			return fmt.Errorf("failed to process projectID: %w", err)
		}

		model, err := tui.RunImportExcelUI()
		if err != nil {
			return err
		}

		if model.FilePath == "" {
			return fmt.Errorf("no file selected")
		}

		rows, err := parseExcel(model.FilePath)
		if err != nil {
			return fmt.Errorf("failed to parse Excel file: %w", err)
		}

		testCases := collectedTestCases(rows)
		if len(testCases) == 0 {
			fmt.Println("No test cases to import.")
			return nil
		}

		payload := schema.BulkCreateTestCaseRequest{
			ProjectID: projectID,
			TestCases: testCases,
		}
		body, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to process the body :%w", err)
		}
		resp, err := client.Default().Post("v1/test-cases/bulk", body)
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

func parseExcel(path string) ([][]string, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	return f.GetRows("Sheet1")
}

func collectedTestCases(rows [][]string) []schema.ExcelTestCase {
	var cases []schema.ExcelTestCase
	start := findHeaderIndex(rows)
	for _, row := range rows[start:] { // skip header
		if len(row) < 7 {
			continue // skip incomplete rows
		}
		cases = append(cases, schema.ExcelTestCase{
			Title:           row[0],
			Description:     row[1],
			Kind:            row[2],
			Code:            row[3],
			FeatureOrModule: row[4],
			Tags:            strings.Split(row[5], ","),
			IsDraft:         strings.ToLower(row[6]) == "true",
		})
	}
	return cases
}

func findHeaderIndex(rows [][]string) int {
	for i, row := range rows {
		if len(row) >= 7 && strings.ToLower(row[0]) == "title" && strings.ToLower(row[2]) == "kind" {
			return i + 1
		}
	}
	return 1
}

func init() {
	importExcelCmd.Flags().Int64("project", 0, "project ID")
	rootCmd.AddCommand(importExcelCmd)
}
