package embedded

import "embed"

//go:embed static
var FS embed.FS

//go:embed template
var Templates embed.FS
