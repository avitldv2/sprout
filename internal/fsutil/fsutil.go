package fsutil

import (
	"os"
	"path/filepath"
)

func WalkDir(root string, fn func(path string, info os.FileInfo) error) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		return fn(path, info)
	})
}

func WriteFileAtomic(path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	tmpFile := path + ".tmp"
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return err
	}
	defer os.Remove(tmpFile)

	if err := os.Rename(tmpFile, path); err != nil {
		return err
	}
	return nil
}

func MkdirAll(path string) error {
	return os.MkdirAll(path, 0755)
}

type FileFingerprint struct {
	MTime int64
	Size  int64
}

func Fingerprint(path string) (FileFingerprint, error) {
	info, err := os.Stat(path)
	if err != nil {
		return FileFingerprint{}, err
	}
	return FileFingerprint{
		MTime: info.ModTime().Unix(),
		Size:  info.Size(),
	}, nil
}
