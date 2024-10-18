package vault

import (
	"fmt"
	"io/fs"
	"log"
	"mdvault/config"
	"mdvault/parser"
	"os"
	"path/filepath"
	"sync"
)

type Vault struct {
	Files     map[string]*parser.Metadata
	BackLinks map[string]*parser.Metadata
	Tagged    map[string][]*parser.Metadata
}

func (v *Vault) Update(metadata *parser.Metadata) {
	v.Files[metadata.Path] = metadata
}

func Extract(root string) {
	parser, err := config.ConfigureParser()
	if err != nil {
		log.Fatal(err)
	}

	err = walkMarkdownFiles(root, func(filename string) {
		source, err := os.ReadFile(filename)
		if err != nil {
			log.Printf("Error reading file %s: %v", filename, err)
			return
		}

		metadata, err := parser.Parse(source)
		if err != nil {
			log.Printf("Error extracting metadata for file %s: %v", filename, err)
			return
		}

		path, err := filepath.Rel(root, filename)
		if err != nil {
			path = filename
		}

		metadata.SetPath(path)

		fmt.Println(metadata)
	})

	if err != nil {
		log.Fatal(err)
	}
}

func walkMarkdownFiles(root string, fn func(filename string)) error {
	var wg sync.WaitGroup

	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("Error walking path %s: %v", path, err)
			return err
		}

		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".md" {
			wg.Add(1)

			go func(filename string) {
				defer wg.Done()

				fn(filename)
			}(path)
		}

		return nil
	})

	wg.Wait()

	return err
}
