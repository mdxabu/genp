/*
Copyright 2025 - github.com/mdxabu
*/

package strength

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// RunInteractive starts the real-time password strength checker.
// It reads each keystroke in raw terminal mode and updates the
// strength bar and roast message live as the user types.
func RunInteractive() (string, error) {
	fd := int(os.Stdin.Fd())

	// Save and restore terminal state
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return "", fmt.Errorf("failed to set raw terminal mode: %w", err)
	}
	defer term.Restore(fd, oldState)

	password := []byte{}
	buf := make([]byte, 1)

	// Print the initial prompt and empty strength line
	writeStr("\r\033[K\033[1;37mEnter the password: \033[0m")
	writeStr("\r\n")
	renderStrength("", 0)
	// Move cursor back up to the password input line, after the prompt
	writeStr("\033[1A\033[20C")

	for {
		_, err := os.Stdin.Read(buf)
		if err != nil {
			return string(password), err
		}

		ch := buf[0]

		switch {
		// Enter key
		case ch == '\r' || ch == '\n':
			// Move to end, print newlines to exit cleanly
			writeStr("\r\n\r\n")
			return string(password), nil

		// Ctrl+C / Ctrl+D
		case ch == 3 || ch == 4:
			writeStr("\r\n\r\n")
			return "", fmt.Errorf("cancelled")

		// Backspace / Delete
		case ch == 127 || ch == 8:
			if len(password) > 0 {
				password = password[:len(password)-1]
			}

		// Ignore other control characters
		case ch < 32:
			continue

		// Normal printable character
		default:
			password = append(password, ch)
		}

		// Redraw the password line (masked)
		redraw(password)
	}
}

func redraw(password []byte) {
	masked := strings.Repeat("*", len(password))

	// Move to the password line, clear it, rewrite
	writeStr("\r\033[K")
	writeStr(fmt.Sprintf("\033[1;37mEnter the password: \033[1;36m%s\033[0m", masked))

	// Save cursor position on password line
	writeStr("\033[s")

	// Move down to the strength line, clear it, render strength
	writeStr("\r\n")
	renderStrength(string(password), len(password))

	// Restore cursor position back to password line
	writeStr("\033[u")
}

func renderStrength(password string, length int) {
	writeStr("\r\033[K")

	if length == 0 {
		writeStr("\033[2;37m  [----------] Type something...\033[0m")
		return
	}

	result := Evaluate(password)

	// Build the bar
	filled := result.BarFill
	if filled > 10 {
		filled = 10
	}
	empty := 10 - filled

	barColor := ansiColor(result.Color)
	bar := fmt.Sprintf("%s%s\033[0m%s",
		barColor,
		strings.Repeat("=", filled),
		strings.Repeat("-", empty),
	)

	label := LevelLabel(result.Level)
	labelColored := fmt.Sprintf("%s%s\033[0m", barColor, label)

	writeStr(fmt.Sprintf("  [%s] %s - %s", bar, labelColored, result.Roast))
}

func ansiColor(color string) string {
	switch color {
	case "red":
		return "\033[1;31m"
	case "yellow":
		return "\033[1;33m"
	case "cyan":
		return "\033[1;36m"
	case "green":
		return "\033[1;32m"
	case "white":
		return "\033[1;37m"
	default:
		return "\033[0m"
	}
}

func writeStr(s string) {
	os.Stdout.WriteString(s)
}
