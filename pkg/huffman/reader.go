package huffman

import (
	"bufio"
	"bytes"
	"io"
)

//#############################################################################
//# Exported
//#############################################################################

// A Reader reads encoded data from an underlying bufio.Reader and decodes it.
type Reader struct {
	headerRead bool
	r          *bufio.Reader
	h          *Header
	decoder    decoder
	codeBuf    int
	bitsRead   int
	totalBytes int
	currByte   byte
	bitPos     uint
	readNext   bool
	err        error
}

// NewReader initializes and returns a new Reader.
func NewReader(r *bufio.Reader) *Reader {
	return &Reader{
		headerRead: false,
		r:          r,
		h:          nil,
		decoder:    nil,
		codeBuf:    0,
		bitsRead:   0,
		totalBytes: 0,
		currByte:   0,
		bitPos:     0,
		readNext:   true,
		err:        nil,
	}
}

// Read reads encoded data from hr's underlying bufio.Reader, decodes it, and
// reads the decoded result into p.
func (hr *Reader) Read(p []byte) (n int, err error) {
	if hr.err != nil {
		return 0, hr.err
	}

	if !hr.headerRead {
		hr.headerRead = true
		var headerDesc string
		if headerDesc, hr.err = hr.readHeader(); hr.err != nil {
			return 0, hr.err
		}
		if hr.h, hr.err = DeriveHeader(headerDesc); hr.err != nil {
			return 0, hr.err
		}
		hr.decoder = hr.h.extractDecoder()
	}

	n = 0
	bi := 0
	for {
		if hr.readNext {
			var b byte
			if b, hr.err = hr.r.ReadByte(); hr.err != nil {
				return 0, hr.err
			}
			hr.currByte = b
		}
		for ; hr.bitPos < 8; hr.bitPos++ {
			if n == len(p) {
				hr.readNext = false
				return n, hr.err
			} else if hr.totalBytes >= hr.h.numUnits {
				hr.err = io.EOF
				return n, hr.err
			}
			currBit := int((hr.currByte >> (7 - hr.bitPos)) & 1)
			hr.codeBuf = (hr.codeBuf << 1) | currBit
			hr.bitsRead++
			if _, ok := hr.decoder[hr.codeBuf]; ok {
				if hr.bitsRead == hr.decoder[hr.codeBuf].codeLen {
					p[bi] = hr.decoder[hr.codeBuf].unit
					bi++
					n++
					hr.totalBytes++
					hr.codeBuf = 0
					hr.bitsRead = 0
				}
			}
		}
		hr.readNext = true
		hr.bitPos = 0
	}
}

//#############################################################################
//# Unexported
//#############################################################################

// readHeader extracts the Header description from an encoded file.
func (hr *Reader) readHeader() (string, error) {
	var buf bytes.Buffer
	for {
		// Error handling later
		line, err := hr.r.ReadString(headerDelim[len(headerDelim)-1])
		if err != nil {
			return "", err
		}
		buf.WriteString(line)
		if len(line) == len(headerDelim) && bytes.HasSuffix(buf.Bytes(), []byte(headerDelim)) {
			bufStr := buf.String()
			return bufStr[0 : len(bufStr)-len(headerDelim)], nil
		}
	}
}
