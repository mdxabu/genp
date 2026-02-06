package auth

import (
	"fmt"
	"os"

	"github.com/cli/oauth"
	"github.com/mdxabu/genp/internal/store"
)

const (
	githubHost = "github.com"
)

func getClientID() (string, error) {
	clientID := os.Getenv("GITHUB_CLIENT_ID")
	if clientID == "" {
		return "", fmt.Errorf("GITHUB_CLIENT_ID environment variable is not set.\n  Please set it before logging in:\n    export GITHUB_CLIENT_ID=your_client_id")
	}
	return clientID, nil
}

func Login() (string, error) {
	clientID, err := getClientID()
	if err != nil {
		return "", err
	}

	flow := &oauth.Flow{
		Hostname: githubHost,
		ClientID: clientID,
		Scopes:   []string{"repo"},
	}

	accessToken, err := flow.DeviceFlow()
	if err != nil {
		return "", err
	}

	token := accessToken.Token

	// Store the token in genp.yaml alongside the passwords
	if err := store.SaveToken(token); err != nil {
		return "", err
	}

	return token, nil
}

func Logout() error {
	return store.RemoveToken()
}

func GetToken() (string, error) {
	return store.GetToken()
}
