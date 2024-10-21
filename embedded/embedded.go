package embedded

import (
	"embed"
	"io/fs"
	"log"
)

//go:embed fs
var fsys embed.FS
var FS = sub(fsys, "fs")

//go:embed templates
var templates embed.FS
var Templates = sub(templates, "templates")

func sub(fsys fs.FS, dir string) fs.FS {
	subFsys, err := fs.Sub(fsys, dir)
	if err != nil {
		log.Fatal(err)
	}

	return subFsys
}
