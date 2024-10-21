package goldmark

import (
	"io"
	"mdvault/markdown"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/frontmatter"
	"go.abhg.dev/goldmark/hashtag"
	"go.abhg.dev/goldmark/wikilink"
)

type Parser struct {
	md goldmark.Markdown
}

func (parser Parser) Parse(reader io.Reader) (*markdown.Metadata, error) {
	source, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	document, source, err := parseMarkdown(source, parser.md)
	if err != nil {
		return nil, err
	}

	metadata, err := extractMetadata(document, source)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

func NewParser() markdown.Parser {
	return Parser{md: goldmark.New(
		goldmark.WithExtensions(
			extension.TaskList,
			&frontmatter.Extender{Mode: frontmatter.SetMetadata},
			&hashtag.Extender{Variant: hashtag.ObsidianVariant},
			&wikilink.Extender{}))}
}

func parseMarkdown(source []byte, md goldmark.Markdown) (*ast.Document, []byte, error) {
	reader := text.NewReader(source)
	root := md.Parser().Parse(reader)
	doc := root.OwnerDocument()

	return doc, source, nil
}

func extractMetadata(document *ast.Document, source []byte) (*markdown.Metadata, error) {
	metadata := markdown.Metadata{}
	metadata.SetProperties(document.Meta())

	title := ""

	err := ast.Walk(document, func(node ast.Node, enter bool) (ast.WalkStatus, error) {
		if n, ok := node.(*ast.Heading); ok && enter {
			section := string(n.Text(source))

			metadata.AddSection(section)

			if title == "" && n.Level == 1 {
				title = section
				metadata.SetTitle(title)
			}
		}

		if n, ok := node.(*ast.Link); ok && enter {
			metadata.AddURL(string(n.Destination))
		}

		if n, ok := node.(*wikilink.Node); ok && enter {
			link := string(n.Target)
			if len(n.Fragment) > 0 {
				link = link + "#" + string(n.Fragment)
			}

			metadata.AddURL(link)
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

	return &metadata, nil
}
