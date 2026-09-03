package history

import (
	"bufio"
	"os"
)

func Load(path string) ([]string, error) {
	var entries []string
	file, err := os.Open(path)
	if err != nil {
		return entries, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		entries = Append(entries, scanner.Text())
	}
	return entries, scanner.Err()
}
