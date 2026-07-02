package goda

import "math"

func setFlexStartLayoutPosition(parent, child *Node, direction Direction, axis FlexDirection, containingBlockWidth float32) {
	position := child.Style().ComputeFlexStartMargin(axis, direction, containingBlockWidth) +
		parent.GetLayout().Border(flexStartEdge(axis))

	if !child.HasErrata(ErrataAbsolutePositionWithoutInsetsExcludesPadding) && parent.Style().Display() != DisplayGrid {
		position += parent.GetLayout().Padding(flexStartEdge(axis))
	}
	child.SetLayoutPosition(position, flexStartEdge(axis))
}

func setFlexEndLayoutPosition(parent, child *Node, direction Direction, axis FlexDirection, containingBlockWidth float32) {
	flexEndPos := parent.GetLayout().Border(flexEndEdge(axis)) +
		child.Style().ComputeFlexEndMargin(axis, direction, containingBlockWidth)

	if !child.HasErrata(ErrataAbsolutePositionWithoutInsetsExcludesPadding) && parent.Style().Display() != DisplayGrid {
		flexEndPos += parent.GetLayout().Padding(flexEndEdge(axis))
	}
	child.SetLayoutPosition(getPositionOfOppositeEdge(flexEndPos, axis, parent, child), flexStartEdge(axis))
}

func setCenterLayoutPosition(parent, child *Node, direction Direction, axis FlexDirection, containingBlockWidth float32) {
	parentContentBoxSize := parent.GetLayout().MeasuredDimension(dimension(axis)) -
		parent.GetLayout().Border(flexStartEdge(axis)) -
		parent.GetLayout().Border(flexEndEdge(axis))

	if !child.HasErrata(ErrataAbsolutePositionWithoutInsetsExcludesPadding) && parent.Style().Display() != DisplayGrid {
		parentContentBoxSize -= parent.GetLayout().Padding(flexStartEdge(axis))
		parentContentBoxSize -= parent.GetLayout().Padding(flexEndEdge(axis))
	}

	childOuterSize := child.GetLayout().MeasuredDimension(dimension(axis)) +
		child.Style().ComputeMarginForAxis(axis, containingBlockWidth)

	position := (parentContentBoxSize-childOuterSize)/2.0 +
		parent.GetLayout().Border(flexStartEdge(axis)) +
		child.Style().ComputeFlexStartMargin(axis, direction, containingBlockWidth)

	if !child.HasErrata(ErrataAbsolutePositionWithoutInsetsExcludesPadding) && parent.Style().Display() != DisplayGrid {
		position += parent.GetLayout().Padding(flexStartEdge(axis))
	}
	child.SetLayoutPosition(position, flexStartEdge(axis))
}

func justifyAbsoluteChild(parent, child *Node, direction Direction, mainAxis FlexDirection, containingBlockWidth float32) {
	justify := parent.Style().JustifyContent()
	if parent.Style().Display() == DisplayGrid {
		justify = resolveChildJustification(parent, child)
	}
	switch justify {
	case JustifyStart, JustifyAuto, JustifyStretch, JustifyFlexStart, JustifySpaceBetween:
		setFlexStartLayoutPosition(parent, child, direction, mainAxis, containingBlockWidth)
	case JustifyEnd, JustifyFlexEnd:
		setFlexEndLayoutPosition(parent, child, direction, mainAxis, containingBlockWidth)
	case JustifyCenter, JustifySpaceAround, JustifySpaceEvenly:
		setCenterLayoutPosition(parent, child, direction, mainAxis, containingBlockWidth)
	}
}

func alignAbsoluteChild(parent, child *Node, direction Direction, crossAxis FlexDirection, containingBlockWidth float32) {
	itemAlign := resolveChildAlignment(parent, child)
	parentWrap := parent.Style().FlexWrap()
	if parentWrap == WrapWrapReverse {
		if itemAlign == AlignFlexEnd {
			itemAlign = AlignFlexStart
		} else if itemAlign != AlignCenter {
			itemAlign = AlignFlexEnd
		}
	}

	switch itemAlign {
	case AlignStart, AlignAuto, AlignFlexStart, AlignBaseline, AlignSpaceAround, AlignSpaceBetween, AlignStretch, AlignSpaceEvenly:
		setFlexStartLayoutPosition(parent, child, direction, crossAxis, containingBlockWidth)
	case AlignEnd, AlignFlexEnd:
		setFlexEndLayoutPosition(parent, child, direction, crossAxis, containingBlockWidth)
	case AlignCenter:
		setCenterLayoutPosition(parent, child, direction, crossAxis, containingBlockWidth)
	}
}

