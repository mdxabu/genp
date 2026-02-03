/*
Copyright © 2026 @mdxabu

*/

package crypto

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/fatih/color"
	"golang.org/x/term"
)

// PromptForMasterPassword prompts the user to enter a master password
// without echoing it to the terminal. Returns the password or an error.
func PromptForMasterPassword(promptText string) (string, error) {
	if promptText == "" {
		promptText = "Enter master password: "
	}

	color.Magenta(promptText)
	
	// Read password without echoing
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println() // Add newline after password input
	
	if err != nil {
		return "", fmt.Errorf("failed to read password: %w", err)
	}

	password := strings.TrimSpace(string(bytePassword))
	if password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}

	return password, nil
}

// PromptForMasterPasswordWithConfirm prompts for a master password twice
// to confirm it matches. Used when setting a new master password.
func PromptForMasterPasswordWithConfirm() (string, error) {
	password1, err := PromptForMasterPassword("Enter master password: ")
	if err != nil {
		return "", err
	}

	password2, err := PromptForMasterPassword("Confirm master password: ")
	if err != nil {
		return "", err
	}

	if password1 != password2 {
		return "", fmt.Errorf("passwords do not match")
	}

	return password1, nil
}

// CheckMasterPasswordExists checks if a master password has been set
// by checking if the config file exists
func CheckMasterPasswordExists(configPath string) bool {
	_, err := os.Stat(configPath)
	return err == nil
}
