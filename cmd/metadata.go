package cmd

import (
	"mdvault/vault"

	"github.com/spf13/cobra"
)

var metadataCmd = &cobra.Command{
	Use:     "metadata",
	Aliases: []string{"meta", "md"},
	Short:   "Extract vault metadata",
	Long:    "Extract vault metadata",
	Run: func(cmd *cobra.Command, args []string) {
		vault.Parse(".")
	},
}

func init() {
	rootCmd.AddCommand(metadataCmd)
}