func positionAbsoluteChild(containingNode, parent, child *Node, direction Direction, axis FlexDirection, isMainAxis bool, containingBlockWidth, containingBlockHeight float32) {
	isAxisRow := isRow(axis)
	containingBlockSize := containingBlockWidth
	if !isAxisRow {
		containingBlockSize = containingBlockHeight
	}

	if child.Style().IsInlineStartPositionDefined(axis, direction) && !child.Style().IsInlineStartPositionAuto(axis, direction) {
		posRelativeToInlineStart := child.Style().ComputeInlineStartPosition(axis, direction, containingBlockSize) +
			containingNode.Style().ComputeInlineStartBorder(axis, direction) +
			child.Style().ComputeInlineStartMargin(axis, direction, containingBlockSize)
		posRelativeToFlexStart := posRelativeToInlineStart
		if inlineStartEdge(axis, direction) != flexStartEdge(axis) {
			posRelativeToFlexStart = getPositionOfOppositeEdge(posRelativeToInlineStart, axis, containingNode, child)
		}
		child.SetLayoutPosition(posRelativeToFlexStart, flexStartEdge(axis))
	} else if child.Style().IsInlineEndPositionDefined(axis, direction) && !child.Style().IsInlineEndPositionAuto(axis, direction) {
		posRelativeToInlineStart := containingNode.GetLayout().MeasuredDimension(dimension(axis)) -
			child.GetLayout().MeasuredDimension(dimension(axis)) -
			containingNode.Style().ComputeInlineEndBorder(axis, direction) -
			child.Style().ComputeInlineEndMargin(axis, direction, containingBlockSize) -
			child.Style().ComputeInlineEndPosition(axis, direction, containingBlockSize)
		posRelativeToFlexStart := posRelativeToInlineStart
		if inlineStartEdge(axis, direction) != flexStartEdge(axis) {
			posRelativeToFlexStart = getPositionOfOppositeEdge(posRelativeToInlineStart, axis, containingNode, child)
		}
		child.SetLayoutPosition(posRelativeToFlexStart, flexStartEdge(axis))
	} else {
		if isMainAxis {
			justifyAbsoluteChild(parent, child, direction, axis, containingBlockWidth)
		} else {
			alignAbsoluteChild(parent, child, direction, axis, containingBlockWidth)
		}
	}
}

