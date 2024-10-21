package cmd

import (
	"io/fs"
	"log"
	"mdvault/embedded"
	"mdvault/vault"
	"net/http"

	"github.com/rs/cors"
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

		addr := ":8080"
		mux := http.NewServeMux()
		mux.HandleFunc("/", server.Handler)

		if err := handleFS(mux, embedded.FS); err != nil {
			log.Fatal(err)
		}

		log.Printf("Listening on %s...\n", addr)

		handler := logger(cors.Default().Handler(mux))
		err = http.ListenAndServe(addr, handler)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func handleFS(mux *http.ServeMux, fsys fs.FS) error {
	entries, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return err
	}

	fsysHandler := http.FileServer(http.FS(fsys))
	for _, entry := range entries {
		if entry.IsDir() {
			mux.Handle("/"+entry.Name()+"/", fsysHandler)
		} else {
			mux.Handle("/"+entry.Name(), fsysHandler)
		}
	}

	return nil
}

func logger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("%s %s %s\n", request.RemoteAddr, request.Method, request.URL)
		handler.ServeHTTP(writer, request)
	})
}
