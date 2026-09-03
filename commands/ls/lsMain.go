// Adam Pannell
// COSC 3750
// 03/05/2026

package ls

import (
	"final/commands/grep"
	"final/commands/ls/functions"
	"fmt"
	"io"
	"strings"
)

func LsMain(args []string, stdin io.Reader, stdout io.Writer) error {
	// Get the file args and flags
	files, flags, err := filterFlags(args)
	if err != nil {
		return err
	}

	// Call ls
	err = functions.LS(stdout, files, flags, grep.IsTerminal(stdout))
	if err != nil {
		return fmt.Errorf("ls: %v", err)
	}

	return nil
}

func filterFlags(args []string) ([]string, []bool, error) {
	// Set all flags to false
	a := false
	l := false
	n := false
	h := false
	r := false
	flags := []bool{a, l, n, h, r}

	// Slice for all non-flag args
	var files []string

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			// Set the corresponding flag state
			switch arg {
			case "-a":
				flags[0] = true
			case "-l":
				flags[1] = true
			case "-n":
				flags[2] = true
			case "-h":
				flags[3] = true
			case "-R":
				flags[4] = true
			default:
				return nil, nil, fmt.Errorf("invalid flag: %s", arg)
			}
		} else {
			// Add non-flag args
			files = append(files, arg)
		}
	}
	return files, flags, nil
}
