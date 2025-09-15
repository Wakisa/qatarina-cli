package cmd

import (
	"encoding/json"
	"fmt"

	"qatarina-cli/internal/auth"
	"qatarina-cli/internal/client"
	"qatarina-cli/internal/schema"

	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with the Qatarina API",
	RunE: func(cmd *cobra.Command, args []string) error {
		email, _ := cmd.Flags().GetString("email")
		password, _ := cmd.Flags().GetString("password")

		if email == "" || password == "" {
			return fmt.Errorf("email and password are required")
		}

		payload := schema.LoginResquest{
			Email:    email,
			Password: password,
		}
		body, _ := json.Marshal(payload)

		resp, err := client.Default().Post("v1/auth/login", body)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		var result schema.LoginResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		token := result.Token
		if token == "" {
			return fmt.Errorf("login failed: no token received")
		}

		if err := auth.SaveToken(token); err != nil {
			return fmt.Errorf("failed to save token: %w", err)
		}
		fmt.Println("Logged in successfully!")
		return nil
	},
}

func init() {
	loginCmd.Flags().String("email", "", "your email")
	loginCmd.Flags().String("password", "", "your password")
	rootCmd.AddCommand(loginCmd)
}
