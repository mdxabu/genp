/*
Copyright © 2026 @mdxabu

*/

package store

import (
	"fmt"
	"runtime"

	"github.com/fatih/color"
	"github.com/mdxabu/genp/internal/crypto"
)

func StorepasswordLocally(password string) {
	var passwordName string
	color.New(color.FgCyan).Print("Enter a name for the password: ")
	fmt.Scanln(&passwordName)

	OSName := runtime.GOOS

	// Get the config path to check if this is first time setup
	baseDir, err := ConfigBaseDir("genp", OSName)
	if err != nil {
		color.Red("Failed to determine config directory: %v\n", err)
		return
	}

	// Prompt for master password
	var masterPassword string
	if crypto.CheckMasterPasswordExists(baseDir + "/genp.yaml") {
		// Existing user - just prompt for password
		masterPassword, err = crypto.PromptForMasterPassword("Enter master password: ")
		if err != nil {
			color.Red("Failed to read master password: %v\n", err)
			return
		}
	} else {
		// First time - prompt with confirmation
		color.Yellow("First time setup: Please create a master password to encrypt your passwords.\n")
		masterPassword, err = crypto.PromptForMasterPasswordWithConfirm()
		if err != nil {
			color.Red("Failed to set master password: %v\n", err)
			return
		}
	}

	// Encrypt the password
	encryptedPassword, err := crypto.Encrypt(password, masterPassword)
	if err != nil {
		color.Red("Failed to encrypt password: %v\n", err)
		return
	}

	confPath, err := StoreLocalConfig(passwordName, encryptedPassword, OSName)
	if err != nil {
		color.Red("Failed to store password locally: %v\n", err)
		return
	}

	color.Green("Password encrypted and stored locally at: %s\n", confPath)
}
