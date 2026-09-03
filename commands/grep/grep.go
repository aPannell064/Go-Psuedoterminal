package grep

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
)

const reset = "\x1b[0m"
const red = "\x1b[31m"
const magenta = "\033[35m"

func Grep(args []string, stdin io.Reader, stdout io.Writer) error {
	// Ensure arguments are provided
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: grep PATTERN [FILE]...\n")
	} else {
		// Get pattern and regexp for it
		pattern := args[0]
		re, err := regexp.Compile(pattern)

		// Check for stdout is terminal for color usage
		useColor := IsTerminal(stdout)
		if err != nil {
			return fmt.Errorf("grep: failed to compile regexp '%s': %v", pattern, err)
		}

		// Read from stdin if no file args
		if len(args) == 1 {
			lineScanner := bufio.NewScanner(stdin)
			for lineScanner.Scan() {
				line := lineScanner.Text()

				// Print (in red) if there is a match
				if re.MatchString(line) {
					if useColor {
						line = re.ReplaceAllStringFunc(line, func(match string) string {
							return red + match + reset
						})
					}

					// Write line (and newline)
					_, err = stdout.Write([]byte(line + "\n"))
					if err != nil {
						if errors.Is(err, io.ErrClosedPipe) {
							return err
						}
						return fmt.Errorf("grep: failed to write to stdout: %v", err)
					}
				}
			}
		} else {
			// Get file args
			files := args[1:]
			multipleFiles := len(files) > 1
			for _, fName := range files {
				// Check file's existence
				info, err := os.Lstat(fName)
				if err != nil {
					fmt.Fprintf(os.Stderr, "grep: %s: No such file or directory\n", fName)
					continue
				}

				// Check if it is a directory
				if info.IsDir() {
					fmt.Fprintf(os.Stderr, "grep: %s Is a directory\n", fName)
					continue
				}

				// Open file
				file, err := os.Open(fName)
				if err != nil {
					fmt.Fprintf(os.Stderr, "grep: failed to open file: '%s': %v\n", fName, err)
					continue
				}

				// Scan file
				lineScanner := bufio.NewScanner(file)
				for lineScanner.Scan() {

					// Print (in red) if there is a match
					line := lineScanner.Text()
					if re.MatchString(line) {
						if useColor {
							line = re.ReplaceAllStringFunc(line, func(match string) string {
								return red + match + reset
							})
						}

						if multipleFiles {
							// Add prefix if multiple files
							prefix := fName + ":"
							if useColor {
								prefix = magenta + prefix + reset
							}

							line = prefix + line
						}

						// Write line (and newline)
						_, err = stdout.Write([]byte(line + "\n"))
						if err != nil {
							return fmt.Errorf("grep: failed to write to stdout: %v", err)
						}
					}
				}
				// Close file
				file.Close()
			}
		}

	}
	return nil
}
