package rm

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func RM(args []string, stdin io.Reader, stdout io.Writer) error {
	// Filter the flags
	files, r, err := filterFlags(args)
	if err != nil {
		return err
	}

	// Scan from os.Stdin specifically (not stdin argument)
	lineScanner := bufio.NewScanner(os.Stdin)
	for _, file := range files {
		if r {
			// If -r prompt user in os.Stdout
			prompt := fmt.Sprintf("Remove directory '%s' and all contents? [y/N]: ", file)
			_, err = os.Stdout.WriteString(prompt)
			if err != nil {
				return fmt.Errorf("rm: Failed to write to terminal: %v", err)
			}

			// Only remove if 'y' or 'Y'
			if lineScanner.Scan() && (lineScanner.Text() == "y" || lineScanner.Text() == "Y") {
				err = os.RemoveAll(file)
				if err != nil {
					fmt.Fprintf(os.Stderr, "rm: Failed to remove '%s' recursively: %v\n", file, err)
				}
			}
		} else {
			// Check for existence and if it's a directory or not
			info, err := os.Lstat(file)
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					fmt.Fprintf(os.Stderr, "rm: no such file or directory: %s\n", file)
				} else {
					fmt.Fprintf(os.Stderr, "rm: could not access file: %s\n", file)
				}
				continue
			} else if info.IsDir() {
				fmt.Fprintf(os.Stderr, "rm: cannot remove '%s': Is a directory\n", file)
				continue
			}

			// If not -r, remove file
			err = os.Remove(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "rm: Failed to remove 'file': %v\n", err)
			}
		}
	}
	return nil
}

func filterFlags(args []string) ([]string, bool, error) {
	// Files slice for file args
	var files []string
	r := false
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			// Set the corresponding flag state
			switch arg {
			case "-r":
				r = true
			default:
				return nil, r, fmt.Errorf("invalid flag: %s", arg)
			}
		} else {
			// Add all non-flag-args to files
			files = append(files, arg)
		}
	}
	return files, r, nil
}
