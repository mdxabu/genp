/*
Copyright © 2026 @mdxabu
*/
package cmd

import (
	"github.com/fatih/color"
	"github.com/mdxabu/genp/internal/auth"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to GitHub via device flow",
	Long:  `Authenticate with GitHub using the OAuth device flow. The token is stored locally in genp.yaml alongside your passwords.`,
	Run: func(cmd *cobra.Command, args []string) {
		color.Cyan("Logging in with GitHub...\n")
		token, err := auth.Login()
		if err != nil {
			color.Red("Login failed: %v\n", err)
			return
		}
		_ = token
		color.Green("Successfully logged in! Token stored in genp.yaml.\n")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
