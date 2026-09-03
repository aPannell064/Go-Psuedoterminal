package functions

import (
	"bufio"
	"fmt"
	"io"
)

func reader(r io.Reader, writer *bufio.Writer, flags []bool, linenum *int, newline *bool) error {
	reader := bufio.NewReader(r)
	lastEmpty := false
	for {
		b, err := reader.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		isEmpty := *newline && b == '\n'

		// Handle -s flag
		if flags[2] && isEmpty {
			if lastEmpty {
				continue
			}
			lastEmpty = true
		} else if b != '\n' {
			lastEmpty = false
		}

		if *newline {
			// Handle -n and -b flags
			if flags[1] {
				if !isEmpty {
					_, err = writer.WriteString(fmt.Sprintf("%6d\t", *linenum))
					if err != nil {
						return err
					}
					*linenum++
				}
			} else if flags[0] && *newline {
				_, err = writer.WriteString(fmt.Sprintf("%6d\t", *linenum))
				if err != nil {
					return err
				}
				*linenum++
			}
			*newline = false
		}

		// Handle -E flag
		if flags[3] && b == '\n' {
			_, err = writer.WriteString("$")
			if err != nil {
				return err
			}
		}

		// Handle -T flag
		if flags[4] && b == '\t' {
			_, err = writer.WriteString("^I")
			if err != nil {
				return err
			}
			continue
		}

		// Handle -v flag
		if flags[5] {
			newByte := b
			if newByte >= 128 {
				newByte &= 0x7F
				_, err = writer.WriteString("M-")
				if err != nil {
					return err
				}
			}
			if newByte < 32 && newByte != '\n' && newByte != '\t' {
				newByte += 64
				_, err = writer.WriteString("^")
			} else if newByte == 127 {
				_, err = writer.WriteString("^?")
			}
			if err != nil {
				return err
			}

			err = writer.WriteByte(newByte)
			if err != nil {
				return err
			}
		} else {
			err = writer.WriteByte(b)
			if err != nil {
				return err
			}
		}

		if b == '\n' {
			*newline = true
			if err = writer.Flush(); err != nil {
				return err
			}
		}
	}
	return writer.Flush()
}
