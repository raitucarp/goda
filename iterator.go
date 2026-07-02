package goda

import "math"

// LayoutableIterator iterates over children of a node that participate in layout,
// transparently flattening DisplayContents children.
type LayoutableIterator struct {
	node      *Node
	childIdx  int
	backtrack []backtrackEntry
}

type backtrackEntry struct {
	node     *Node
	childIdx int
}

func NewLayoutableIterator(n *Node) *LayoutableIterator {
	return &LayoutableIterator{node: n, childIdx: -1}
}

func (it *LayoutableIterator) enterContents(contentsNode *Node) {
	it.backtrack = append(it.backtrack, backtrackEntry{it.node, it.childIdx})
	it.node = contentsNode
	it.childIdx = -1
}

func (it *LayoutableIterator) Reset(n *Node) {
	it.node = n
	it.childIdx = -1
	it.backtrack = nil
}

func (it *LayoutableIterator) Current() *Node {
	if it.node == nil || it.childIdx < 0 || it.childIdx >= len(it.node.children) {
		return nil
	}
	return it.node.children[it.childIdx]
}

func (it *LayoutableIterator) Next() bool {
	if it.node == nil {
		return false
	}
	it.childIdx++
	for {
		if it.childIdx >= len(it.node.children) {
			if len(it.backtrack) == 0 {
				it.node = nil
				return false
			}
			entry := it.backtrack[len(it.backtrack)-1]
			it.backtrack = it.backtrack[:len(it.backtrack)-1]
			it.node = entry.node
			it.childIdx = entry.childIdx
			it.childIdx++
			if it.childIdx >= len(it.node.children) {
				continue
			}
		}
		child := it.node.children[it.childIdx]
		if child.style.Display() == DisplayContents && len(child.children) > 0 {
			it.enterContents(child)
			it.childIdx++
			continue
		}
		return true
	}
}

func (n *Node) LayoutableSlice() []*Node {
	var result []*Node
	iter := NewLayoutableIterator(n)
	for iter.Next() {
		result = append(result, iter.Current())
	}
	return result
}

func (n *Node) HasLayoutableChildren() bool {
	return n.GetLayoutChildCount() > 0
}

func init() {
	_ = math.NaN()
}
