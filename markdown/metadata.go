package markdown

import (
	"net/url"
	"path"
	"time"
)

type Metadata struct {
	Names      map[string]struct{}
	Dates      map[time.Time]struct{}
	Links      map[string]struct{}
	URLs       map[string]struct{}
	Sections   map[string]struct{}
	Tags       map[string]struct{}
	Tasks      map[string]struct{}
	Properties map[string]interface{}
}

func (metadata *Metadata) AddName(name string) {
	if metadata.Names == nil {
		metadata.Names = make(map[string]struct{})
	}

	metadata.Names[name] = struct{}{}
}

func (metadata *Metadata) AddDate(date time.Time) {
	if metadata.Dates == nil {
		metadata.Dates = make(map[time.Time]struct{})
	}

	metadata.Dates[date] = struct{}{}
}

func (metadata *Metadata) AddURL(rawURL string) {
	metadata.addURL(rawURL)

	u, err := url.Parse(rawURL)
	if err != nil || u.IsAbs() {
		return
	}

	u.Fragment = ""
	rawURL = u.String()
	if rawURL == "" {
		return
	}

	if path.Ext(rawURL) == "" {
		rawURL = rawURL + ".md"
	}

	metadata.addLink(rawURL)
}

func (metadata *Metadata) addLink(link string) {
	if metadata.Links == nil {
		metadata.Links = make(map[string]struct{})
	}

	metadata.Links[link] = struct{}{}
}

func (metadata *Metadata) addURL(url string) {
	if metadata.URLs == nil {
		metadata.URLs = make(map[string]struct{})
	}

	metadata.URLs[url] = struct{}{}
}

func (metadata *Metadata) AddSection(section string) {
	if metadata.Sections == nil {
		metadata.Sections = make(map[string]struct{})
	}

	metadata.Sections[section] = struct{}{}
}

func (metadata *Metadata) AddTag(tag string) {
	if metadata.Tags == nil {
		metadata.Tags = make(map[string]struct{})
	}

	metadata.Tags[tag] = struct{}{}
}

func (metadata *Metadata) AddTask(task string) {
	if metadata.Tasks == nil {
		metadata.Tasks = make(map[string]struct{})
	}

	metadata.Tasks[task] = struct{}{}
}

func (metadata *Metadata) SetPath(path string) {
	// TODO: Extract date/name from path
	metadata.AddName(path)
}

func (metadata *Metadata) SetProperties(properties map[string]interface{}) {
	metadata.Properties = properties
}

func (metadata *Metadata) ExtractCommonProperties() {
	if metadata.Properties == nil {
		return
	}

	extractProperty(metadata.Properties, "id", metadata.AddName)
	extractProperty(metadata.Properties, "Id", metadata.AddName)
	extractProperty(metadata.Properties, "ID", metadata.AddName)

	extractProperty(metadata.Properties, "alias", metadata.AddName)
	extractProperty(metadata.Properties, "Alias", metadata.AddName)
	extractProperty(metadata.Properties, "ALIAS", metadata.AddName)
	extractProperty(metadata.Properties, "aliases", metadata.AddName)
	extractProperty(metadata.Properties, "Aliases", metadata.AddName)
	extractProperty(metadata.Properties, "ALIASES", metadata.AddName)

	extractProperty(metadata.Properties, "date", metadata.AddDate)
	extractProperty(metadata.Properties, "Date", metadata.AddDate)
	extractProperty(metadata.Properties, "DATE", metadata.AddDate)

	extractProperty(metadata.Properties, "time", metadata.AddDate)
	extractProperty(metadata.Properties, "Time", metadata.AddDate)
	extractProperty(metadata.Properties, "TIME", metadata.AddDate)

	extractProperty(metadata.Properties, "tag", metadata.AddTag)
	extractProperty(metadata.Properties, "Tag", metadata.AddTag)
	extractProperty(metadata.Properties, "TAG", metadata.AddTag)
	extractProperty(metadata.Properties, "tags", metadata.AddTag)
	extractProperty(metadata.Properties, "Tags", metadata.AddTag)
	extractProperty(metadata.Properties, "TAGS", metadata.AddTag)
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
