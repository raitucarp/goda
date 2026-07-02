package goda

func getPositionOfOppositeEdge(position float32, axis FlexDirection, containingNode, node *Node) float32 {
	return containingNode.GetLayout().MeasuredDimension(dimension(axis)) -
		node.GetLayout().MeasuredDimension(dimension(axis)) - position
}

func setChildTrailingPosition(node, child *Node, axis FlexDirection) {
	child.SetLayoutPosition(
		getPositionOfOppositeEdge(child.GetLayout().Position(flexStartEdge(axis)), axis, node, child),
		flexEndEdge(axis))
}

func needsTrailingPosition(axis FlexDirection) bool {
	return axis == FlexDirectionRowReverse || axis == FlexDirectionColumnReverse
}
