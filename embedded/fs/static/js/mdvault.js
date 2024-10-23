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
			alert('Error creating file: ' + error);
		});
	}

	mdvault.save = function(markdown) {
		fetch(window.location.href, {
			method: 'PUT',
			headers: { 'Content-Type': 'text/markdown' },
			body: markdown
		})
		.catch(error => {
			alert('Error saving file: ' + error);
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
			alert('Error deleting file: ' + error);
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
			alert('Error rendering file: ' + error);
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
