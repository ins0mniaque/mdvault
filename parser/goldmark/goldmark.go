package goldmark

import (
	"mdvault/parser"
	"net/url"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/frontmatter"
	"go.abhg.dev/goldmark/hashtag"
	"go.abhg.dev/goldmark/wikilink"
)

type GoldmarkParser struct {
	md goldmark.Markdown
}

func (r GoldmarkParser) Parse(source []byte) (*parser.Metadata, error) {
	document, source, err := parseMarkdown(source, r.md)
	if err != nil {
		return nil, err
	}

	metadata, err := extractMetadata(document, source)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

func Create() (parser.Parser, error) {
	return GoldmarkParser{md: goldmark.New(
		goldmark.WithExtensions(
			extension.TaskList,
			&frontmatter.Extender{Mode: frontmatter.SetMetadata},
			&hashtag.Extender{Variant: hashtag.ObsidianVariant},
			&wikilink.Extender{}))}, nil
}

func parseMarkdown(source []byte, md goldmark.Markdown) (*ast.Document, []byte, error) {
	reader := text.NewReader(source)
	root := md.Parser().Parse(reader)
	doc := root.OwnerDocument()

	return doc, source, nil
}

func extractMetadata(document *ast.Document, source []byte) (*parser.Metadata, error) {
	metadata := parser.Metadata{
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

	// TODO: Move to entry
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
