package content

import (
	"bytes"
	"strings"

	"sprout/internal/model"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

func newMarkdownRenderer(unsafeHTML bool) goldmark.Markdown {
	var opts []goldmark.Option
	
	opts = append(opts, goldmark.WithExtensions(extension.GFM))
	
	if unsafeHTML {
		opts = append(opts, goldmark.WithRendererOptions(html.WithUnsafe()))
	}

	return goldmark.New(opts...)
}

func RenderMarkdown(markdown []byte, unsafeHTML bool) ([]byte, string, error) {
	md := newMarkdownRenderer(unsafeHTML)
	var buf bytes.Buffer
	if err := md.Convert(markdown, &buf); err != nil {
		return nil, "", err
	}

	html := buf.Bytes()
	title := extractFirstH1(markdown)

	return html, title, nil
}

func extractFirstH1(markdown []byte) string {
	lines := bytes.Split(markdown, []byte("\n"))
	for _, line := range lines {
		line = bytes.TrimSpace(line)
		if bytes.HasPrefix(line, []byte("# ")) {
			title := string(bytes.TrimPrefix(line, []byte("# ")))
			return strings.TrimSpace(title)
		}
	}
	return ""
}

func ParseAndRender(sourcePath string, raw []byte, unsafeHTML bool) (*model.FrontMatter, []byte, string, error) {
	fm, markdown, err := ParsePage(sourcePath, raw)
	if err != nil {
		return nil, nil, "", err
	}

	html, derivedTitle, err := RenderMarkdown(markdown, unsafeHTML)
	if err != nil {
		return nil, nil, "", err
	}

	if fm.Title == "" {
		fm.Title = derivedTitle
	}

	if fm.Slug == "" {
		fm.Slug = DeriveSlug(sourcePath)
	}

	return fm, html, derivedTitle, nil
}
