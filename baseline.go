package goda

import "math"

func calculateBaseline(node *Node) float32 {
	if node.HasBaselineFunc() {
		baseline := node.Baseline(
			node.GetLayout().MeasuredDimension(DimensionWidth),
			node.GetLayout().MeasuredDimension(DimensionHeight))
		if math.IsNaN(float64(baseline)) {
			panic("Expect custom baseline function to not return NaN")
		}
		return baseline
	}

	var baselineChild *Node
	iter := NewLayoutableIterator(node)
	for iter.Next() {
		child := iter.Current()
		if child.GetLineIndex() > 0 {
			break
		}
		if child.Style().PositionType() == PositionTypeAbsolute {
			continue
		}
		if resolveChildAlignment(node, child) == AlignBaseline || child.IsReferenceBaseline() {
			baselineChild = child
			break
		}
		if baselineChild == nil {
			baselineChild = child
		}
	}

	if baselineChild == nil {
		return node.GetLayout().MeasuredDimension(DimensionHeight)
	}

	return calculateBaseline(baselineChild) + baselineChild.GetLayout().Position(PhysicalEdgeTop)
}

func isBaselineLayout(node *Node) bool {
	if isColumn(node.Style().FlexDirection()) {
		return false
	}
	if node.Style().AlignItems() == AlignBaseline {
		return true
	}
	iter := NewLayoutableIterator(node)
	for iter.Next() {
		child := iter.Current()
		if child.Style().PositionType() != PositionTypeAbsolute &&
			child.Style().AlignSelf() == AlignBaseline {
			return true
		}
	}
	return false
}
