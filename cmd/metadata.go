package cmd

import (
	"encoding/json"
	"log"
	"mdvault/vault"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var metadataFormat string

var metadataCmd = &cobra.Command{
	Use:     "metadata",
	Aliases: []string{"meta", "md"},
	Short:   "Extract vault metadata",
	Long:    "Extract vault metadata in JSON, YAML or TOML format",
	Run: func(cmd *cobra.Command, args []string) {
		v := vault.NewVault(vaultDir)

		if err := v.Load(); err != nil {
			log.Fatalf("Error loading vault %s: %s", v.Dir(), err)
		}

		if metadataFormat == "json" {
			json, err := json.Marshal(v.Entries())
			if err != nil {
				log.Fatal(err)
			}

			println(string(json))
		} else if metadataFormat == "yaml" {
			yaml, err := yaml.Marshal(v.Entries())
			if err != nil {
				log.Fatal(err)
			}

			println(string(yaml))
		} else if metadataFormat == "toml" {
			err := toml.NewEncoder(os.Stdout).Encode(v.Entries())
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatalf("Invalid format: %s. Available formats: json|yaml|toml", metadataFormat)
		}
	},
}

func init() {
	metadataCmd.Flags().StringVarP(&metadataFormat, "format", "f", "yaml", "Output format: json|yaml|toml")

	rootCmd.AddCommand(metadataCmd)
}
