package goda

import "math"

// calculateLayoutImpl is the main flexbox layout routine.
func calculateLayoutImpl(node *Node, availableWidth, availableHeight float32, ownerDirection Direction,
	widthMode, heightMode SizingMode, ownerWidth, ownerHeight float32, performLayout bool,
	reason int, layoutMarkerData *LayoutData, depth int, generationCount uint32) {

	if isUndefined(availableWidth) && widthMode != SizingModeMaxContent {
		if widthMode != SizingModeMaxContent {
			widthMode = SizingModeMaxContent
		}
	}
	if isUndefined(availableHeight) && heightMode != SizingModeMaxContent {
		if heightMode != SizingModeMaxContent {
			heightMode = SizingModeMaxContent
		}
	}

	if performLayout {
		layoutMarkerData.Layouts++
	} else {
		layoutMarkerData.Measures++
	}

	direction := node.ResolveDirection(ownerDirection)
	node.SetLayoutDirection(direction)

	flexRowDir := resolveDirection(FlexDirectionRow, direction)
	flexColDir := resolveDirection(FlexDirectionColumn, direction)

	startEdge := PhysicalEdgeLeft
	endEdge := PhysicalEdgeRight
	if direction == DirectionRTL {
		startEdge = PhysicalEdgeRight
		endEdge = PhysicalEdgeLeft
	}

	marginRowLeading := node.Style().ComputeInlineStartMargin(flexRowDir, direction, ownerWidth)
	node.SetLayoutMargin(marginRowLeading, startEdge)
	marginRowTrailing := node.Style().ComputeInlineEndMargin(flexRowDir, direction, ownerWidth)
	node.SetLayoutMargin(marginRowTrailing, endEdge)
	marginColLeading := node.Style().ComputeInlineStartMargin(flexColDir, direction, ownerWidth)
	node.SetLayoutMargin(marginColLeading, PhysicalEdgeTop)
	marginColTrailing := node.Style().ComputeInlineEndMargin(flexColDir, direction, ownerWidth)
	node.SetLayoutMargin(marginColTrailing, PhysicalEdgeBottom)

	marginAxisRow := marginRowLeading + marginRowTrailing
	marginAxisCol := marginColLeading + marginColTrailing

	node.SetLayoutBorder(node.Style().ComputeInlineStartBorder(flexRowDir, direction), startEdge)
	node.SetLayoutBorder(node.Style().ComputeInlineEndBorder(flexRowDir, direction), endEdge)
	node.SetLayoutBorder(node.Style().ComputeInlineStartBorder(flexColDir, direction), PhysicalEdgeTop)
	node.SetLayoutBorder(node.Style().ComputeInlineEndBorder(flexColDir, direction), PhysicalEdgeBottom)

	node.SetLayoutPadding(node.Style().ComputeInlineStartPadding(flexRowDir, direction, ownerWidth), startEdge)
	node.SetLayoutPadding(node.Style().ComputeInlineEndPadding(flexRowDir, direction, ownerWidth), endEdge)
	node.SetLayoutPadding(node.Style().ComputeInlineStartPadding(flexColDir, direction, ownerWidth), PhysicalEdgeTop)
	node.SetLayoutPadding(node.Style().ComputeInlineEndPadding(flexColDir, direction, ownerWidth), PhysicalEdgeBottom)

	if node.HasMeasureFunc() {
		measureNodeWithMeasureFunc(node, direction, availableWidth-marginAxisRow, availableHeight-marginAxisCol,
			widthMode, heightMode, ownerWidth, ownerHeight, layoutMarkerData, reason)
		cleanupContentsNodesRecursively(node, performLayout)
		return
	}

	childCount := node.GetLayoutChildCount()
	if childCount == 0 {
		measureNodeWithoutChildren(node, direction, availableWidth-marginAxisRow, availableHeight-marginAxisCol,
			widthMode, heightMode, ownerWidth, ownerHeight)
		cleanupContentsNodesRecursively(node, performLayout)
		return
	}

	if !performLayout && measureNodeWithFixedSize(node, direction, availableWidth-marginAxisRow, availableHeight-marginAxisCol,
		widthMode, heightMode, ownerWidth, ownerHeight) {
		cleanupContentsNodesRecursively(node, false)
		return
	}

	node.CloneChildrenIfNeeded()
	node.SetLayoutHadOverflow(false)
	cleanupContentsNodesRecursively(node, performLayout)

	mainAxis := resolveDirection(node.Style().FlexDirection(), direction)
	crossAxis := resolveCrossDirection(mainAxis, direction)
	isMainAxisRow := isRow(mainAxis)
	isNodeFlexWrap := node.Style().FlexWrap() != WrapNoWrap

	mainAxisOwnerSize := ownerWidth
	crossAxisOwnerSize := ownerHeight
	if !isMainAxisRow {
		mainAxisOwnerSize, crossAxisOwnerSize = ownerHeight, ownerWidth
	}

	pbAxisMain := paddingAndBorderForAxis(node, mainAxis, direction, ownerWidth)
	pbAxisCross := paddingAndBorderForAxis(node, crossAxis, direction, ownerWidth)
	leadingPBCross := node.Style().ComputeFlexStartPaddingAndBorder(crossAxis, direction, ownerWidth)

	sizingModeMainDim := widthMode
	sizingModeCrossDim := heightMode
	if !isMainAxisRow {
		sizingModeMainDim, sizingModeCrossDim = heightMode, widthMode
	}

	pbAxisRow := pbAxisMain
	pbAxisCol := pbAxisCross
	if !isMainAxisRow {
		pbAxisRow, pbAxisCol = pbAxisCross, pbAxisMain
	}

	availInnerW := calculateAvailableInnerDimension(node, direction, DimensionWidth, availableWidth-marginAxisRow, pbAxisRow, ownerWidth, ownerWidth)
	availInnerH := calculateAvailableInnerDimension(node, direction, DimensionHeight, availableHeight-marginAxisCol, pbAxisCol, ownerHeight, ownerWidth)

	availInnerMainDim := availInnerW
	availInnerCrossDim := availInnerH
	if !isMainAxisRow {
		availInnerMainDim, availInnerCrossDim = availInnerH, availInnerW
	}

	ownerWidthForChildren := availInnerW
	ownerHeightForChildren := availInnerH

	if node.GetConfig().IsExperimentalFeatureEnabled(ExperimentalFeatureFixFlexBasisFitContent) {
		owner := node.GetOwner()
		isChildOfScroll := owner != nil && owner.Style().Overflow() == OverflowScroll
		if !isChildOfScroll {
			if isUndefined(ownerWidthForChildren) && isDefined(ownerWidth) {
				ownerWidthForChildren = calculateAvailableInnerDimension(node, direction, DimensionWidth, ownerWidth-marginAxisRow, pbAxisRow, ownerWidth, ownerWidth)
			}
			if isUndefined(ownerHeightForChildren) && isDefined(ownerHeight) {
				ownerHeightForChildren = calculateAvailableInnerDimension(node, direction, DimensionHeight, ownerHeight-marginAxisCol, pbAxisCol, ownerHeight, ownerWidth)
			}
		}
	}

	totalMainDim := computeFlexBasisForChildren(node, availInnerW, availInnerH, ownerWidthForChildren, ownerHeightForChildren,
		widthMode, heightMode, direction, mainAxis, performLayout, layoutMarkerData, depth, generationCount)

	if childCount > 1 {
		totalMainDim += node.Style().ComputeGapForAxis(mainAxis, availInnerMainDim) * float32(childCount-1)
	}

	mainAxisOverflows := (sizingModeMainDim != SizingModeMaxContent) && totalMainDim > availInnerMainDim

	if isNodeFlexWrap && mainAxisOverflows && sizingModeMainDim == SizingModeFitContent {
		sizingModeMainDim = SizingModeStretchFit
	}

	startIterator := NewLayoutableIterator(node)
	lineCount := 0
	totalLineCrossDim := float32(0.0)
	crossAxisGap := node.Style().ComputeGapForAxis(crossAxis, availInnerCrossDim)
	maxLineMainDim := float32(0.0)

	for startIterator.Next() || lineCount > 0 {
		flexLine := calculateFlexLine(node, ownerDirection, ownerWidth, mainAxisOwnerSize, availInnerW, availInnerMainDim, startIterator, lineCount)

		if len(flexLine.ItemsInFlow) == 0 {
			break
		}

		canSkipFlex := !performLayout && sizingModeCrossDim == SizingModeStretchFit

		sizeBasedOnContent := false
		if sizingModeMainDim != SizingModeStretchFit {
			minInnerW := node.Style().ResolvedMinDimension(direction, DimensionWidth, ownerWidth, ownerWidth).Unwrap() - pbAxisRow
			maxInnerW := node.Style().ResolvedMaxDimension(direction, DimensionWidth, ownerWidth, ownerWidth).Unwrap() - pbAxisRow
			minInnerH := node.Style().ResolvedMinDimension(direction, DimensionHeight, ownerHeight, ownerWidth).Unwrap() - pbAxisCol
			maxInnerH := node.Style().ResolvedMaxDimension(direction, DimensionHeight, ownerHeight, ownerWidth).Unwrap() - pbAxisCol

			minInnerMain := minInnerW
			maxInnerMain := maxInnerW
			if !isMainAxisRow {
				minInnerMain, maxInnerMain = minInnerH, maxInnerH
			}

			if isDefined(minInnerMain) && flexLine.SizeConsumed < minInnerMain {
				availInnerMainDim = minInnerMain
			} else if isDefined(maxInnerMain) && flexLine.SizeConsumed > maxInnerMain {
				availInnerMainDim = maxInnerMain
			} else {
				useLegacy := node.HasErrata(ErrataStretchFlexBasis)
				if !useLegacy &&
					((isDefined(flexLine.Layout.TotalFlexGrowFactors) && flexLine.Layout.TotalFlexGrowFactors == 0) ||
						(isDefined(node.ResolveFlexGrow()) && node.ResolveFlexGrow() == 0)) {
					availInnerMainDim = flexLine.SizeConsumed
				}
				sizeBasedOnContent = !useLegacy
			}
		}

		if !sizeBasedOnContent && isDefined(availInnerMainDim) {
			flexLine.Layout.RemainingFreeSpace = availInnerMainDim - flexLine.SizeConsumed
		} else if flexLine.SizeConsumed < 0 {
			flexLine.Layout.RemainingFreeSpace = -flexLine.SizeConsumed
		}

		if !canSkipFlex {
			resolveFlexibleLength(node, &flexLine, mainAxis, crossAxis, direction, ownerWidth,
				mainAxisOwnerSize, availInnerMainDim, availInnerCrossDim, availInnerW, availInnerH,
				mainAxisOverflows, sizingModeCrossDim, performLayout, layoutMarkerData, depth, generationCount)
		}

		node.SetLayoutHadOverflow(node.layout.HadOverflow() || (flexLine.Layout.RemainingFreeSpace < 0))

		justifyMainAxis(node, &flexLine, mainAxis, crossAxis, direction, sizingModeMainDim, sizingModeCrossDim,
			mainAxisOwnerSize, ownerWidth, availInnerMainDim, availInnerCrossDim, availInnerW, performLayout)

		containerCrossAxis := availInnerCrossDim
		if sizingModeCrossDim == SizingModeMaxContent || sizingModeCrossDim == SizingModeFitContent {
			containerCrossAxis = boundAxis(node, crossAxis, direction, flexLine.Layout.CrossDim+pbAxisCross, crossAxisOwnerSize, ownerWidth) - pbAxisCross
		}

		if !isNodeFlexWrap && sizingModeCrossDim == SizingModeStretchFit {
			flexLine.Layout.CrossDim = availInnerCrossDim
		}

		if !isNodeFlexWrap {
			flexLine.Layout.CrossDim = boundAxis(node, crossAxis, direction, flexLine.Layout.CrossDim+pbAxisCross, crossAxisOwnerSize, ownerWidth) - pbAxisCross
		}

		if performLayout {
			for _, child := range flexLine.ItemsInFlow {
				leadingCrossDim := leadingPBCross
				alignItem := resolveChildAlignment(node, child)

				if alignItem == AlignStretch &&
					!child.Style().FlexStartMarginIsAuto(crossAxis, direction) &&
					!child.Style().FlexEndMarginIsAuto(crossAxis, direction) {
					if !child.HasDefiniteLength(dimension(crossAxis), availInnerCrossDim) {
						childMainSize := child.layout.MeasuredDimension(dimension(mainAxis))
						childCrossSize := flexLine.Layout.CrossDim
						if child.Style().AspectRatio().IsDefined() {
							childCrossSize = child.Style().ComputeMarginForAxis(crossAxis, availInnerW)
							if isMainAxisRow {
								childCrossSize += childMainSize / child.Style().AspectRatio().Unwrap()
							} else {
								childCrossSize += childMainSize * child.Style().AspectRatio().Unwrap()
							}
						}

						childMainSize += child.Style().ComputeMarginForAxis(mainAxis, availInnerW)
						childMainSizingMode := SizingModeStretchFit
						childCrossSizingMode := SizingModeStretchFit
						constrainMaxSizeForMode(child, direction, mainAxis, availInnerMainDim, availInnerW, &childMainSizingMode, &childMainSize)
						constrainMaxSizeForMode(child, direction, crossAxis, availInnerCrossDim, availInnerW, &childCrossSizingMode, &childCrossSize)

						cW := childMainSize
						cH := childCrossSize
						cWSizing := childMainSizingMode
						cHSizing := childCrossSizingMode
						if !isMainAxisRow {
							cW, cH = childCrossSize, childMainSize
							cWSizing, cHSizing = childCrossSizingMode, childMainSizingMode
						}

						alignContent := node.Style().AlignContent()
						crossAxisNoGrow := alignContent != AlignStretch && isNodeFlexWrap
						if isUndefined(cW) || (!isMainAxisRow && crossAxisNoGrow) {
							cWSizing = SizingModeMaxContent
						}
						if isUndefined(cH) || (isMainAxisRow && crossAxisNoGrow) {
							cHSizing = SizingModeMaxContent
						}

						CalculateLayoutInternal(child, cW, cH, direction, cWSizing, cHSizing,
							availInnerW, availInnerH, true, LayoutPassStretch, layoutMarkerData, depth, generationCount)
					}
				} else {
					remainingCrossDim := containerCrossAxis - child.DimensionWithMargin(crossAxis, availInnerW)
					if child.Style().FlexStartMarginIsAuto(crossAxis, direction) &&
						child.Style().FlexEndMarginIsAuto(crossAxis, direction) {
						leadingCrossDim += maxOrDefined(0.0, remainingCrossDim/2)
					} else if child.Style().FlexStartMarginIsAuto(crossAxis, direction) {
						leadingCrossDim += maxOrDefined(0.0, remainingCrossDim)
					} else if alignItem == AlignCenter {
						leadingCrossDim += remainingCrossDim / 2
					} else if alignItem == AlignFlexEnd || alignItem == AlignEnd {
						leadingCrossDim += remainingCrossDim
					}
				}

				child.SetLayoutPosition(child.layout.Position(flexStartEdge(crossAxis))+totalLineCrossDim+leadingCrossDim, flexStartEdge(crossAxis))
			}
		}

		appliedCrossGap := float32(0.0)
		if lineCount != 0 {
			appliedCrossGap = crossAxisGap
		}
		totalLineCrossDim += flexLine.Layout.CrossDim + appliedCrossGap
		maxLineMainDim = maxOrDefined(maxLineMainDim, flexLine.Layout.MainDim)
		lineCount++
	}

	// Multi-line alignment
	if performLayout && (isNodeFlexWrap || isBaselineLayout(node)) {
		leadPerLine := float32(0.0)
		currentLead := leadingPBCross
		extraSpacePerLine := float32(0.0)

		unclampedCrossDim := availInnerCrossDim + pbAxisCross
		if sizingModeCrossDim != SizingModeStretchFit {
			if node.HasDefiniteLength(dimension(crossAxis), crossAxisOwnerSize) {
				unclampedCrossDim = node.GetResolvedDimension(direction, dimension(crossAxis), crossAxisOwnerSize, ownerWidth).Unwrap()
			} else {
				unclampedCrossDim = totalLineCrossDim + pbAxisCross
			}
		}

		innerCrossDim := boundAxis(node, crossAxis, direction, unclampedCrossDim, crossAxisOwnerSize, ownerWidth) - pbAxisCross
		remainingAlignContentDim := innerCrossDim - totalLineCrossDim

		alignContent := node.Style().AlignContent()
		if remainingAlignContentDim < 0 {
			alignContent = fallbackAlignmentAlign(node.Style().AlignContent())
		}

		switch alignContent {
		case AlignFlexEnd, AlignEnd:
			currentLead += remainingAlignContentDim
		case AlignCenter:
			currentLead += remainingAlignContentDim / 2
		case AlignStretch:
			extraSpacePerLine = remainingAlignContentDim / float32(lineCount)
		case AlignSpaceAround:
			currentLead += remainingAlignContentDim / (2 * float32(lineCount))
			leadPerLine = remainingAlignContentDim / float32(lineCount)
		case AlignSpaceEvenly:
			currentLead += remainingAlignContentDim / float32(lineCount+1)
			leadPerLine = remainingAlignContentDim / float32(lineCount+1)
		case AlignSpaceBetween:
			if lineCount > 1 {
				leadPerLine = remainingAlignContentDim / float32(lineCount-1)
			}
		}

		endIter := NewLayoutableIterator(node)
		for i := 0; i < lineCount; i++ {
			startIter := endIter
			iterator := startIter

			lineHeight := float32(0.0)
			maxAscentForLine := float32(0.0)
			maxDescentForLine := float32(0.0)
			for iterator.Next() {
				child := iterator.Current()
				if child.Style().Display() == DisplayNone {
					continue
				}
				if child.Style().PositionType() != PositionTypeAbsolute {
					if child.GetLineIndex() != i {
						break
					}
					if child.IsLayoutDimensionDefined(crossAxis) {
						lineHeight = maxOrDefined(lineHeight, child.layout.MeasuredDimension(dimension(crossAxis))+
							child.Style().ComputeMarginForAxis(crossAxis, availInnerW))
					}
					if resolveChildAlignment(node, child) == AlignBaseline {
						ascent := calculateBaseline(child) +
							child.Style().ComputeFlexStartMargin(FlexDirectionColumn, direction, availInnerW)
						descent := child.layout.MeasuredDimension(DimensionHeight) +
							child.Style().ComputeMarginForAxis(FlexDirectionColumn, availInnerW) - ascent
						maxAscentForLine = maxOrDefined(maxAscentForLine, ascent)
						maxDescentForLine = maxOrDefined(maxDescentForLine, descent)
						lineHeight = maxOrDefined(lineHeight, maxAscentForLine+maxDescentForLine)
					}
				}
			}
			endIter = iterator

			if i != 0 {
				currentLead += crossAxisGap
			}
			lineHeight += extraSpacePerLine

			rIter := NewLayoutableIterator(node)
			for rIter.Next() {
				child := rIter.Current()
				if child.Style().Display() == DisplayNone {
					continue
				}
				if child.Style().PositionType() != PositionTypeAbsolute && child.GetLineIndex() == i {
					switch resolveChildAlignment(node, child) {
					case AlignFlexStart:
						child.SetLayoutPosition(currentLead+child.Style().ComputeFlexStartPosition(crossAxis, direction, availInnerW), flexStartEdge(crossAxis))
					case AlignFlexEnd:
						child.SetLayoutPosition(currentLead+lineHeight-
							child.Style().ComputeFlexEndMargin(crossAxis, direction, availInnerW)-
							child.layout.MeasuredDimension(dimension(crossAxis)), flexStartEdge(crossAxis))
					case AlignCenter:
						childHeight := child.layout.MeasuredDimension(dimension(crossAxis))
						child.SetLayoutPosition(currentLead+(lineHeight-childHeight)/2, flexStartEdge(crossAxis))
					case AlignStretch:
						child.SetLayoutPosition(currentLead+child.Style().ComputeFlexStartMargin(crossAxis, direction, availInnerW), flexStartEdge(crossAxis))
						if !child.HasDefiniteLength(dimension(crossAxis), availInnerCrossDim) {
							cW := child.layout.MeasuredDimension(DimensionWidth) + child.Style().ComputeMarginForAxis(mainAxis, availInnerW)
							cH := child.layout.MeasuredDimension(DimensionHeight) + child.Style().ComputeMarginForAxis(crossAxis, availInnerW)
							if isMainAxisRow {
								cW = child.layout.MeasuredDimension(DimensionWidth) + child.Style().ComputeMarginForAxis(mainAxis, availInnerW)
							} else {
								cW = leadPerLine + lineHeight
							}
							if !isMainAxisRow {
								cH = child.layout.MeasuredDimension(DimensionHeight) + child.Style().ComputeMarginForAxis(crossAxis, availInnerW)
							} else {
								cH = leadPerLine + lineHeight
							}
							if !(inexactEqualsFloat(cW, child.layout.MeasuredDimension(DimensionWidth)) &&
								inexactEqualsFloat(cH, child.layout.MeasuredDimension(DimensionHeight))) {
								CalculateLayoutInternal(child, cW, cH, direction, SizingModeStretchFit, SizingModeStretchFit,
									availInnerW, availInnerH, true, LayoutPassMultilineStretch, layoutMarkerData, depth, generationCount)
							}
						}
					case AlignBaseline:
						child.SetLayoutPosition(currentLead+maxAscentForLine-calculateBaseline(child)+
							child.Style().ComputeFlexStartPosition(FlexDirectionColumn, direction, availInnerCrossDim), PhysicalEdgeTop)
					}
				}
			}

			currentLead += leadPerLine + lineHeight
		}
	}

	// Final dimensions
	node.SetLayoutMeasuredDimension(boundAxis(node, FlexDirectionRow, direction, availableWidth-marginAxisRow, ownerWidth, ownerWidth), DimensionWidth)
	node.SetLayoutMeasuredDimension(boundAxis(node, FlexDirectionColumn, direction, availableHeight-marginAxisCol, ownerHeight, ownerWidth), DimensionHeight)

	if sizingModeMainDim == SizingModeMaxContent ||
		(node.Style().Overflow() != OverflowScroll && sizingModeMainDim == SizingModeFitContent) {
		node.SetLayoutMeasuredDimension(boundAxis(node, mainAxis, direction, maxLineMainDim, mainAxisOwnerSize, ownerWidth), dimension(mainAxis))
	} else if sizingModeMainDim == SizingModeFitContent && node.Style().Overflow() == OverflowScroll {
		node.SetLayoutMeasuredDimension(maxOrDefined(minOrDefined(availInnerMainDim+pbAxisMain,
			boundAxisWithinMinAndMax(node, direction, mainAxis, NewFloatOptional(maxLineMainDim), mainAxisOwnerSize, ownerWidth).Unwrap()),
			pbAxisMain), dimension(mainAxis))
	}

	if sizingModeCrossDim == SizingModeMaxContent ||
		(node.Style().Overflow() != OverflowScroll && sizingModeCrossDim == SizingModeFitContent) {
		node.SetLayoutMeasuredDimension(boundAxis(node, crossAxis, direction, totalLineCrossDim+pbAxisCross, crossAxisOwnerSize, ownerWidth), dimension(crossAxis))
	} else if sizingModeCrossDim == SizingModeFitContent && node.Style().Overflow() == OverflowScroll {
		node.SetLayoutMeasuredDimension(maxOrDefined(minOrDefined(availInnerCrossDim+pbAxisCross,
			boundAxisWithinMinAndMax(node, direction, crossAxis, NewFloatOptional(totalLineCrossDim+pbAxisCross), crossAxisOwnerSize, ownerWidth).Unwrap()),
			pbAxisCross), dimension(crossAxis))
	}

	// Wrap-reverse
	if performLayout && node.Style().FlexWrap() == WrapWrapReverse {
		iter := NewLayoutableIterator(node)
		for iter.Next() {
			child := iter.Current()
			if child.Style().PositionType() != PositionTypeAbsolute {
				child.SetLayoutPosition(node.layout.MeasuredDimension(dimension(crossAxis))-
					child.layout.Position(flexStartEdge(crossAxis))-
					child.layout.MeasuredDimension(dimension(crossAxis)), flexStartEdge(crossAxis))
			}
		}
	}

	if performLayout {
		needsMainTrailing := needsTrailingPosition(mainAxis)
		needsCrossTrailing := needsTrailingPosition(crossAxis)
		if needsMainTrailing || needsCrossTrailing {
			iter := NewLayoutableIterator(node)
			for iter.Next() {
				child := iter.Current()
				if child.Style().Display() == DisplayNone || child.Style().PositionType() == PositionTypeAbsolute {
					continue
				}
				if needsMainTrailing {
					setChildTrailingPosition(node, child, mainAxis)
				}
				if needsCrossTrailing {
					setChildTrailingPosition(node, child, crossAxis)
				}
			}
		}

		if node.Style().PositionType() != PositionTypeStatic || node.AlwaysFormsContainingBlock() || depth == 1 {
			sizingModeForAbs := sizingModeMainDim
			if !isMainAxisRow {
				sizingModeForAbs = sizingModeCrossDim
			}
			layoutAbsoluteDescendants(node, node, sizingModeForAbs, direction, layoutMarkerData, depth, generationCount, 0, 0, availInnerW, availInnerH)
		}
	}
}