func layoutAbsoluteChild(containingNode, node, child *Node, containingBlockWidth, containingBlockHeight float32,
	widthMode SizingMode, direction Direction, layoutMarkerData *LayoutData, depth int, generationCount uint32) {

	mainAxis := resolveDirection(node.Style().FlexDirection(), direction)
	crossAxis := resolveCrossDirection(mainAxis, direction)
	if node.Style().Display() == DisplayGrid {
		mainAxis = resolveDirection(FlexDirectionRow, direction)
		crossAxis = FlexDirectionColumn
	}

	childWidth := float32(math.NaN())
	childHeight := float32(math.NaN())
	childWidthMode := SizingModeMaxContent
	childHeightMode := SizingModeMaxContent

	marginRow := child.Style().ComputeMarginForAxis(FlexDirectionRow, containingBlockWidth)
	marginColumn := child.Style().ComputeMarginForAxis(FlexDirectionColumn, containingBlockWidth)

	if child.HasDefiniteLength(DimensionWidth, containingBlockWidth) {
		childWidth = child.GetResolvedDimension(direction, DimensionWidth, containingBlockWidth, containingBlockWidth).Unwrap() + marginRow
	} else {
		if child.Style().IsFlexStartPositionDefined(FlexDirectionRow, direction) &&
			child.Style().IsFlexEndPositionDefined(FlexDirectionRow, direction) &&
			!child.Style().IsFlexStartPositionAuto(FlexDirectionRow, direction) &&
			!child.Style().IsFlexEndPositionAuto(FlexDirectionRow, direction) {
			childWidth = containingNode.GetLayout().MeasuredDimension(DimensionWidth) -
				(containingNode.Style().ComputeFlexStartBorder(FlexDirectionRow, direction) +
					containingNode.Style().ComputeFlexEndBorder(FlexDirectionRow, direction)) -
				(child.Style().ComputeFlexStartPosition(FlexDirectionRow, direction, containingBlockWidth) +
					child.Style().ComputeFlexEndPosition(FlexDirectionRow, direction, containingBlockWidth))
			childWidth = boundAxis(child, FlexDirectionRow, direction, childWidth, containingBlockWidth, containingBlockWidth)
		}
	}

	if child.HasDefiniteLength(DimensionHeight, containingBlockHeight) {
		childHeight = child.GetResolvedDimension(direction, DimensionHeight, containingBlockHeight, containingBlockWidth).Unwrap() + marginColumn
	} else {
		if child.Style().IsFlexStartPositionDefined(FlexDirectionColumn, direction) &&
			child.Style().IsFlexEndPositionDefined(FlexDirectionColumn, direction) &&
			!child.Style().IsFlexStartPositionAuto(FlexDirectionColumn, direction) &&
			!child.Style().IsFlexEndPositionAuto(FlexDirectionColumn, direction) {
			childHeight = containingNode.GetLayout().MeasuredDimension(DimensionHeight) -
				(containingNode.Style().ComputeFlexStartBorder(FlexDirectionColumn, direction) +
					containingNode.Style().ComputeFlexEndBorder(FlexDirectionColumn, direction)) -
				(child.Style().ComputeFlexStartPosition(FlexDirectionColumn, direction, containingBlockHeight) +
					child.Style().ComputeFlexEndPosition(FlexDirectionColumn, direction, containingBlockHeight))
			childHeight = boundAxis(child, FlexDirectionColumn, direction, childHeight, containingBlockHeight, containingBlockWidth)
		}
	}

	childStyle := child.Style()
	if isUndefined(childWidth) != isUndefined(childHeight) {
		if childStyle.AspectRatio().IsDefined() {
			if isUndefined(childWidth) {
				childWidth = marginRow + (childHeight-marginColumn)*childStyle.AspectRatio().Unwrap()
			} else {
				childHeight = marginColumn + (childWidth-marginRow)/childStyle.AspectRatio().Unwrap()
			}
		}
	}

	if isUndefined(childWidth) || isUndefined(childHeight) {
		if isUndefined(childWidth) {
			childWidthMode = SizingModeMaxContent
		} else {
			childWidthMode = SizingModeStretchFit
		}
		if isUndefined(childHeight) {
			childHeightMode = SizingModeMaxContent
		} else {
			childHeightMode = SizingModeStretchFit
		}

		if !isRow(mainAxis) && isUndefined(childWidth) && widthMode != SizingModeMaxContent &&
			isDefined(containingBlockWidth) && containingBlockWidth > 0 {
			childWidth = containingBlockWidth
			childWidthMode = SizingModeFitContent
		}

		CalculateLayoutInternal(child, childWidth, childHeight, direction, childWidthMode, childHeightMode,
			containingBlockWidth, containingBlockHeight, false, LayoutPassAbsMeasureChild,
			layoutMarkerData, depth, generationCount)

		childWidth = child.GetLayout().MeasuredDimension(DimensionWidth) +
			child.Style().ComputeMarginForAxis(FlexDirectionRow, containingBlockWidth)
		childHeight = child.GetLayout().MeasuredDimension(DimensionHeight) +
			child.Style().ComputeMarginForAxis(FlexDirectionColumn, containingBlockWidth)
	}

	CalculateLayoutInternal(child, childWidth, childHeight, direction, SizingModeStretchFit, SizingModeStretchFit,
		containingBlockWidth, containingBlockHeight, true, LayoutPassAbsLayout,
		layoutMarkerData, depth, generationCount)

	positionAbsoluteChild(containingNode, node, child, direction, mainAxis, true, containingBlockWidth, containingBlockHeight)
	positionAbsoluteChild(containingNode, node, child, direction, crossAxis, false, containingBlockWidth, containingBlockHeight)
}

