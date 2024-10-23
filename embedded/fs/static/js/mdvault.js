var mdvault = (function() {
	var mdvault = { };

	mdvault.create = function(markdown = '') {
		fetch(window.location.href, {
			method: 'PUT',
			headers: { 'Content-Type': 'text/markdown' },
			body: markdown
		})
		.then(() => {
			window.location.reload();
		})
		.catch(error => {
			Swal.fire({ icon: 'error', title: 'Error creating file', text: error });
		});
	}

	mdvault.save = function(markdown) {
		fetch(window.location.href, {
			method: 'PUT',
			headers: { 'Content-Type': 'text/markdown' },
			body: markdown
		})
		.then(() => {
			Swal.fire({ icon: 'success', title: 'Saved', position: 'top-end', showConfirmButton: false, timer: 1500 })
		})
		.catch(error => {
			Swal.fire({ icon: 'error', title: 'Error saving file', text: error });
		});
	}

	mdvault.delete = function() {
		fetch(window.location.href, {
			method: 'DELETE'
		})
		.then(() => {
			window.location.reload();
		})
		.catch(error => {
			Swal.fire({ icon: 'error', title: 'Error deleting file', text: error, icon: 'error' });
		});
	}

	mdvault.render = function(markdown, onrendered) {
		fetch(window.location.origin + '?render', {
			method: 'POST',
			headers: { 'Content-Type': 'text/markdown' },
			body: markdown
		})
		.then(response => response.text())
		.then(html => onrendered(html))
		.catch(error => {
			Swal.fire({ icon: 'error', title: 'Error rendering markdown', text: error, position: 'top-end', showConfirmButton: false, timer: 1500 })
		});
	}

	var mermaid_initialized = false

	function initializeMermaid() {
		if(mermaid_initialized)
			return;

		mermaid.initialize({ startOnLoad: false });
		mermaid_initialized = true;
	}

	mdvault.assign = function(element, html) {
		element.innerHTML = html;

		MathJax.typeset();

		initializeMermaid();
		mermaid.run();
	}

	return mdvault;
})();
