package config

import (
	"mdvault/markdown"
	"mdvault/markdown/goldmark"
)

// TODO: Expose options through config
func ConfigureParser() (markdown.Parser, error) {
	return goldmark.NewParser(), nil
}
