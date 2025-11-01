package cmd

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"github.com/wakisa/qatarina-cli/internal/client"
	"github.com/wakisa/qatarina-cli/internal/schema"
)

var moduleCmd = &cobra.Command{
	Use:   "module",
	Short: "Module commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var createModuleCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new moudle",
	RunE: func(cmd *cobra.Command, args []string) error {
		projectID, _ := cmd.Flags().GetInt32("project-id")
		name, _ := cmd.Flags().GetString("name")
		code, _ := cmd.Flags().GetString("code")
		priority, _ := cmd.Flags().GetInt32("priority")
		moduleType, _ := cmd.Flags().GetString("type")
		description, _ := cmd.Flags().GetString("description")

		payload := schema.ModuleRequest{
			ProjectID:   projectID,
			Name:        name,
			Code:        code,
			Priority:    priority,
			Type:        moduleType,
			Description: description,
		}

		body, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshal payload: %w", err)
		}

		resp, err := client.Default().Post("v1/modules", body)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		bodyBytes, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != 200 {
			return fmt.Errorf("API error: %s", string(bodyBytes))
		}

		fmt.Println("Module created succcessfully.")
		return nil
	},
}

func init() {
	createModuleCmd.Flags().Int32("project-id", 0, "Project ID")
	createModuleCmd.Flags().String("name", "", "Module name")
	createModuleCmd.Flags().String("code", "", "Module code")
	createModuleCmd.Flags().Int32("priority", 0, "Module priority")
	createModuleCmd.Flags().String("type", "", "Module type")
	createModuleCmd.Flags().String("description", "", "Module description")

	moduleCmd.AddCommand(createModuleCmd)
	rootCmd.AddCommand(moduleCmd)
}
