package huffman

//#############################################################################
//# Exported
//#############################################################################

type Writer struct {
	Header      huffmanHeader
	w           io.Writter
	encTable    *encodingTable
	codebook    *Codebook
	currByte    byte
	bitsWritten int
}

func NewWriter(w io.Writer, cb *Codebook) *Writter {
	encTable := make(map[byte]string)
	for _, entry := range cb {
		encTable[entry.unit] = entry.code
	}
	return &Writer{
		Header:      "",
		w:           w,
		encTable:    encTable,
		codebook:    cb,
		currByte:    0,
		bitsWritten: 0,
	}
}

//#############################################################################
//# Unexported
//#############################################################################
type huffmanHeader struct {
	bitLens  []byte
	alphabet []byte
}

type encodingTable map[byte]string
