package cat

import (
	"final/commands/cat/functions"
	"fmt"
	"io"
	"strings"
)

func any(flags []bool) bool {
	for _, flag := range flags {
		if flag {
			return true
		}
	}
	return false
}

func CatMain(args []string, stdin io.Reader, stdout io.Writer) error {
	// Filter flags from file args
	files, flags, err := filterFlags(args)
	if err != nil {
		return fmt.Errorf("cat: %v", err)
	}

	// If no files, set files to "-" for stdin
	if len(files) == 0 {
		files = []string{"-"}
	}

	// Cat if flags, dumbcat if not
	if anyFlags(flags) {
		err = functions.Cat(files, flags, stdin, stdout)
	} else {
		err = functions.DumbCat(files, stdin, stdout)
	}
	if err != nil {
		return fmt.Errorf("cat: %v", err)
	}

	return nil
}

func filterFlags(args []string) ([]string, []bool, error) {
	// Set all flags to false
	n := false
	b := false
	s := false
	e := false
	t := false
	v := false
	flags := []bool{n, b, s, e, t, v}

	// Slice for all non-flag args
	var files []string

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			// Set the corresponding flag state or append if it is stdin
			switch arg {
			case "-n":
				flags[0] = true
			case "-b":
				flags[1] = true
			case "-s":
				flags[2] = true
			case "-E":
				flags[3] = true
			case "-T":
				flags[4] = true
			case "-v":
				flags[5] = true
			case "-e":
				flags[3] = true
				flags[5] = true
			case "-t":
				flags[4] = true
				flags[5] = true
			case "-A":
				flags[3] = true
				flags[4] = true
				flags[5] = true
			case "-":
				files = append(files, arg)
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

// Check if any flags are marked
func anyFlags(flags []bool) bool {
	for _, flag := range flags {
		if flag {
			return true
		}
	}
	return false
}
