package history

import (
	"os"
)

func Save(path string, entries []string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, entry := range entries {
		_, err := file.WriteString(entry + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}
