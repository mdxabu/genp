/*
Copyright 2025 - github.com/mdxabu
*/
package cmd

import (
	"github.com/fatih/color"
	"github.com/mdxabu/genp/internal/strength"
	"github.com/spf13/cobra"
)

// strengthCmd represents the strength command
var strengthCmd = &cobra.Command{
	Use:   "strength",
	Short: "Check password strength in real-time",
	Long: `Interactively check the strength of a password as you type it.

Each character you type updates the strength meter and roast message
in real-time. Press Enter when done.

Example:
  genp strength`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := strength.RunInteractive()
		if err != nil {
			if err.Error() == "cancelled" {
				color.Yellow("Password strength check cancelled.\n")
				return
			}
			color.Red("Error: %v\n", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(strengthCmd)
}
