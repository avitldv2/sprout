package assets

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"sprout/internal/fsutil"
)

func CopyAll(staticDir, publicDir string) error {
	if staticDir == "" || publicDir == "" {
		return fmt.Errorf("static and public directories cannot be empty")
	}

	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		return nil
	}

	return fsutil.WalkDir(staticDir, func(path string, info os.FileInfo) error {
		relPath, err := filepath.Rel(staticDir, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		if strings.Contains(relPath, "..") {
			return fmt.Errorf("invalid relative path: %s", relPath)
		}

		destPath := filepath.Join(publicDir, relPath)
		if !strings.HasPrefix(destPath, publicDir) {
			return fmt.Errorf("invalid destination path: %s", destPath)
		}

		if err := fsutil.MkdirAll(filepath.Dir(destPath)); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		if err := copyFile(path, destPath); err != nil {
			return fmt.Errorf("failed to copy file %s: %w", path, err)
		}

		return nil
	})
}

func CopyChanged(staticDir, publicDir string, changedFiles []string) error {
	if staticDir == "" || publicDir == "" {
		return fmt.Errorf("static and public directories cannot be empty")
	}

	for _, relPath := range changedFiles {
		if relPath == "" {
			continue
		}
		if strings.Contains(relPath, "..") {
			return fmt.Errorf("invalid relative path: %s", relPath)
		}

		srcPath := filepath.Join(staticDir, relPath)
		destPath := filepath.Join(publicDir, relPath)

		if !strings.HasPrefix(destPath, publicDir) {
			return fmt.Errorf("invalid destination path: %s", destPath)
		}

		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			if err := os.Remove(destPath); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("failed to remove deleted file %s: %w", destPath, err)
			}
			continue
		}

		if err := fsutil.MkdirAll(filepath.Dir(destPath)); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		if err := copyFile(srcPath, destPath); err != nil {
			return fmt.Errorf("failed to copy file %s: %w", srcPath, err)
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	data, err := io.ReadAll(srcFile)
	if err != nil {
		return err
	}

	return fsutil.WriteFileAtomic(dst, data)
}
