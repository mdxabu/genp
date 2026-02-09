/*
Copyright 2025 - github.com/mdxabu
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

// PromptForMasterPassword prompts the user to enter their system lock screen
// password once (without echoing) and verifies it against the OS.
// On success it returns the verified password for use as the encryption key.
func PromptForMasterPassword(promptText string) (string, error) {
	if promptText == "" {
		promptText = "Enter system password: "
	}

	color.New(color.FgMagenta).Print(promptText)

	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()

	if err != nil {
		return "", fmt.Errorf("failed to read password: %w", err)
	}

	password := strings.TrimSpace(string(bytePassword))
	if password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}

	// Verify the password against the operating system
	if err := VerifySystemPassword(password); err != nil {
		return "", fmt.Errorf("authentication failed: %w", err)
	}

	return password, nil
}

// CheckMasterPasswordExists checks if a config file already exists at the
// given path, which indicates passwords have been stored before.
func CheckMasterPasswordExists(configPath string) bool {
	// No longer used for branching into a double-prompt flow, but kept for
	// backward compatibility with callers that check whether the config
	// file is present (e.g. to print different informational messages).
	return fileExists(configPath)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
