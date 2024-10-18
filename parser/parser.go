package parser

import (
	"time"
)

type Parser interface {
	Parse(source []byte) (*Metadata, error)
}

type Metadata struct {
	Path       string                 `json:"path"`
	Names      map[string]struct{}    `json:"names"`
	Dates      map[time.Time]struct{} `json:"dates"`
	Links      map[string]struct{}    `json:"links"`
	Tags       map[string]struct{}    `json:"tags"`
	Tasks      map[string]struct{}    `json:"tasks"`
	Properties map[string]interface{} `json:"properties"`
}

func (m *Metadata) AddName(name string) {
	m.Names[name] = struct{}{}
}

func (m *Metadata) AddDate(date time.Time) {
	m.Dates[date] = struct{}{}
}

func (m *Metadata) AddLink(link string) {
	m.Links[link] = struct{}{}
}

func (m *Metadata) AddTag(tag string) {
	m.Tags[tag] = struct{}{}
}

func (m *Metadata) AddTask(task string) {
	m.Tasks[task] = struct{}{}
}

func (m *Metadata) SetPath(path string) {
	// TODO: Extract date/name from path
	m.Path = path
	m.Names[path] = struct{}{}
}
