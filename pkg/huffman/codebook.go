package huffman

import (
	"container/heap"
	"sort"
)

//#############################################################################
//# Exported
//#############################################################################
// A Codebook is an ordered collection detailing how symbols should be encoded.
type Codebook []*codebookEntry

// NewCanonicalCB returns a canonical Huffman codebook.
func NewCanonicalCB(weightMap map[byte]float64) Codebook {
	leaves := newHuffTree(weightMap)
	cb := make(Codebook, len(leaves))
	for i, leaf := range leaves {
		//cb[i] = &codebookEntry{leaf.val, "", 0}
		cb[i] = &codebookEntry{leaf.val, 0, 0}
		for leaf.parent != nil {
			cb[i].codeLen++
			leaf = leaf.parent
		}
	}
	// Entries are sorted first accoding to code length, then
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
		//cb[i].code = fmt.Sprintf("%.*b", cb[i].codeLen, code)
		//code = (code + 1) << uint(cb[i+1].codeLen-cb[i].codeLen)
	}
	cb[i].code = code
	//cb[i].code = fmt.Sprintf("%.*b", cb[i].codeLen, code)
	return cb
}

//#############################################################################
//# Unexported
//#############################################################################
//-----------------------------------------------------------------------------
//- Codebook
//-----------------------------------------------------------------------------
type codebookEntry struct {
	unit byte
	code int
	//code    string
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
