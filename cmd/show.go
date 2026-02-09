/*
Copyright © 2025 @mdxabu
*/
package cmd

import (
	"github.com/fatih/color"
	"github.com/mdxabu/genp/internal/crypto"
	"github.com/mdxabu/genp/internal/store"
	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Display stored passwords",
	Long: `Display all stored passwords after decrypting them with your master password.

This command will prompt you for your master password and then display
all stored passwords in decrypted form.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get all encrypted passwords
		passwords, err := store.GetAllPasswords()
		if err != nil {
			color.Red("Error: %v\n", err)
			return
		}

		// Prompt for master password
		masterPassword, err := crypto.PromptForMasterPassword("Enter system password: ")
		if err != nil {
			color.Red("Error reading master password: %v\n", err)
			return
		}

		// Decrypt and display all passwords
		color.Cyan("\n=== Stored Passwords ===\n")
		hasError := false
		for name, encrypted := range passwords {
			decrypted, err := store.DecryptPassword(encrypted, masterPassword)
			if err != nil {
				color.Red("%s: [Failed to decrypt - incorrect master password or corrupted data]\n", name)
				hasError = true
				continue
			}
			color.New(color.FgGreen).Printf("%s: ", name)
			color.Yellow("%s\n", decrypted)
		}

		if hasError {
			color.Red("\nNote: Some passwords could not be decrypted. Please check your master password.\n")
		}
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
