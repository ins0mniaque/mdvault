package cmd

import (
	"log"
	"mdvault/vault"
	"net/http"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:     "serve",
	Aliases: []string{"start"},
	Short:   "Start vault server",
	Long:    "Start vault server rendering markdown files as HTML",
	Run: func(cmd *cobra.Command, args []string) {
		v := vault.NewVault(vaultDir)
		if err := v.Load(); err != nil {
			log.Fatalf("Error loading vault %s: %s", v.Dir(), err)
		}

		server, err := vault.NewServer(v)
		if err != nil {
			log.Fatalf("Error crearing server for vault %s: %s", v.Dir(), err)
		}

		http.HandleFunc("/", server.Handler)

		log.Println("Listening on :8080...")
		err = http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
