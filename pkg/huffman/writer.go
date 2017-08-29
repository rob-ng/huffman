package huffman

import (
	"io"
)

//#############################################################################
//# Exported
//#############################################################################

// A Writer
type Writer struct {
	wroteHeader bool
	w           io.Writer
	h           *Header
	encoder     encoder
	currByte    byte
	bitsWritten int
}

// NewWriter acts as a constructor for Writer.
//func NewWriter(w io.Writer, cb Codebook) *Writer {
func NewWriter(w io.Writer, h *Header) *Writer {
	enc := h.ExtractEncoder()
	return &Writer{
		w:           w,
		wroteHeader: false,
		h:           h,
		encoder:     enc,
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
		//io.WriteString(hw.w, hw.h.String())
		hw.w.Write([]byte(hw.h.String()))
	}
	for _, b := range p {
		enc := hw.encoder[b]
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
