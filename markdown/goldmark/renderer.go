package goldmark

import (
	"io"
	"mdvault/markdown"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	mathjax "github.com/litao91/goldmark-mathjax"
	fences "github.com/stefanfritsch/goldmark-fences"
	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/frontmatter"
	"go.abhg.dev/goldmark/hashtag"
	"go.abhg.dev/goldmark/mermaid"
	"go.abhg.dev/goldmark/wikilink"
)

type Renderer struct {
	md goldmark.Markdown
}

func (parser Renderer) Render(reader io.Reader, writer io.Writer) error {
	source, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	return parser.md.Convert(source, writer)
}

func NewRenderer() markdown.Renderer {
	return Renderer{md: goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.DefinitionList,
			extension.Footnote,
			extension.CJK,
			emoji.Emoji,
			&fences.Extender{},
			&frontmatter.Extender{Mode: frontmatter.SetMetadata},
			&hashtag.Extender{Variant: hashtag.ObsidianVariant},
			highlighting.NewHighlighting(
				highlighting.WithStyle("monokai"),
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(true),
				),
			),
			mathjax.MathJax,
			&mermaid.Extender{RenderMode: mermaid.RenderModeClient, NoScript: true},
			&wikilink.Extender{}),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		))}
}
