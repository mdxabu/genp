/*
Copyright © 2026 @mdxabu

*/

package crypto

import (
	"strings"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	tests := []struct {
		name      string
		plaintext string
		password  string
	}{
		{
			name:      "Simple password",
			plaintext: "mySecretPassword123!",
			password:  "masterKey",
		},
		{
			name:      "Long password with special chars",
			plaintext: "p@ssw0rd!#$%^&*()_+-=[]{}|;:,.<>?",
			password:  "veryStrongMasterPassword123!@#",
		},
		{
			name:      "Unicode password",
			plaintext: "пароль密码🔐",
			password:  "мастер-ключ-密钥",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encrypt
			encrypted, err := Encrypt(tt.plaintext, tt.password)
			if err != nil {
				t.Fatalf("Encrypt failed: %v", err)
			}

			// Verify encrypted is not empty and is different from plaintext
			if encrypted == "" {
				t.Fatal("Encrypted string is empty")
			}
			if encrypted == tt.plaintext {
				t.Fatal("Encrypted string is the same as plaintext")
			}

			// Decrypt
			decrypted, err := Decrypt(encrypted, tt.password)
			if err != nil {
				t.Fatalf("Decrypt failed: %v", err)
			}

			// Verify decrypted matches original
			if decrypted != tt.plaintext {
				t.Fatalf("Decrypted text doesn't match original. Got %q, want %q", decrypted, tt.plaintext)
			}
		})
	}
}

func TestDecryptWithWrongPassword(t *testing.T) {
	plaintext := "mySecret"
	correctPassword := "correct"
	wrongPassword := "wrong"

	encrypted, err := Encrypt(plaintext, correctPassword)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	_, err = Decrypt(encrypted, wrongPassword)
	if err == nil {
		t.Fatal("Expected error when decrypting with wrong password, got nil")
	}
	if !strings.Contains(err.Error(), "failed to decrypt") {
		t.Fatalf("Expected 'failed to decrypt' error, got: %v", err)
	}
}

func TestEncryptEmptyPlaintext(t *testing.T) {
	_, err := Encrypt("", "password")
	if err == nil {
		t.Fatal("Expected error for empty plaintext, got nil")
	}
	if !strings.Contains(err.Error(), "plaintext cannot be empty") {
		t.Fatalf("Expected 'plaintext cannot be empty' error, got: %v", err)
	}
}

func TestEncryptEmptyPassword(t *testing.T) {
	_, err := Encrypt("plaintext", "")
	if err == nil {
		t.Fatal("Expected error for empty password, got nil")
	}
	if !strings.Contains(err.Error(), "password cannot be empty") {
		t.Fatalf("Expected 'password cannot be empty' error, got: %v", err)
	}
}

func TestDecryptEmptyData(t *testing.T) {
	_, err := Decrypt("", "password")
	if err == nil {
		t.Fatal("Expected error for empty encrypted data, got nil")
	}
	if !strings.Contains(err.Error(), "encrypted data cannot be empty") {
		t.Fatalf("Expected 'encrypted data cannot be empty' error, got: %v", err)
	}
}

func TestDecryptInvalidBase64(t *testing.T) {
	_, err := Decrypt("not-valid-base64!!!", "password")
	if err == nil {
		t.Fatal("Expected error for invalid base64, got nil")
	}
	if !strings.Contains(err.Error(), "failed to decode base64") {
		t.Fatalf("Expected 'failed to decode base64' error, got: %v", err)
	}
}

func TestDecryptTooShortData(t *testing.T) {
	// Create a base64 string that's too short
	shortData := "YWJjZA==" // "abcd" in base64, which is too short
	_, err := Decrypt(shortData, "password")
	if err == nil {
		t.Fatal("Expected error for too short data, got nil")
	}
	if !strings.Contains(err.Error(), "encrypted data is too short") {
		t.Fatalf("Expected 'encrypted data is too short' error, got: %v", err)
	}
}

func TestEncryptionProducesDifferentOutputs(t *testing.T) {
	plaintext := "testPassword"
	password := "masterKey"

	encrypted1, err := Encrypt(plaintext, password)
	if err != nil {
		t.Fatalf("First encryption failed: %v", err)
	}

	encrypted2, err := Encrypt(plaintext, password)
	if err != nil {
		t.Fatalf("Second encryption failed: %v", err)
	}

	// Due to random salt and nonce, two encryptions should produce different outputs
	if encrypted1 == encrypted2 {
		t.Fatal("Two encryptions of the same plaintext produced identical ciphertext (expected different due to random salt/nonce)")
	}

	// But both should decrypt to the same plaintext
	decrypted1, _ := Decrypt(encrypted1, password)
	decrypted2, _ := Decrypt(encrypted2, password)

	if decrypted1 != plaintext || decrypted2 != plaintext {
		t.Fatal("Decrypted values don't match original plaintext")
	}
}
