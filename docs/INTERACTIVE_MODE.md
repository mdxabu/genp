# GenP Interactive Mode

## Overview

GenP now features an interactive mode with a Claude Code-like command palette interface. This provides a modern, user-friendly way to interact with the password manager.

## Starting Interactive Mode

```bash
genp interactive
# or use the short alias
genp i
```

## UI Preview

### Initial Screen

When you start interactive mode, you'll see:

```
Starting GenP Interactive Mode...
Press Ctrl+C to exit

GenP Interactive Mode
                     
Type / to see available commands or type a command directly

> █
```

### Command Palette (typing `/`)

When you type `/`, the command palette appears:

```
GenP Interactive Mode
                     
Type / to see available commands or type a command directly

> /█

╭───────────────────────────────────────╮
│                                       │
│ Available Commands:                   │
│                                       │
│   ▶ /create - Generate a new password │
│     /show - Display stored passwords  │
│     /help - Show available commands   │
│     /exit - Exit interactive mode     │
│                                       │
│                                       │
╰───────────────────────────────────────╯

Use ↑↓ to navigate, Enter to select, Esc to cancel
```

### Navigation

- **Arrow Keys (↑↓)**: Navigate through commands
- **Enter**: Execute the selected command
- **Esc**: Cancel and clear the palette
- **Ctrl+C**: Exit interactive mode

### Command Execution

When you execute a command, the output appears below:

```
> █

┌───────────────────────────────────────────────────┐
│                                                   │
│ Generated Password: XeTunUdWVmsa                  │
│                                                   │
│ Do you want to store this password?               │
│ (Run 'genp create' command with flags for custom │
│  options)                                         │
│                                                   │
└───────────────────────────────────────────────────┘
```

## Available Commands

| Command  | Description                      | Notes                                      |
|----------|----------------------------------|--------------------------------------------|
| `/create` | Generate a new password         | Uses default settings (12 chars, all types)|
| `/show`   | Display stored passwords         | Shows count only (use CLI for decryption)  |
| `/help`   | Show available commands          | Lists all commands with descriptions       |
| `/exit`   | Exit interactive mode            | Alternative to Ctrl+C                      |

## Features

- **Visual Command Palette**: Clean, bordered UI similar to Claude Code
- **Keyboard Navigation**: Full keyboard control with arrow keys
- **Instant Feedback**: Commands execute immediately upon selection
- **Styled Output**: Color-coded output for better readability
- **Non-intrusive**: Doesn't affect existing CLI functionality

## When to Use Interactive Mode vs CLI

### Use Interactive Mode When:
- You want a visual, guided experience
- You're generating multiple passwords in a session
- You prefer keyboard navigation over typing full commands
- You want to explore available commands

### Use CLI Mode When:
- You need advanced options (custom password lengths, character sets)
- You're scripting or automating tasks
- You need to decrypt and view stored passwords (requires secure input)
- You're integrating with other command-line tools

## Examples

### Quick Password Generation

1. Start interactive mode: `genp interactive`
2. Type `/`
3. Press Enter (create is selected by default)
4. View generated password

### Checking Stored Passwords

1. Start interactive mode: `genp interactive`
2. Type `/`
3. Press ↓ to select "show"
4. Press Enter
5. See password count (use `genp show` in CLI to decrypt)

### Getting Help

1. Start interactive mode: `genp interactive`
2. Type `/help` and press Enter
3. Or type `/` and navigate to help with arrow keys
