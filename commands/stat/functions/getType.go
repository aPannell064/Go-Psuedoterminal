package functions

import (
	"io/fs"
	"os"
)

// Enum for file type
const (
	reg     = "regular file"
	dir     = "directory"
	symlink = "symlink"
	oth     = "other"
)

// Returns file type as a string
func getFileType(info fs.FileInfo, mode fs.FileMode) string {
	if info.IsDir() {
		// Directory
		return dir
	} else if mode.IsRegular() {
		// Regular file
		return reg
	} else if mode&os.ModeSymlink != 0 {
		// Symbolic link
		return symlink
	}

	// Other
	return oth
}
