package config

import (
	"html/template"
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
func ConfigureEditorTemplate() (*template.Template, error) {
	return template.New("editor").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>{{ .Title }}</title>
	<link rel="stylesheet" href="https://cdn.jsdelivr.net/simplemde/latest/simplemde.min.css">
	<script src="https://cdn.jsdelivr.net/simplemde/latest/simplemde.min.js"></script>
</head>
<body>
	<textarea></textarea>
	<script>var editor = new SimpleMDE(); editor.value({{ .Markdown }});</script>
</body>
</html>
`)
}

// TODO: Read template from vault config
func ConfigureRenderTemplate() (*template.Template, error) {
	return template.New("render").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>{{ .Title }}</title>
</head>
<body>
	{{ .Markdown }}
</body>
</html>
`)
}
