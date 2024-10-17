package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/frontmatter"
	"go.abhg.dev/goldmark/hashtag"
	"go.abhg.dev/goldmark/wikilink"
)

type Metadata struct {
	Names      map[string]string      `json:"names"`
	Links      map[string]struct{}    `json:"links"`
	Tags       map[string]struct{}    `json:"tags"`
	Properties map[string]interface{} `json:"properties"`
}

func extractMetadata(document *ast.Document) (*Metadata, error) {
	links := make(map[string]struct{})
	tags := make(map[string]struct{})
	properties := document.Meta()

	// TODO: Extract title from first title header
	err := ast.Walk(document, func(node ast.Node, enter bool) (ast.WalkStatus, error) {
		if n, ok := node.(*ast.Link); ok && enter {
			links[string(n.Destination)] = struct{}{}
		}

		if n, ok := node.(*wikilink.Node); ok && enter {
			if len(n.Target) > 0 {
				links[string(n.Target)] = struct{}{}
			}
		}

		if n, ok := node.(*hashtag.Node); ok && enter {
			tags[string(n.Tag)] = struct{}{}
		}

		return ast.WalkContinue, nil
	})

	if err != nil {
		return nil, err
	}

	extractTags := func(key string) {
		if metatags, ok := properties[key].([]interface{}); ok {
			for _, v := range metatags {
				if s, ok := v.(string); ok {
					tags[s] = struct{}{}
				}
			}
		}
	}

	extractTags("Tags")
	extractTags("tags")

	// TODO: Extract aliases, title and id into names

	return &Metadata{Links: links, Tags: tags, Properties: properties}, nil
}

func parseMarkdown(filename string, md goldmark.Markdown) (*ast.Document, error) {
	contents, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("Error reading file %s: %v", filename, err)
		return nil, err
	}

	reader := text.NewReader(contents)
	root := md.Parser().Parse(reader)
	doc := root.OwnerDocument()

	return doc, nil
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

// TODO: Expose options through config
func configureParser() (goldmark.Markdown, error) {
	return goldmark.New(
		goldmark.WithExtensions(
			&frontmatter.Extender{Mode: frontmatter.SetMetadata},
			&hashtag.Extender{Variant: hashtag.ObsidianVariant},
			&wikilink.Extender{})), nil
}

func main() {
	root := "."

	md, err := configureParser()
	if err != nil {
		log.Fatal(err)
	}

	err = walkMarkdownFiles(root, func(filename string) {
		document, err := parseMarkdown(filename, md)
		if err != nil {
			log.Printf("Error processing file %s: %v", filename, err)
			return
		}

		metadata, err := extractMetadata(document)
		if err != nil {
			log.Printf("Error extracting metadata for file %s: %v", filename, err)
			return
		}

		fmt.Println(filename, "=>", metadata)
	})

	if err != nil {
		log.Fatal(err)
	}
}
