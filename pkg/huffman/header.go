package huffman

//=============================================================================
//= Exported
//=============================================================================

type Header struct {
	cb       Codebook
	numUnits int
}

// NewHeader creates and returns a pointer to a new Header.
// To do this, it needs to be given the number of units (bytes) in the file, as
// well as a weight for each such unit (usually its frequency within the file).
func NewHeader(unitWeights map[byte]float64, numUnits int) *Header {

}

// DeriveHeader recreates the Header described by headerDesc.
// In nearly all cases, headerDesc should be the first n lines of an encoded
// file, the last ending with 'headerDelim'.
func DeriveHeader(headerDesc string) *Header {

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
type Codebook []*cbEntry
