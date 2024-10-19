package markdown

type Parser interface {
	Parse(source []byte) (*Metadata, error)
}
