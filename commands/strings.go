package commands

import (
	"fmt"
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

// IsIncluded returns whether path is included
func IsIncluded(path string, paths map[string]bool) (string, bool) {
	for k := range paths {
		if StartsWith(path, k) {
			return k, true
		}
	}

	return "", false
}

// FilterIncluded removes already included paths from map
func FilterIncluded(path string, paths map[string]bool) {
	toRemove := []string{}
	for p := range paths {
		if StartsWith(p, path) {
			toRemove = append(toRemove, p)
		}
	}

	for _, v := range toRemove {
		delete(paths, v)
		fmt.Println("Redundant path " + v + " removed.")
	}
}

// IsGoFile ...
func IsGoFile(name string) bool {
	return filepath.Ext(name) == ".go"
}
