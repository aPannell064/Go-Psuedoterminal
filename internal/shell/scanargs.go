// Based off Dann's example pty code

package shell

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

func isQuote(r rune) bool {
	return r == '"' || r == '\''
}

func isOp(r rune) bool {
	return r == '|' || r == '>'
}

func ScanArgs(data []byte, atEOF bool) (advance int, token []byte, err error) {
	//start is the index of the start of the token,
	// first is the first character of the token
	start, first := 0, rune(0)

	for width := 0; start < len(data); start += width {
		// Iterate over the input data, looking for the start of a token
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		// Decode the next rune (character) from the input data
		if !unicode.IsSpace(r) {
			// If the character is not a space, it is the start of a token
			first = r
			break
		}
	}

	if isOp(first) {
		// Check for ">>"
		if first == '>' && start+1 < len(data) && data[start+1] == '>' {
			return start + 2, data[start : start+2], nil
		}
		return start + 1, data[start : start+1], nil
	}

	if isQuote(first) {
		start++ //skip the opening quote
	}

	for width, i := 0, start; i < len(data); i += width {
		// Iterate over the input data, looking for the end of the token
		var r rune

		// Decode the next rune (character) from the input data
		r, width = utf8.DecodeRune(data[i:])

		// If not in quote, and we get a space or a operator,
		// or we are in a quoted token and we encounter the closing quote,
		// then end token
		if ok := isQuote(first); (!ok && (unicode.IsSpace(r) || isOp(r))) || (ok && r == first) {

			// If operator, then don't skip it
			if !ok && isOp(r) {
				return i, data[start:i], nil
			}

			// Return the index of the next token, the current token, and no error
			return i + width, data[start:i], nil
		}
	}
	if atEOF && len(data) > start {
		// If we have reached the end of the input data and there is still a token to return,
		// return it
		if isQuote(first) {
			err = fmt.Errorf("unterminated quote: %q", first)
		}
		return len(data), data[start:], err
	}
	if isQuote(first) {
		start-- // if we are at the end of the data and we have an unterminated quote,
		// we want to include the opening quote in the token
	}
	return start, nil, nil // Return the index of the next token, no token, and no error
}
