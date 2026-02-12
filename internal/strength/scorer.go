/*
Copyright 2025 - github.com/mdxabu
*/

package strength

import (
	"strings"
	"unicode"
)

// Level represents the password strength tier
type Level int

const (
	LevelEmpty Level = iota
	LevelPathetic
	LevelWeak
	LevelFair
	LevelDecent
	LevelStrong
	LevelFortKnox
)

// Result holds the evaluated password strength details
type Result struct {
	Score   int
	Max     int
	Level   Level
	Roast   string
	BarFill int // out of 10
	Color   string
}

// roasts maps each level to a list of meme-level roast messages
var roasts = map[Level][]string{
	LevelEmpty: {
		"Go on, type something... I dare you.",
		"The password field is lonelier than you on a Friday night.",
		"Even 'password' would be an upgrade right now.",
	},
	LevelPathetic: {
		"My grandma's WiFi password is stronger than this.",
		"A toddler mashing a keyboard would do better.",
		"This isn't a password, it's a cry for help.",
		"Hackers wouldn't even bother, it's too easy to be fun.",
		"Did you just slam your face on the keyboard? Try harder.",
	},
	LevelWeak: {
		"This password is on every hacker's 'first try' list.",
		"Bro, my pet goldfish could crack this.",
		"You locked the front door but left the window wide open.",
		"This password has the strength of wet tissue paper.",
		"Even a script kiddie is laughing at this.",
	},
	LevelFair: {
		"You're trying, I'll give you that... barely.",
		"It's giving 'minimum effort to pass the class'.",
		"Mediocre, like gas station sushi.",
		"Not terrible, but I wouldn't trust it with my Netflix.",
		"Your password is the participation trophy of security.",
	},
	LevelDecent: {
		"Okay okay, now we're getting somewhere.",
		"Solid. Like a B+ student who could try harder.",
		"Most hackers just sighed and moved on.",
		"Your password actually has a backbone now.",
		"Decent enough to survive a coffee shop WiFi.",
	},
	LevelStrong: {
		"Fort Knox called, they want your password strategy.",
		"Hackers saw this and chose a different career.",
		"This password drinks black coffee and deadlifts.",
		"Even the NSA just raised an eyebrow. Respect.",
		"Your password has more layers than an onion.",
	},
	LevelFortKnox: {
		"This password doesn't need a firewall, it IS the firewall.",
		"Brute force? More like brute FARCE against this beast.",
		"CIA, FBI, KGB... none of them are getting in.",
		"You didn't create a password, you created a LEGEND.",
		"This password bench presses other passwords for fun.",
	},
}

// Evaluate scores a password and returns a Result with roast message
func Evaluate(password string) Result {
	if len(password) == 0 {
		return Result{
			Score:   0,
			Max:     100,
			Level:   LevelEmpty,
			Roast:   pickRoast(LevelEmpty, 0),
			BarFill: 0,
			Color:   "red",
		}
	}

	score := 0

	// Length scoring (up to 30 points)
	length := len(password)
	switch {
	case length >= 20:
		score += 30
	case length >= 16:
		score += 25
	case length >= 12:
		score += 20
	case length >= 8:
		score += 15
	case length >= 6:
		score += 8
	case length >= 4:
		score += 4
	default:
		score += 1
	}

	// Character variety (up to 40 points)
	hasLower := false
	hasUpper := false
	hasDigit := false
	hasSpecial := false

	for _, ch := range password {
		switch {
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsDigit(ch):
			hasDigit = true
		default:
			hasSpecial = true
		}
	}

	varieties := 0
	if hasLower {
		score += 10
		varieties++
	}
	if hasUpper {
		score += 10
		varieties++
	}
	if hasDigit {
		score += 10
		varieties++
	}
	if hasSpecial {
		score += 10
		varieties++
	}

	// Bonus for using all 4 character types (10 points)
	if varieties == 4 {
		score += 10
	}

	// Penalty for common patterns (up to -20)
	score -= penaltyForPatterns(password)

	// Clamp score
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	level := scoreToLevel(score)
	barFill := score / 10
	if score > 0 && barFill == 0 {
		barFill = 1
	}

	return Result{
		Score:   score,
		Max:     100,
		Level:   level,
		Roast:   pickRoast(level, length),
		BarFill: barFill,
		Color:   levelColor(level),
	}
}

func scoreToLevel(score int) Level {
	switch {
	case score >= 90:
		return LevelFortKnox
	case score >= 70:
		return LevelStrong
	case score >= 55:
		return LevelDecent
	case score >= 40:
		return LevelFair
	case score >= 20:
		return LevelWeak
	default:
		return LevelPathetic
	}
}

func levelColor(level Level) string {
	switch level {
	case LevelEmpty, LevelPathetic:
		return "red"
	case LevelWeak:
		return "yellow"
	case LevelFair:
		return "yellow"
	case LevelDecent:
		return "cyan"
	case LevelStrong:
		return "green"
	case LevelFortKnox:
		return "green"
	default:
		return "white"
	}
}

func pickRoast(level Level, passwordLen int) string {
	messages := roasts[level]
	if len(messages) == 0 {
		return ""
	}
	// Pick based on password length to give variety as user types
	idx := passwordLen % len(messages)
	return messages[idx]
}

func penaltyForPatterns(password string) int {
	penalty := 0
	lower := strings.ToLower(password)

	// Check for common weak passwords
	commonPasswords := []string{
		"password", "123456", "qwerty", "abc123", "letmein",
		"admin", "welcome", "monkey", "master", "dragon",
		"login", "princess", "football", "shadow", "sunshine",
		"trustno1", "iloveyou", "batman", "access", "hello",
		"charlie", "donald", "passw0rd",
	}
	for _, common := range commonPasswords {
		if lower == common {
			penalty += 20
			break
		}
	}

	// Repeated characters (e.g., "aaaa")
	if len(password) > 2 {
		repeats := 0
		for i := 1; i < len(password); i++ {
			if password[i] == password[i-1] {
				repeats++
			}
		}
		if float64(repeats) > float64(len(password))*0.5 {
			penalty += 10
		}
	}

	// Sequential characters (e.g., "abcd", "1234")
	if len(password) > 3 {
		sequential := 0
		for i := 2; i < len(password); i++ {
			if password[i]-password[i-1] == 1 && password[i-1]-password[i-2] == 1 {
				sequential++
			}
		}
		if sequential > 2 {
			penalty += 10
		}
	}

	// All same case with no variety
	if len(password) > 4 {
		allSame := true
		for i := 1; i < len(password); i++ {
			if password[i] != password[0] {
				allSame = false
				break
			}
		}
		if allSame {
			penalty += 15
		}
	}

	return penalty
}

// LevelLabel returns a human-readable label for the strength level
func LevelLabel(level Level) string {
	switch level {
	case LevelEmpty:
		return "EMPTY"
	case LevelPathetic:
		return "PATHETIC"
	case LevelWeak:
		return "WEAK"
	case LevelFair:
		return "FAIR"
	case LevelDecent:
		return "DECENT"
	case LevelStrong:
		return "STRONG"
	case LevelFortKnox:
		return "FORT KNOX"
	default:
		return "UNKNOWN"
	}
}
