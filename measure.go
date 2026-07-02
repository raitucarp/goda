package goda

import "math"

func isFixedSizeFunc(dim float32, mode SizingMode) bool {
	return mode == SizingModeStretchFit || (isDefined(dim) && mode == SizingModeFitContent && dim <= 0.0)
}

func measureNodeWithMeasureFunc(node *Node, direction Direction, availableWidth, availableHeight float32,
	widthMode, heightMode SizingMode, ownerWidth, ownerHeight float32,
	layoutMarkerData *LayoutData, reason int) {

	if !node.HasMeasureFunc() {
		panic("Expected node to have custom measure function")
	}

	avW := availableWidth
	avH := availableHeight
	if widthMode == SizingModeMaxContent {
		avW = float32(math.NaN())
	}
	if heightMode == SizingModeMaxContent {
		avH = float32(math.NaN())
	}

	layout := node.GetLayout()
	pbRow := layout.Padding(PhysicalEdgeLeft) + layout.Padding(PhysicalEdgeRight) +
		layout.Border(PhysicalEdgeLeft) + layout.Border(PhysicalEdgeRight)
	pbCol := layout.Padding(PhysicalEdgeTop) + layout.Padding(PhysicalEdgeBottom) +
		layout.Border(PhysicalEdgeTop) + layout.Border(PhysicalEdgeBottom)

	innerW := avW
	if isDefined(avW) {
		innerW = maxOrDefined(0.0, avW-pbRow)
	}
	innerH := avH
	if isDefined(avH) {
		innerH = maxOrDefined(0.0, avH-pbCol)
	}

	if widthMode == SizingModeStretchFit && heightMode == SizingModeStretchFit {
		node.SetLayoutMeasuredDimension(boundAxis(node, FlexDirectionRow, direction, avW, ownerWidth, ownerWidth), DimensionWidth)
		node.SetLayoutMeasuredDimension(boundAxis(node, FlexDirectionColumn, direction, avH, ownerHeight, ownerWidth), DimensionHeight)
	} else {
		measuredSize := node.Measure(innerW, measureMode(widthMode), innerH, measureMode(heightMode))
		layoutMarkerData.MeasureCallbacks++
		if reason >= 0 && reason < LayoutPassCount {
			layoutMarkerData.MeasureCallbackReasons[reason]++
		}

		finalW := avW
		if widthMode == SizingModeMaxContent || widthMode == SizingModeFitContent {
			finalW = measuredSize.Width + pbRow
		}
		finalH := avH
		if heightMode == SizingModeMaxContent || heightMode == SizingModeFitContent {
			finalH = measuredSize.Height + pbCol
		}

		node.SetLayoutMeasuredDimension(boundAxis(node, FlexDirectionRow, direction, finalW, ownerWidth, ownerWidth), DimensionWidth)
		node.SetLayoutMeasuredDimension(boundAxis(node, FlexDirectionColumn, direction, finalH, ownerHeight, ownerWidth), DimensionHeight)
	}
}

func measureNodeWithoutChildren(node *Node, direction Direction, availableWidth, availableHeight float32,
	widthMode, heightMode SizingMode, ownerWidth, ownerHeight float32) {

	layout := node.GetLayout()
	w := availableWidth
	if widthMode == SizingModeMaxContent || widthMode == SizingModeFitContent {
		w = layout.Padding(PhysicalEdgeLeft) + layout.Padding(PhysicalEdgeRight) +
			layout.Border(PhysicalEdgeLeft) + layout.Border(PhysicalEdgeRight)
	}
	node.SetLayoutMeasuredDimension(boundAxis(node, FlexDirectionRow, direction, w, ownerWidth, ownerWidth), DimensionWidth)

	h := availableHeight
	if heightMode == SizingModeMaxContent || heightMode == SizingModeFitContent {
		h = layout.Padding(PhysicalEdgeTop) + layout.Padding(PhysicalEdgeBottom) +
			layout.Border(PhysicalEdgeTop) + layout.Border(PhysicalEdgeBottom)
	}
	node.SetLayoutMeasuredDimension(boundAxis(node, FlexDirectionColumn, direction, h, ownerHeight, ownerWidth), DimensionHeight)
}

// Debug variables for last measured node.
var (
	lastMeasuredWidth, lastMeasuredHeight                                         float32
	lastMeasuredNode                                                              *Node
	lastMeasuredAvailableW, lastMeasuredAvailableH                                float32
	lastMeasuredWidthMode, lastMeasuredHeightMode                                 SizingMode
	lastMeasuredOwnerW, lastMeasuredOwnerH                                        float32
)

func measureNodeWithFixedSize(node *Node, direction Direction, availableWidth, availableHeight float32,
	widthMode, heightMode SizingMode, ownerWidth, ownerHeight float32) bool {

	if isFixedSizeFunc(availableWidth, widthMode) && isFixedSizeFunc(availableHeight, heightMode) {
		w := availableWidth
		if isUndefined(availableWidth) || (widthMode == SizingModeFitContent && availableWidth < 0) {
			w = 0.0
		}
		node.SetLayoutMeasuredDimension(boundAxis(node, FlexDirectionRow, direction, w, ownerWidth, ownerWidth), DimensionWidth)

		h := availableHeight
		if isUndefined(availableHeight) || (heightMode == SizingModeFitContent && availableHeight < 0) {
			h = 0.0
		}
		node.SetLayoutMeasuredDimension(boundAxis(node, FlexDirectionColumn, direction, h, ownerHeight, ownerWidth), DimensionHeight)
		return true
	}
	return false
}

func zeroOutLayoutRecursively(node *Node) {
	node.SetLayout(NewLayoutResults())
	node.SetLayoutDimension(0, DimensionWidth)
	node.SetLayoutDimension(0, DimensionHeight)
	node.SetHasNewLayout(true)
	node.CloneChildrenIfNeeded()
	for _, child := range node.GetChildren() {
		zeroOutLayoutRecursively(child)
	}
}

func cleanupContentsNodesRecursively(node *Node, didPerformLayout bool) {
	if node.HasContentsChildren() {
		node.CloneContentsChildrenIfNeeded()
		for _, child := range node.GetChildren() {
			if child.Style().Display() == DisplayContents {
				child.SetLayout(NewLayoutResults())
				child.SetLayoutDimension(0, DimensionWidth)
				child.SetLayoutDimension(0, DimensionHeight)
				if didPerformLayout {
					child.SetHasNewLayout(true)
				}
				child.SetDirty(false)
				child.CloneChildrenIfNeeded()
				cleanupContentsNodesRecursively(child, didPerformLayout)
			}
		}
	}
}

func calculateAvailableInnerDimension(node *Node, direction Direction, dim Dimension,
	availableDim, paddingAndBorder, ownerDim, ownerWidth float32) float32 {

	availableInnerDim := availableDim - paddingAndBorder
	if isDefined(availableInnerDim) {
		minDimOpt := node.Style().ResolvedMinDimension(direction, dim, ownerDim, ownerWidth)
		minInner := float32(0.0)
		if minDimOpt.IsDefined() {
			minInner = minDimOpt.Unwrap() - paddingAndBorder
		}

		maxDimOpt := node.Style().ResolvedMaxDimension(direction, dim, ownerDim, ownerWidth)
		maxInner := float32(math.MaxFloat32)
		if maxDimOpt.IsDefined() {
			maxInner = maxDimOpt.Unwrap() - paddingAndBorder
		}

		availableInnerDim = maxOrDefined(minOrDefined(availableInnerDim, maxInner), minInner)
	}
	return availableInnerDim
}

func init() {
	_ = math.NaN()
}
