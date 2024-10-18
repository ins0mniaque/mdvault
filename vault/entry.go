package vault

import (
	"maps"
	"mdvault/markdown"
	"slices"
	"time"
)

type Entry struct {
	Names      []string               `json:"names,omitempty"      yaml:"names,omitempty"`
	Dates      []time.Time            `json:"dates,omitempty"      yaml:"dates,omitempty"`
	Links      []string               `json:"links,omitempty"      yaml:"links,omitempty"`
	Backlinks  []string               `json:"backlinks,omitempty"  yaml:"backlinks,omitempty"`
	Tags       []string               `json:"tags,omitempty"       yaml:"tags,omitempty"`
	Tasks      []string               `json:"tasks,omitempty"      yaml:"tasks,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty" yaml:"properties,omitempty"`
}

func NewEntry(metadata *markdown.Metadata) *Entry {
	if metadata == nil {
		return &Entry{}
	}

	return &Entry{
		Names:      slices.Collect(maps.Keys(metadata.Names)),
		Dates:      slices.Collect(maps.Keys(metadata.Dates)),
		Links:      slices.Collect(maps.Keys(metadata.Links)),
		Tags:       slices.Collect(maps.Keys(metadata.Tags)),
		Tasks:      slices.Collect(maps.Keys(metadata.Tasks)),
		Properties: metadata.Properties,
	}
}
