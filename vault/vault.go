package vault

import (
	"io/fs"
	"log"
	"mdvault/config"
	"mdvault/markdown"
	"os"
	"path/filepath"
	"sync"
)

type Vault struct {
	dir     string
	entries map[string]*markdown.Metadata
}

func NewVault(dir string) *Vault {
	return &Vault{
		dir: dir,
	}
}

func (v *Vault) Dir() string {
	return v.dir
}

func (v *Vault) Entries() map[string]*Entry {
	entries := make(map[string]*Entry, len(v.entries))
	for k, v := range v.entries {
		entries[k] = NewEntry(v)
	}

	return entries
}

func (v *Vault) IsLoaded() bool {
	return v.entries != nil
}

func (v *Vault) Load() error {
	var mutex sync.Mutex
	var wg sync.WaitGroup

	files := make(map[string]*markdown.Metadata)

	parser, err := config.ConfigureParser()
	if err != nil {
		log.Printf("Error configuring parser for vault %s: %v", v.dir, err)
		return err
	}

	err = filepath.WalkDir(v.dir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if skip, err := shouldSkip(entry); skip {
			return err
		}

		key, err := filepath.Rel(v.dir, path)
		if err != nil {
			key = path
		}

		ext := filepath.Ext(path)
		if ext == ".md" || ext == ".MD" {
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
		log.Printf("Error loading vault %s: %v", v.dir, err)
	}

	wg.Wait()

	v.entries = files

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
