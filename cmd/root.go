package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mdvault",
	Short: "mdvault is a markdown knowledge base command-line tool",
	Long:  "mdvault is a markdown knowledge base command-line tool",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
