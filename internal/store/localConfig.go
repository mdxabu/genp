/*
Copyright © 2026 @mdxabu

*/

package store

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mdxabu/genp/internal/config"
	"github.com/mdxabu/genp/internal/crypto"
	"gopkg.in/yaml.v3"
)

// ConfigFile represents the top-level structure of genp.yaml
type ConfigFile struct {
	Password map[string]string `yaml:"password"`
}

// StoreLocalConfig creates a cross-platform config directory and writes a credentials file
// with restrictive permissions. It avoids OS/env shadowing and uses standard per-OS locations.
//
// - Windows: %LOCALAPPDATA%\genp\genp.yaml (fallback to %APPDATA% if LOCALAPPDATA is empty)
// - macOS: ~/Library/Application Support/genp/genp.yaml
// - Linux/Other Unix: $XDG_CONFIG_HOME/genp/genp.yaml (fallback to ~/.config/genp/genp.yaml)
//
// The file uses proper YAML marshaling to avoid duplicate key issues.

func StoreLocalConfig(passwordName string, password string, osName string) (string, error) {
	if passwordName == "" {
		return "", errors.New("passwordName must not be empty")
	}

	baseDir, err := ConfigBaseDir("genp", osName)
	if err != nil {
		return "", err
	}

	// Ensure directory exists with restrictive permissions where supported
	// Unix: 0700, Windows ACLs are handled by OS; mode is best-effort
	if err := os.MkdirAll(baseDir, 0o700); err != nil {
		return "", fmt.Errorf("failed to create config directory %s: %w", baseDir, err)
	}

	confPath := filepath.Join(baseDir, "genp.yaml")

	// Load existing config or create a new one
	cfg, err := loadConfigFile(confPath)
	if err != nil {
		return "", fmt.Errorf("failed to load existing config: %w", err)
	}

	// Add or update the password entry
	cfg.Password[passwordName] = password

	// Marshal back to YAML
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return "", fmt.Errorf("failed to marshal config to YAML: %w", err)
	}

	// Write file with restrictive permissions
	if err := os.WriteFile(confPath, data, 0o600); err != nil {
		return "", fmt.Errorf("failed to write config file %s: %w", confPath, err)
	}

	return confPath, nil
}

// loadConfigFile reads and parses the genp.yaml file.
// If the file does not exist, it returns a new empty config.
// If the file has duplicate YAML keys (from older versions of genp),
// it falls back to a line-based dedup parser that keeps the last value for each key.
func loadConfigFile(confPath string) (*ConfigFile, error) {
	cfg := &ConfigFile{
		Password: make(map[string]string),
	}

	data, err := os.ReadFile(confPath)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("failed to read config file %s: %w", confPath, err)
	}

	if len(data) == 0 {
		return cfg, nil
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		// If standard YAML parsing fails (e.g. duplicate keys from old genp versions),
		// fall back to a line-based dedup parser that keeps the last value per key.
		dedupedCfg, dedupErr := parseDuplicateKeyYAML(data)
		if dedupErr != nil {
			// Return the original YAML error if fallback also fails
			return nil, fmt.Errorf("failed to parse config file %s: %w", confPath, err)
		}

		// Repair the file on disk so future reads don't hit this path
		repaired, marshalErr := yaml.Marshal(dedupedCfg)
		if marshalErr == nil {
			_ = os.WriteFile(confPath, repaired, 0o600)
		}

		return dedupedCfg, nil
	}

	// Ensure the map is initialized even if YAML had no password entries
	if cfg.Password == nil {
		cfg.Password = make(map[string]string)
	}

	return cfg, nil
}

// parseDuplicateKeyYAML handles YAML files that have duplicate mapping keys
// (produced by older versions of genp that used string concatenation).
// It parses line-by-line under the "password:" section and keeps the last
// value for each key, deduplicating them.
func parseDuplicateKeyYAML(data []byte) (*ConfigFile, error) {
	cfg := &ConfigFile{
		Password: make(map[string]string),
	}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	inPasswordSection := false

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Skip empty lines and comments
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		// Detect top-level "password:" header
		if trimmed == "password:" {
			inPasswordSection = true
			continue
		}

		// If we hit another top-level key (no leading whitespace), leave password section
		if !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t") && strings.Contains(trimmed, ":") {
			inPasswordSection = false
			continue
		}

		// Parse indented entries under "password:"
		if inPasswordSection {
			colonIdx := strings.Index(trimmed, ":")
			if colonIdx > 0 {
				key := strings.TrimSpace(trimmed[:colonIdx])
				val := strings.TrimSpace(trimmed[colonIdx+1:])
				// Remove surrounding quotes if present
				val = strings.Trim(val, "\"")
				if key != "" && val != "" {
					// Last value wins for duplicate keys
					cfg.Password[key] = val
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan config file: %w", err)
	}

	return cfg, nil
}

// GetConfigFilePath returns the full path to genp.yaml for the current OS
func GetConfigFilePath() (string, error) {
	osName := runtime.GOOS
	return config.ConfigFilePath(osName)
}

// ConfigBaseDir determines the per-OS base config directory.
// appName should be a stable identifier for your application.
// This delegates to the shared config package.
func ConfigBaseDir(appName string, osName string) (string, error) {
	return config.BaseDirForApp(appName, osName)
}

// PasswordEntry represents a stored password entry
type PasswordEntry struct {
	Name      string
	Encrypted string
}

// GetAllPasswords reads all stored passwords from the config file
// Returns a map of password names to their encrypted values
func GetAllPasswords() (map[string]string, error) {
	confPath, err := GetConfigFilePath()
	if err != nil {
		return nil, fmt.Errorf("failed to determine config file path: %w", err)
	}

	// Check if file exists
	if _, err := os.Stat(confPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("no passwords stored yet. Config file does not exist at: %s", confPath)
	}

	cfg, err := loadConfigFile(confPath)
	if err != nil {
		return nil, err
	}

	if len(cfg.Password) == 0 {
		return nil, fmt.Errorf("no passwords found in config file")
	}

	return cfg.Password, nil
}

// DecryptPassword decrypts a single password using the master password
func DecryptPassword(encryptedPassword string, masterPassword string) (string, error) {
	return crypto.Decrypt(encryptedPassword, masterPassword)
}
