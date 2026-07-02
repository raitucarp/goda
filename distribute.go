package goda

import "math"

func distributeFreeSpaceSecondPass(
	flexLine *FlexLine, node *Node, mainAxis, crossAxis FlexDirection,
	direction Direction, ownerWidth, mainAxisOwnerSize, availableInnerMainDim,
	availableInnerCrossDim, availableInnerWidth, availableInnerHeight float32,
	mainAxisOverflows bool, sizingModeCrossDim SizingMode, performLayout bool,
	layoutMarkerData *LayoutData, depth int, generationCount uint32) float32 {

	deltaFreeSpace := float32(0.0)
	isMainAxisRow := isRow(mainAxis)
	isNodeFlexWrap := node.Style().FlexWrap() != WrapNoWrap

	for _, child := range flexLine.ItemsInFlow {
		childFlexBasis := boundAxisWithinMinAndMax(child, direction, mainAxis,
			child.layout.computedFlexBasis, mainAxisOwnerSize, ownerWidth).Unwrap()
		updatedMainSize := childFlexBasis

		if isDefined(flexLine.Layout.RemainingFreeSpace) && flexLine.Layout.RemainingFreeSpace < 0 {
			flexShrinkScaled := -child.ResolveFlexShrink() * childFlexBasis
			if flexShrinkScaled != 0 {
				childSize := float32(math.NaN())
				if isDefined(flexLine.Layout.TotalFlexShrinkScaledFactors) && flexLine.Layout.TotalFlexShrinkScaledFactors == 0 {
					childSize = childFlexBasis + flexShrinkScaled
				} else {
					childSize = childFlexBasis + (flexLine.Layout.RemainingFreeSpace/flexLine.Layout.TotalFlexShrinkScaledFactors)*flexShrinkScaled
				}
				updatedMainSize = boundAxisWithAutoMin(child, mainAxis, direction, childSize, availableInnerMainDim, availableInnerWidth)
			}
		} else if isDefined(flexLine.Layout.RemainingFreeSpace) && flexLine.Layout.RemainingFreeSpace > 0 {
			flexGrowFactor := child.ResolveFlexGrow()
			if !math.IsNaN(float64(flexGrowFactor)) && flexGrowFactor != 0 {
				updatedMainSize = boundAxisWithAutoMin(child, mainAxis, direction,
					childFlexBasis+flexLine.Layout.RemainingFreeSpace/flexLine.Layout.TotalFlexGrowFactors*flexGrowFactor,
					availableInnerMainDim, availableInnerWidth)
			}
		}

		deltaFreeSpace += updatedMainSize - childFlexBasis

		marginMain := child.Style().ComputeMarginForAxis(mainAxis, availableInnerWidth)
		marginCross := child.Style().ComputeMarginForAxis(crossAxis, availableInnerWidth)

		childMainSize := updatedMainSize + marginMain
		childCrossSize := float32(math.NaN())
		childMainSizingMode := SizingModeStretchFit
		childCrossSizingMode := SizingModeStretchFit

		childStyle := child.Style()
		if childStyle.AspectRatio().IsDefined() {
			if isMainAxisRow {
				childCrossSize = (childMainSize-marginMain)/childStyle.AspectRatio().Unwrap() + marginCross
			} else {
				childCrossSize = (childMainSize-marginMain)*childStyle.AspectRatio().Unwrap() + marginCross
			}
			childCrossSizingMode = SizingModeStretchFit
		} else if !math.IsNaN(float64(availableInnerCrossDim)) &&
			!child.HasDefiniteLength(dimension(crossAxis), availableInnerCrossDim) &&
			sizingModeCrossDim == SizingModeStretchFit &&
			!(isNodeFlexWrap && mainAxisOverflows) &&
			resolveChildAlignment(node, child) == AlignStretch &&
			!child.Style().FlexStartMarginIsAuto(crossAxis, direction) &&
			!child.Style().FlexEndMarginIsAuto(crossAxis, direction) {
			childCrossSize = availableInnerCrossDim
			childCrossSizingMode = SizingModeStretchFit
		} else if !child.HasDefiniteLength(dimension(crossAxis), availableInnerCrossDim) {
			childCrossSize = availableInnerCrossDim
			if isUndefined(childCrossSize) {
				childCrossSizingMode = SizingModeMaxContent
			} else {
				childCrossSizingMode = SizingModeFitContent
			}
		} else {
			childCrossSize = child.GetResolvedDimension(direction, dimension(crossAxis), availableInnerCrossDim, availableInnerWidth).Unwrap() + marginCross
			isLoosePct := child.GetProcessedDimension(dimension(crossAxis)).IsPercent() && sizingModeCrossDim != SizingModeStretchFit
			if isUndefined(childCrossSize) || isLoosePct {
				childCrossSizingMode = SizingModeMaxContent
			} else {
				childCrossSizingMode = SizingModeStretchFit
			}
		}

		constrainMaxSizeForMode(child, direction, mainAxis, availableInnerMainDim, availableInnerWidth, &childMainSizingMode, &childMainSize)
		constrainMaxSizeForMode(child, direction, crossAxis, availableInnerCrossDim, availableInnerWidth, &childCrossSizingMode, &childCrossSize)

		requiresStretch := !child.HasDefiniteLength(dimension(crossAxis), availableInnerCrossDim) &&
			resolveChildAlignment(node, child) == AlignStretch &&
			!child.Style().FlexStartMarginIsAuto(crossAxis, direction) &&
			!child.Style().FlexEndMarginIsAuto(crossAxis, direction)

		cW := childMainSize
		cH := childCrossSize
		cWSizing := childMainSizingMode
		cHSizing := childCrossSizingMode
		if !isMainAxisRow {
			cW, cH = childCrossSize, childMainSize
			cWSizing, cHSizing = childCrossSizingMode, childMainSizingMode
		}
		if isUndefined(cW) {
			cWSizing = SizingModeMaxContent
		}
		if isUndefined(cH) {
			cHSizing = SizingModeMaxContent
		}

		isLayoutPass := performLayout && !requiresStretch
		reason := LayoutPassFlexMeasure
		if isLayoutPass {
			reason = LayoutPassFlexLayout
		}

		CalculateLayoutInternal(child, cW, cH, node.layout.Direction(), cWSizing, cHSizing,
			availableInnerWidth, availableInnerHeight, isLayoutPass, reason,
			layoutMarkerData, depth, generationCount)

		node.SetLayoutHadOverflow(node.layout.HadOverflow() || child.layout.HadOverflow())
	}
	return deltaFreeSpace
}

