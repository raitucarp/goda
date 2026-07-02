package goda

import "math"

func roundValueToPixelGrid(value, pointScaleFactor float64, forceCeil, forceFloor bool) float32 {
	if math.IsNaN(value) || math.IsNaN(pointScaleFactor) || pointScaleFactor == 0 {
		return float32(math.NaN())
	}

	scaledValue := value * pointScaleFactor
	fractional := math.Mod(scaledValue, 1.0)
	if fractional < 0 {
		fractional++
	}
	if inexactEqualsDouble(fractional, 0) {
		scaledValue -= fractional
	} else if inexactEqualsDouble(fractional, 1.0) {
		scaledValue -= fractional + 1.0
	} else if forceCeil {
		scaledValue -= fractional + 1.0
	} else if forceFloor {
		scaledValue -= fractional
	} else {
		add := 0.0
		if !math.IsNaN(fractional) && (fractional > 0.5 || inexactEqualsDouble(fractional, 0.5)) {
			add = 1.0
		}
		scaledValue -= fractional + add
	}
	if math.IsNaN(scaledValue) || math.IsNaN(pointScaleFactor) {
		return float32(math.NaN())
	}
	return float32(scaledValue / pointScaleFactor)
}

func roundLayoutResultsToPixelGrid(node *Node, absoluteLeft, absoluteTop float64) {
	pointScaleFactor := float64(node.GetConfig().GetPointScaleFactor())

	nodeLeft := float64(node.GetLayout().Position(PhysicalEdgeLeft))
	nodeTop := float64(node.GetLayout().Position(PhysicalEdgeTop))
	nodeWidth := float64(node.GetLayout().Dimension(DimensionWidth))
	nodeHeight := float64(node.GetLayout().Dimension(DimensionHeight))

	absNodeLeft := absoluteLeft + nodeLeft
	absNodeTop := absoluteTop + nodeTop
	absNodeRight := absNodeLeft + nodeWidth
	absNodeBottom := absNodeTop + nodeHeight

	if pointScaleFactor != 0.0 {
		textRounding := node.GetNodeType() == NodeTypeText

		node.SetLayoutPosition(
			roundValueToPixelGrid(nodeLeft, pointScaleFactor, false, textRounding),
			PhysicalEdgeLeft)
		node.SetLayoutPosition(
			roundValueToPixelGrid(nodeTop, pointScaleFactor, false, textRounding),
			PhysicalEdgeTop)

		scaledW := nodeWidth * pointScaleFactor
		hasFracW := !inexactEqualsDouble(math.Round(scaledW), scaledW)
		scaledH := nodeHeight * pointScaleFactor
		hasFracH := !inexactEqualsDouble(math.Round(scaledH), scaledH)

		width := roundValueToPixelGrid(absNodeRight, pointScaleFactor, textRounding && hasFracW, textRounding && !hasFracW) -
			roundValueToPixelGrid(absNodeLeft, pointScaleFactor, false, textRounding)
		height := roundValueToPixelGrid(absNodeBottom, pointScaleFactor, textRounding && hasFracH, textRounding && !hasFracH) -
			roundValueToPixelGrid(absNodeTop, pointScaleFactor, false, textRounding)

		node.GetLayout().SetDimension(DimensionWidth, width)
		node.GetLayout().SetDimension(DimensionHeight, height)
	}

	for _, child := range node.GetChildren() {
		roundLayoutResultsToPixelGrid(child, absNodeLeft, absNodeTop)
	}
}