var (
	debugLayoutInternalCalls, debugCalcLayoutImplCalls   int
	debugLastAvailableWidth, debugLastAvailableHeight     float32
	debugLastWidthMode, debugLastHeightMode               SizingMode
)

// CalculateLayoutInternal is the caching wrapper around the layout implementation.
func CalculateLayoutInternal(node *Node, availableWidth, availableHeight float32, ownerDirection Direction,
	widthMode, heightMode SizingMode, ownerWidth, ownerHeight float32, performLayout bool,
	reason int, layoutMarkerData *LayoutData, depth int, generationCount uint32) bool {

	layout := node.GetLayout()
	depth++

	needToVisit := (node.IsDirty() && layout.generationCount != generationCount) ||
		layout.configVersion != node.GetConfig().GetVersion() ||
		layout.lastOwnerDirection != ownerDirection

	if needToVisit {
		layout.nextCachedMeasurementsIndex = 0
		layout.cachedLayout.availableWidth = -1
		layout.cachedLayout.availableHeight = -1
		layout.cachedLayout.widthSizingMode = SizingModeMaxContent
		layout.cachedLayout.heightSizingMode = SizingModeMaxContent
		layout.cachedLayout.computedWidth = -1
		layout.cachedLayout.computedHeight = -1
	}

	var cachedResults *cachedMeasurement

	if node.HasMeasureFunc() {
		marginAxisRow := node.Style().ComputeMarginForAxis(FlexDirectionRow, ownerWidth)
		marginAxisCol := node.Style().ComputeMarginForAxis(FlexDirectionColumn, ownerWidth)

		if canUseCachedMeasurement(widthMode, availableWidth, heightMode, availableHeight,
			layout.cachedLayout.widthSizingMode, layout.cachedLayout.availableWidth,
			layout.cachedLayout.heightSizingMode, layout.cachedLayout.availableHeight,
			layout.cachedLayout.computedWidth, layout.cachedLayout.computedHeight,
			marginAxisRow, marginAxisCol, node.GetConfig()) {
			cachedResults = &layout.cachedLayout
		} else {
			for i := uint32(0); i < layout.nextCachedMeasurementsIndex; i++ {
				if canUseCachedMeasurement(widthMode, availableWidth, heightMode, availableHeight,
					layout.cachedMeasurements[i].widthSizingMode, layout.cachedMeasurements[i].availableWidth,
					layout.cachedMeasurements[i].heightSizingMode, layout.cachedMeasurements[i].availableHeight,
					layout.cachedMeasurements[i].computedWidth, layout.cachedMeasurements[i].computedHeight,
					marginAxisRow, marginAxisCol, node.GetConfig()) {
					cachedResults = &layout.cachedMeasurements[i]
					break
				}
			}
		}
	} else if performLayout {
		if inexactEqualsFloat(layout.cachedLayout.availableWidth, availableWidth) &&
			inexactEqualsFloat(layout.cachedLayout.availableHeight, availableHeight) &&
			layout.cachedLayout.widthSizingMode == widthMode &&
			layout.cachedLayout.heightSizingMode == heightMode {
			cachedResults = &layout.cachedLayout
		}
	} else {
		for i := uint32(0); i < layout.nextCachedMeasurementsIndex; i++ {
			if inexactEqualsFloat(layout.cachedMeasurements[i].availableWidth, availableWidth) &&
				inexactEqualsFloat(layout.cachedMeasurements[i].availableHeight, availableHeight) &&
				layout.cachedMeasurements[i].widthSizingMode == widthMode &&
				layout.cachedMeasurements[i].heightSizingMode == heightMode {
				cachedResults = &layout.cachedMeasurements[i]
				break
			}
		}
	}

	if !needToVisit && cachedResults != nil {
		layout.SetMeasuredDimension(DimensionWidth, cachedResults.computedWidth)
		layout.SetMeasuredDimension(DimensionHeight, cachedResults.computedHeight)
		if performLayout {
			layoutMarkerData.CachedLayouts++
		} else {
			layoutMarkerData.CachedMeasures++
		}
	} else {
		calculateLayoutImpl(node, availableWidth, availableHeight, ownerDirection, widthMode, heightMode,
			ownerWidth, ownerHeight, performLayout, reason, layoutMarkerData, depth, generationCount)

		layout.lastOwnerDirection = ownerDirection
		layout.configVersion = node.GetConfig().GetVersion()

		if cachedResults == nil {
			if uint32(layout.nextCachedMeasurementsIndex)+1 > layoutMarkerData.MaxMeasureCache {
				layoutMarkerData.MaxMeasureCache = uint32(layout.nextCachedMeasurementsIndex) + 1
			}

			if layout.nextCachedMeasurementsIndex == maxCachedMeasurements {
				layout.nextCachedMeasurementsIndex = 0
			}

			var newCache *cachedMeasurement
			if performLayout {
				newCache = &layout.cachedLayout
			} else {
				newCache = &layout.cachedMeasurements[layout.nextCachedMeasurementsIndex]
				layout.nextCachedMeasurementsIndex++
			}

			newCache.availableWidth = availableWidth
			newCache.availableHeight = availableHeight
			newCache.widthSizingMode = widthMode
			newCache.heightSizingMode = heightMode
			newCache.computedWidth = layout.MeasuredDimension(DimensionWidth)
			newCache.computedHeight = layout.MeasuredDimension(DimensionHeight)
		}
	}

	if performLayout {
		node.SetLayoutDimension(layout.MeasuredDimension(DimensionWidth), DimensionWidth)
		node.SetLayoutDimension(layout.MeasuredDimension(DimensionHeight), DimensionHeight)
		node.SetHasNewLayout(true)
		node.SetDirty(false)
	}

	layout.generationCount = generationCount

	return (needToVisit || cachedResults == nil)
}

