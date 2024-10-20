package vault

import (
	"html/template"
	"log"
	"mdvault/config"
	"mdvault/markdown"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
	filename := filepath.Join(server.Vault.Dir(), request.URL.Path)
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
			http.NotFound(writer, request)
		} else {
			log.Printf("Error reading file: %v", err)
			http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	if render {
		title := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
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
		title := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
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
