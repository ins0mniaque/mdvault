package cmd

import (
	"fmt"
	"log"
	"mdvault/vault"

	"github.com/spf13/cobra"
)

var linksBack bool
var linksDead bool

var linksCmd = &cobra.Command{
	Use:     "links",
	Aliases: []string{"link", "ln"},
	Short:   "List links, backlinks and dead links",
	Long:    "List links, backlinks and dead links",
	Run: func(cmd *cobra.Command, args []string) {
		v := vault.NewVault(vaultDir)

		if err := v.Load(); err != nil {
			log.Fatalf("Error loading vault %s: %s", v.Dir(), err)
		}

		if linksBack {
			if linksDead {
				return
			}

			for path, metadata := range v.Entries() {
				for _, backlink := range metadata.Backlinks {
					fmt.Printf("%s <-- %s\n", path, backlink)
				}
			}
		} else if !linksDead {
			for path, metadata := range v.Entries() {
				for _, link := range metadata.Links {
					fmt.Printf("%s --> %s\n", path, link)
				}
			}
		} else {
			entries := v.Entries()

			for path, metadata := range entries {
				for _, link := range metadata.Links {
					if _, ok := entries[link]; !ok {
						fmt.Printf("%s --> %s\n", path, link)
					}
				}
			}
		}
	},
}

func init() {
	linksCmd.Flags().BoolVarP(&linksBack, "back", "b", false, "List backlinks instead")
	linksCmd.Flags().BoolVarP(&linksDead, "dead", "d", false, "List dead links only")

	rootCmd.AddCommand(linksCmd)
}
