# GenP
[![Release](https://github.com/mdxabu/genp/actions/workflows/release.yml/badge.svg)](https://github.com/mdxabu/genp/actions/workflows/release.yml)

 GenP is a command-line tool for generating passwords and storing them in E2EE (End-to-End Encrypted) mode. It provides a secure way to manage your passwords and ensures that only you have access to them.

## Installation

### Go

```bash
go install github.com/mdxabu/genp@latest
```

### From Source

```bash
git clone https://github.com/mdxabu/genp.git
cd genp
go build -o genp .
```

## Features

- Generate strong, random passwords with customizable length and character sets
- Store passwords securely with end-to-end encryption
- Local storage with encryption ensuring your passwords never leave your device unencrypted
- Simple and intuitive commands for password management

## About

GenP focuses on security and simplicity. All passwords are encrypted locally using industry-standard encryption before being stored. The encryption keys are derived from your master password, meaning that your data remains secure and only accessible to you. No passwords or encryption keys are ever transmitted to external servers.

The tool is designed for users who prefer command-line utilities and want full control over their password management without relying on third-party cloud services.

## Usage

#### Generate a Password

```bash
# Generate a basic password
genp create

# Generate a password with specific options
genp create -0 -A -$ -l 16
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
