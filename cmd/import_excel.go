package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wakisa/qatarina-cli/internal/client"
	"github.com/wakisa/qatarina-cli/internal/schema"
	"github.com/xuri/excelize/v2"
)

var importFileCmd = &cobra.Command{
	Use:   "import-file",
	Short: "Import test cases from an Excel or CSV file",
	RunE: func(cmd *cobra.Command, args []string) error {
		projectID, err := cmd.Flags().GetInt64("project")
		if err != nil || projectID == 0 {
			return fmt.Errorf("invalid or missing projectID: %w", err)
		}
		filePath, err := cmd.Flags().GetString("file")
		if err != nil || filePath == "" {
			return fmt.Errorf("file path is required")
		}

		ext := strings.ToLower(filepath.Ext(filePath))
		var rows [][]string

		switch ext {
		case ".xlsx":
			rows, err = parseExcel(filePath)
		case ".csv":
			rows, err = parseCSV(filePath)
		default:
			return fmt.Errorf("unsupported file type: %s", ext)
		}
		if err != nil {
			return fmt.Errorf("failed to parse file: %w", err)
		}

		testCases := collectedTestCases(rows)
		if len(testCases) == 0 {
			fmt.Println("No valid test cases found.")
			return nil
		}

		payload := schema.BulkCreateTestCaseRequest{
			ProjectID: projectID,
			TestCases: testCases,
		}
		body, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to encode payload: %w", err)
		}

		resp, err := client.Default().Post("v1/test-cases/bulk", body)
		if err != nil {
			return fmt.Errorf("API error: %w", err)
		}
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil || resp.StatusCode != 200 {
			return fmt.Errorf("API error: %s", string(bodyBytes))
		}

		var message schema.MessageResponse
		if err := json.Unmarshal(bodyBytes, &message); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}

		fmt.Println(message.Message)
		fmt.Printf("Imported %d test cases to project %d\n", len(testCases), projectID)
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

func parseCSV(path string) ([][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	return reader.ReadAll()
}
func collectedTestCases(rows [][]string) []schema.ExcelTestCase {
	var cases []schema.ExcelTestCase
	start := findHeaderIndex(rows)
	for _, row := range rows[start:] { // skip header
		if len(row) < 7 {
			continue // skip incomplete rows
		}
		tags := strings.Split(row[5], ",")
		for i := range tags {
			tags[i] = strings.TrimSpace(tags[i])
		}
		cases = append(cases, schema.ExcelTestCase{
			Title:           row[0],
			Description:     row[1],
			Kind:            row[2],
			Code:            row[3],
			FeatureOrModule: row[4],
			Tags:            tags,
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
	importFileCmd.Flags().Int64("project", 0, "Project ID")
	importFileCmd.Flags().String("file", "", "Path to excel or CSV file")
	importFileCmd.MarkFlagRequired("project")
	importFileCmd.MarkFlagRequired("file")
	rootCmd.AddCommand(importFileCmd)
}
