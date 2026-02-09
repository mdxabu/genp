/*
Copyright 2025 - github.com/mdxabu
*/

package github

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	vaultRepoName  = "genp-vault"
	vaultFileName  = "genp.yaml"
	maxRetries     = 3
	initialBackoff = 2 * time.Second
)

// RepoInfo represents basic GitHub repository information
type RepoInfo struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Private  bool   `json:"private"`
	HTMLURL  string `json:"html_url"`
}

// FileContent represents GitHub API file content response
type FileContent struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	SHA     string `json:"sha"`
	Content string `json:"content"`
}

// isRetryableStatus returns true if the HTTP status code indicates a transient
// server-side error that is worth retrying.
func isRetryableStatus(statusCode int) bool {
	switch statusCode {
	case http.StatusBadGateway, // 502
		http.StatusServiceUnavailable,  // 503
		http.StatusGatewayTimeout,      // 504
		http.StatusInternalServerError: // 500
		return true
	default:
		return false
	}
}

// doRequestWithRetry executes an HTTP request with exponential backoff retries
// for transient 5xx errors.  It returns the response body bytes, the final
// status code, and any hard error.  The caller is responsible for interpreting
// the status code after retries are exhausted.
func doRequestWithRetry(buildReq func() (*http.Request, error)) ([]byte, int, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	backoff := initialBackoff

	var lastBody []byte
	var lastStatus int

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(backoff)
			backoff *= 2
		}

		req, err := buildReq()
		if err != nil {
			return nil, 0, fmt.Errorf("failed to build request: %w", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			// Network-level errors are retryable
			if attempt < maxRetries {
				continue
			}
			return nil, 0, fmt.Errorf("failed to connect to GitHub API after %d attempts: %w", maxRetries+1, err)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			if attempt < maxRetries {
				continue
			}
			return nil, resp.StatusCode, fmt.Errorf("failed to read response: %w", err)
		}

		lastBody = body
		lastStatus = resp.StatusCode

		if !isRetryableStatus(resp.StatusCode) {
			return body, resp.StatusCode, nil
		}

		// Retryable status – log and loop
		if attempt < maxRetries {
			fmt.Printf("  [retry] GitHub returned %d, retrying in %v (%d/%d)...\n",
				resp.StatusCode, backoff, attempt+1, maxRetries)
		}
	}

	return lastBody, lastStatus, nil
}

// CreateOrGetVaultRepo ensures the genp-vault private repo exists on the user's GitHub account.
// If it already exists, it returns the existing repo info. Otherwise, it creates a new one.
func CreateOrGetVaultRepo(token string) (*RepoInfo, error) {
	// First, check if the repo already exists
	repo, err := getRepo(token, vaultRepoName)
	if err == nil {
		return repo, nil
	}

	// Repo doesn't exist, create it
	return createRepo(token, vaultRepoName)
}

// SyncConfigToVault pushes the local genp.yaml file to the genp-vault GitHub repo.
// It handles both creating and updating the file.
func SyncConfigToVault(configPath string) error {
	tokenInfo, err := LoadToken()
	if err != nil {
		// Not logged in, skip sync silently
		return nil
	}

	// Read the local config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// Ensure the vault repo exists
	_, err = CreateOrGetVaultRepo(tokenInfo.Token)
	if err != nil {
		return fmt.Errorf("failed to ensure vault repo exists: %w", err)
	}

	// Push the file to the repo
	return pushFile(tokenInfo.Token, tokenInfo.Username, vaultRepoName, vaultFileName, data)
}

// SyncConfigToVaultIfLoggedIn is a convenience wrapper that only syncs if the user
// is logged in to GitHub. Returns nil if not logged in (non-blocking).
func SyncConfigToVaultIfLoggedIn(configPath string) error {
	if !IsLoggedIn() {
		return nil
	}
	return SyncConfigToVault(configPath)
}

// getRepo fetches repository info from GitHub
func getRepo(token string, repoName string) (*RepoInfo, error) {
	// Get the authenticated user first to build the URL
	user, err := getAuthenticatedUser(token)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/repos/%s/%s", githubAPIBase, user.Login, repoName)

	body, status, err := doRequestWithRetry(func() (*http.Request, error) {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
		return req, nil
	})
	if err != nil {
		return nil, err
	}

	if status == http.StatusNotFound {
		return nil, fmt.Errorf("repository %s not found", repoName)
	}

	if status != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d: %s", status, string(body))
	}

	var repo RepoInfo
	if err := json.Unmarshal(body, &repo); err != nil {
		return nil, fmt.Errorf("failed to parse repo info: %w", err)
	}

	return &repo, nil
}

