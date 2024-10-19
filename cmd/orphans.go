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

var orphansConfirm bool
var orphansDelete bool

var orphansCmd = &cobra.Command{
	Use:     "orphans",
	Aliases: []string{"orphan", "orphaned"},
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

		if orphansDelete && len(orphans) > 0 {
			if orphansConfirm && !confirm(fmt.Sprintf("\nDelete %d orphan files?", len(orphans))) {
				return
			}

			fmt.Printf("\nDeleting %d orphan files...\n", len(orphans))

			for _, path := range orphans {
				err := os.Remove(filepath.Join(v.Dir(), path))
				if err != nil {
					log.Println(err)
				}
			}

			fmt.Printf("\n%d orphan files deleted\n", len(orphans))
		} else if len(orphans) > 0 {
			fmt.Printf("\n%d orphan files\n", len(orphans))
		} else {
			println("No orphan files")
		}
	},
}

func init() {
	orphansCmd.Flags().BoolVarP(&orphansConfirm, "confirm", "c", true, "Confirm orphan files deletion")
	orphansCmd.Flags().BoolVarP(&orphansDelete, "delete", "d", false, "Delete orphan files")

	rootCmd.AddCommand(orphansCmd)
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
