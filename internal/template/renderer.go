package template

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"sprout/internal/model"
)

type Renderer struct {
	templates *template.Template
}

func NewRenderer(templatesDir string) (*Renderer, error) {
	if templatesDir == "" {
		return nil, fmt.Errorf("templates directory cannot be empty")
	}

	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("templates directory does not exist: %s", templatesDir)
	}

	tmpl := template.New("base")

	var templateFiles []string
	err := filepath.Walk(templatesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".html") {
			templateFiles = append(templateFiles, path)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk templates directory: %w", err)
	}

	if len(templateFiles) == 0 {
		return nil, fmt.Errorf("no template files found in %s", templatesDir)
	}

	for _, path := range templateFiles {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read template %s: %w", path, err)
		}

		if len(data) == 0 {
			return nil, fmt.Errorf("template file is empty: %s", path)
		}

		if _, err := tmpl.Parse(string(data)); err != nil {
			return nil, fmt.Errorf("failed to parse template %s: %w", path, err)
		}
	}

	return &Renderer{templates: tmpl}, nil
}

func (r *Renderer) RenderPage(site model.Site, page model.Page) ([]byte, error) {
	templateName := page.Layout
	if templateName == "" {
		templateName = "page"
	}

	var t *template.Template
	if r.templates.Lookup("base.html") != nil {
		t = r.templates.Lookup("base.html")
	} else {
		layoutTemplateName := templateName + ".html"
		if r.templates.Lookup(layoutTemplateName) != nil {
			t = r.templates.Lookup(layoutTemplateName)
		} else if r.templates.Lookup("page.html") != nil {
			t = r.templates.Lookup("page.html")
		} else {
			t = r.templates
		}
	}

	var buf strings.Builder
	if err := t.Execute(&buf, map[string]interface{}{
		"Site": site,
		"Page": page,
	}); err != nil {
		return nil, fmt.Errorf("failed to render template: %w", err)
	}

	return []byte(buf.String()), nil
}
