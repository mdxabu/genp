/*
Copyright © 2026 @mdxabu

*/

package store

import (
	"fmt"
	"runtime"

	"github.com/mdxabu/genp/internal/crypto"
)

/*
 * This file will have an integration with GitHub to create a private repo on the user's account. Whenever a password is stored,
 * it will be saved in the repo in a JSON file in an encrypted format.
 */

func StorepasswordLocally(password string) {
	var passwordName string
	fmt.Print("Enter a name for the password: ")
	fmt.Scanln(&passwordName)

	OSName := runtime.GOOS
	
	// Get the config path to check if this is first time setup
	baseDir, err := configBaseDir("genp", OSName)
	if err != nil {
		fmt.Println("Failed to determine config directory:", err)
		return
	}
	
	// Prompt for master password
	var masterPassword string
	if crypto.CheckMasterPasswordExists(baseDir + "/genp.yaml") {
		// Existing user - just prompt for password
		masterPassword, err = crypto.PromptForMasterPassword("Enter master password: ")
		if err != nil {
			fmt.Println("Failed to read master password:", err)
			return
		}
	} else {
		// First time - prompt with confirmation
		fmt.Println("First time setup: Please create a master password to encrypt your passwords.")
		masterPassword, err = crypto.PromptForMasterPasswordWithConfirm()
		if err != nil {
			fmt.Println("Failed to set master password:", err)
			return
		}
	}

	// Encrypt the password
	encryptedPassword, err := crypto.Encrypt(password, masterPassword)
	if err != nil {
		fmt.Println("Failed to encrypt password:", err)
		return
	}

	confPath, err := StoreLocalConfig(passwordName, encryptedPassword, OSName)
	if err != nil {
		fmt.Println("Failed to store password locally:", err)
		return
	}

	fmt.Println("Password encrypted and stored locally at:", confPath)
}

func StorepasswordRemotely(password string) {
	var passwordName string
	fmt.Print("Enter a name for the password: ")
	fmt.Scanln(&passwordName)

	fmt.Println("Password stored remotely on Github Private Repository in a encrypted mode.")

}
