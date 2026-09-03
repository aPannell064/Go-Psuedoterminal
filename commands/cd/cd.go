package cd

import (
	"errors"
	"final/internal/shell/stateInterface"
	"fmt"
	"io"
	"os"
	"strings"
)

// Changes directory
func Cd(state stateInterface.State) func(args []string, stdin io.Reader, stdout io.Writer) error {
	// Closure so that state can be passed in
	return func(args []string, stdin io.Reader, stdout io.Writer) error {
		var path string
		var err error

		// Ensure 1 or 0 args
		if len(args) > 1 {
			return fmt.Errorf("cd: Too many arguments")
		} else if len(args) == 0 {
			// Go to home directory if no args
			path, err = os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("cd: Could not get home directory: %v", err)
			}
		} else {
			// Go to specified path if no args
			path = args[0]

			// Resolve user using '~' as home directory
			if strings.HasPrefix(path, "~") {
				home, err := os.UserHomeDir()
				if err != nil {
					return fmt.Errorf("cd: Could not change directory: %v", err)
				}
				path = strings.Replace(path, "~", home, 1)
			}

		}

		// Get stat to ensure existance
		info, err := os.Lstat(path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("cd: Directory '%s' does not exist", path)
			}
			return fmt.Errorf("cd: could not get directory info for '%s': %v", path, err)
		}

		// Also ensure it is a directory
		if !info.IsDir() {
			return fmt.Errorf("cd: '%s' is not a directory", path)
		}

		// Actually change the directory
		err = os.Chdir(path)
		if err != nil {
			return fmt.Errorf("cd: Could not change directory: %v", err)
		}

		// Get new working directory
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("cd: Could not update directory: %v", err)
		}

		// Update shell
		state.SetCWD(wd)

		return nil
	}
}
