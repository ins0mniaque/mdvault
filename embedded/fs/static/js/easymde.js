var editor = new EasyMDE({
	autofocus: true,
	element: document.getElementById('editor'),
	toolbar: [
		{
			name: "save",
			action: function(editor) { saveFile(editor.value()); },
			className: "fa fa-save",
			title: "Save",
		},
		{
			name: "delete",
			action: function() { deleteFile(); },
			className: "fa fa-trash",
			title: "Delete",
		},
		"|",
		"bold", "italic", "heading", "|",
		"quote", "unordered-list", "ordered-list", "|",
		"link", "image", "|",
		"preview", "side-by-side", "fullscreen", "|",
		"guide"
	],
	previewRender: function(markdown, preview) {
		clearTimeout(editor.previewRenderTimeout)
		editor.previewRenderTimeout = setTimeout(function() {
			renderFile(markdown, function(html) {
				preview.innerHTML = html;
			});
		}, 250);

		return preview.innerHTML;
	}
});

document.addEventListener('keydown', function(e) {
	if ((e.ctrlKey || e.metaKey) && e.key === 's') {
		e.preventDefault();

		saveFile(editor.value());
	}
});
