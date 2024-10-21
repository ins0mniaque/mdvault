var editor = new SimpleMDE({
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
	]
});

document.addEventListener('keydown', function(e) {
	if ((e.ctrlKey || e.metaKey) && e.key === 's') {
		e.preventDefault();

		saveFile(editor.value());
	}
});
