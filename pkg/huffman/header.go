package huffman

import (
	"container/heap"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
)

//=============================================================================
//= Exported
//=============================================================================

// A Header
type Header struct {
	cb       codebook
	numUnits int
}

// String returns a string "representation" of the header.
// This string is intended to be writen to the top of the output file during
// encoding. The string contains enough information so that during decoding,
// the original Header can be recreated.
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

	return fmt.Sprintf("%s%s%s%s%d%s%s",
		line1, headerDelim, line2, headerDelim, line3, headerDelim, headerDelim)
}

// ExtractEncoder uses the Header's codebook to create a map between units and
// their associated code and code length.
func (h *Header) ExtractEncoder() encoder {
	enc := make(encoder)
	for _, e := range h.cb {
		v := enc[e.unit]
		v.code = e.code
		v.codeLen = e.codeLen
		enc[e.unit] = v
	}
	return enc
}

// ExtractDecoder uses the Header's codebook to create a map between codes and
// their associated unit and code length.
func (h *Header) ExtractDecoder() decoder {
	dec := make(decoder)
	for _, e := range h.cb {
		v := dec[e.code]
		v.unit = e.unit
		v.codeLen = e.codeLen
		dec[e.code] = v
	}
	return dec
}

// NewHeader creates and returns a pointer to a new Header.
// To do this, it needs to be given the number of units (bytes) in the file, as
// well as a weight for each such unit (usually its frequency within the file).
func NewHeader(unitWeights map[byte]float64, numUnits int) *Header {
	leaves := newHuffTree(unitWeights)
	cb := make(codebook, len(leaves))
	for i, leaf := range leaves {
		cb[i] = &cbEntry{
			unit:    leaf.val,
			code:    0,
			codeLen: 0,
		}
		for leaf.parent != nil {
			cb[i].codeLen++
			leaf = leaf.parent
		}
	}
	// Codebook entries are sorted first according to code length, then
	// alphabetically within same length.
	sort.Slice(cb, func(i, j int) bool {
		if cb[i].codeLen < cb[j].codeLen {
			return true
		} else if cb[i].codeLen > cb[j].codeLen {
			return false
		}
		return cb[i].unit < cb[j].unit
	})
	code := 0
	i := 0
	for ; i < len(cb)-1; i++ {
		cb[i].code = code
		code = (code + 1) << uint(cb[i+1].codeLen-cb[i].codeLen)
	}
	cb[i].code = code

	header := &Header{
		cb:       cb,
		numUnits: numUnits,
	}
	return header
}

// NewHeaderFromReader reads data from a data stream, recording the values in
// the stream and their frequency. This information is then used to create and
// return a Header.
func NewHeaderFromReader(in io.Reader) (header *Header, err error) {
	unitWeights := make(map[byte]float64)
	var numUnits float64
	currUnit := make([]byte, 1)
	for {
		_, err = in.Read(currUnit)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
		if _, ok := unitWeights[currUnit[0]]; !ok {
			unitWeights[currUnit[0]] = 0
		}
		unitWeights[currUnit[0]]++
		numUnits++
	}
	for i, _ := range unitWeights {
		unitWeights[i] = unitWeights[i] / numUnits
	}

	return NewHeader(unitWeights, int(numUnits)), nil
}

// DeriveHeader recreates the Header described by headerDesc.
// In nearly all cases, headerDesc should be the first n lines of an encoded
// file, the last ending with 'headerDelim'.
func DeriveHeader(headerDesc string) (header *Header, err error) {
	lines := strings.Split(headerDesc, headerDelim)
	line0 := strings.Split(lines[0], " ")
	totals := make([]int, len(line0))
	for i, t := range line0 {
		if totals[i], err = strconv.Atoi(t); err != nil {
			return nil, err
		}
	}
	units := []byte(lines[1])
	var numUnits int
	if numUnits, err = strconv.Atoi(lines[2]); err != nil {
		return nil, err
	}

	code := 0
	i := 0
	cb := make(codebook, len(units))
	for ti := 0; ti < len(totals); ti++ {
		for totals[ti] > 0 {
			cb[i] = &cbEntry{units[i], code, ti + 1}
			totals[ti]--
			tiNext := ti
			for tiNext < len(totals) && totals[tiNext] == 0 {
				tiNext++
			}
			code = (code + 1) << uint(tiNext-ti)
			i++
		}
	}

	header = &Header{
		cb:       cb,
		numUnits: numUnits,
	}
	return header, nil
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

type encoder map[byte]struct {
	code    int
	codeLen int
}

type decoder map[int]struct {
	unit    byte
	codeLen int
}

//-----------------------------------------------------------------------------
//- Huffman Tree
//-----------------------------------------------------------------------------
type huffTree struct {
	val    byte
	weight float64
	left   *huffTree
	right  *huffTree
	parent *huffTree
}
type huffTreePQ []*huffTree

func (htpq huffTreePQ) Len() int { return len(htpq) }

func (htpq huffTreePQ) Less(i, j int) bool {
	return htpq[i].weight < htpq[j].weight
}

func (htpq huffTreePQ) Swap(i, j int) {
	htpq[i], htpq[j] = htpq[j], htpq[i]
}

func (htpq *huffTreePQ) Push(x interface{}) {
	item := x.(*huffTree)
	*htpq = append(*htpq, item)
}

func (htpq *huffTreePQ) Pop() interface{} {
	old := *htpq
	n := len(old)
	item := old[n-1]
	*htpq = old[0 : n-1]
	return item
}

// newHuffTree creates a new Huffman tree.
// The function returns the leaves of the tree, rather than its root.
// This was done because most callers of this function want access to the
// tree's leaves.
func newHuffTree(weightMap map[byte]float64) (leaves []*huffTree) {
	treePQ := make(huffTreePQ, len(weightMap))
	leaves = make([]*huffTree, len(weightMap))
	i := 0
	for b, w := range weightMap {
		leaf := &huffTree{b, w, nil, nil, nil}
		treePQ[i] = leaf
		leaves[i] = leaf
		i++
	}
	heap.Init(&treePQ)

	for treePQ.Len() > 1 {
		min1 := heap.Pop(&treePQ).(*huffTree)
		min2 := heap.Pop(&treePQ).(*huffTree)
		comb := &huffTree{0, min1.weight + min2.weight, min1, min2, nil}
		min1.parent = comb
		min2.parent = comb
		heap.Push(&treePQ, comb)
	}

	return
}
