package router

import (
	"path/filepath"
	"strings"
)

func RelPermalink(sourcePath, slug string) string {
	if slug == "" {
		relPath := sourcePath
		if strings.HasPrefix(relPath, "content/") {
			relPath = strings.TrimPrefix(relPath, "content/")
		}
		dir := filepath.Dir(relPath)
		if dir == "." || dir == "" {
			return "/"
		}
		dir = filepath.ToSlash(dir)
		return "/" + dir + "/"
	}

	relPath := sourcePath
	if strings.HasPrefix(relPath, "content/") {
		relPath = strings.TrimPrefix(relPath, "content/")
	}

	if idx := strings.LastIndex(relPath, "."); idx != -1 {
		relPath = relPath[:idx]
	}

	dir := filepath.Dir(relPath)
	dir = filepath.ToSlash(dir)

	if dir == "." || dir == "" {
		return "/" + slug + "/"
	}

	return "/" + dir + "/" + slug + "/"
}

func OutputPath(publicDir, relPermalink string) string {
	if relPermalink == "/" {
		return filepath.Join(publicDir, "index.html")
	}

	path := strings.Trim(relPermalink, "/")
	return filepath.Join(publicDir, path, "index.html")
}

func SimpleRelPermalink(slug string) string {
	if slug == "" {
		return "/"
	}
	return "/" + slug + "/"
}
