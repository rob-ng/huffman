package huffman

import (
	"fmt"
	"strconv"
	"strings"
)

//=============================================================================
//= Exported
//=============================================================================

type Header struct {
	cb       codebook
	numUnits int
}

func (h *Header) String() string {
	maxBitLen := 0
	for _, e := range h.cb {
		if e.codeLen > maxBitLen {
			maxBitLen = e.codeLen
		}
	}
	totals := make([]int, maxBitLen)
	units := make([]byte, len(h.cb))
	for i, e := range h.cb {
		totals[e.codeLen-1]++
		units[i] = e.unit
	}
	totalsStr := make([]string, maxBitLen)
	for i, _ := range totals {
		totalsStr[i] = strconv.Itoa(totals[i])
	}
	numUnits := h.numUnits

	line1 := strings.Join(totalsStr, " ")
	line2 := units
	line3 := numUnits

	return fmt.Sprintf("%s\n%s\n%d\n%s", line1, line2, line3, headerDelim)
}

// NewHeader creates and returns a pointer to a new Header.
// To do this, it needs to be given the number of units (bytes) in the file, as
// well as a weight for each such unit (usually its frequency within the file).
func NewHeader(unitWeights map[byte]float64, numUnits int) *Header {
	return nil
}

// DeriveHeader recreates the Header described by headerDesc.
// In nearly all cases, headerDesc should be the first n lines of an encoded
// file, the last ending with 'headerDelim'.
func DeriveHeader(headerDesc string) *Header {
	return nil
}

//=============================================================================
//= Unexported
//=============================================================================

const headerDelim = "\r\n"

type cbEntry struct {
	unit    byte
	code    int
	codeLen int
}
type codebook []*cbEntry
