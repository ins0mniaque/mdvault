package config

import (
	"mdvault/parser"
	"mdvault/parser/goldmark"
)

// TODO: Expose options through config
func ConfigureParser() (parser.Parser, error) {
	return goldmark.Create()
}
