package goldmark

import (
	"mdvault/markdown"
	"net/url"
	"path"

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

func (parser Parser) Parse(source []byte) (*markdown.Metadata, error) {
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

func parseLink(source []byte) string {
	u, err := url.Parse(string(source))
	if err != nil || u.IsAbs() {
		return ""
	}

	u.Fragment = ""
	link := u.String()
	if link == "" {
		return ""
	}

	if path.Ext(link) == "" {
		link = link + ".md"
	}
	return link
}

func extractMetadata(document *ast.Document, source []byte) (*markdown.Metadata, error) {
	metadata := markdown.Metadata{}
	metadata.SetProperties(document.Meta())

	title := ""

	err := ast.Walk(document, func(node ast.Node, enter bool) (ast.WalkStatus, error) {
		if title == "" {
			if n, ok := node.(*ast.Heading); ok && n.Level == 1 && enter {
				title = string(n.Text(source))
				metadata.AddName(title)
			}
		}

		if n, ok := node.(*ast.Link); ok && enter {
			if link := parseLink(n.Destination); link != "" {
				metadata.AddLink(link)
			}
		}

		if n, ok := node.(*wikilink.Node); ok && enter {
			if link := parseLink(n.Target); link != "" {
				metadata.AddLink(link)
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

	return &metadata, nil
}
