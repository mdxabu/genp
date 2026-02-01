package store

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// StoreLocalConfig creates a cross-platform config directory and writes a credentials file
// with restrictive permissions. It avoids OS/env shadowing and uses standard per-OS locations.
//
// - Windows: %LOCALAPPDATA%\genp\genp.yaml (fallback to %APPDATA% if LOCALAPPDATA is empty)
// - macOS: ~/Library/Application Support/genp/genp.yaml
// - Linux/Other Unix: $XDG_CONFIG_HOME/genp/genp.yaml (fallback to ~/.config/genp/genp.yaml)
//
// The file content is a simple "name=value" line. Caller is responsible for not providing plaintext
// secrets if that is inappropriate; prefer using OS keychain APIs for real secrets.

func StoreLocalConfig(passwordName string, password string, osName string) (string, error) {
	if passwordName == "" {
		return "", errors.New("passwordName must not be empty")
	}

	baseDir, err := configBaseDir("genp", osName)
	if err != nil {
		return "", err
	}

	// Ensure base config directory exists

	// Ensure directory exists with restrictive permissions where supported
	// Unix: 0700, Windows ACLs are handled by OS; mode is best-effort
	if err := os.MkdirAll(baseDir, 0o700); err != nil {
		return "", fmt.Errorf("failed to create config directory %s: %w", baseDir, err)
	}

	// Create genp.yaml file inside the base config directory
	confPath := filepath.Join(baseDir, "genp.yaml")

	// Prepare YAML content under a top-level "password:" map supporting multiple entries.
	// If the file exists and already has content, append a new entry under "password:".
	// This is a simple line-based approach without full YAML parsing.
	var content string
	existing, readErr := os.ReadFile(confPath)
	if readErr == nil {
		existingStr := string(existing)
		if existingStr == "" {
			content = fmt.Sprintf("password:\n  %s: %q\n", passwordName, password)
		} else {
			// Ensure the header exists
			if !strings.Contains(existingStr, "\npassword:\n") && !strings.HasPrefix(existingStr, "password:\n") {
				// Prepend the header if missing
				if existingStr[len(existingStr)-1] != '\n' {
					existingStr += "\n"
				}
				existingStr = "password:\n" + existingStr
			}
			// Ensure newline at the end, then append the new entry
			if existingStr[len(existingStr)-1] != '\n' {
				existingStr += "\n"
			}
			content = existingStr + fmt.Sprintf("  %s: %q\n", passwordName, password)
		}
	} else {
		// If the file doesn't exist or can't be read, start fresh
		content = fmt.Sprintf("password:\n  %s: %q\n", passwordName, password)
	}

	// Write file with restrictive permissions:
	// - Use os.WriteFile for brevity; set mode 0600 on Unix, best-effort on Windows.
	if err := os.WriteFile(confPath, []byte(content), 0o600); err != nil {
		return "", fmt.Errorf("failed to write config file %s: %w", confPath, err)
	}

	return confPath, nil
}

// configBaseDir determines the per-OS base config directory.
// appName should be a stable identifier for your application.
// Returns an absolute path like:
// - Windows: %LOCALAPPDATA%\appName
// - macOS: ~/Library/Application Support/appName
// - Linux/Unix: $XDG_CONFIG_HOME/appName or ~/.config/appName
func configBaseDir(appName string, osName string) (string, error) {
	switch osName {
	case "windows":
		// Prefer LOCALAPPDATA, fallback to APPDATA
		if local := os.Getenv("LOCALAPPDATA"); local != "" {
			return filepath.Join(local, appName), nil
		}
		if roaming := os.Getenv("APPDATA"); roaming != "" {
			return filepath.Join(roaming, appName), nil
		}
		// Fallback to user home if envs missing
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("cannot determine user home: %w", err)
		}
		return filepath.Join(home, "AppData", "Local", appName), nil

	case "darwin":
		// macOS: ~/Library/Application Support/appName
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("cannot determine user home: %w", err)
		}
		return filepath.Join(home, "Library", "Application Support", appName), nil

	default:
		// Unix/Linux: XDG_CONFIG_HOME/appName or ~/.config/appName
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
