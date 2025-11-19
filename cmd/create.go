/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/mdxabu/genp/internal"
	"github.com/spf13/cobra"
)

var (
	includeNumbers   bool
	includeUppercase bool
	includeSpecial   bool
	passwordLength   int
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Generate a secure password",
	Long: `Generate a secure password with customizable options.

You can specify which character types to include in your password:
  -0 : Include numbers (0-9)
  -A : Include uppercase letters (A-Z)
  -$ : Include special characters (!@#$&)

Example:
  genp create -0 -A -$ --length 16`,
	Run: func(cmd *cobra.Command, args []string) {
		password := internal.GeneratePassword(passwordLength, includeNumbers, includeUppercase, includeSpecial)
		fmt.Println("Generated Password:", password)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().BoolVarP(&includeNumbers, "numbers", "0", false, "Include numbers in password")
	createCmd.Flags().BoolVarP(&includeUppercase, "uppercase", "A", false, "Include uppercase letters in password")
	createCmd.Flags().BoolVarP(&includeSpecial, "special", "$", false, "Include special characters in password")
	createCmd.Flags().IntVarP(&passwordLength, "length", "l", 12, "Length of the password")
}
