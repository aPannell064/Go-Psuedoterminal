package pwd

import (
	"errors"
	"fmt"
	"io"
	"os"
)

func Pwd(args []string, stdin io.Reader, stdout io.Writer) error {
	// Get state working directory and print
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cd: Could not get working directory: %v\n", err)
	}

	wd += "\n"

	// Write to stdout
	_, err = stdout.Write([]byte(wd))
	if err != nil {
		if errors.Is(err, io.ErrClosedPipe) {
			return err
		}
		return fmt.Errorf("cd: Could not write to stdout: %v\n", err)
	}
	return nil
}
