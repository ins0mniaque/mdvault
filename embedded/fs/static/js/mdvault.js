function createFile(markdown = '') {
	fetch(window.location.href, {
		method: 'PUT',
		headers: { 'Content-Type': 'text/markdown' },
		body: markdown
	})
	.then(() => {
		window.location.reload();
	})
	.catch((error) => {
		alert('Error creating file: ' + error);
	});
}

function saveFile(markdown) {
	fetch(window.location.href, {
		method: 'PUT',
		headers: { 'Content-Type': 'text/markdown' },
		body: markdown
	})
	.catch((error) => {
		alert('Error saving file: ' + error);
	});
}

function deleteFile(editor) {
	fetch(window.location.href, {
		method: 'DELETE'
	})
	.then(() => {
		window.location.reload();
	})
	.catch((error) => {
		alert('Error deleting file: ' + error);
	});
}
