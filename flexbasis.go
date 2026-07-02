package goda

import "math"

func constrainMaxSizeForMode(node *Node, dir Direction, axis FlexDirection, ownerAxisSize, ownerWidth float32, mode *SizingMode, size *float32) {
	maxSize := node.Style().ResolvedMaxDimension(dir, dimension(axis), ownerAxisSize, ownerWidth)
	if !maxSize.IsDefined() {
		return
	}
	mg := node.Style().ComputeMarginForAxis(axis, ownerWidth)
	maxVal := NewFloatOptional(maxSize.Unwrap() + mg)

	switch *mode {
	case SizingModeStretchFit, SizingModeFitContent:
		if maxVal.IsDefined() && *size > maxVal.Unwrap() {
			*size = maxVal.Unwrap()
		}
	case SizingModeMaxContent:
		if maxVal.IsDefined() {
			*mode = SizingModeFitContent
			*size = maxVal.Unwrap()
		}
	}
}

func computeFlexBasisForChild(node, child *Node, width float32, widthMode SizingMode,
	height, ownerWidth, ownerHeight float32, heightMode SizingMode,
	direction Direction, layoutMarkerData *LayoutData, depth int, generationCount uint32) {

	mainAxis := resolveDirection(node.Style().FlexDirection(), direction)
	isMainAxisRow := isRow(mainAxis)
	mainAxisSize := width
	mainAxisOwnerSize := ownerWidth
	if !isMainAxisRow {
		mainAxisSize = height
		mainAxisOwnerSize = ownerHeight
	}

	childWidth := float32(math.NaN())
	childHeight := float32(math.NaN())
	childWidthSizingMode := SizingModeMaxContent
	childHeightSizingMode := SizingModeMaxContent

	resolvedFlexBasis := child.ResolveFlexBasis(direction, mainAxis, mainAxisOwnerSize, ownerWidth)
	isRowStyleDefined := child.HasDefiniteLength(DimensionWidth, ownerWidth)
	isColStyleDefined := child.HasDefiniteLength(DimensionHeight, ownerHeight)

	fixFlexBasisFitContent := node.GetConfig().IsExperimentalFeatureEnabled(ExperimentalFeatureFixFlexBasisFitContent)

	useResolvedFlexBasis := resolvedFlexBasis.IsDefined() && isDefined(mainAxisSize)
	if fixFlexBasisFitContent && resolvedFlexBasis.IsDefined() && resolvedFlexBasis.Unwrap() > 0 {
		useResolvedFlexBasis = true
	}

	if useResolvedFlexBasis {
		if child.layout.computedFlexBasis.IsUndefined() ||
			(child.GetConfig().IsExperimentalFeatureEnabled(ExperimentalFeatureWebFlexBasis) &&
				child.layout.computedFlexBasisGeneration != generationCount) {
			pb := NewFloatOptional(paddingAndBorderForAxis(child, mainAxis, direction, ownerWidth))
			child.SetLayoutComputedFlexBasis(NewFloatOptional(maxOrDefined(resolvedFlexBasis.Unwrap(), pb.Unwrap())))
		}
	} else if isMainAxisRow && isRowStyleDefined {
		pb := NewFloatOptional(paddingAndBorderForAxis(child, FlexDirectionRow, direction, ownerWidth))
		resolved := child.GetResolvedDimension(direction, DimensionWidth, ownerWidth, ownerWidth)
		child.SetLayoutComputedFlexBasis(NewFloatOptional(maxOrDefined(resolved.Unwrap(), pb.Unwrap())))
	} else if !isMainAxisRow && isColStyleDefined {
		pb := NewFloatOptional(paddingAndBorderForAxis(child, FlexDirectionColumn, direction, ownerWidth))
		resolved := child.GetResolvedDimension(direction, DimensionHeight, ownerHeight, ownerWidth)
		child.SetLayoutComputedFlexBasis(NewFloatOptional(maxOrDefined(resolved.Unwrap(), pb.Unwrap())))
	} else {
		marginRow := child.Style().ComputeMarginForAxis(FlexDirectionRow, ownerWidth)
		marginColumn := child.Style().ComputeMarginForAxis(FlexDirectionColumn, ownerWidth)

		if isRowStyleDefined {
			childWidth = child.GetResolvedDimension(direction, DimensionWidth, ownerWidth, ownerWidth).Unwrap() + marginRow
			childWidthSizingMode = SizingModeStretchFit
		}
		if isColStyleDefined {
			childHeight = child.GetResolvedDimension(direction, DimensionHeight, ownerHeight, ownerWidth).Unwrap() + marginColumn
			childHeightSizingMode = SizingModeStretchFit
		}

		if (!isMainAxisRow && node.Style().Overflow() == OverflowScroll) || node.Style().Overflow() != OverflowScroll {
			if isUndefined(childWidth) && isDefined(width) {
				childWidth = width
				childWidthSizingMode = SizingModeFitContent
			}
		}

		applyHeightFitContent := isMainAxisRow || node.Style().Overflow() != OverflowScroll
		if fixFlexBasisFitContent {
			nodeHasScrollAncestor := false
			for owner := node.GetOwner(); owner != nil; owner = owner.GetOwner() {
				if owner.Style().Overflow() == OverflowScroll {
					nodeHasScrollAncestor = true
					break
				}
			}
			applyHeightFitContent = isMainAxisRow ||
				((child.HasMeasureFunc() || !nodeHasScrollAncestor) && node.Style().Overflow() != OverflowScroll)
		}
		if applyHeightFitContent && isUndefined(childHeight) && isDefined(height) {
			childHeight = height
			childHeightSizingMode = SizingModeFitContent
		}

		childStyle := child.Style()
		if childStyle.AspectRatio().IsDefined() {
			if !isMainAxisRow && childWidthSizingMode == SizingModeStretchFit {
				childHeight = marginColumn + (childWidth-marginRow)/childStyle.AspectRatio().Unwrap()
				childHeightSizingMode = SizingModeStretchFit
			} else if isMainAxisRow && childHeightSizingMode == SizingModeStretchFit {
				childWidth = marginRow + (childHeight-marginColumn)*childStyle.AspectRatio().Unwrap()
				childWidthSizingMode = SizingModeStretchFit
			}
		}

		hasExactWidth := isDefined(width) && widthMode == SizingModeStretchFit
		childWidthStretch := resolveChildAlignment(node, child) == AlignStretch && childWidthSizingMode != SizingModeStretchFit
		if !isMainAxisRow && !isRowStyleDefined && hasExactWidth && childWidthStretch {
			childWidth = width
			childWidthSizingMode = SizingModeStretchFit
			if childStyle.AspectRatio().IsDefined() {
				childHeight = (childWidth - marginRow) / childStyle.AspectRatio().Unwrap()
				childHeightSizingMode = SizingModeStretchFit
			}
		}

		hasExactHeight := isDefined(height) && heightMode == SizingModeStretchFit
		childHeightStretch := resolveChildAlignment(node, child) == AlignStretch && childHeightSizingMode != SizingModeStretchFit
		if isMainAxisRow && !isColStyleDefined && hasExactHeight && childHeightStretch {
			childHeight = height
			childHeightSizingMode = SizingModeStretchFit
			if childStyle.AspectRatio().IsDefined() {
				childWidth = (childHeight - marginColumn) * childStyle.AspectRatio().Unwrap()
				childWidthSizingMode = SizingModeStretchFit
			}
		}

		constrainMaxSizeForMode(child, direction, FlexDirectionRow, ownerWidth, ownerWidth, &childWidthSizingMode, &childWidth)
		constrainMaxSizeForMode(child, direction, FlexDirectionColumn, ownerHeight, ownerWidth, &childHeightSizingMode, &childHeight)

		CalculateLayoutInternal(child, childWidth, childHeight, direction, childWidthSizingMode, childHeightSizingMode,
			ownerWidth, ownerHeight, false, LayoutPassMeasureChild, layoutMarkerData, depth, generationCount)

		measured := child.layout.MeasuredDimension(dimension(mainAxis))
		pb := paddingAndBorderForAxis(child, mainAxis, direction, ownerWidth)
		child.SetLayoutComputedFlexBasis(NewFloatOptional(maxOrDefined(measured, pb)))
	}
	child.SetLayoutComputedFlexBasisGeneration(generationCount)
}

