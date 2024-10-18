package vault

import (
	"fmt"
	"io/fs"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/frontmatter"
	"go.abhg.dev/goldmark/hashtag"
	"go.abhg.dev/goldmark/wikilink"
)

type Vault struct {
	Files     map[string]*Metadata
	BackLinks map[string]*Metadata
	Tagged    map[string][]*Metadata
}

func (v *Vault) Update(metadata *Metadata) {
	v.Files[metadata.Path] = metadata
}

type Metadata struct {
	Path       string                 `json:"path"`
	Names      map[string]struct{}    `json:"names"`
	Dates      map[time.Time]struct{} `json:"dates"`
	Links      map[string]struct{}    `json:"links"`
	Tags       map[string]struct{}    `json:"tags"`
	Tasks      map[string]struct{}    `json:"tasks"`
	Properties map[string]interface{} `json:"properties"`
}

func (m *Metadata) AddName(name string) {
	m.Names[name] = struct{}{}
}

func (m *Metadata) AddDate(date time.Time) {
	m.Dates[date] = struct{}{}
}

func (m *Metadata) AddLink(link string) {
	m.Links[link] = struct{}{}
}

func (m *Metadata) AddTag(tag string) {
	m.Tags[tag] = struct{}{}
}

func (m *Metadata) AddTask(task string) {
	m.Tasks[task] = struct{}{}
}

func (m *Metadata) SetPath(path string) {
	// TODO: Extract date/name from path
	m.Path = path
	m.Names[path] = struct{}{}
}

func extractMetadata(document *ast.Document, source []byte) (*Metadata, error) {
	metadata := Metadata{
		Names:      make(map[string]struct{}),
		Links:      make(map[string]struct{}),
		Tags:       make(map[string]struct{}),
		Tasks:      make(map[string]struct{}),
		Properties: document.Meta(),
	}

	title := ""

	err := ast.Walk(document, func(node ast.Node, enter bool) (ast.WalkStatus, error) {
		if title == "" {
			if n, ok := node.(*ast.Heading); ok && n.Level == 1 && enter {
				title = string(n.Text(source))
				metadata.AddName(title)
			}
		}

		if n, ok := node.(*ast.Link); ok && enter {
			u, _ := url.Parse(string(n.Destination))
			u.Fragment = ""
			link := u.String()
			if len(link) > 0 {
				metadata.AddLink(link)
			}
		}

		if n, ok := node.(*wikilink.Node); ok && enter {
			if len(n.Target) > 0 {
				metadata.AddLink(string(n.Target))
			}
		}

		if n, ok := node.(*hashtag.Node); ok && enter {
			metadata.AddTag(string(n.Tag))
		}

		if n, ok := node.(*east.TaskCheckBox); ok && enter {
			metadata.AddTask(string(n.Parent().Text(source)))
		}

		return ast.WalkContinue, nil
	})

	if err != nil {
		return nil, err
	}

	extractTags := func(key string) {
		if metatags, ok := metadata.Properties[key].([]interface{}); ok {
			for _, v := range metatags {
				if s, ok := v.(string); ok {
					metadata.AddTag(s)
				}
			}
		}
	}

	extractTags("Tags")
	extractTags("tags")

	// TODO: Extract name, aliases, title and id into names
	//       e.g. alias, ALIAS and ALIASES, date, time

	return &metadata, nil
}

func parseMarkdown(filename string, md goldmark.Markdown) (*ast.Document, []byte, error) {
	source, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("Error reading file %s: %v", filename, err)
		return nil, nil, err
	}

	reader := text.NewReader(source)
	root := md.Parser().Parse(reader)
	doc := root.OwnerDocument()

	return doc, source, nil
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
			extension.TaskList,
			&frontmatter.Extender{Mode: frontmatter.SetMetadata},
			&hashtag.Extender{Variant: hashtag.ObsidianVariant},
			&wikilink.Extender{})), nil
}

func Parse(root string) {
	md, err := configureParser()
	if err != nil {
		log.Fatal(err)
	}

	err = walkMarkdownFiles(root, func(filename string) {
		document, source, err := parseMarkdown(filename, md)
		if err != nil {
			log.Printf("Error processing file %s: %v", filename, err)
			return
		}

		metadata, err := extractMetadata(document, source)
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
