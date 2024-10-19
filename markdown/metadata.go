package markdown

import "time"

type Metadata struct {
	Names      map[string]struct{}    `json:"names"`
	Dates      map[time.Time]struct{} `json:"dates"`
	Links      map[string]struct{}    `json:"links"`
	Tags       map[string]struct{}    `json:"tags"`
	Tasks      map[string]struct{}    `json:"tasks"`
	Properties map[string]interface{} `json:"properties"`
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

func (metadata *Metadata) AddLink(link string) {
	if metadata.Links == nil {
		metadata.Links = make(map[string]struct{})
	}

	metadata.Links[link] = struct{}{}
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
