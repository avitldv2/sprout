package cache

import (
	"os"
	"path/filepath"

	"sprout/internal/fsutil"
)

func CreateSnapshot(contentDir, templatesDir, staticDir string) (*Snapshot, error) {
	snapshot := &Snapshot{
		ContentFiles:  make(map[string]FileEntry),
		TemplateFiles: make(map[string]FileEntry),
		StaticFiles:   make(map[string]FileEntry),
		ConfigFile:    nil,
	}

	if err := snapshotDir(contentDir, snapshot.ContentFiles); err != nil {
		return nil, err
	}

	if err := snapshotDir(templatesDir, snapshot.TemplateFiles); err != nil {
		return nil, err
	}

	if err := snapshotDir(staticDir, snapshot.StaticFiles); err != nil {
		return nil, err
	}

	return snapshot, nil
}

func CreateSnapshotWithConfig(contentDir, templatesDir, staticDir, configPath string) (*Snapshot, error) {
	snapshot, err := CreateSnapshot(contentDir, templatesDir, staticDir)
	if err != nil {
		return nil, err
	}
	
	configFile, err := SnapshotConfigFile(configPath)
	if err != nil {
		return nil, err
	}
	snapshot.ConfigFile = configFile
	
	return snapshot, nil
}

func snapshotDir(dir string, entries map[string]FileEntry) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil
	}
	return fsutil.WalkDir(dir, func(path string, info os.FileInfo) error {
		fp, err := fsutil.Fingerprint(path)
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		relPath = filepath.ToSlash(relPath)

		entries[relPath] = FileEntry{
			MTime: fp.MTime,
			Size:  fp.Size,
		}

		return nil
	})
}

func Diff(cache *Cache, snapshot *Snapshot) *Plan {
	plan := &Plan{
		PagesToRebuild:  []string{},
		AssetsToCopy:     []string{},
		TemplatesChanged: false,
	}

	// Check if config file changed
	configChanged := false
	if snapshot.ConfigFile != nil {
		if cache.ConfigFile == nil {
			configChanged = true
		} else if cache.ConfigFile.MTime != snapshot.ConfigFile.MTime || cache.ConfigFile.Size != snapshot.ConfigFile.Size {
			configChanged = true
		}
	} else if cache.ConfigFile != nil {
		configChanged = true
	}

	for path, entry := range snapshot.TemplateFiles {
		oldEntry, exists := cache.TemplateFiles[path]
		if !exists || oldEntry.MTime != entry.MTime || oldEntry.Size != entry.Size {
			plan.TemplatesChanged = true
			break
		}
	}

	for path := range cache.TemplateFiles {
		if _, exists := snapshot.TemplateFiles[path]; !exists {
			plan.TemplatesChanged = true
			break
		}
	}

	// If config or templates changed, rebuild all pages
	if plan.TemplatesChanged || configChanged {
		for path := range snapshot.ContentFiles {
			plan.PagesToRebuild = append(plan.PagesToRebuild, path)
		}
	} else {
		for path, entry := range snapshot.ContentFiles {
			oldEntry, exists := cache.ContentFiles[path]
			if !exists || oldEntry.MTime != entry.MTime || oldEntry.Size != entry.Size {
				plan.PagesToRebuild = append(plan.PagesToRebuild, path)
			}
		}
	}

	for path, entry := range snapshot.StaticFiles {
		oldEntry, exists := cache.StaticFiles[path]
		if !exists || oldEntry.MTime != entry.MTime || oldEntry.Size != entry.Size {
			plan.AssetsToCopy = append(plan.AssetsToCopy, path)
		}
	}

	for path := range cache.StaticFiles {
		if _, exists := snapshot.StaticFiles[path]; !exists {
			plan.AssetsToCopy = append(plan.AssetsToCopy, path)
		}
	}

	return plan
}
