package commands

import (
	"path/filepath"
	"strings"
)

// CleanTemplateKey ...
func CleanTemplateKey(key string) string {
	return strings.Replace(key, " ", "", -1)
}

// GetDirName can also figure out file name from path
func GetDirName(path string) string {
	for i := len(path) - 2; i >= 0; i-- {
		if path[i] == '/' || path[i] == '\\' {
			return path[i+1:]
		}
	}

	return path
}

// IsGoFile ...
func IsGoFile(name string) bool {
	return filepath.Ext(name) == ".go"
}
