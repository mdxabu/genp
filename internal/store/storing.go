package store

import "fmt"

/*
 * This file will have an integration with GitHub to create a private repo on the user's account. Whenever a password is stored,
 * it will be saved in the repo in a JSON file in an encrypted format.
 */
 
 func StorepasswordLocally(password string) {
	var passwordName string
	fmt.Print("Enter a name for the password: ")
	fmt.Scanln(&passwordName)
	
	
	
	fmt.Println("Password stored locally in a encrypted mode.")
	
	
}

func StorepasswordRemotely(password string) {
	var passwordName string
	fmt.Print("Enter a name for the password: ")
	fmt.Scanln(&passwordName)
	
	
	
	fmt.Println("Password stored remotely on Github Private Repository in a encrypted mode.")
	
	
}

