package vault

import (
	"html/template"
	"io"
	"log"
	"mdvault/config"
	"mdvault/markdown"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

type Server struct {
	Vault          *Vault
	renderer       markdown.Renderer
	editorTemplate *template.Template
	renderTemplate *template.Template
}

func NewServer(vault *Vault) (*Server, error) {
	renderer, err := config.ConfigureRenderer()
	if err != nil {
		log.Printf("Error configuring renderer for vault %s: %v\n", vault.Dir(), err)
		return nil, err
	}

	editorTemplate, err := config.ConfigureEditorTemplate()
	if err != nil {
		log.Printf("Error configuring editor template for vault %s: %v\n", vault.Dir(), err)
		return nil, err
	}

	renderTemplate, err := config.ConfigureRenderTemplate()
	if err != nil {
		log.Printf("Error configuring render template for vault %s: %v\n", vault.Dir(), err)
		return nil, err
	}

	return &Server{
		Vault:          vault,
		renderer:       renderer,
		editorTemplate: editorTemplate,
		renderTemplate: renderTemplate}, nil
}

type EditorPage struct {
	Title    string
	Markdown template.JSStr
}

type RenderPage struct {
	Title    string
	Markdown template.HTML
}

func (server *Server) Handler(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "HEAD":
		server.head(writer, request)
	case "GET":
		server.get(writer, request)
	case "PUT":
		server.put(writer, request)
	case "DELETE":
		server.delete(writer, request)
	case "PATCH":
		server.patch(writer, request)
	default:
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (server *Server) head(writer http.ResponseWriter, request *http.Request) {
	path := filepath.Join(server.Vault.Dir(), request.URL.Path)
	ext := strings.ToLower(filepath.Ext(path))

	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) && ext == ".html" {
		path = path[:len(path)-len(ext)] + ".md"
		ext = ".md"
		_, err = os.Stat(path)
	}

	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(writer, request)
		} else {
			log.Printf("Error checking file: %v", err)
			http.Error(writer, "Failed to check file", http.StatusInternalServerError)
		}
		return
	}
}

func (server *Server) get(writer http.ResponseWriter, request *http.Request) {
	path := filepath.Join(server.Vault.Dir(), request.URL.Path)
	ext := strings.ToLower(filepath.Ext(path))
	render := false

	data, err := os.ReadFile(path)
	if err != nil && os.IsNotExist(err) && ext == ".html" {
		path = path[:len(path)-len(ext)] + ".md"
		ext = ".md"
		render = true

		data, err = os.ReadFile(path)
	}

	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(writer, request)
		} else {
			log.Printf("Error reading file: %v", err)
			http.Error(writer, "Failed to read file", http.StatusInternalServerError)
		}
		return
	}

	if render {
		title := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		html := new(strings.Builder)
		err := server.renderer.Render(data, html)
		if err != nil {
			log.Printf("Error rendering markdown: %v", err)
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}

		page := RenderPage{
			Title:    title,
			Markdown: template.HTML(html.String())}

		err = server.renderTemplate.Execute(writer, page)
		if err != nil {
			log.Printf("Error rendering template: %v", err)
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	} else if ext == ".md" {
		title := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		js := template.JSEscapeString(string(data))

		page := EditorPage{
			Title:    title,
			Markdown: template.JSStr(js)}

		err := server.editorTemplate.Execute(writer, page)
		if err != nil {
			log.Printf("Error rendering template: %v", err)
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	} else {
		_, err := writer.Write(data)
		if err != nil {
			log.Printf("Error rendering resource: %v", err)
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (server *Server) put(writer http.ResponseWriter, request *http.Request) {
	path := filepath.Join(server.Vault.Dir(), request.URL.Path)
	ext := strings.ToLower(filepath.Ext(path))

	if ext != ".md" {
		http.Error(writer, "Only markdown files are allowed", http.StatusForbidden)
		return
	}

	file, err := os.Create(path)
	if err != nil {
		log.Printf("Error creating file: %v", err)
		http.Error(writer, "Failed to create file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, request.Body)
	if err != nil {
		log.Printf("Error writing file: %v", err)
		http.Error(writer, "Failed to write file", http.StatusInternalServerError)
		return
	}
}

func (server *Server) delete(writer http.ResponseWriter, request *http.Request) {
	path := filepath.Join(server.Vault.Dir(), request.URL.Path)
	ext := strings.ToLower(filepath.Ext(path))

	if ext != ".md" {
		http.Error(writer, "Only markdown files are allowed", http.StatusForbidden)
		return
	}

	err := os.Remove(path)
	if err != nil {
		http.Error(writer, "File not found or failed to delete file", http.StatusNotFound)
		return
	}
}

func (server *Server) patch(writer http.ResponseWriter, request *http.Request) {
	path := filepath.Join(server.Vault.Dir(), request.URL.Path)
	ext := strings.ToLower(filepath.Ext(path))

	if ext != ".md" {
		http.Error(writer, "Only markdown files are allowed", http.StatusForbidden)
		return
	}

	builder := new(strings.Builder)
	_, err := io.Copy(builder, request.Body)
	if err != nil {
		log.Printf("Error reading patch: %v", err)
		http.Error(writer, "Failed to read patch", http.StatusInternalServerError)
		return
	}

	dmp := diffmatchpatch.New()
	patches, err := dmp.PatchFromText(builder.String())
	if err != nil {
		log.Printf("Invalid patch: %v", err)
		http.Error(writer, "Invalid patch", http.StatusInternalServerError)
		return
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(writer, request)
		} else {
			log.Printf("Error reading file: %v", err)
			http.Error(writer, "Failed to read file", http.StatusInternalServerError)
		}
		return
	}

	patched, applieds := dmp.PatchApply(patches, string(data))
	for _, applied := range applieds {
		if !applied {
			log.Printf("Error applying patch: %v", err)
			http.Error(writer, "Failed to apply patch", http.StatusInternalServerError)
			return
		}
	}

	file, err := os.Create(path)
	if err != nil {
		log.Printf("Error creating file: %v", err)
		http.Error(writer, "Failed to create file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	_, err = file.Write([]byte(patched))
	if err != nil {
		log.Printf("Error writing file: %v", err)
		http.Error(writer, "Failed to write file", http.StatusInternalServerError)
		return
	}
}
