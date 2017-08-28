package huffman

import (
	"io"
	"strconv"
	"strings"
)

//#############################################################################
//# Exported
//#############################################################################

// A Writer
type Writer struct {
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
		encTable[entry.unit] = encodingTableEntry{entry.code, entry.codeLen}
	}
	return &Writer{
		w:           w,
		wroteHeader: false,
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
		hw.writeHeader()
	}
	for _, b := range p {
		enc := hw.encTable[b]
		bitArray := make([]int, enc.codeLen)
		for i := 0; i < enc.codeLen; i++ {
			// Get least significant bit
			bitArray[enc.codeLen-1-i] = enc.code & -enc.code
			enc.code >>= 1
		}
		for _, bit := range bitArray {
			hw.currByte = (hw.currByte << 1) | byte(bit)
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
// writeHeader writes a description of the codebook to the underlying writer.
// The first line of the header describes how many units have a particular
// code length (starting at 1).
// The second line lists the units in order of increasing code length.
func (hw *Writer) writeHeader() error {
	curr := 0
	totals := make([]string, len(hw.codebook))
	units := make([]byte, len(hw.codebook))
	for i, e := range hw.codebook {
		if curr < len(hw.codebook) {
			total := 0
			for curr < len(hw.codebook) && hw.codebook[curr].codeLen == (i+1) {
				total++
				curr++
			}
			totals[i] = strconv.Itoa(total)
		}
		units[i] = e.unit
	}
	io.WriteString(hw.w, strings.Trim(strings.Join(totals, " "), " "))
	hw.w.Write([]byte{'\n'})
	hw.w.Write(units)
	hw.w.Write([]byte{'\n'})
	return nil
}

//type encodingTable map[byte]string
type encodingTable map[byte]encodingTableEntry

type encodingTableEntry struct {
	code    int
	codeLen int
}
