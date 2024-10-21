package markdown

import "io"

type Parser interface {
	Parse(reader io.Reader) (*Metadata, error)
}
