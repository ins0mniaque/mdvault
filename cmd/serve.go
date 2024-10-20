package cmd

import (
	"fmt"
	"log"
	"mdvault/config"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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
			ext := strings.ToLower(filepath.Ext(filename))
			render := false

			data, err := os.ReadFile(filename)
			if err != nil && os.IsNotExist(err) && ext == ".html" {
				filename = filename[:len(filename)-len(ext)] + ".md"
				ext = ".md"
				render = true

				data, err = os.ReadFile(filename)
			}

			if err != nil {
				if os.IsNotExist(err) {
					http.NotFound(w, r)
				} else {
					log.Printf("Error reading file: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
				return
			}

			if render {
				fmt.Fprint(w, "<html><body>")
				renderer.Render(data, w)
				fmt.Fprint(w, "</body></html>")
			} else if ext == ".md" {
				fmt.Fprint(w, "<html><body>")
				fmt.Fprint(w, "<head>")
				fmt.Fprint(w, `<link rel="stylesheet" href="https://cdn.jsdelivr.net/simplemde/latest/simplemde.min.css">`)
				fmt.Fprint(w, `<script src="https://cdn.jsdelivr.net/simplemde/latest/simplemde.min.js"></script>`)
				fmt.Fprint(w, "</head>")
				fmt.Fprint(w, "<textarea></textarea>")
				fmt.Fprintf(w, "<script>var editor = new SimpleMDE(); editor.value(%q); </script>", data)
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
