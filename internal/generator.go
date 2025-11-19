package internal

import (
	"crypto/rand"
	"math/big"
)

const (
	lowercaseBytes = "abcdefghijklmnopqrstuvwxyz"
	uppercaseBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numberBytes    = "0123456789"
	specialBytes   = "!@#$&"
)

func GeneratePassword(length int, includeNumbers, includeUppercase, includeSpecial bool) string {
	charset := lowercaseBytes

	if includeUppercase {
		charset += uppercaseBytes
	}
	if includeNumbers {
		charset += numberBytes
	}
	if includeSpecial {
		charset += specialBytes
	}

	password := make([]byte, length)
	charsetLength := big.NewInt(int64(len(charset)))

	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			password[i] = charset[0]
			continue
		}
		password[i] = charset[randomIndex.Int64()]
	}

	return string(password)
}
