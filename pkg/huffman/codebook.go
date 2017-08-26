package huffman

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

func (htpq huffTreePQ) Less(i, j int) {
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
