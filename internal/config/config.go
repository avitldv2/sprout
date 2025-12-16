package config

import (
	"fmt"
	"os"
	"path/filepath"

	"sprout/internal/model"

	"github.com/pelletier/go-toml/v2"
)

func Load(root string) (*model.Config, error) {
	configPath := filepath.Join(root, "sprout.toml")

	cfg := &model.Config{
		BaseURL:    "http://localhost:1313",
		PrettyURLs: true,
		UnsafeHTML: false,
	}
	cfg.Paths.Content = "content"
	cfg.Paths.Templates = "templates"
	cfg.Paths.Static = "static"
	cfg.Paths.Public = "public"

	data, err := os.ReadFile(configPath)
	if os.IsNotExist(err) {
		return cfg, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	if err := toml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return cfg, nil
}

func Resolve(cfg *model.Config, root string) (*model.ResolvedConfig, error) {
	resolved := &model.ResolvedConfig{
		BaseURL:    cfg.BaseURL,
		PrettyURLs: cfg.PrettyURLs,
		UnsafeHTML: cfg.UnsafeHTML,
		Paths: model.Paths{
			Content:   filepath.Join(root, cfg.Paths.Content),
			Templates: filepath.Join(root, cfg.Paths.Templates),
			Static:    filepath.Join(root, cfg.Paths.Static),
			Public:    filepath.Join(root, cfg.Paths.Public),
		},
	}

	if err := os.MkdirAll(resolved.Paths.Public, 0755); err != nil {
		return nil, fmt.Errorf("failed to create public directory: %w", err)
	}

	return resolved, nil
}
