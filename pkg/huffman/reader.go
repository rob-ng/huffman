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
	readBuf    int
	bitsRead   int
}

func NewReader(r *bufio.Reader) *Reader {
	return &Reader{
		headerRead: false,
		r:          r,
		h:          nil,
		decoder:    nil,
		readBuf:    0,
		bitsRead:   0,
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
	for i := 0; i < len(p); i++ {
		currByte, _ := hr.r.ReadByte()
		var j uint
		for ; j < 8; j++ {
			if n >= hr.h.numUnits {
				return 0, io.EOF
			}
			currBit := int((currByte >> (7 - j)) & 1)
			hr.readBuf = (hr.readBuf << 1) | currBit
			hr.bitsRead++
			if _, ok := hr.decoder[hr.readBuf]; ok {
				if hr.bitsRead == hr.decoder[hr.readBuf].codeLen {
					p[bi] = hr.decoder[hr.readBuf].unit
					bi++
					n++
					hr.bitsRead = 0
					hr.readBuf = 0
				}
			}
		}
	}

	return n, nil
}

//#############################################################################
//# Unexported
//#############################################################################

// Before read second line via ReadBytes
func (hr *Reader) readHeader() string {
	var buf bytes.Buffer
	for {
		// Error handling later
		line, _ := hr.r.ReadString('\n')
		//line, _ := hr.r.ReadBytes('\n')
		if len(line) <= 1 {
			break
		}
		buf.WriteString(line)
		//buf.WriteBytes(line)

	}
	return buf.String()
}
