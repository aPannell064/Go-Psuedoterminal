package rmdir

import (
	"fmt"
	"io"
	"os"
)

func RMDir(args []string, stdin io.Reader, stdout io.Writer) error {
	// For all arguments
	for _, dir := range args {
		// Get file info to check for existence
		info, err := os.Lstat(dir)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "rmdir: Directory '%s' does not exist\n", dir)
			} else {
				fmt.Fprintf(os.Stderr, "rmdir: could not get directory info for '%s': %v\n", dir, err)
			}
			continue
		} else if !info.IsDir() {
			// Don't remove if not a directory
			fmt.Fprintf(os.Stderr, "rmdir: '%s' is not a directory\n", dir)
			continue
		}

		// Attempt to remove
		err = os.Remove(dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "rmdir: '%s': Directory not empty\n", dir)
			continue
		}
	}
	return nil
}
