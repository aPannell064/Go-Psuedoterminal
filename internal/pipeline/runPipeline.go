package pipeline

import (
	"errors"
	"final/internal/dispatch"
	"fmt"
	"io"
	"os"
	"sync"
)

func runPipeline(stages []Stage, d *dispatch.Dispatcher, stdin io.Reader, stdout io.Writer) {
	var wg sync.WaitGroup
	var prevReader io.Reader

	// Make error channel
	errCh := make(chan error, len(stages))
	// For each stage
	for i, stage := range stages {
		var in io.Reader
		var out io.Writer

		// Set input
		if i == 0 {
			in = stdin
		} else {
			in = prevReader
		}

		// Set output
		if i == len(stages)-1 {
			out = stdout
		} else {
			r, w := io.Pipe()
			out = w
			prevReader = r
		}

		// Increment wait counter
		wg.Add(1)

		// Goroutine to dispatch
		go func(stage Stage, in io.Reader, out io.Writer) {
			defer wg.Done()

			// Close read end
			if pr, ok := in.(*io.PipeReader); ok {
				defer pr.Close()
			}

			// Close write end
			if pw, ok := out.(*io.PipeWriter); ok {
				defer pw.Close()
			}

			// Handle cd differently
			if stage.Cmd == "cd" && len(stages) > 1 {
				return
			}

			if err := d.Dispatch(stage.Cmd, stage.Args, in, out); err != nil {
				if !errors.Is(err, io.ErrClosedPipe) {
					errCh <- err
				}

			}

		}(stage, in, out)
	}

	// Goroutine to wait then close error channel
	go func() {
		wg.Wait()
		close(errCh)
	}()

	var firstErr error
	for err := range errCh {
		if firstErr == nil {
			firstErr = err
		}
	}

	if firstErr != nil {
		fmt.Fprintln(os.Stderr, firstErr)
	}
}
