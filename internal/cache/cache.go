package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"sprout/internal/fsutil"
)

type FileEntry struct {
	MTime int64 `json:"mtime"`
	Size  int64 `json:"size"`
}

type Cache struct {
	ContentFiles  map[string]FileEntry `json:"content_files"`
	TemplateFiles map[string]FileEntry `json:"template_files"`
	StaticFiles   map[string]FileEntry `json:"static_files"`
	BuildTime     time.Time             `json:"build_time"`
}

type Snapshot struct {
	ContentFiles  map[string]FileEntry
	TemplateFiles map[string]FileEntry
	StaticFiles   map[string]FileEntry
}

type Plan struct {
	PagesToRebuild  []string
	AssetsToCopy    []string
	TemplatesChanged bool
}

func LoadCache(path string) (*Cache, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &Cache{
			ContentFiles:  make(map[string]FileEntry),
			TemplateFiles: make(map[string]FileEntry),
			StaticFiles:   make(map[string]FileEntry),
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read cache: %w", err)
	}

	var cache Cache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, fmt.Errorf("failed to parse cache: %w", err)
	}

	if cache.ContentFiles == nil {
		cache.ContentFiles = make(map[string]FileEntry)
	}
	if cache.TemplateFiles == nil {
		cache.TemplateFiles = make(map[string]FileEntry)
	}
	if cache.StaticFiles == nil {
		cache.StaticFiles = make(map[string]FileEntry)
	}

	return &cache, nil
}

func SaveCache(path string, cache *Cache) error {
	cache.BuildTime = time.Now()
	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cache: %w", err)
	}

	return fsutil.WriteFileAtomic(path, data)
}
