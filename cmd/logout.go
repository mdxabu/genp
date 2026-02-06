/*
Copyright © 2026 @mdxabu
*/
package cmd

import (
	"github.com/fatih/color"
	"github.com/mdxabu/genp/internal/auth"
	"github.com/spf13/cobra"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout and remove stored GitHub token",
	Long:  `Remove the stored GitHub OAuth token from genp.yaml.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := auth.Logout(); err != nil {
			color.Red("Logout failed: %v\n", err)
			return
		}
		color.Green("Successfully logged out. Token removed from genp.yaml.\n")
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
