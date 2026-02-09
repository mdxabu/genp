/*
Copyright © 2025 - github.com/mdxabu
*/

package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	// AppName is the application identifier used for config directories
	AppName = "genp"
	// ConfigFileName is the name of the main config file
	ConfigFileName = "genp.yaml"
	// GitHubTokenFileName is the name of the GitHub token file
	GitHubTokenFileName = "github_token"
)

// BaseDir determines the per-OS base config directory.
// Returns an absolute path like:
//   - Windows: %LOCALAPPDATA%\genp
//   - macOS: ~/Library/Application Support/genp
//   - Linux/Unix: $XDG_CONFIG_HOME/genp or ~/.config/genp
func BaseDir(osName string) (string, error) {
	switch osName {
	case "windows":
		if local := os.Getenv("LOCALAPPDATA"); local != "" {
			return filepath.Join(local, AppName), nil
		}
		if roaming := os.Getenv("APPDATA"); roaming != "" {
			return filepath.Join(roaming, AppName), nil
		}
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("cannot determine user home: %w", err)
		}
		return filepath.Join(home, "AppData", "Local", AppName), nil

	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("cannot determine user home: %w", err)
		}
		return filepath.Join(home, "Library", "Application Support", AppName), nil

	default:
		if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
			return filepath.Join(xdg, AppName), nil
		}
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("cannot determine user home: %w", err)
		}
		return filepath.Join(home, ".config", AppName), nil
	}
}

// BaseDirForApp determines the per-OS base config directory for an arbitrary app name.
// This is useful for testing with a different app name.
func BaseDirForApp(appName string, osName string) (string, error) {
	switch osName {
	case "windows":
		if local := os.Getenv("LOCALAPPDATA"); local != "" {
			return filepath.Join(local, appName), nil
		}
		if roaming := os.Getenv("APPDATA"); roaming != "" {
			return filepath.Join(roaming, appName), nil
		}
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("cannot determine user home: %w", err)
		}
		return filepath.Join(home, "AppData", "Local", appName), nil

	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("cannot determine user home: %w", err)
		}
		return filepath.Join(home, "Library", "Application Support", appName), nil

	default:
		if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
			return filepath.Join(xdg, appName), nil
		}
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("cannot determine user home: %w", err)
		}
		return filepath.Join(home, ".config", appName), nil
	}
}

// ConfigFilePath returns the full path to genp.yaml for the given OS
func ConfigFilePath(osName string) (string, error) {
	baseDir, err := BaseDir(osName)
	if err != nil {
		return "", err
	}
	return filepath.Join(baseDir, ConfigFileName), nil
}

// GitHubTokenPath returns the full path to the GitHub token file for the given OS
func GitHubTokenPath(osName string) (string, error) {
	baseDir, err := BaseDir(osName)
	if err != nil {
		return "", err
	}
	return filepath.Join(baseDir, GitHubTokenFileName), nil
}
