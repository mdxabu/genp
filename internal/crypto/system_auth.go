/*
Copyright 2025 - github.com/mdxabu
*/

package crypto

import (
	"fmt"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
)

// VerifySystemPassword verifies the given password against the operating system's
// user account password (the same password used to unlock the lock screen).
// Returns nil if the password is correct, or an error describing the failure.
func VerifySystemPassword(password string) error {
	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	switch runtime.GOOS {
	case "darwin":
		return verifyMacOS(password)
	case "linux":
		return verifyLinux(password)
	case "windows":
		return verifyWindows(password)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

// verifyMacOS uses dscl to authenticate the current user against the local directory.
// dscl . -authonly <username> <password> exits 0 on success, non-zero on failure.
func verifyMacOS(password string) error {
	currentUser, err := user.Current()
	if err != nil {
		return fmt.Errorf("failed to determine current user: %w", err)
	}

	cmd := exec.Command("dscl", ".", "-authonly", currentUser.Username, password)
	output, err := cmd.CombinedOutput()
	if err != nil {
		outStr := strings.TrimSpace(string(output))
		if outStr != "" {
			return fmt.Errorf("system password verification failed: %s", outStr)
		}
		return fmt.Errorf("system password verification failed: incorrect password")
	}

	return nil
}

// verifyLinux uses sudo to validate the current user's password.
// It pipes the password into `sudo -k -S -v` which validates without running a command.
// -k invalidates the cached credentials first so it always prompts.
// -S reads the password from stdin.
// -v updates the cached credentials (validates) without running a command.
func verifyLinux(password string) error {
	cmd := exec.Command("sudo", "-k", "-S", "-v")
	cmd.Stdin = strings.NewReader(password + "\n")
	output, err := cmd.CombinedOutput()
	if err != nil {
		outStr := strings.TrimSpace(string(output))
		if strings.Contains(outStr, "incorrect password") || strings.Contains(outStr, "Sorry") {
			return fmt.Errorf("system password verification failed: incorrect password")
		}
		if outStr != "" {
			return fmt.Errorf("system password verification failed: %s", outStr)
		}
		return fmt.Errorf("system password verification failed: incorrect password")
	}

	// Immediately invalidate the sudo timestamp so we don't leave an open sudo session
	_ = exec.Command("sudo", "-k").Run()

	return nil
}

// verifyWindows uses PowerShell and the .NET DirectoryServices to validate the
// current user's password against the local machine account database.
func verifyWindows(password string) error {
	currentUser, err := user.Current()
	if err != nil {
		return fmt.Errorf("failed to determine current user: %w", err)
	}

	// Extract just the username if it comes as DOMAIN\user
	username := currentUser.Username
	if parts := strings.SplitN(username, `\`, 2); len(parts) == 2 {
		username = parts[1]
	}

	// Use PowerShell to validate against local machine
	script := fmt.Sprintf(
		`Add-Type -AssemblyName System.DirectoryServices.AccountManagement; `+
			`$ctx = New-Object System.DirectoryServices.AccountManagement.PrincipalContext([System.DirectoryServices.AccountManagement.ContextType]::Machine); `+
			`$result = $ctx.ValidateCredentials('%s', '%s'); `+
			`if (-not $result) { exit 1 }`,
		escapePowerShellString(username),
		escapePowerShellString(password),
	)

	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		outStr := strings.TrimSpace(string(output))
		if outStr != "" {
			return fmt.Errorf("system password verification failed: %s", outStr)
		}
		return fmt.Errorf("system password verification failed: incorrect password")
	}

	return nil
}

// escapePowerShellString escapes single quotes in a string for safe use
// inside a PowerShell single-quoted string literal.
func escapePowerShellString(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}
