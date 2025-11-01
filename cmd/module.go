package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

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

		payload := schema.CreateModuleRequest{
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

var updateModuleCmd = &cobra.Command{
	Use:   "update <moduleID>",
	Short: "Update a module",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid module ID: %w", err)
		}

		name, _ := cmd.Flags().GetString("name")
		code, _ := cmd.Flags().GetString("code")
		priority, _ := cmd.Flags().GetInt32("priority")
		moduleType, _ := cmd.Flags().GetString("type")
		description, _ := cmd.Flags().GetString("description")

		payload := schema.UpdateModuleRequest{
			ID:          int32(id),
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

		resp, err := client.Default().Post("v1/modules/"+args[0], body)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		bodyBytes, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != 200 {
			return fmt.Errorf("API error: %s", string(bodyBytes))
		}

		fmt.Println("Module updated successfully.")
		return nil
	},
}

var listModulesCmd = &cobra.Command{
	Use:   "list",
	Short: "List all modules",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := client.Default().Get("v1/modules")
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		var wrapper struct {
			Modules []schema.ModulesResponse `json:"modules"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
			return fmt.Errorf("failed to decode reponse: %w", err)
		}

		for _, m := range wrapper.Modules {
			fmt.Printf("• [%d] %s — %s\n", m.ID, m.Name, m.Description)
		}
		return nil
	},
}

var viewModuleCmd = &cobra.Command{
	Use:   "view <moduleID>",
	Short: "View module details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		resp, err := client.Default().Get("v1/modules/" + id)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}

		if resp.StatusCode != 200 {
			if len(bodyBytes) == 0 {
				return fmt.Errorf("module not found (ID: %s)", id)
			}
			return fmt.Errorf("module not found (ID: %s): %s", id, string(bodyBytes))
		}

		var module schema.ModulesResponse
		if err := json.Unmarshal(bodyBytes, &module); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}

		fmt.Printf("Module: %s\nID: %d\nDescription: %s\n", module.Name, module.ID, module.Description)
		return nil
	},
}

var deleteModuleCmd = &cobra.Command{
	Use:   "delete <moduleID>",
	Short: "Delete a module",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		resp, err := client.Default().Delete("v1/modules/" + id)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		bodyBytes, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != 200 {
			return fmt.Errorf("failed to delete module (ID: %s): %s", id, string(bodyBytes))
		}

		fmt.Println("Module deleted successfully.")
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

	updateModuleCmd.Flags().String("name", "", "Module name")
	updateModuleCmd.Flags().String("code", "", "Module code")
	updateModuleCmd.Flags().Int32("priority", 0, "Module priority")
	updateModuleCmd.Flags().String("type", "", "Module type")
	updateModuleCmd.Flags().String("description", "", "Module description")

	moduleCmd.AddCommand(createModuleCmd)
	moduleCmd.AddCommand(updateModuleCmd)
	moduleCmd.AddCommand(listModulesCmd)
	moduleCmd.AddCommand(viewModuleCmd)
	moduleCmd.AddCommand(deleteModuleCmd)
	rootCmd.AddCommand(moduleCmd)
}
