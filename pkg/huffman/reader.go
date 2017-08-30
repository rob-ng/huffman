package huffman

import (
	"bufio"
	"bytes"
	"io"
)

//#############################################################################
//# Exported
//#############################################################################

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
}

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
	}
}

func (hr *Reader) Read(p []byte) (n int, err error) {
	if !hr.headerRead {
		hr.headerRead = true
		headerDesc := hr.readHeader()
		hr.h = DeriveHeader(headerDesc)
		hr.decoder = hr.h.ExtractDecoder()
	}

	n = 0
	bi := 0
	for {
		if hr.readNext {
			b, err := hr.r.ReadByte()
			if err != nil {
				return 0, nil
			}
			hr.currByte = b
		}
		for ; hr.bitPos < 8; hr.bitPos++ {
			if n == len(p) {
				hr.readNext = false
				return n, nil
			} else if hr.totalBytes >= hr.h.numUnits {
				return n, io.EOF
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
	return
}

//#############################################################################
//# Unexported
//#############################################################################

func (hr *Reader) readHeader() string {
	var buf bytes.Buffer
	for {
		// Error handling later
		line, _ := hr.r.ReadString(headerDelim[len(headerDelim)-1])
		buf.WriteString(line)
		if len(line) == len(headerDelim) && bytes.HasSuffix(buf.Bytes(), []byte(headerDelim)) {
			bufStr := buf.String()
			return bufStr[0 : len(bufStr)-len(headerDelim)]
		}
	}
}
