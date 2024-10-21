package goldmark

import (
	"io"
	"mdvault/markdown"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"go.abhg.dev/goldmark/frontmatter"
	"go.abhg.dev/goldmark/hashtag"
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
			extension.TaskList,
			&frontmatter.Extender{Mode: frontmatter.SetMetadata},
			&hashtag.Extender{Variant: hashtag.ObsidianVariant},
			&wikilink.Extender{}))}
}
