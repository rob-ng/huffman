package huffman

import (
	"bufio"
	"io"
	"strings"
	"testing"
)

func TestReadMalformedHeader(t *testing.T) {
	malformedHeaderInput := []struct {
		input string
	}{
		// Missing header entirely
		{input: "This is a body without a header"},
		// Missing line
		{input: "1 2 3 4\r\n50\r\n\r\n"},
		// Line contains values of wrong type
		{input: "a b c d\r\nabcdef\r\n10\r\n\r\n"},
		// Missing final delimiter
		{input: "0 1 2 0 4\r\nabcdefg\r\n20\r\n"},
	}

	for _, test := range malformedHeaderInput {
		read := make([]byte, len(test.input))
		r := NewReader(bufio.NewReader(strings.NewReader(test.input)))
		n, err := r.Read(read)
		if n != 0 {
			t.Errorf("Should have read 0 bytes. Instead was %d.",
				n)
		}

		if err == nil {
			t.Errorf("Should have returned an error.")
		}
	}
}

func TestReadNoBody(t *testing.T) {
	headerNoBody := "0 1 2 4\r\nabcdefg\r\n20\r\n\r\n"

	read := make([]byte, len(headerNoBody))
	r := NewReader(bufio.NewReader(strings.NewReader(headerNoBody)))
	n, err := r.Read(read)
	if n != 0 {
		t.Errorf("Should have read 0 bytes. Instead was %d", n)
	}

	if err != io.EOF {
		t.Errorf("Should have returned EOF. Instead returned %v",
			err)
	}
}
