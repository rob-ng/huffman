package huffman

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

type Reader struct {
	headerRead bool
	r          *bufio.Reader
	readBuf    int
	bitsRead   int
	decTable   decodeTable
}

type Header struct {
	lenCounts []int
	alphabet  []byte
}

func NewReader(r *bufio.Reader) *Reader {
	dt := make(decodeTable)
	return &Reader{
		headerRead: false,
		r:          r,
		readBuf:    0,
		bitsRead:   0,
		decTable:   dt,
	}
}

func (hr *Reader) Read(p []byte) (n int, err error) {
	if !hr.headerRead {
		hr.headerRead = true
		header := hr.readHeader()
		cb := DeriveCanonicalCB(header)
		for _, entry := range cb {
			hr.decTable[entry.code] = &decodeTableEntry{
				unit:    entry.unit,
				codeLen: entry.codeLen,
			}
		}
	}

	n = 0
	bi := 0
	for i := 0; i < len(p); i++ {
		currByte, _ := hr.r.ReadByte()
		var j uint
		for ; j < 8; j++ {
			currBit := int((currByte >> (7 - j)) & 1)
			hr.readBuf = (hr.readBuf << 1) | currBit
			hr.bitsRead++
			if n >= len(p) {
				return 0, io.EOF
			}
			if _, ok := hr.decTable[hr.readBuf]; ok {
				if hr.bitsRead == hr.decTable[hr.readBuf].codeLen {
					p[bi] = hr.decTable[hr.readBuf].unit
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

type decodeTable map[int]*decodeTableEntry

type decodeTableEntry struct {
	unit    byte
	codeLen int
}

func (hr *Reader) readHeader() *Header {
	h1, _ := hr.r.ReadString('\n')
	//h1, _ := hr.r.ReadBytes('\n')
	h2, _ := hr.r.ReadBytes('\n')

	h1Split := strings.Split(strings.Trim(h1, "\n"), " ")
	h2 = h2[0 : len(h2)-1]

	lenCounts := make([]int, len(h1Split))
	for i, c := range h1Split {
		lenCounts[i], _ = strconv.Atoi(c)
	}

	return &Header{
		lenCounts,
		h2,
	}
	return nil
}
