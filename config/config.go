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
func ConfigureCreatorTemplate() (*template.Template, error) {
	return template.New("creator").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>{{ .Title }}</title>
</head>
<body>
	<button id="create-button">Create {{ .Title }}</button>
	<script>
		document.getElementById('create-button').addEventListener('click', function() {
			fetch(window.location.href, {
				method: 'PUT',
				headers: { 'Content-Type': 'text/markdown' },
				body: ''
			})
			.then(() => {
				window.location.reload();
			})
			.catch((error) => {
				alert('Error creating file: ' + error);
			});
		});
	</script>
</body>
</html>
`)
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
	<button id="save-button">Save</button>
	<button id="delete-button">Delete</button>
	<script>var editor = new SimpleMDE(); editor.value({{ .Markdown }});</script>
	<script>
		document.getElementById('save-button').addEventListener('click', function() {
			fetch(window.location.href, {
				method: 'PUT',
				headers: { 'Content-Type': 'text/markdown' },
				body: editor.value()
			})
			.catch((error) => {
				alert('Error saving file: ' + error);
			});
		});
	</script>
	<script>
		document.getElementById('delete-button').addEventListener('click', function() {
			fetch(window.location.href, {
				method: 'DELETE'
			})
			.then(() => {
				window.location.reload();
			})
			.catch((error) => {
				alert('Error deleting file: ' + error);
			});
		});
	</script>
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
