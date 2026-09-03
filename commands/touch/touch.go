package touch

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

func Touch(args []string, stdin io.Reader, stdout io.Writer) error {
	// For each argument
	for _, path := range args {
		// Check for existence
		_, err := os.Lstat(path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				// Create file if it doesn't exist
				f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					return fmt.Errorf("touch: Could not create file: '%s': %v", path, err)
				}
				f.Close()
			} else {
				return fmt.Errorf("touch: Could not determine if '%s' exists: %v", path, err)
			}
		} else {
			// Update time if it exists
			tm := time.Now()
			err = os.Chtimes(path, tm, tm)
			if err != nil {
				return fmt.Errorf("touch: Could not update timestamps on file: '%s': %v", path, err)
			}
		}
	}
	return nil
}
