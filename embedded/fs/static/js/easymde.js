var editor = new EasyMDE({
	autofocus: true,
	element: document.getElementById('editor'),
	toolbar: [
		{
			name: "save",
			action: function(editor) { mdvault.save(editor.value()); },
			className: "fa fa-save",
			title: "Save",
		},
		{
			name: "delete",
			action: function() { mdvault.delete(); },
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
			mdvault.render(markdown, function(html) {
				mdvault.assign(preview, html);
			});
		}, 250);

		return null;
	}
});

document.addEventListener('keydown', function(e) {
	if ((e.ctrlKey || e.metaKey) && e.key === 's') {
		e.preventDefault();

		mdvault.save(editor.value());
	}
});
