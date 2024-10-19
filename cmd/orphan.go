package cmd

import (
	"fmt"
	"log"
	"mdvault/vault"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

var orphanConfirm bool
var orphanDelete bool

var orphanCmd = &cobra.Command{
	Use:     "orphan",
	Aliases: []string{"orphans", "orphaned"},
	Short:   "List orphan files",
	Long:    "List and optionally delete orphan files",
	Run: func(cmd *cobra.Command, args []string) {
		v := vault.NewVault(vaultDir)
		v.Load()

		orphans := make([]string, 0)

		for path, metadata := range v.Entries() {
			if len(metadata.Names) == 0 && len(metadata.Backlinks) == 0 {
				orphans = append(orphans, path)
			}
		}

		sort.Strings(orphans)

		for _, path := range orphans {
			println(path)
		}

		if orphanDelete && len(orphans) > 0 {
			if orphanConfirm && !confirm(fmt.Sprintf("\nDelete %d orphans?", len(orphans))) {
				return
			}

			fmt.Printf("\nDeleting %d orphans...\n", len(orphans))

			for _, path := range orphans {
				err := os.Remove(filepath.Join(v.Dir(), path))
				if err != nil {
					log.Println(err)
				}
			}

			fmt.Printf("\n%d orphans deleted\n", len(orphans))
		} else if len(orphans) > 0 {
			fmt.Printf("\n%d orphans\n", len(orphans))
		} else {
			println("No orphans")
		}
	},
}

func init() {
	orphanCmd.Flags().BoolVarP(&orphanConfirm, "confirm", "c", true, "Confirm orphan files deletion")
	orphanCmd.Flags().BoolVarP(&orphanDelete, "delete", "d", false, "Delete orphan files")

	rootCmd.AddCommand(orphanCmd)
}

func confirm(prompt string) bool {
	fmt.Printf("%s [y|n]: ", prompt)

	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		log.Fatal(err)
	}

	input = strings.ToLower(input)

	return input == "y" || input == "yes"
}
