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

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage users in Qatarina",
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new user via wizard",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runCreateUser()
	},
}

func runCreateUser() error {
	m := tui.NewUserCreateModel()
	prog := tea.NewProgram(m)
	final, err := prog.Run()
	if err != nil {
		return err
	}
	um, ok := final.(*tui.UserCreateModel)
	if !ok {
		return fmt.Errorf("unexpected model type: %T", final)
	}
	a := um.Answers()

	required := []string{"FirstName", "LastName", "DisplayName", "Email", "Password"}
	for _, key := range required {
		if strings.TrimSpace(a[key]) == "" {
			return fmt.Errorf("missing value for field %s", key)
		}
	}

	payload := schema.NewUserRequest{
		FirstName:   a["FirstName"],
		LastName:    a["LastName"],
		DisplayName: a["DisplayName"],
		Email:       a["Email"],
		Password:    a["Password"],
		OrgID:       0, // Optional: pass via flag
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := client.Default().Post("v1/users", body)
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

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all users",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := client.Default().Get("v1/users")
		if err != nil {
			fmt.Println("Failed to fetch users:", err)
			return
		}
		defer resp.Body.Close()

		var result schema.CompactUserListResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			fmt.Println("Failed to parse response:", err)
			return
		}

		fmt.Printf("Number of Users: %d\n", len(result.Users))
		for _, u := range result.Users {
			fmt.Printf("• ID: %d | Name: %s | Email: %s | Created: %s\n", u.ID, u.DisplayName, u.Email, u.CreatedAt)
		}
	},
}

var getCmd = &cobra.Command{
	Use:   "get [userID]",
	Short: "Get user by ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		resp, err := client.Default().Get("v1/users/" + id)
		if err != nil {
			fmt.Println("Failed to fetch user:", err)
			return
		}
		defer resp.Body.Close()

		var user map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&user); err != nil {
			fmt.Println("Failed to parse response:", err)
			return
		}

		fmt.Printf("User: %+v\n", user)
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete [userID]",
	Short: "Delete user by ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		resp, err := client.Default().Delete("v1/users/" + id)
		if err != nil {
			fmt.Println("Failed to delete user:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			bodyBytes, _ := io.ReadAll(resp.Body)
			fmt.Printf("API error: %s\n", string(bodyBytes))
			return
		}

		fmt.Println("User deleted successfully.")
	},
}

func init() {
	userCmd.AddCommand(createCmd)
	userCmd.AddCommand(listCmd)
	userCmd.AddCommand(getCmd)
	userCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(userCmd)
}
