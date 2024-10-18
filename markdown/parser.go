package markdown

import (
	"time"
)

type Parser interface {
	Parse(source []byte) (*Metadata, error)
}

type Metadata struct {
	Names      map[string]struct{}    `json:"names"`
	Dates      map[time.Time]struct{} `json:"dates"`
	Links      map[string]struct{}    `json:"links"`
	Tags       map[string]struct{}    `json:"tags"`
	Tasks      map[string]struct{}    `json:"tasks"`
	Properties map[string]interface{} `json:"properties"`
}

func (m *Metadata) AddName(name string) {
	if m.Names == nil {
		m.Names = make(map[string]struct{})
	}

	m.Names[name] = struct{}{}
}

func (m *Metadata) AddDate(date time.Time) {
	if m.Dates == nil {
		m.Dates = make(map[time.Time]struct{})
	}

	m.Dates[date] = struct{}{}
}

func (m *Metadata) AddLink(link string) {
	if m.Links == nil {
		m.Links = make(map[string]struct{})
	}

	m.Links[link] = struct{}{}
}

func (m *Metadata) AddTag(tag string) {
	if m.Tags == nil {
		m.Tags = make(map[string]struct{})
	}

	m.Tags[tag] = struct{}{}
}

func (m *Metadata) AddTask(task string) {
	if m.Tasks == nil {
		m.Tasks = make(map[string]struct{})
	}

	m.Tasks[task] = struct{}{}
}

func (m *Metadata) SetPath(path string) {
	// TODO: Extract date/name from path
	m.AddName(path)
}

func (m *Metadata) SetProperties(properties map[string]interface{}) {
	m.Properties = properties
}

func (m *Metadata) ExtractCommonProperties() {
	if m.Properties == nil {
		return
	}

	extractProperty(m.Properties, "id", m.AddName)
	extractProperty(m.Properties, "Id", m.AddName)
	extractProperty(m.Properties, "ID", m.AddName)

	extractProperty(m.Properties, "alias", m.AddName)
	extractProperty(m.Properties, "Alias", m.AddName)
	extractProperty(m.Properties, "ALIAS", m.AddName)
	extractProperty(m.Properties, "aliases", m.AddName)
	extractProperty(m.Properties, "Aliases", m.AddName)
	extractProperty(m.Properties, "ALIASES", m.AddName)

	extractProperty(m.Properties, "date", m.AddDate)
	extractProperty(m.Properties, "Date", m.AddDate)
	extractProperty(m.Properties, "DATE", m.AddDate)

	extractProperty(m.Properties, "time", m.AddDate)
	extractProperty(m.Properties, "Time", m.AddDate)
	extractProperty(m.Properties, "TIME", m.AddDate)

	extractProperty(m.Properties, "tag", m.AddTag)
	extractProperty(m.Properties, "Tag", m.AddTag)
	extractProperty(m.Properties, "TAG", m.AddTag)
	extractProperty(m.Properties, "tags", m.AddTag)
	extractProperty(m.Properties, "Tags", m.AddTag)
	extractProperty(m.Properties, "TAGS", m.AddTag)
}

func extractProperty[V any](properties map[string]interface{}, key string, fn func(s V)) {
	value := properties[key]
	if tag, ok := value.(V); ok {
		fn(tag)
	} else if tags, ok := value.([]interface{}); ok {
		for _, v := range tags {
			if tag, ok := v.(V); ok {
				fn(tag)
			}
		}
	}
}
