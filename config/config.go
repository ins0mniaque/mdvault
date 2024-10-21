package config

import (
	"html/template"
	"mdvault/embedded"
	"mdvault/markdown"
	"mdvault/markdown/goldmark"
)

// TODO: Expose options through config
func ConfigureParser() (markdown.Parser, error) {
	return goldmark.NewParser(), nil
}

// TODO: Expose options through config
func ConfigureRenderer() (markdown.Renderer, error) {
	return goldmark.NewRenderer(), nil
}

// TODO: Read template from vault config
func ConfigureCreatorTemplate() (*template.Template, error) {
	return template.ParseFS(embedded.FS, "template/creator.tmpl")
}

// TODO: Read template from vault config
func ConfigureEditorTemplate() (*template.Template, error) {
	return template.ParseFS(embedded.FS, "template/editor.tmpl")
}

// TODO: Read template from vault config
func ConfigureRenderTemplate() (*template.Template, error) {
	return template.ParseFS(embedded.FS, "template/render.tmpl")
}
