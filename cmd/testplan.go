package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
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

		return RunAssignUI(projectID, planID)
	},
}

func init() {
	assignCasesCmd.Flags().Int64("project", 0, "Project ID")
	assignCasesCmd.Flags().Int64("plan", 0, "Test Plan ID")

	rootCmd.AddCommand(assignCasesCmd)
}