func distributeFreeSpaceFirstPass(
	flexLine *FlexLine, direction Direction, mainAxis FlexDirection,
	ownerWidth, mainAxisOwnerSize, availableInnerMainDim, availableInnerWidth float32) {

	deltaFreeSpace := float32(0.0)

	for _, child := range flexLine.ItemsInFlow {
		childFlexBasis := boundAxisWithinMinAndMax(child, direction, mainAxis,
			child.layout.computedFlexBasis, mainAxisOwnerSize, ownerWidth).Unwrap()

		if flexLine.Layout.RemainingFreeSpace < 0 {
			flexShrinkScaled := -child.ResolveFlexShrink() * childFlexBasis
			if isDefined(flexShrinkScaled) && flexShrinkScaled != 0 {
				baseMainSize := childFlexBasis + flexLine.Layout.RemainingFreeSpace/flexLine.Layout.TotalFlexShrinkScaledFactors*flexShrinkScaled
				boundMainSize := boundAxisWithAutoMin(child, mainAxis, direction, baseMainSize, availableInnerMainDim, availableInnerWidth)
				if isDefined(baseMainSize) && isDefined(boundMainSize) && baseMainSize != boundMainSize {
					deltaFreeSpace += boundMainSize - childFlexBasis
					flexLine.Layout.TotalFlexShrinkScaledFactors -= -child.ResolveFlexShrink() * child.layout.computedFlexBasis.Unwrap()
				}
			}
		} else if isDefined(flexLine.Layout.RemainingFreeSpace) && flexLine.Layout.RemainingFreeSpace > 0 {
			flexGrowFactor := child.ResolveFlexGrow()
			if isDefined(flexGrowFactor) && flexGrowFactor != 0 {
				baseMainSize := childFlexBasis + flexLine.Layout.RemainingFreeSpace/flexLine.Layout.TotalFlexGrowFactors*flexGrowFactor
				boundMainSize := boundAxis(child, mainAxis, direction, baseMainSize, availableInnerMainDim, availableInnerWidth)
				if isDefined(baseMainSize) && isDefined(boundMainSize) && baseMainSize != boundMainSize {
					deltaFreeSpace += boundMainSize - childFlexBasis
					flexLine.Layout.TotalFlexGrowFactors -= flexGrowFactor
				}
			}
		}
	}
	flexLine.Layout.RemainingFreeSpace -= deltaFreeSpace
}

func resolveFlexibleLength(node *Node, flexLine *FlexLine, mainAxis, crossAxis FlexDirection,
	direction Direction, ownerWidth, mainAxisOwnerSize, availableInnerMainDim,
	availableInnerCrossDim, availableInnerWidth, availableInnerHeight float32,
	mainAxisOverflows bool, sizingModeCrossDim SizingMode, performLayout bool,
	layoutMarkerData *LayoutData, depth int, generationCount uint32) {

	originalFreeSpace := flexLine.Layout.RemainingFreeSpace

	if !node.HasErrata(ErrataMinSizeUndefinedInsteadOfAuto) {
		for _, child := range flexLine.ItemsInFlow {
			child.layout.computedAutoMinMainSize = computeAutoMinMainSize(child, mainAxis, direction, mainAxisOwnerSize, availableInnerWidth, availableInnerHeight)
		}
	} else {
		for _, child := range flexLine.ItemsInFlow {
			child.layout.computedAutoMinMainSize = FloatOptional{}
		}
	}

	distributeFreeSpaceFirstPass(flexLine, direction, mainAxis, ownerWidth, mainAxisOwnerSize, availableInnerMainDim, availableInnerWidth)

	distributedFreeSpace := distributeFreeSpaceSecondPass(flexLine, node, mainAxis, crossAxis,
		direction, ownerWidth, mainAxisOwnerSize, availableInnerMainDim, availableInnerCrossDim,
		availableInnerWidth, availableInnerHeight, mainAxisOverflows, sizingModeCrossDim,
		performLayout, layoutMarkerData, depth, generationCount)

	flexLine.Layout.RemainingFreeSpace = originalFreeSpace - distributedFreeSpace
}
