package functions

import (
	"errors"
	"fmt"
	"io"
	"os"
)

// Basic cat for no flags
func DumbCat(files []string, stdin io.Reader, stdout io.Writer) error {
	var f *os.File
	var err error
	// For each file
	for _, file := range files {
		if file == "-" {
			// If file is "-" read from stdin
			_, err := io.Copy(stdout, stdin)
			if err != nil && !errors.Is(err, io.ErrClosedPipe) {
				fmt.Fprintf(os.Stderr, "Failed to copy from stdin: %v\n", err)
			}
		} else {
			// Otherwise, read from file
			f, err = os.Open(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: No such file or directory\n", file)
				continue
			}

			_, err = io.Copy(stdout, f)
			f.Close()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to copy from %q: %v\n", file, err)
			}
		}

	}
	return nil
}
