/*
Copyright © 2025 - github.com/mdxabu
*/
package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "genp",
	Short: "Generate Password, store and encrypted in CLI",
	Long: `Generate`,
	
	Run: func(cmd *cobra.Command, args []string) {
		asciiBanner := `
  /$$$$$$                      /$$$$$$$ 
 /$$__  $$                    | $$__  $$
| $$  \__/  /$$$$$$  /$$$$$$$ | $$  \ $$
| $$ /$$$$ /$$__  $$| $$__  $$| $$$$$$$/
| $$|_  $$| $$$$$$$$| $$  \ $$| $$____/ 
| $$  \ $$| $$_____/| $$  | $$| $$      
|  $$$$$$/|  $$$$$$$| $$  | $$| $$      
 \______/  \_______/|__/  |__/|__/      
                                        
                                        
                                        
		`
		
		fmt.Println(asciiBanner)
		fmt.Println("Welcome to GenP, to create and store password E2EE :)")
		
		fmt.Printf("GenP is running on %s\n", runtime.GOOS)
	},

}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