func computeFlexBasisForChildren(node *Node, availableInnerWidth, availableInnerHeight, ownerWidth, ownerHeight float32,
	widthMode, heightMode SizingMode, direction Direction, mainAxis FlexDirection,
	performLayout bool, layoutMarkerData *LayoutData, depth int, generationCount uint32) float32 {

	totalOuterFlexBasis := float32(0.0)
	var singleFlexChild *Node

	sizingModeMainDim := widthMode
	if !isRow(mainAxis) {
		sizingModeMainDim = heightMode
	}

	if sizingModeMainDim == SizingModeStretchFit {
		iter := NewLayoutableIterator(node)
		for iter.Next() {
			child := iter.Current()
			if child.IsNodeFlexible() {
				if singleFlexChild != nil ||
					inexactEqualsFloat(child.ResolveFlexGrow(), 0) ||
					inexactEqualsFloat(child.ResolveFlexShrink(), 0) {
					singleFlexChild = nil
					break
				}
				singleFlexChild = child
			}
		}
	}

	iter := NewLayoutableIterator(node)
	for iter.Next() {
		child := iter.Current()
		child.ProcessDimensions()

		if child.Style().Display() == DisplayNone {
			if performLayout {
				zeroOutLayoutRecursively(child)
				child.SetHasNewLayout(true)
				child.SetDirty(false)
			}
			continue
		}

		if performLayout {
			childDir := child.ResolveDirection(direction)
			child.SetPosition(childDir, availableInnerWidth, availableInnerHeight)
		}

		if child.Style().PositionType() == PositionTypeAbsolute {
			continue
		}

		if child == singleFlexChild {
			child.SetLayoutComputedFlexBasisGeneration(generationCount)
			child.SetLayoutComputedFlexBasis(NewFloatOptional(0))
		} else {
			computeFlexBasisForChild(node, child, availableInnerWidth, widthMode,
				availableInnerHeight, ownerWidth, ownerHeight, heightMode,
				direction, layoutMarkerData, depth, generationCount)
		}

		totalOuterFlexBasis += child.layout.computedFlexBasis.Unwrap() +
			child.Style().ComputeMarginForAxis(mainAxis, availableInnerWidth)
	}

	return totalOuterFlexBasis
}

