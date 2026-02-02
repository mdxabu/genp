/*
Copyright © 2026 @mdxabu

*/

package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

const (
	// SaltSize is the size of the salt in bytes
	SaltSize = 16
	// NonceSize is the size of the nonce in bytes for GCM
	NonceSize = 12
	// KeySize is the size of the encryption key in bytes (256 bits)
	KeySize = 32
	// Iterations for PBKDF2
	Iterations = 100000
)

// Encrypt encrypts the plaintext using AES-256-GCM with a password-derived key.
// The output is base64 encoded and contains: salt + nonce + ciphertext
func Encrypt(plaintext string, password string) (string, error) {
	if plaintext == "" {
		return "", errors.New("plaintext cannot be empty")
	}
	if password == "" {
		return "", errors.New("password cannot be empty")
	}

	// Generate a random salt
	salt := make([]byte, SaltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Derive key from password using PBKDF2
	key := pbkdf2.Key([]byte(password), salt, Iterations, KeySize, sha256.New)

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate nonce
	nonce := make([]byte, NonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt the plaintext
	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)

	// Combine salt + nonce + ciphertext
	result := append(salt, nonce...)
	result = append(result, ciphertext...)

	// Encode to base64
	return base64.StdEncoding.EncodeToString(result), nil
}

// Decrypt decrypts the base64-encoded ciphertext using AES-256-GCM with a password-derived key.
// The input should contain: salt + nonce + ciphertext
func Decrypt(encryptedData string, password string) (string, error) {
	if encryptedData == "" {
		return "", errors.New("encrypted data cannot be empty")
	}
	if password == "" {
		return "", errors.New("password cannot be empty")
	}

	// Decode from base64
	data, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	// Check minimum size: salt + nonce + at least some ciphertext
	minSize := SaltSize + NonceSize + 16 // 16 is the GCM tag size
	if len(data) < minSize {
		return "", errors.New("encrypted data is too short")
	}

	// Extract salt, nonce, and ciphertext
	salt := data[:SaltSize]
	nonce := data[SaltSize : SaltSize+NonceSize]
	ciphertext := data[SaltSize+NonceSize:]

	// Derive key from password using PBKDF2
	key := pbkdf2.Key([]byte(password), salt, Iterations, KeySize, sha256.New)

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Decrypt the ciphertext
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: incorrect password or corrupted data")
	}

	return string(plaintext), nil
}
