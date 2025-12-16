package app

import (
	"fmt"
	"os"
	"path/filepath"

	"sprout/internal/fsutil"
)

func Init(root string) error {
	configPath := filepath.Join(root, "sprout.toml")
	if _, err := os.Stat(configPath); err == nil {
		return fmt.Errorf("sprout.toml already exists")
	}

	configContent := `base_url = "http://localhost:1313"
pretty_urls = true

[paths]
content = "content"
templates = "templates"
static = "static"
public = "public"
`

	if err := fsutil.WriteFileAtomic(configPath, []byte(configContent)); err != nil {
		return fmt.Errorf("failed to create config: %w", err)
	}

	dirs := []string{
		filepath.Join(root, "content"),
		filepath.Join(root, "templates"),
		filepath.Join(root, "static"),
		filepath.Join(root, "public"),
	}

	for _, dir := range dirs {
		if err := fsutil.MkdirAll(dir); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	baseHTML := `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>{{.Page.Title}}</title>
</head>
<body>
    <main>
        <h1>{{.Page.Title}}</h1>
        <div>{{.Page.ContentHTML}}</div>
    </main>
</body>
</html>
`

	pageHTML := `{{define "content"}}
<h1>{{.Page.Title}}</h1>
<div>{{.Page.ContentHTML}}</div>
{{end}}
`

	basePath := filepath.Join(root, "templates", "base.html")
	if err := fsutil.WriteFileAtomic(basePath, []byte(baseHTML)); err != nil {
		return fmt.Errorf("failed to create base.html: %w", err)
	}

	pagePath := filepath.Join(root, "templates", "page.html")
	if err := fsutil.WriteFileAtomic(pagePath, []byte(pageHTML)); err != nil {
		return fmt.Errorf("failed to create page.html: %w", err)
	}

	indexContent := `+++
title = "Home"
+++

# Welcome

This is your new Sprout site. Edit this file to get started.
`

	indexPath := filepath.Join(root, "content", "index.md")
	if err := fsutil.WriteFileAtomic(indexPath, []byte(indexContent)); err != nil {
		return fmt.Errorf("failed to create index.md: %w", err)
	}

	return nil
}

