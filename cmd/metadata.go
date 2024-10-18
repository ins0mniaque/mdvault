package cmd

import (
	"fmt"
	"mdvault/vault"

	"github.com/spf13/cobra"
)

var metadataCmd = &cobra.Command{
	Use:     "metadata",
	Aliases: []string{"meta", "md"},
	Short:   "Extract vault metadata",
	Long:    "Extract vault metadata",
	Run: func(cmd *cobra.Command, args []string) {
		v := vault.NewVault(vaultDir)
		v.Load()

		for path, metadata := range v.Files() {
			if metadata != nil {
				fmt.Println(path, *metadata)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(metadataCmd)
}