func ternaryFloat(cond bool, a, b float32) float32 {
	if cond {
		return a
	}
	return b
}

func ternaryMeasure(cond bool, a, b MeasureMode) MeasureMode {
	if cond {
		return a
	}
	return b
}

// computeMinContentMainSize computes the min-content size per CSS Flexbox 4.5
func computeMinContentMainSize(node *Node, requestedAxis FlexDirection, ownerDirection Direction, ownerWidth, ownerHeight float32) float32 {
	wantRow := isRow(requestedAxis)

	staticMin := node.GetMinContentWidth()
	if !wantRow {
		staticMin = node.GetMinContentHeight()
	}
	if staticMin.IsDefined() {
		return staticMin.Unwrap()
	}

	if node.HasMeasureFunc() {
		var size Size
		if node.HasMinContentMeasureFunc() {
			size = node.MeasureMinContent(
				ternaryFloat(wantRow, 0.0, float32(math.NaN())),
				ternaryMeasure(wantRow, MeasureModeAtMost, MeasureModeUndefined),
				ternaryFloat(wantRow, float32(math.NaN()), 0.0),
				ternaryMeasure(wantRow, MeasureModeUndefined, MeasureModeAtMost))
		} else {
			size = node.Measure(
				ternaryFloat(wantRow, 0.0, float32(math.NaN())),
				ternaryMeasure(wantRow, MeasureModeAtMost, MeasureModeUndefined),
				ternaryFloat(wantRow, float32(math.NaN()), 0.0),
				ternaryMeasure(wantRow, MeasureModeUndefined, MeasureModeAtMost))
		}
		leafDir := node.ResolveDirection(ownerDirection)
		pb := node.Style().ComputeFlexStartPaddingAndBorder(requestedAxis, leafDir, ownerWidth) +
			node.Style().ComputeFlexEndPaddingAndBorder(requestedAxis, leafDir, ownerWidth)
		result := size.Width
		if !wantRow {
			result = size.Height
		}
		return result + pb
	}

	if len(node.GetChildren()) == 0 {
		return 0.0
	}

	direction := node.ResolveDirection(ownerDirection)
	nodeMainAxis := resolveDirection(node.Style().FlexDirection(), direction)
	nodeCrossAxis := resolveCrossDirection(nodeMainAxis, direction)

	mainTotal := float32(0.0)
	crossMax := float32(0.0)

	for _, child := range node.GetChildren() {
		if child.Style().Display() == DisplayNone || child.Style().PositionType() == PositionTypeAbsolute {
			continue
		}
		childMain := computeMinContentMainSize(child, nodeMainAxis, direction, ownerWidth, ownerHeight)
		childMain += child.Style().ComputeMarginForAxis(nodeMainAxis, ownerWidth)

		childCross := computeMinContentMainSize(child, nodeCrossAxis, direction, ownerWidth, ownerHeight)
		childCross += child.Style().ComputeMarginForAxis(nodeCrossAxis, ownerWidth)

		mainTotal += childMain
		crossMax = max(crossMax, childCross)
	}

	mainTotal += node.Style().ComputeFlexStartPaddingAndBorder(nodeMainAxis, direction, ownerWidth) +
		node.Style().ComputeFlexEndPaddingAndBorder(nodeMainAxis, direction, ownerWidth)
	crossMax += node.Style().ComputeFlexStartPaddingAndBorder(nodeCrossAxis, direction, ownerWidth) +
		node.Style().ComputeFlexEndPaddingAndBorder(nodeCrossAxis, direction, ownerWidth)

	nodeMainIsRow := isRow(nodeMainAxis)
	widthMin := mainTotal
	heightMin := mainTotal
	if nodeMainIsRow {
		heightMin = crossMax
	} else {
		widthMin = crossMax
	}

	if wantRow {
		return widthMin
	}
	return heightMin
}

