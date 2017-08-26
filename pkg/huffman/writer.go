package huffman

import (
	"fmt"
	"io"
)

//#############################################################################
//# Exported
//#############################################################################

// A Writer
type Writer struct {
	Header      huffmanHeader
	wroteHeader bool
	w           io.Writer
	encTable    encodingTable
	codebook    Codebook
	currByte    byte
	bitsWritten int
}

// NewWriter acts as a constructor for Writer.
func NewWriter(w io.Writer, cb Codebook) *Writer {
	encTable := make(encodingTable)
	for _, entry := range cb {
		encTable[entry.unit] = entry.code
		fmt.Printf("unit: %b, code: %s\n", entry.unit, entry.code)
	}
	return &Writer{
		Header:      huffmanHeader{},
		wroteHeader: false,
		w:           w,
		encTable:    encTable,
		codebook:    cb,
		currByte:    0,
		bitsWritten: 0,
	}
}

// Write writes an encoded form of p to the underlying io.Writer.
// Note that final value of currByte is not guaranteed to be written and hence
// calls to Write() should be followed with a call to Flush().
func (hw *Writer) Write(p []byte) (n int, err error) {
	n = 0
	if !hw.wroteHeader {
		hw.wroteHeader = true
		fmt.Printf("writing header...")
	}
	for _, b := range p {
		bits := hw.encTable[b]
		for _, bit := range bits {
			hw.currByte = (hw.currByte << 1) | byte(bit-'0')
			hw.bitsWritten++
			if hw.bitsWritten == 8 {
				hw.w.Write([]byte{hw.currByte})
				hw.currByte = 0
				hw.bitsWritten = 0
				n++
			}
		}
	}
	return
}

// Flush fills the rest of currByte with 0 bits before writing it to the
// underlying io.Writter.
func (hw *Writer) Flush() error {
	if hw.bitsWritten != 0 {
		for hw.bitsWritten < 8 {
			hw.currByte = (hw.currByte << 1) | 0
			hw.bitsWritten++
		}
		hw.w.Write([]byte{hw.currByte})
	}
	return nil
}

//#############################################################################
//# Unexported
//#############################################################################
type huffmanHeader struct {
	bitLens  []byte
	alphabet []byte
}

type encodingTable map[byte]string