// createRepo creates a new private repository on GitHub
func createRepo(token string, repoName string) (*RepoInfo, error) {
	payload := map[string]interface{}{
		"name":        repoName,
		"description": "GenP password vault - encrypted password storage",
		"private":     true,
		"auto_init":   true, // Initialize with a README so we have a default branch
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	body, status, err := doRequestWithRetry(func() (*http.Request, error) {
		req, err := http.NewRequest("POST", githubAPIBase+"/user/repos", bytes.NewBuffer(jsonPayload))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
		return req, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create repository: %w", err)
	}

	if status == http.StatusCreated {
		var repo RepoInfo
		if err := json.Unmarshal(body, &repo); err != nil {
			return nil, fmt.Errorf("failed to parse repo info: %w", err)
		}
		return &repo, nil
	}

	// Check if repo already exists (422 = name already taken / race condition)
	if status == http.StatusUnprocessableEntity {
		existing, getErr := getRepo(token, repoName)
		if getErr == nil {
			return existing, nil
		}
	}

	return nil, fmt.Errorf("failed to create repository (status %d): %s", status, string(body))
}

// pushFile creates or updates a file in a GitHub repository using the Contents API
func pushFile(token, owner, repoName, filePath string, content []byte) error {
	url := fmt.Sprintf("%s/repos/%s/%s/contents/%s", githubAPIBase, owner, repoName, filePath)

	// Check if the file already exists to get its SHA (needed for updates)
	existingSHA, err := getFileSHA(token, owner, repoName, filePath)
	if err != nil {
		// File doesn't exist yet, that's fine - we'll create it
		existingSHA = ""
	}

	// Build the payload
	payloadMap := map[string]interface{}{
		"message": fmt.Sprintf("vault: sync %s", filePath),
		"content": base64.StdEncoding.EncodeToString(content),
	}

	if existingSHA != "" {
		payloadMap["sha"] = existingSHA
	}

	jsonPayload, err := json.Marshal(payloadMap)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	body, status, err := doRequestWithRetry(func() (*http.Request, error) {
		req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonPayload))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
		return req, nil
	})
	if err != nil {
		return fmt.Errorf("failed to push file: %w", err)
	}

	// 200 = updated, 201 = created
	if status != http.StatusOK && status != http.StatusCreated {
		return fmt.Errorf("failed to push file (status %d): %s", status, string(body))
	}

	return nil
}

// getFileSHA retrieves the SHA of an existing file in a GitHub repository.
// Returns an error if the file does not exist.
func getFileSHA(token, owner, repoName, filePath string) (string, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/contents/%s", githubAPIBase, owner, repoName, filePath)

	body, status, err := doRequestWithRetry(func() (*http.Request, error) {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
		return req, nil
	})
	if err != nil {
		return "", err
	}

	if status == http.StatusNotFound {
		return "", fmt.Errorf("file not found")
	}

	if status != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned status %d: %s", status, string(body))
	}

	var fileContent FileContent
	if err := json.Unmarshal(body, &fileContent); err != nil {
		return "", fmt.Errorf("failed to parse file content: %w", err)
	}

	return fileContent.SHA, nil
}

// PullConfigFromVault downloads the genp.yaml from the vault repo and returns its content.
// This can be used to restore passwords from the cloud backup.
func PullConfigFromVault(token, username string) ([]byte, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/contents/%s", githubAPIBase, username, vaultRepoName, vaultFileName)

	body, status, err := doRequestWithRetry(func() (*http.Request, error) {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
		return req, nil
	})
	if err != nil {
		return nil, err
	}

	if status == http.StatusNotFound {
		return nil, fmt.Errorf("vault file not found in genp-vault repo")
	}

	if status != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d: %s", status, string(body))
	}

	var fileContent FileContent
	if err := json.Unmarshal(body, &fileContent); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Decode base64 content
	decoded, err := base64.StdEncoding.DecodeString(fileContent.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to decode file content: %w", err)
	}

	return decoded, nil
}
