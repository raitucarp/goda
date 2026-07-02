package goda

// InsertChildNode inserts a child node at the given index.
// Returns the parent node for chaining.
func (n *Node) InsertChildNode(child *Node, index int) *Node {
	if child.GetOwner() != nil {
		panic("Child already has an owner, it must be removed first.")
	}
	if n.HasMeasureFunc() {
		panic("Cannot add child: Nodes with measure functions cannot have children.")
	}
	n.InsertChild(child, index)
	child.SetOwner(n)
	n.MarkDirtyAndPropagate()
	return n
}

// SwapChildNode swaps a child node at the given index.
// The old child at index is properly detached (owner cleared, layout reset, marked dirty).
// The new child must not already have an owner.
func (n *Node) SwapChildNode(child *Node, index int) *Node {
	if index < 0 || index >= len(n.children) {
		return n
	}
	if child.GetOwner() != nil {
		panic("Child already has an owner, it must be removed first.")
	}
	old := n.children[index]
	old.SetLayout(NewLayoutResults())
	old.SetOwner(nil)
	dirtiedFunc := old.GetDirtiedFunc()
	old.SetDirtiedFunc(nil)
	old.SetDirty(true)
	old.SetDirtiedFunc(dirtiedFunc)
	n.ReplaceChildAt(child, index)
	child.SetOwner(n)
	n.MarkDirtyAndPropagate()
	return n
}

// RemoveChildNode removes a child node. Returns the parent for chaining.
func (n *Node) RemoveChildNode(child *Node) *Node {
	if len(n.children) == 0 {
		return n
	}
	childOwner := child.GetOwner()
	if n.RemoveChild(child) {
		if n == childOwner {
			child.SetLayout(NewLayoutResults())
			child.SetOwner(nil)
			dirtiedFunc := child.GetDirtiedFunc()
			child.SetDirtiedFunc(nil)
			child.SetDirty(true)
			child.SetDirtiedFunc(dirtiedFunc)
		}
		n.MarkDirtyAndPropagate()
	}
	return n
}

// RemoveAllChildren removes all children from the node.
func (n *Node) RemoveAllChildren() *Node {
	childCount := len(n.children)
	if childCount == 0 {
		return n
	}
	firstChild := n.children[0]
	if firstChild.GetOwner() == n {
		for _, child := range n.children {
			child.SetLayout(NewLayoutResults())
			child.SetOwner(nil)
			df := child.GetDirtiedFunc()
			child.SetDirtiedFunc(nil)
			child.SetDirty(true)
			child.SetDirtiedFunc(df)
		}
		n.ClearChildren()
		n.MarkDirtyAndPropagate()
		return n
	}
	n.SetChildren(nil)
	n.MarkDirtyAndPropagate()
	return n
}

// SetChildrenList replaces all children with the given list.
func (n *Node) SetChildrenList(children []*Node) *Node {
	if n == nil {
		return nil
	}
	if len(n.children) > 0 {
		for _, old := range n.children {
			found := false
			for _, newChild := range children {
				if old == newChild {
					found = true
					break
				}
			}
		if !found {
			old.SetLayout(NewLayoutResults())
			old.SetOwner(nil)
			df := old.GetDirtiedFunc()
			old.SetDirtiedFunc(nil)
			old.SetDirty(true)
			old.SetDirtiedFunc(df)
		}
		}
	}
	n.SetChildren(children)
	for _, child := range children {
		child.SetOwner(n)
	}
	n.MarkDirtyAndPropagate()
	return n
}

// MarkDirty marks the node as dirty. Only valid for leaf nodes with measure functions.
func (n *Node) MarkDirty() *Node {
	if !n.HasMeasureFunc() {
		panic("Only leaf nodes with custom measure functions should manually mark themselves as dirty")
	}
	n.MarkDirtyAndPropagate()
	return n
}

// CopyStyleFrom copies all style properties from the source node (deep copy).
func (n *Node) CopyStyleFrom(src *Node) *Node {
	if !n.style.Equals(&src.style) {
		n.style = src.style.Copy()
		n.MarkDirtyAndPropagate()
	}
	return n
}
