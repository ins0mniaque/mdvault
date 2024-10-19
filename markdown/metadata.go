package markdown

import (
	"net/url"
	"path"
	"path/filepath"
	"regexp"
	"strings"
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

var dateRegexp = regexp.MustCompile(`\d{4}-\d{2}-\d{2}|\d{8}`)

func parseDateMatch(match string) (time.Time, error) {
	layout := "2006-01-02"
	if len(match) == 8 {
		layout = "20060102"
	}

	return time.Parse(layout, match)
}

func (metadata *Metadata) addDatesFrom(value string) {
	for _, match := range dateRegexp.FindAllString(value, -1) {
		if date, err := parseDateMatch(match); err == nil {
			metadata.AddDate(date)
		}
	}
}

func (metadata *Metadata) AddURL(rawURL string) {
	metadata.addURL(rawURL)

	URL, err := url.Parse(rawURL)
	if err != nil || URL.IsAbs() {
		return
	}

	URL.Fragment = ""
	rawURL = URL.String()
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
	name := filepath.Base(path)
	name = strings.TrimSuffix(name, filepath.Ext(name))

	metadata.AddName(name)
	metadata.addDatesFrom(path)
}

func (metadata *Metadata) SetTitle(title string) {
	metadata.AddName(title)
	metadata.addDatesFrom(title)
}

func (metadata *Metadata) SetProperties(properties map[string]interface{}) {
	metadata.Properties = properties
}

func (metadata *Metadata) ExtractCommonProperties() {
	properties := metadata.Properties
	if len(properties) == 0 {
		return
	}

	extractProperty(properties, "id", metadata.SetTitle)
	extractProperty(properties, "Id", metadata.SetTitle)
	extractProperty(properties, "ID", metadata.SetTitle)

	extractProperty(properties, "name", metadata.SetTitle)
	extractProperty(properties, "Name", metadata.SetTitle)
	extractProperty(properties, "NAME", metadata.SetTitle)

	extractProperty(properties, "title", metadata.SetTitle)
	extractProperty(properties, "Title", metadata.SetTitle)
	extractProperty(properties, "TITLE", metadata.SetTitle)

	extractProperty(properties, "alias", metadata.AddName)
	extractProperty(properties, "Alias", metadata.AddName)
	extractProperty(properties, "ALIAS", metadata.AddName)
	extractProperty(properties, "aliases", metadata.AddName)
	extractProperty(properties, "Aliases", metadata.AddName)
	extractProperty(properties, "ALIASES", metadata.AddName)

	extractProperty(properties, "date", metadata.AddDate)
	extractProperty(properties, "Date", metadata.AddDate)
	extractProperty(properties, "DATE", metadata.AddDate)

	extractProperty(properties, "time", metadata.AddDate)
	extractProperty(properties, "Time", metadata.AddDate)
	extractProperty(properties, "TIME", metadata.AddDate)

	extractProperty(properties, "tag", metadata.AddTag)
	extractProperty(properties, "Tag", metadata.AddTag)
	extractProperty(properties, "TAG", metadata.AddTag)
	extractProperty(properties, "tags", metadata.AddTag)
	extractProperty(properties, "Tags", metadata.AddTag)
	extractProperty(properties, "TAGS", metadata.AddTag)
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