var gCurrentGenerationCount uint32

// CalculateLayout is the public entry point for performing layout on a node tree.
func CalculateLayout(node *Node, ownerWidth, ownerHeight float32, ownerDirection Direction) {
	layoutMarkerData := &LayoutData{}

	gCurrentGenerationCount++
	generationCount := gCurrentGenerationCount

	node.ProcessDimensions()
	direction := node.ResolveDirection(ownerDirection)

	width := float32(math.NaN())
	widthMode := SizingModeMaxContent
	style := node.Style()
	if node.HasDefiniteLength(DimensionWidth, ownerWidth) {
		width = node.GetResolvedDimension(direction, dimension(FlexDirectionRow), ownerWidth, ownerWidth).Unwrap() +
			style.ComputeMarginForAxis(FlexDirectionRow, ownerWidth)
		widthMode = SizingModeStretchFit
	} else if style.ResolvedMaxDimension(direction, DimensionWidth, ownerWidth, ownerWidth).IsDefined() {
		width = style.ResolvedMaxDimension(direction, DimensionWidth, ownerWidth, ownerWidth).Unwrap()
		widthMode = SizingModeFitContent
	} else {
		width = ownerWidth
		if isUndefined(width) {
			widthMode = SizingModeMaxContent
		} else {
			widthMode = SizingModeStretchFit
		}
	}

	height := float32(math.NaN())
	heightMode := SizingModeMaxContent
	if node.HasDefiniteLength(DimensionHeight, ownerHeight) {
		height = node.GetResolvedDimension(direction, dimension(FlexDirectionColumn), ownerHeight, ownerWidth).Unwrap() +
			style.ComputeMarginForAxis(FlexDirectionColumn, ownerWidth)
		heightMode = SizingModeStretchFit
	} else if style.ResolvedMaxDimension(direction, DimensionHeight, ownerHeight, ownerWidth).IsDefined() {
		height = style.ResolvedMaxDimension(direction, DimensionHeight, ownerHeight, ownerWidth).Unwrap()
		heightMode = SizingModeFitContent
	} else {
		height = ownerHeight
		if isUndefined(height) {
			heightMode = SizingModeMaxContent
		} else {
			heightMode = SizingModeStretchFit
		}
	}

	if CalculateLayoutInternal(node, width, height, ownerDirection, widthMode, heightMode,
		ownerWidth, ownerHeight, true, LayoutPassInitial, layoutMarkerData, 0, generationCount) {
		node.SetPosition(node.layout.Direction(), ownerWidth, ownerHeight)
		roundLayoutResultsToPixelGrid(node, 0.0, 0.0)
	}
}

func init() {
	_ = math.NaN()
}
