# GenP

```bash

  /$$$$$$                      /$$$$$$$ 
 /$$__  $$                    | $$__  $$
| $$  \__/  /$$$$$$  /$$$$$$$ | $$  \ $$
| $$ /$$$$ /$$__  $$| $$__  $$| $$$$$$$/
| $$|_  $$| $$$$$$$$| $$  \ $$| $$____/ 
| $$  \ $$| $$_____/| $$  | $$| $$      
|  $$$$$$/|  $$$$$$$| $$  | $$| $$      
 \______/  \_______/|__/  |__/|__/     
 
 ```
 
 GenP is a command-line tool for generating passwords and storing them in E2EE (End-to-End Encrypted) mode. It provides a secure way to manage your passwords and ensures that only you have access to them.

## Features

- Generate strong, random passwords with customizable length and character sets
- Store passwords securely with end-to-end encryption
- **Interactive mode with Claude Code-like command palette** - Type `/` to see all available commands
- Command-line interface for easy integration into workflows
- Local storage with encryption ensuring your passwords never leave your device unencrypted
- Simple and intuitive commands for password management

## About

GenP focuses on security and simplicity. All passwords are encrypted locally using industry-standard encryption before being stored. The encryption keys are derived from your master password, meaning that your data remains secure and only accessible to you. No passwords or encryption keys are ever transmitted to external servers.

The tool is designed for users who prefer command-line utilities and want full control over their password management without relying on third-party cloud services.

## Usage

### Interactive Mode (Recommended)

GenP now features an interactive mode with a Claude Code-like command palette:

```bash
genp interactive
# or use the short alias
genp i
```

In interactive mode:
- Type `/` to see all available commands
- Use ↑↓ arrow keys to navigate the command list
- Press Enter to execute a selected command
- Type commands directly (e.g., `/create`, `/show`, `/help`)
- Press Esc to cancel or Ctrl+C to exit

### Traditional CLI Mode

You can also use GenP with traditional command-line arguments:

#### Generate a Password

```bash
# Generate a basic password
genp create

# Generate a password with specific options
genp create -0 -A -$ --length 16
```

Options:
- `-0` or `--numbers`: Include numbers (0-9)
- `-A` or `--uppercase`: Include uppercase letters (A-Z)
- `-$` or `--special`: Include special characters (!@#$&)
- `-l` or `--length`: Set password length (default: 12)

#### Show Stored Passwords

```bash
genp show
```

This will prompt for your master password and display all stored passwords.
 