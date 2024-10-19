package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:     "init",
	Aliases: []string{},
	Short:   "Initialize vault",
	Long:    "Initialize vault configuration folder (.mdvault)",
	Run: func(cmd *cobra.Command, args []string) {
		vaultConfigDir := filepath.Join(vaultDir, ".mdvault")

		info, err := os.Stat(vaultConfigDir)
		if err == nil {
			if info.IsDir() {
				log.Fatalf("Vault %s is already initialized", vaultDir)
			}

			log.Fatalf("%s already exists and is not a directory", vaultConfigDir)
		}

		if os.IsNotExist(err) {
			err = os.Mkdir(vaultConfigDir, os.ModePerm)
		}

		if err != nil {
			log.Fatalf("Error initializing %s: %s", vaultDir, err)
		}

		fmt.Printf("Vault %s initialized\n", vaultDir)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
