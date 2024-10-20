package vault

import (
	"io/fs"
	"log"
	"mdvault/config"
	"mdvault/markdown"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/rjeczalik/notify"
)

type Vault struct {
	dir     string
	entries map[string]*markdown.Metadata
	watched bool
}

func NewVault(dir string) *Vault {
	return &Vault{
		dir: dir,
	}
}

func (vault *Vault) Dir() string {
	return vault.dir
}

func (vault *Vault) Backlinks() map[string][]string {
	entries := make(map[string][]string, len(vault.entries))

	for link := range vault.entries {
		backlinks := make([]string, 0)
		for backlink, metadata := range vault.entries {
			if metadata == nil {
				continue
			}

			if _, ok := metadata.Links[link]; ok {
				backlinks = append(backlinks, backlink)
			}
		}

		entries[link] = backlinks
	}

	return entries
}

func (vault *Vault) Entries() map[string]*Entry {
	backlinks := vault.Backlinks()
	entries := make(map[string]*Entry, len(vault.entries))

	for k, v := range vault.entries {
		entry := NewEntry(v)
		entry.Backlinks = backlinks[k]
		entries[k] = entry
	}

	return entries
}

func (vault *Vault) IsLoaded() bool {
	return vault.entries != nil
}

func (vault *Vault) Load() error {
	var mutex sync.Mutex
	var wg sync.WaitGroup

	files := make(map[string]*markdown.Metadata)

	parser, err := config.ConfigureParser()
	if err != nil {
		log.Printf("Error configuring parser for vault %s: %v", vault.dir, err)
		return err
	}

	err = filepath.WalkDir(vault.dir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if skip, err := shouldSkip(entry); skip {
			return err
		}

		key, err := filepath.Rel(vault.dir, path)
		if err != nil {
			key = path
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".md" {
			wg.Add(1)

			go func(markdownPath string) {
				defer wg.Done()

				metadata := parse(parser, key, markdownPath)

				mutex.Lock()
				defer mutex.Unlock()

				files[key] = metadata
			}(path)
		} else {
			mutex.Lock()
			defer mutex.Unlock()

			files[key] = nil
		}

		return nil
	})

	if err != nil {
		log.Printf("Error loading vault %s: %v", vault.dir, err)
	}

	wg.Wait()

	vault.entries = files

	return err
}

func parse(parser markdown.Parser, key string, path string) *markdown.Metadata {
	source, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Error reading file %s: %v", key, err)
		return nil
	}

	metadata, err := parser.Parse(source)
	if err != nil {
		log.Printf("Error extracting metadata for file %s: %v", key, err)
		return nil
	}

	metadata.ExtractCommonProperties()
	metadata.SetPath(key)

	return metadata
}

func shouldSkip(entry fs.DirEntry) (bool, error) {
	name := entry.Name()

	if entry.IsDir() {
		if name[0] == '.' {
			return true, filepath.SkipDir
		}

		return true, nil
	}

	if name[0] == '.' {
		return true, nil
	}

	return false, nil
}

type Event int

const (
	Create Event = iota
	Remove
	Write
	Rename
)

func (vault *Vault) Watch(onEvent func(event Event, path string)) (func(), error) {
	if vault.watched {
		log.Fatal("Vault is already being watched")
	}

	events := make(chan notify.EventInfo, 1)
	if err := notify.Watch(filepath.Join(vault.dir, "..."), events, notify.All); err != nil {
		return nil, err
	}

	vault.watched = true

	go func() {
		defer notify.Stop(events)

		for event := range events {
			switch event.Event() {
			case notify.Create:
				onEvent(Create, event.Path())
			case notify.Write:
				onEvent(Write, event.Path())
			case notify.Remove:
				onEvent(Remove, event.Path())
			case notify.Rename:
				onEvent(Remove, event.Path())
			}
		}
	}()

	return func() {
		close(events)

		vault.watched = false
	}, nil
}
