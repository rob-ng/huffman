package huffman

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

type Reader struct {
	headerRead bool
	r          *bufio.Reader
	readBuf    byte
	cb         Codebook
}

type Header struct {
	lenCounts []int
	alphabet  []byte
}

func NewReader(r *bufio.Reader) *Reader {
	return &Reader{
		headerRead: false,
		r:          r,
		readBuf:    0,
		cb:         nil,
	}
}

func (hr *Reader) Read(p []byte) (n int, err error) {
	if !hr.headerRead {
		hr.headerRead = true
		header := hr.readHeader()
		fmt.Printf("%v\n", header)
	}

	return 0, nil
}

func (hr *Reader) readHeader() *Header {
	h1, _ := hr.r.ReadString('\n')
	h2, _ := hr.r.ReadBytes('\n')

	h1Split := strings.Split(strings.Trim(h1, "\n"), " ")

	lenCounts := make([]int, len(h1Split))
	for i, c := range h1Split {
		lenCounts[i], _ = strconv.Atoi(c)
	}

	return &Header{
		lenCounts,
		h2,
	}
}
