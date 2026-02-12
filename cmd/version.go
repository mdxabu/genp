/*
Copyright 2025 - github.com/mdxabu
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of genp",
	Long:  `Display the current version of the genp CLI tool.`,
	Run: func(cmd *cobra.Command, args []string) {
		version := "v0.1.1"
		fmt.Println("genp " + version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
