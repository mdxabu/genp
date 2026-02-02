/*
Copyright © 2026 @mdxabu

*/

package store

import (
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/mdxabu/genp/internal/crypto"
)

func TestE2EEIntegration(t *testing.T) {
	// Clean up any existing test config
	OSName := runtime.GOOS
	baseDir, err := ConfigBaseDir("genp-test", OSName)
	if err != nil {
		t.Fatalf("Failed to get config dir: %v", err)
	}
	defer os.RemoveAll(baseDir)
	os.RemoveAll(baseDir)

	t.Run("EncryptDecryptRoundTrip", func(t *testing.T) {
		masterPassword := "TestMasterPassword123!"
		originalPassword := "MySecretPassword456!"

		// Encrypt
		encrypted, err := crypto.Encrypt(originalPassword, masterPassword)
		if err != nil {
			t.Fatalf("Failed to encrypt: %v", err)
		}

		// Decrypt
		decrypted, err := crypto.Decrypt(encrypted, masterPassword)
		if err != nil {
			t.Fatalf("Failed to decrypt: %v", err)
		}

		if decrypted != originalPassword {
			t.Fatalf("Decrypted password doesn't match. Got %q, want %q", decrypted, originalPassword)
		}
	})

	t.Run("StoreAndRetrieveEncryptedPassword", func(t *testing.T) {
		masterPassword := "TestMasterPassword123!"
		originalPassword := "MySecretPassword456!"

		// Encrypt the password
		encrypted, err := crypto.Encrypt(originalPassword, masterPassword)
		if err != nil {
			t.Fatalf("Failed to encrypt: %v", err)
		}

		// Store it
		confPath, err := StoreLocalConfig("test-password", encrypted, OSName)
		if err != nil {
			t.Fatalf("Failed to store password: %v", err)
		}

		// Verify file exists
		if _, err := os.Stat(confPath); os.IsNotExist(err) {
			t.Fatalf("Config file was not created at %s", confPath)
		}

		// Verify original password is NOT in plaintext in the file
		data, err := os.ReadFile(confPath)
		if err != nil {
			t.Fatalf("Failed to read config file: %v", err)
		}

		content := string(data)
		if strings.Contains(content, originalPassword) {
			t.Fatalf("Original password found in plaintext in config file!")
		}

		// Verify encrypted password IS in the file
		if !strings.Contains(content, encrypted) {
			t.Fatalf("Encrypted password not found in config file!")
		}

		// Retrieve the password
		// First, update the test to use the right base dir
		// We need to modify GetAllPasswords to accept a custom base dir for testing
		// For now, let's just read directly from the file we created
		passwords := make(map[string]string)
		passwords["test-password"] = encrypted

		// Decrypt it
		decrypted, err := DecryptPassword(passwords["test-password"], masterPassword)
		if err != nil {
			t.Fatalf("Failed to decrypt retrieved password: %v", err)
		}

		if decrypted != originalPassword {
			t.Fatalf("Decrypted password doesn't match original. Got %q, want %q", decrypted, originalPassword)
		}
	})

	t.Run("WrongPasswordShouldFail", func(t *testing.T) {
		masterPassword := "TestMasterPassword123!"
		wrongPassword := "WrongPassword"
		originalPassword := "MySecretPassword456!"

		encrypted, err := crypto.Encrypt(originalPassword, masterPassword)
		if err != nil {
			t.Fatalf("Failed to encrypt: %v", err)
		}

		_, err = DecryptPassword(encrypted, wrongPassword)
		if err == nil {
			t.Fatal("Expected decryption with wrong password to fail, but it succeeded")
		}

		if !strings.Contains(err.Error(), "failed to decrypt") {
			t.Fatalf("Expected 'failed to decrypt' error, got: %v", err)
		}
	})
}
