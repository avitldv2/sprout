package app

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"sprout/internal/assets"
	"sprout/internal/cache"
	"sprout/internal/config"
	"sprout/internal/content"
	"sprout/internal/fsutil"
	"sprout/internal/logx"
	"sprout/internal/model"
	"sprout/internal/router"
	tmpl "sprout/internal/template"
)

type BuildResult struct {
	BuiltPages   int
	CopiedAssets int
}

func Build(root string, clean bool) (*BuildResult, error) {
	if root == "" {
		return nil, fmt.Errorf("root directory cannot be empty")
	}

	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("invalid root path: %w", err)
	}

	if _, err := os.Stat(rootAbs); os.IsNotExist(err) {
		return nil, fmt.Errorf("root directory does not exist: %s", rootAbs)
	}

	cfg, err := config.Load(rootAbs)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	resolved, err := config.Resolve(cfg, rootAbs)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve config: %w", err)
	}

	if clean {
		if err := os.RemoveAll(resolved.Paths.Public); err != nil {
			return nil, fmt.Errorf("failed to clean public directory: %w", err)
		}
		if err := fsutil.MkdirAll(resolved.Paths.Public); err != nil {
			return nil, fmt.Errorf("failed to create public directory: %w", err)
		}
	}

	if _, err := os.Stat(resolved.Paths.Templates); os.IsNotExist(err) {
		return nil, fmt.Errorf("templates directory does not exist: %s", resolved.Paths.Templates)
	}

	renderer, err := tmpl.NewRenderer(resolved.Paths.Templates)
	if err != nil {
		return nil, fmt.Errorf("failed to create renderer: %w", err)
	}

	site := model.Site{
		BaseURL: resolved.BaseURL,
	}

	result := &BuildResult{}

	cachePath := filepath.Join(resolved.Paths.Public, ".sprout-cache.json")
	oldCache, err := cache.LoadCache(cachePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load cache: %w", err)
	}

	configPath := filepath.Join(rootAbs, "sprout.toml")
	snapshot, err := cache.CreateSnapshotWithConfig(resolved.Paths.Content, resolved.Paths.Templates, resolved.Paths.Static, configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create snapshot: %w", err)
	}

	var plan *cache.Plan
	if len(oldCache.ContentFiles) == 0 && len(oldCache.TemplateFiles) == 0 {
		plan = &cache.Plan{
			PagesToRebuild:  make([]string, 0),
			AssetsToCopy:     make([]string, 0),
			TemplatesChanged: false,
		}
		for path := range snapshot.ContentFiles {
			plan.PagesToRebuild = append(plan.PagesToRebuild, path)
		}
		for path := range snapshot.StaticFiles {
			plan.AssetsToCopy = append(plan.AssetsToCopy, path)
		}
	} else {
		plan = cache.Diff(oldCache, snapshot)
	}

	logx.Infof("Templates changed: %v, Pages to rebuild: %d, Assets to copy: %d",
		plan.TemplatesChanged, len(plan.PagesToRebuild), len(plan.AssetsToCopy))

	if len(plan.PagesToRebuild) > 0 {
		for _, relPath := range plan.PagesToRebuild {
			if relPath == "" {
				continue
			}
			contentPath := filepath.Join(resolved.Paths.Content, relPath)
			contentPath = filepath.Clean(contentPath)
			if !strings.HasPrefix(contentPath, resolved.Paths.Content) {
				return nil, fmt.Errorf("invalid content path: %s", contentPath)
			}
			if err := buildPage(contentPath, resolved, site, renderer); err != nil {
				return nil, fmt.Errorf("failed to build page %s: %w", contentPath, err)
			}
			result.BuiltPages++
		}
	}

	if len(plan.AssetsToCopy) > 0 {
		if err := assets.CopyChanged(resolved.Paths.Static, resolved.Paths.Public, plan.AssetsToCopy); err != nil {
			return nil, fmt.Errorf("failed to copy assets: %w", err)
		}
		result.CopiedAssets = len(plan.AssetsToCopy)
	}

	newCache := &cache.Cache{
		ContentFiles:  snapshot.ContentFiles,
		TemplateFiles: snapshot.TemplateFiles,
		StaticFiles:   snapshot.StaticFiles,
		ConfigFile:    snapshot.ConfigFile,
	}
	if err := cache.SaveCache(cachePath, newCache); err != nil {
		return nil, fmt.Errorf("failed to save cache: %w", err)
	}

	logx.Infof("Built %d pages, copied %d assets", result.BuiltPages, result.CopiedAssets)

	return result, nil
}

func buildPage(contentPath string, resolved *model.ResolvedConfig, site model.Site, renderer *tmpl.Renderer) error {
	info, err := os.Stat(contentPath)
	if err != nil {
		return fmt.Errorf("failed to stat content file: %w", err)
	}
	if info.Size() == 0 {
		return fmt.Errorf("content file is empty: %s", contentPath)
	}

	raw, err := os.ReadFile(contentPath)
	if err != nil {
		return fmt.Errorf("failed to read content: %w", err)
	}

	if len(raw) == 0 {
		return fmt.Errorf("content file is empty: %s", contentPath)
	}

	fm, html, _, err := content.ParseAndRender(contentPath, raw, resolved.UnsafeHTML)
	if err != nil {
		return fmt.Errorf("failed to parse content: %w", err)
	}

	relPath, err := filepath.Rel(resolved.Paths.Content, contentPath)
	if err != nil {
		return fmt.Errorf("failed to get relative path: %w", err)
	}
	relPermalink := router.RelPermalink(relPath, fm.Slug)
	if relPermalink == "" {
		relPermalink = "/"
	}

	page := model.Page{
		Title:        fm.Title,
		Slug:         fm.Slug,
		Layout:       fm.Layout,
		RelPermalink: relPermalink,
		ContentHTML:  template.HTML(html),
		SourcePath:   contentPath,
	}

	output, err := renderer.RenderPage(site, page)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	outputPath := router.OutputPath(resolved.Paths.Public, relPermalink)

	if err := fsutil.WriteFileAtomic(outputPath, output); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	logx.Infof("Built: %s -> %s", contentPath, relPermalink)

	return nil
}
