/*
Copyright 2025 - github.com/mdxabu
*/

package github

import (
	"bytes"
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
	githubAPIBase       = "https://api.github.com"
	githubDeviceCodeURL = "https://github.com/login/device/code"
	githubAccessToken   = "https://github.com/login/oauth/access_token"
	// GenP's OAuth App Client ID - users can replace this with their own
	defaultClientID = "GENP_GITHUB_CLIENT_ID"
	tokenFileName   = "github_token"
)

// TokenInfo stores the GitHub authentication token and metadata
type TokenInfo struct {
	Token     string `json:"token"`
	LoginType string `json:"login_type"` // "token" or "oauth"
	Username  string `json:"username"`
}

// DeviceCodeResponse represents the response from GitHub's device code endpoint
type DeviceCodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

// OAuthTokenResponse represents the response from GitHub's access token endpoint
type OAuthTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	Error       string `json:"error,omitempty"`
	ErrorDesc   string `json:"error_description,omitempty"`
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

// LoginWithOAuth initiates the OAuth device flow for GitHub authentication
func LoginWithOAuth(clientID string) (*TokenInfo, error) {
	if clientID == "" {
		clientID = os.Getenv("GENP_GITHUB_CLIENT_ID")
		if clientID == "" {
			return nil, fmt.Errorf("GitHub OAuth Client ID not set. Set GENP_GITHUB_CLIENT_ID environment variable or pass --client-id flag")
		}
	}

	// Step 1: Request device code
	deviceCode, err := requestDeviceCode(clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to request device code: %w", err)
	}

	fmt.Printf("\nGitHub Device Authentication\n")
	fmt.Printf("-------------------------------\n")
	fmt.Printf("1. Open this URL in your browser: %s\n", deviceCode.VerificationURI)
	fmt.Printf("2. Enter this code: %s\n", deviceCode.UserCode)
	fmt.Printf("-------------------------------\n")
	fmt.Printf("Waiting for authorization...\n\n")

	// Step 2: Poll for the access token
	token, err := pollForToken(clientID, deviceCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	// Step 3: Validate and get user info
	user, err := getAuthenticatedUser(token)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}

	info := &TokenInfo{
		Token:     token,
		LoginType: "oauth",
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

// requestDeviceCode initiates the OAuth device flow
func requestDeviceCode(clientID string) (*DeviceCodeResponse, error) {
	payload := map[string]string{
		"client_id": clientID,
		"scope":     "repo",
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", githubDeviceCodeURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to GitHub: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub returned status %d: %s", resp.StatusCode, string(body))
	}

	var deviceCode DeviceCodeResponse
	if err := json.Unmarshal(body, &deviceCode); err != nil {
		return nil, fmt.Errorf("failed to parse device code response: %w", err)
	}

	return &deviceCode, nil
}

// pollForToken polls GitHub for the access token after device code authorization
func pollForToken(clientID string, deviceCode *DeviceCodeResponse) (string, error) {
	interval := deviceCode.Interval
	if interval < 5 {
		interval = 5
	}

	deadline := time.Now().Add(time.Duration(deviceCode.ExpiresIn) * time.Second)

	for time.Now().Before(deadline) {
		time.Sleep(time.Duration(interval) * time.Second)

		payload := map[string]string{
			"client_id":   clientID,
			"device_code": deviceCode.DeviceCode,
			"grant_type":  "urn:ietf:params:oauth:grant-type:device_code",
		}

		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			return "", err
		}

		req, err := http.NewRequest("POST", githubAccessToken, bytes.NewBuffer(jsonPayload))
		if err != nil {
			return "", err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			continue // Retry on network errors
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			continue
		}

		var tokenResp OAuthTokenResponse
		if err := json.Unmarshal(body, &tokenResp); err != nil {
			continue
		}

		switch tokenResp.Error {
		case "":
			// Success
			if tokenResp.AccessToken != "" {
				return tokenResp.AccessToken, nil
			}
		case "authorization_pending":
			// User hasn't authorized yet, keep polling
			continue
		case "slow_down":
			// We're polling too fast, increase interval
			interval += 5
			continue
		case "expired_token":
			return "", fmt.Errorf("device code expired. Please try again")
		case "access_denied":
			return "", fmt.Errorf("authorization denied by user")
		default:
			return "", fmt.Errorf("OAuth error: %s - %s", tokenResp.Error, tokenResp.ErrorDesc)
		}
	}

	return "", fmt.Errorf("authorization timed out. Please try again")
}
