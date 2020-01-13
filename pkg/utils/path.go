package utils

import (
	"path/filepath"
)

func PrefixedPath(prefix, path string) string {
	filePath, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	return filepath.Join(filepath.Dir(filePath), prefix+filepath.Base(filePath))
}
