/*
Copyright 2025 - github.com/mdxabu
*/

package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/mdxabu/genp/internal/config"
)

const (
	githubAPIBase = "https://api.github.com"
	tokenFileName = "github_token"
)

// TokenInfo stores the GitHub authentication token and metadata
type TokenInfo struct {
	Token     string `json:"token"`
	LoginType string `json:"login_type"`
	Username  string `json:"username"`
}

// GitHubUser represents basic GitHub user info
type GitHubUser struct {
	Login string `json:"login"`
	Name  string `json:"name"`
}

// getTokenFilePath returns the path to the stored GitHub token file
func getTokenFilePath() (string, error) {
	osName := runtime.GOOS
	return config.GitHubTokenPath(osName)
}

// GetTokenStorePath returns the path where the GitHub token is stored on disk.
func GetTokenStorePath() (string, error) {
	return getTokenFilePath()
}

// LoginWithToken authenticates using a personal access token
func LoginWithToken(token string) (*TokenInfo, error) {
	// Validate the token by making an API call
	user, err := getAuthenticatedUser(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	info := &TokenInfo{
		Token:     token,
		LoginType: "token",
		Username:  user.Login,
	}

	// Save the token
	if err := saveToken(info); err != nil {
		return nil, fmt.Errorf("failed to save token: %w", err)
	}

	return info, nil
}

// Logout removes the stored GitHub token
func Logout() error {
	tokenPath, err := getTokenFilePath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(tokenPath); os.IsNotExist(err) {
		return fmt.Errorf("not logged in to GitHub")
	}

	return os.Remove(tokenPath)
}

// LoadToken reads the stored GitHub token from disk
func LoadToken() (*TokenInfo, error) {
	tokenPath, err := getTokenFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(tokenPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("not logged in to GitHub. Run 'genp login' first")
		}
		return nil, fmt.Errorf("failed to read token file: %w", err)
	}

	var info TokenInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, fmt.Errorf("failed to parse token file: %w", err)
	}

	return &info, nil
}

// IsLoggedIn checks if the user is currently logged in to GitHub
func IsLoggedIn() bool {
	_, err := LoadToken()
	return err == nil
}

// saveToken writes the token info to disk
func saveToken(info *TokenInfo) error {
	tokenPath, err := getTokenFilePath()
	if err != nil {
		return err
	}

	// Ensure directory exists
	dir := filepath.Dir(tokenPath)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	return os.WriteFile(tokenPath, data, 0o600)
}

// getAuthenticatedUser fetches the authenticated user's info from GitHub API
func getAuthenticatedUser(token string) (*GitHubUser, error) {
	req, err := http.NewRequest("GET", githubAPIBase+"/user", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to GitHub API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d: %s", resp.StatusCode, string(body))
	}

	var user GitHubUser
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	return &user, nil
}
