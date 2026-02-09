/*
Copyright 2025 - github.com/mdxabu
*/

package store

import (
	"fmt"
	"runtime"

	"github.com/fatih/color"
	"github.com/mdxabu/genp/internal/crypto"
)

func StorepasswordLocally(password string) string {
	var passwordName string
	color.New(color.FgCyan).Print("Enter a name for the password: ")
	fmt.Scanln(&passwordName)

	OSName := runtime.GOOS

	// Prompt for system password (single prompt, verified against OS)
	masterPassword, err := crypto.PromptForMasterPassword("Enter system password: ")
	if err != nil {
		color.Red("Failed to authenticate: %v\n", err)
		return ""
	}

	// Encrypt the password
	encryptedPassword, err := crypto.Encrypt(password, masterPassword)
	if err != nil {
		color.Red("Failed to encrypt password: %v\n", err)
		return ""
	}

	confPath, err := StoreLocalConfig(passwordName, encryptedPassword, OSName)
	if err != nil {
		color.Red("Failed to store password locally: %v\n", err)
		return ""
	}

	color.Green("Password encrypted and stored locally at: %s\n", confPath)
	return confPath
}
