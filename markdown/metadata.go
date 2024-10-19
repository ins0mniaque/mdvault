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
	if len(metadata.Properties) == 0 {
		return
	}

	extractProperty(metadata.Properties, "id", metadata.SetTitle)
	extractProperty(metadata.Properties, "Id", metadata.SetTitle)
	extractProperty(metadata.Properties, "ID", metadata.SetTitle)

	extractProperty(metadata.Properties, "name", metadata.SetTitle)
	extractProperty(metadata.Properties, "Name", metadata.SetTitle)
	extractProperty(metadata.Properties, "NAME", metadata.SetTitle)

	extractProperty(metadata.Properties, "title", metadata.SetTitle)
	extractProperty(metadata.Properties, "Title", metadata.SetTitle)
	extractProperty(metadata.Properties, "TITLE", metadata.SetTitle)

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
