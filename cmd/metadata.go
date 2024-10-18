package cmd

import (
	"encoding/json"
	"log"
	"mdvault/vault"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var format string

var metadataCmd = &cobra.Command{
	Use:     "metadata",
	Aliases: []string{"meta", "md"},
	Short:   "Extract vault metadata",
	Long:    "Extract vault metadata in JSON or YAML format",
	Run: func(cmd *cobra.Command, args []string) {
		v := vault.NewVault(vaultDir)
		v.Load()

		if format == "json" {
			json, err := json.Marshal(v.Entries())
			if err != nil {
				log.Fatal(err)
			}

			println(string(json))
		} else if format == "yaml" {
			yaml, err := yaml.Marshal(v.Entries())
			if err != nil {
				log.Fatal(err)
			}

			println(string(yaml))
		} else {
			log.Fatalf("Invalid format: %s. Format must be json or yaml.\n", format)
		}
	},
}

func init() {
	metadataCmd.Flags().StringVarP(&format, "format", "f", "yaml", "Output format: json|yaml")

	rootCmd.AddCommand(metadataCmd)
}
