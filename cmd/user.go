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
		return runCreateUser(cmd)
	},
}

func getFlag(cmd *cobra.Command, name string) string {
	val, err := cmd.Flags().GetString(name)
	if err != nil {
		return ""
	}
	return val
}

func runCreateUser(cmd *cobra.Command) error {
	a := map[string]string{
		"FirstName":   getFlag(cmd, "first-name"),
		"LastName":    getFlag(cmd, "last-name"),
		"DisplayName": getFlag(cmd, "display-name"),
		"Email":       getFlag(cmd, "email"),
		"Password":    getFlag(cmd, "password"),
	}

	// If all required flags are present, skip wizard
	allFlagsPresent := true
	for _, key := range []string{"FirstName", "LastName", "DisplayName", "Email", "Password"} {
		if strings.TrimSpace(a[key]) == "" {
			allFlagsPresent = false
			break
		}
	}

	if !allFlagsPresent {
		// Launch wizard
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
		a = um.Answers()
	}

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

var viewCmd = &cobra.Command{
	Use:   "view [userID]",
	Short: "View user by ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		resp, err := client.Default().Get("v1/users/" + id)
		if err != nil {
			fmt.Println("Failed to fetch user:", err)
			return
		}
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Failed to read response:", err)
			return
		}

		var rawResponse map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &rawResponse); err != nil {
			fmt.Println("Failed to decode response:", err)
			return
		}

		if rawResponse["ID"] == nil {
			fmt.Println("No user found with that ID.")
			return
		}

		fmt.Println("User Details:")
		if idVal, ok := rawResponse["ID"].(float64); ok {
			fmt.Printf("• ID: %d\n", int(idVal))
		} else {
			fmt.Println("• ID: <not set>")
		}
		fmt.Printf("• Name: %s %s\n", safeResponse(rawResponse["FirstName"]), safeResponse(rawResponse["LastName"]))
		fmt.Printf("• Display Name: %s\n", safeResponse(rawResponse["DisplayName"]))
		fmt.Printf("• Email: %s\n", safeResponse(rawResponse["Email"]))
		fmt.Printf("• Created At: %s\n", safeResponse(rawResponse["CreatedAt"]))
	},
}

func safeResponse(v interface{}) string {
	switch val := v.(type) {
	case string:
		if val != "" {
			return val
		}
	case map[string]interface{}:
		if s, ok := val["String"].(string); ok && s != "" {
			return s
		}
		if t, ok := val["Time"].(string); ok && t != "" {
			return t
		}
	}
	return "<not set>"
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
	createCmd.Flags().String("first-name", "", "First name")
	createCmd.Flags().String("last-name", "", "Last name")
	createCmd.Flags().String("display-name", "", "Display name")
	createCmd.Flags().String("email", "", "Email address")
	createCmd.Flags().String("password", "", "Password")
	createCmd.Flags().String("org", "", "Organization ID (optional)")
	userCmd.AddCommand(createCmd)
	userCmd.AddCommand(listCmd)
	userCmd.AddCommand(viewCmd)
	userCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(userCmd)
}
