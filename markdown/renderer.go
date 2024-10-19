package markdown

import "io"

type Renderer interface {
	Render(source []byte, writer io.Writer) error
}
