package content

import (
	"bytes"
	"strings"

	"sprout/internal/model"

	"github.com/yuin/goldmark"
)

var md = goldmark.New(
	goldmark.WithExtensions(),
)

func RenderMarkdown(markdown []byte) ([]byte, string, error) {
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

func ParseAndRender(sourcePath string, raw []byte) (*model.FrontMatter, []byte, string, error) {
	fm, markdown, err := ParsePage(sourcePath, raw)
	if err != nil {
		return nil, nil, "", err
	}

	html, derivedTitle, err := RenderMarkdown(markdown)
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
