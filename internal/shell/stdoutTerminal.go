package shell

import "os"

func stdoutIsTerminal(stdout *os.File) bool {
	fi, err := stdout.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}
