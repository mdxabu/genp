/*
Copyright © 2025 - github.com/mdxabu
*/
package cmd

import (
	"github.com/fatih/color"
	"github.com/mdxabu/genp/internal/github"
	"github.com/mdxabu/genp/internal/store"
	"github.com/spf13/cobra"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync local passwords to GitHub vault",
	Long: `Manually sync your local encrypted passwords to the genp-vault
private repository on GitHub.

This command will:
  1. Verify your GitHub authentication
  2. Create the genp-vault repo if it doesn't exist
  3. Push your local genp.yaml to the repo

You must be logged in first. Use 'genp login' to authenticate.

Examples:
  genp sync`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if logged in
		tokenInfo, err := github.LoadToken()
		if err != nil {
			color.Red("Error: %v\n", err)
			color.Yellow("Run 'genp login --token <token>' or 'genp login --oauth' to authenticate first.\n")
			return
		}

		color.Cyan("Logged in as %s\n", tokenInfo.Username)

		// Ensure vault repo exists
		color.Cyan("Ensuring genp-vault repository exists...\n")
		repo, err := github.CreateOrGetVaultRepo(tokenInfo.Token)
		if err != nil {
			color.Red("[error] Failed to set up vault repository: %v\n", err)
			return
		}
		color.Green("[ok] Vault repository ready: %s (private: %v)\n", repo.FullName, repo.Private)

		// Get local config path
		confPath, err := store.GetConfigFilePath()
		if err != nil {
			color.Red("[error] Failed to determine config file path: %v\n", err)
			return
		}

		// Sync the config file
		color.Cyan("Pushing genp.yaml to vault...\n")
		if err := github.SyncConfigToVault(confPath); err != nil {
			color.Red("[error] Failed to sync: %v\n", err)
			return
		}

		color.Green("[ok] Successfully synced passwords to %s\n", repo.HTMLURL)
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
