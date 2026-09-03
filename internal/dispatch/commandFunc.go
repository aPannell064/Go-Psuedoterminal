package dispatch

import "io"

type CommandFunc func(args []string, stdin io.Reader, stdout io.Writer) error
