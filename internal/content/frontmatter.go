package content

import (
	"bytes"
	"fmt"
	"strings"

	"sprout/internal/model"
	"github.com/pelletier/go-toml/v2"
)

func ParsePage(sourcePath string, raw []byte) (*model.FrontMatter, []byte, error) {
	fm := &model.FrontMatter{
		Layout: "page",
	}

	content := raw
	delimiter := []byte("+++")

	if !bytes.HasPrefix(raw, delimiter) {
		return fm, raw, nil
	}

	endIdx := bytes.Index(raw[3:], delimiter)
	if endIdx == -1 {
		return nil, nil, fmt.Errorf("unclosed front matter delimiter: expected closing +++")
	}

	if endIdx+6 >= len(raw) {
		return nil, nil, fmt.Errorf("invalid front matter: content too short")
	}

	frontMatterBytes := raw[3 : endIdx+3]
	if len(frontMatterBytes) == 0 {
		return nil, nil, fmt.Errorf("empty front matter block")
	}

	content = raw[endIdx+6:]

	if err := toml.Unmarshal(frontMatterBytes, fm); err != nil {
		return nil, nil, fmt.Errorf("failed to parse front matter: %w", err)
	}

	content = bytes.TrimSpace(content)

	return fm, content, nil
}

func DeriveSlug(filepath string) string {
	filename := filepath
	if idx := strings.LastIndex(filename, "/"); idx != -1 {
		filename = filename[idx+1:]
	}
	if idx := strings.LastIndex(filename, "\\"); idx != -1 {
		filename = filename[idx+1:]
	}

	slug := filename
	if idx := strings.LastIndex(slug, "."); idx != -1 {
		slug = slug[:idx]
	}

	if slug == "index" {
		return ""
	}

	return slug
}
