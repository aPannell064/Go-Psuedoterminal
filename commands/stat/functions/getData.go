package functions

import (
	"fmt"
	"os"
)

// Gets the data for stat
func GetData(file string) ([]string, error) {
	// Get file info
	info, err := os.Lstat(file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%s: No such file or directory", file)
		} else {
			return nil, fmt.Errorf("%s: Cannot access file", file)
		}
	}

	// Get mode, and modified time
	mode := info.Mode()
	modified := info.ModTime().Format("2006-01-02 15:04:05")

	// Format strings for each field
	fileString := fmt.Sprintf("%8s: %s\n", "File", file)
	sizeString := fmt.Sprintf("%8s: %d\n", "Size", info.Size())
	modeString := fmt.Sprintf("%8s: %04o (%s)\n", "Mode", mode.Perm(), mode.String())
	modifiedString := fmt.Sprintf("%8s: %s\n", "Modified", modified)
	typeString := fmt.Sprintf("%8s: %s\n", "Type", getFileType(info, mode))

	// Return a string with all the data
	return []string{fileString, sizeString, modeString, modifiedString, typeString}, nil
}
