package main

import (
	"final/internal/shell"
	"os"
)

func main() {
	// Create a new shell instance
	sh, code := shell.New()
	if code == 0 {
		code = sh.Run()
	}

	//Exit with the returned code
	os.Exit(code)
}
