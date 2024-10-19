package cmd

import (
	"fmt"
	"log"
	"mdvault/config"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:     "serve",
	Aliases: []string{"start"},
	Short:   "Start vault server",
	Long:    "Start vault server rendering markdown files as HTML",
	Run: func(cmd *cobra.Command, args []string) {
		renderer, err := config.ConfigureRenderer()
		if err != nil {
			log.Fatalf("Error configuring renderer for vault %s: %v", vaultDir, err)
		}

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			filename := filepath.Join(vaultDir, r.URL.Path)

			data, err := os.ReadFile(filename)
			if err != nil {
				if os.IsNotExist(err) {
					http.NotFound(w, r)
				} else {
					log.Printf("Error reading file: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
				return
			}

			ext := filepath.Ext(filename)

			if ext == ".md" || ext == ".MD" {
				fmt.Fprint(w, "<html><body>")
				renderer.Render(data, w)
				fmt.Fprint(w, "</body></html>")
			} else {
				w.Write(data)
			}
		})

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
