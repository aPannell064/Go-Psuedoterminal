package functions

import (
	"os"
	"strings"
)

// Removes any hidden files from the dir listing
func dirFilter(entries []os.DirEntry) []os.DirEntry {
	filtered := []os.DirEntry{}
	for _, entry := range entries {
		if !strings.HasPrefix(entry.Name(), ".") {
			filtered = append(filtered, entry)
		}
	}

	return filtered
}
