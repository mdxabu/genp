/*
Copyright © 2025 @mdxabu

*/
package cmd

import (
	"fmt"

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
			fmt.Println("Error:", err)
			return
		}

		// Prompt for master password
		masterPassword, err := crypto.PromptForMasterPassword("Enter master password to decrypt: ")
		if err != nil {
			fmt.Println("Error reading master password:", err)
			return
		}

		// Decrypt and display all passwords
		fmt.Println("\n=== Stored Passwords ===")
		hasError := false
		for name, encrypted := range passwords {
			decrypted, err := store.DecryptPassword(encrypted, masterPassword)
			if err != nil {
				fmt.Printf("%s: [Failed to decrypt - incorrect master password or corrupted data]\n", name)
				hasError = true
				continue
			}
			fmt.Printf("%s: %s\n", name, decrypted)
		}

		if hasError {
			fmt.Println("\nNote: Some passwords could not be decrypted. Please check your master password.")
		}
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
