/*
Copyright 2025 - github.com/mdxabu
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/mdxabu/genp/internal/github"
	"github.com/mdxabu/genp/internal/store"
	"github.com/spf13/cobra"
)

var (
	loginToken    string
	loginOAuth    bool
	oauthClientID string
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with GitHub for cloud vault sync",
	Long: `Login to GitHub to enable automatic syncing of your encrypted passwords
to a private repository called 'genp-vault'.

You can authenticate using either:
  1. A personal access token (PAT) with the 'repo' scope
  2. OAuth device flow (browser-based)

Examples:
  genp login --token ghp_xxxxxxxxxxxxxxxxxxxx
  genp login --oauth
  genp login --oauth --client-id YOUR_CLIENT_ID`,
	Run: func(cmd *cobra.Command, args []string) {
		if loginToken == "" && !loginOAuth {
			color.Yellow("Please specify a login method:\n")
			color.Cyan("  genp login --token <your-github-token>\n")
			color.Cyan("  genp login --oauth [--client-id <client-id>]\n")
			fmt.Println()
			color.White("To create a personal access token:\n")
			color.White("  1. Go to https://github.com/settings/tokens\n")
			color.White("  2. Generate a new token (classic) with 'repo' scope\n")
			color.White("  3. Run: genp login --token <your-token>\n")
			return
		}

		if loginToken != "" && loginOAuth {
			color.Red("Error: cannot use both --token and --oauth at the same time\n")
			return
		}

		if loginToken != "" {
			loginWithToken()
		} else if loginOAuth {
			loginWithOAuth()
		}
	},
}

// statusCmd represents the login status subcommand
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check GitHub login status",
	Long:  `Check if you are currently authenticated with GitHub and display account info.`,
	Run: func(cmd *cobra.Command, args []string) {
		tokenInfo, err := github.LoadToken()
		if err != nil {
			color.Red("Not logged in to GitHub.\n")
			color.Yellow("Run 'genp login --token <token>' or 'genp login --oauth' to authenticate.\n")
			return
		}

		color.Green("[ok] Logged in to GitHub\n")
		color.Cyan("  Username:   %s\n", tokenInfo.Username)
		color.Cyan("  Login type: %s\n", tokenInfo.LoginType)
		color.Cyan("  Token:      %s****\n", tokenInfo.Token[:4])
	},
}

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove stored GitHub credentials",
	Long:  `Remove the stored GitHub authentication token from your local machine.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := github.Logout()
		if err != nil {
			color.Red("Error: %v\n", err)
			return
		}
		color.Green("[ok] Successfully logged out from GitHub.\n")
	},
}

func setupVaultAndSync(token string) {
	// Create or get the vault repo
	color.Cyan("Setting up genp-vault repository...\n")
	repo, err := github.CreateOrGetVaultRepo(token)
	if err != nil {
		color.Yellow("[warn] Could not set up vault repository: %v\n", err)
		color.Yellow("  You can try again later with 'genp sync'. Passwords will still be stored locally.\n")
		return
	}

	color.Green("[ok] Vault repository ready: %s (private: %v)\n", repo.FullName, repo.Private)
	color.Green("  URL: %s\n", repo.HTMLURL)

	// Push existing local passwords to the vault
	confPath, err := store.GetConfigFilePath()
	if err != nil {
		color.Yellow("[warn] Could not determine config file path: %v\n", err)
		return
	}

	color.Cyan("Pushing existing passwords to vault...\n")
	if err := github.SyncConfigToVault(confPath); err != nil {
		color.Yellow("[warn] Failed to push existing passwords: %v\n", err)
		color.Yellow("  You can retry with 'genp sync'.\n")
	} else {
		color.Green("[ok] Existing passwords synced to genp-vault repository.\n")
	}
}

func loginWithToken() {
	token := strings.TrimSpace(loginToken)
	if token == "" {
		color.Red("Error: token cannot be empty\n")
		return
	}

	color.Cyan("Authenticating with GitHub...\n")

	info, err := github.LoginWithToken(token)
	if err != nil {
		color.Red("Error: %v\n", err)
		return
	}

	color.Green("[ok] Successfully logged in as %s\n", info.Username)

	tokenPath, err := github.GetTokenStorePath()
	if err == nil {
		color.Cyan("Token stored at: %s\n", tokenPath)
	}

	setupVaultAndSync(token)
}

func loginWithOAuth() {
	color.Cyan("Starting GitHub OAuth device flow...\n")

	info, err := github.LoginWithOAuth(oauthClientID)
	if err != nil {
		color.Red("Error: %v\n", err)
		return
	}

	color.Green("\n[ok] Successfully logged in as %s\n", info.Username)

	tokenPath, err := github.GetTokenStorePath()
	if err == nil {
		color.Cyan("Token stored at: %s\n", tokenPath)
	}

	setupVaultAndSync(info.Token)
}

func init() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
	loginCmd.AddCommand(statusCmd)

	loginCmd.Flags().StringVar(&loginToken, "token", "", "GitHub personal access token with 'repo' scope")
	loginCmd.Flags().BoolVar(&loginOAuth, "oauth", false, "Use OAuth device flow for authentication")
	loginCmd.Flags().StringVar(&oauthClientID, "client-id", "", "OAuth App Client ID (for --oauth flow)")
}