// computeAutoMinMainSize computes CSS Flexbox 4.5 automatic minimum main-axis size
func computeAutoMinMainSize(child *Node, mainAxis FlexDirection, direction Direction, ownerMainAxisSize, ownerWidth, ownerHeight float32) FloatOptional {
	if child.HasErrata(ErrataMinSizeUndefinedInsteadOfAuto) {
		return FloatOptional{}
	}
	if child.Style().Display() == DisplayNone {
		return FloatOptional{}
	}
	if child.Style().MinDimension(dimension(mainAxis)).IsDefined() {
		return FloatOptional{}
	}
	if child.Style().Overflow() != OverflowVisible {
		return NewFloatOptional(0.0)
	}

	mainDim := dimension(mainAxis)
	crossDim := DimensionWidth
	if isRow(mainAxis) {
		crossDim = DimensionHeight
	}
	isMainAxisRow := isRow(mainAxis)

	specifiedMain := child.GetResolvedDimension(direction, mainDim, ownerMainAxisSize, ownerWidth)

	var transferredMain FloatOptional
	aspectRatio := child.Style().AspectRatio()
	if aspectRatio.IsDefined() {
		crossOwner := ownerWidth
		if isMainAxisRow {
			crossOwner = ownerHeight
		}
		crossRes := child.GetResolvedDimension(direction, crossDim, crossOwner, ownerWidth)
		if crossRes.IsDefined() {
			ratio := aspectRatio.Unwrap()
			crossVal := crossRes.Unwrap()
			if isMainAxisRow {
				transferredMain = NewFloatOptional(crossVal * ratio)
			} else {
				transferredMain = NewFloatOptional(crossVal / ratio)
			}
		}
	}

	contentMain := NewFloatOptional(computeMinContentMainSize(child, mainAxis, direction, ownerWidth, ownerHeight))

	floor := contentMain
	if specifiedMain.IsDefined() {
		if floor.IsUndefined() || specifiedMain.Unwrap() < floor.Unwrap() {
			floor = specifiedMain
		}
	} else if transferredMain.IsDefined() {
		if floor.IsUndefined() || transferredMain.Unwrap() < floor.Unwrap() {
			floor = transferredMain
		}
	}

	maxMain := child.Style().ResolvedMaxDimension(direction, mainDim, ownerMainAxisSize, ownerWidth)
	if maxMain.IsDefined() && floor.Unwrap() > maxMain.Unwrap() {
		floor = maxMain
	}

	if floor.IsUndefined() || floor.Unwrap() < 0.0 {
		floor = NewFloatOptional(0.0)
	}
	return floor
}
