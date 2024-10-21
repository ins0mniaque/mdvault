package markdown

import "io"

type Renderer interface {
	Render(reader io.Reader, writer io.Writer) error
}