func layoutAbsoluteDescendants(containingNode, currentNode *Node, widthSizingMode SizingMode,
	currentDirection Direction, layoutMarkerData *LayoutData, currentDepth int, generationCount uint32,
	currentNodeLeft, currentNodeTop, containingNodeAvailInnerWidth, containingNodeAvailInnerHeight float32) bool {

	hasNewLayout := false
	iter := NewLayoutableIterator(currentNode)
	for iter.Next() {
		child := iter.Current()
		if child.Style().Display() == DisplayNone {
			continue
		} else if child.Style().PositionType() == PositionTypeAbsolute {
			absoluteErrata := currentNode.HasErrata(ErrataAbsolutePercentAgainstInnerSize)
			containingBlockWidth := containingNode.GetLayout().MeasuredDimension(DimensionWidth) -
				containingNode.Style().ComputeBorderForAxis(FlexDirectionRow)
			containingBlockHeight := containingNode.GetLayout().MeasuredDimension(DimensionHeight) -
				containingNode.Style().ComputeBorderForAxis(FlexDirectionColumn)
			if absoluteErrata {
				containingBlockWidth = containingNodeAvailInnerWidth
				containingBlockHeight = containingNodeAvailInnerHeight
			}

			layoutAbsoluteChild(containingNode, currentNode, child, containingBlockWidth, containingBlockHeight,
				widthSizingMode, currentDirection, layoutMarkerData, currentDepth, generationCount)

			hasNewLayout = hasNewLayout || child.GetHasNewLayout()

			parentMainAxis := resolveDirection(currentNode.Style().FlexDirection(), currentDirection)
			parentCrossAxis := resolveCrossDirection(parentMainAxis, currentDirection)

			if needsTrailingPosition(parentMainAxis) {
				mainInsetsDefined := child.Style().HorizontalInsetsDefined()
				if !isRow(parentMainAxis) {
					mainInsetsDefined = child.Style().VerticalInsetsDefined()
				}
				target := currentNode
				if mainInsetsDefined {
					target = containingNode
				}
				setChildTrailingPosition(target, child, parentMainAxis)
			}
			if needsTrailingPosition(parentCrossAxis) {
				crossInsetsDefined := child.Style().HorizontalInsetsDefined()
				if !isRow(parentCrossAxis) {
					crossInsetsDefined = child.Style().VerticalInsetsDefined()
				}
				target := currentNode
				if crossInsetsDefined {
					target = containingNode
				}
				setChildTrailingPosition(target, child, parentCrossAxis)
			}

			childLeft := child.GetLayout().Position(PhysicalEdgeLeft)
			childTop := child.GetLayout().Position(PhysicalEdgeTop)

			childLeftOffset := childLeft
			childTopOffset := childTop
			if child.Style().HorizontalInsetsDefined() {
				childLeftOffset = childLeft - currentNodeLeft
			}
			if child.Style().VerticalInsetsDefined() {
				childTopOffset = childTop - currentNodeTop
			}

			child.SetLayoutPosition(childLeftOffset, PhysicalEdgeLeft)
			child.SetLayoutPosition(childTopOffset, PhysicalEdgeTop)
		} else if child.Style().PositionType() == PositionTypeStatic && !child.AlwaysFormsContainingBlock() {
			child.CloneChildrenIfNeeded()
			childDirection := child.ResolveDirection(currentDirection)

			childLeftOffset := currentNodeLeft + child.GetLayout().Position(PhysicalEdgeLeft)
			childTopOffset := currentNodeTop + child.GetLayout().Position(PhysicalEdgeTop)

			newLayout := layoutAbsoluteDescendants(containingNode, child, widthSizingMode, childDirection,
				layoutMarkerData, currentDepth+1, generationCount, childLeftOffset, childTopOffset,
				containingNodeAvailInnerWidth, containingNodeAvailInnerHeight)

			hasNewLayout = newLayout || hasNewLayout
			cleanupContentsNodesRecursively(child, hasNewLayout)
			if hasNewLayout {
				child.SetHasNewLayout(true)
			}
		}
	}
	return hasNewLayout
}

func init() {
	_ = math.NaN()
}
