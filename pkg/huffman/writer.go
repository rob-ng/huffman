package huffman

import (
	"io"
)

//#############################################################################
//# Exported
//#############################################################################

// A Writer writes encoded data to an underlying io.Writer.
type Writer struct {
	wroteHeader bool
	w           io.Writer
	h           *Header
	encoder     encoder
	currByte    byte
	bitsWritten int
	err         error
}

// NewWriter returns a new Writer.
func NewWriter(w io.Writer, h *Header) *Writer {
	enc := h.ExtractEncoder()
	return &Writer{
		w:           w,
		wroteHeader: false,
		h:           h,
		encoder:     enc,
		currByte:    0,
		bitsWritten: 0,
		err:         nil,
	}
}

// Write writes an encoded form of p to the underlying io.Writer.
// Note that final value of currByte is not guaranteed to be written and hence
// calls to Write() should be followed with a call to Flush().
func (hw *Writer) Write(p []byte) (int, error) {
	if hw.err != nil {
		return 0, hw.err
	}

	var n int
	if !hw.wroteHeader {
		hw.wroteHeader = true
		_, hw.err = hw.w.Write([]byte(hw.h.String()))
		if hw.err != nil {
			return 0, hw.err
		}
	}
	for _, b := range p {
		enc := hw.encoder[b]
		bitArray := make([]int, enc.codeLen)
		for i := 0; i < enc.codeLen; i++ {
			// Repeatedly get least significant bit
			bitArray[enc.codeLen-1-i] = enc.code & -enc.code
			enc.code >>= 1
		}
		for _, bit := range bitArray {
			hw.currByte = (hw.currByte << 1) | byte(bit)
			hw.bitsWritten++
			if hw.bitsWritten == 8 {
				var newWritten int
				newWritten, hw.err = hw.w.Write([]byte{hw.currByte})
				if hw.err != nil {
					return n + newWritten, hw.err
				}
				hw.currByte = 0
				hw.bitsWritten = 0
				n += newWritten
			}
		}
	}
	return n, hw.err
}

// Flush fills the rest of currByte with 0 bits before writing it to the
// underlying io.Writter.
func (hw *Writer) Flush() (int, error) {
	if hw.err != nil {
		return 0, hw.err
	}
	var n int
	if hw.bitsWritten != 0 {
		for hw.bitsWritten < 8 {
			hw.currByte = (hw.currByte << 1) | 0
			hw.bitsWritten++
		}
		n, hw.err = hw.w.Write([]byte{hw.currByte})
	}
	return n, hw.err
}
