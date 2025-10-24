package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/wakisa/qatarina-cli/internal/client"
	"github.com/wakisa/qatarina-cli/internal/schema"
	"github.com/wakisa/qatarina-cli/internal/tui"
)

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Project commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var createProjectCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new project (via flags or wizard)",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		version, _ := cmd.Flags().GetString("version")
		websiteURL, _ := cmd.Flags().GetString("website-url")
		githubURL, _ := cmd.Flags().GetString("github-url")

		if name == "" || description == "" || version == "" || websiteURL == "" {
			fmt.Println("Launching interactive wizard...")
			return runCreateProjectWizard()
		}

		payload := schema.NewProjectRequest{
			Name:        name,
			Description: description,
			Version:     version,
			WebsiteURL:  websiteURL,
			GitHubURL:   githubURL,
		}
		return submitProject(payload)

	},
}

func runCreateProjectWizard() error {
	m := tui.NewCreateProjectModel()
	prog := tea.NewProgram(m)
	final, err := prog.Run()
	if err != nil {
		return err
	}

	pm, ok := final.(*tui.CreateProjectModel)
	if !ok {
		return fmt.Errorf("unexpected model type: %T", final)
	}
	a := pm.Answers()

	// Validate
	required := []string{"Name", "Description", "Version", "Website URL"}
	for _, key := range required {
		if strings.TrimSpace(a[key]) == "" {
			return fmt.Errorf("missing value for field %s", key)
		}
	}

	payload := schema.NewProjectRequest{
		Name:        a["Name"],
		Description: a["Description"],
		Version:     a["Version"],
		WebsiteURL:  a["Website URL"],
		GitHubURL:   a["GitHub URL"],
	}
	return submitProject(payload)
}

func submitProject(payload schema.NewProjectRequest) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := client.Default().Post("v1/projects", body)
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

	var wrapper struct {
		Project schema.ProjectResponse `json:"project"`
	}
	if err := json.Unmarshal(bodyBytes, &wrapper); err != nil {
		return fmt.Errorf("failed ot decode response: %w", err)
	}

	fmt.Printf("Project created: %s (ID: %d)\n", wrapper.Project.Title, wrapper.Project.ID)
	return nil
}

func init() {
	createProjectCmd.Flags().String("name", "", "Project name")
	createProjectCmd.Flags().String("description", "", "Project description")
	createProjectCmd.Flags().String("version", "", "Project version")
	createProjectCmd.Flags().String("website-url", "", "Project website URL")
	createProjectCmd.Flags().String("github-url", "", "Project GitHub URL")

	projectCmd.AddCommand(createProjectCmd)
	rootCmd.AddCommand(projectCmd)
}
