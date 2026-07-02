package goda

func justifyMainAxis(node *Node, flexLine *FlexLine, mainAxis, crossAxis FlexDirection,
	direction Direction, sizingModeMainDim, sizingModeCrossDim SizingMode,
	mainAxisOwnerSize, ownerWidth, availableInnerMainDim, availableInnerCrossDim,
	availableInnerWidth float32, performLayout bool) {

	style := node.Style()
	leadingPBMain := style.ComputeFlexStartPaddingAndBorder(mainAxis, direction, ownerWidth)
	trailingPBMain := style.ComputeFlexEndPaddingAndBorder(mainAxis, direction, ownerWidth)
	gap := style.ComputeGapForAxis(mainAxis, availableInnerMainDim)

	if sizingModeMainDim == SizingModeFitContent && flexLine.Layout.RemainingFreeSpace > 0 {
		if style.MinDimension(dimension(mainAxis)).IsDefined() &&
			style.ResolvedMinDimension(direction, dimension(mainAxis), mainAxisOwnerSize, ownerWidth).IsDefined() {
			minAvail := style.ResolvedMinDimension(direction, dimension(mainAxis), mainAxisOwnerSize, ownerWidth).Unwrap() -
				leadingPBMain - trailingPBMain
			occupied := availableInnerMainDim - flexLine.Layout.RemainingFreeSpace
			flexLine.Layout.RemainingFreeSpace = maxOrDefined(0.0, minAvail-occupied)
		} else {
			flexLine.Layout.RemainingFreeSpace = 0
		}
	}

	leadingMainDim := float32(0.0)
	betweenMainDim := gap

	justifyContent := style.JustifyContent()
	if flexLine.Layout.RemainingFreeSpace < 0 {
		justifyContent = fallbackAlignmentJustify(style.JustifyContent())
	}

	if flexLine.NumberOfAutoMargins == 0 {
		switch justifyContent {
		case JustifyCenter:
			leadingMainDim = flexLine.Layout.RemainingFreeSpace / 2
		case JustifyFlexEnd, JustifyEnd:
			leadingMainDim = flexLine.Layout.RemainingFreeSpace
		case JustifySpaceBetween:
			if len(flexLine.ItemsInFlow) > 1 {
				betweenMainDim += flexLine.Layout.RemainingFreeSpace / float32(len(flexLine.ItemsInFlow)-1)
			}
		case JustifySpaceEvenly:
			leadingMainDim = flexLine.Layout.RemainingFreeSpace / float32(len(flexLine.ItemsInFlow)+1)
			betweenMainDim += leadingMainDim
		case JustifySpaceAround:
			leadingMainDim = 0.5 * flexLine.Layout.RemainingFreeSpace / float32(len(flexLine.ItemsInFlow))
			betweenMainDim += leadingMainDim * 2
		}
	}

	flexLine.Layout.MainDim = leadingPBMain + leadingMainDim
	flexLine.Layout.CrossDim = 0

	maxAscent := float32(0.0)
	maxDescent := float32(0.0)
	isBaseLayout := isBaselineLayout(node)

	for i, child := range flexLine.ItemsInFlow {
		childLayout := child.GetLayout()
		if child.Style().FlexStartMarginIsAuto(mainAxis, direction) && flexLine.Layout.RemainingFreeSpace > 0 {
			flexLine.Layout.MainDim += flexLine.Layout.RemainingFreeSpace / float32(flexLine.NumberOfAutoMargins)
		}

		if performLayout {
			child.SetLayoutPosition(childLayout.Position(flexStartEdge(mainAxis))+flexLine.Layout.MainDim, flexStartEdge(mainAxis))
		}

		if i < len(flexLine.ItemsInFlow)-1 {
			flexLine.Layout.MainDim += betweenMainDim
		}

		if child.Style().FlexEndMarginIsAuto(mainAxis, direction) && flexLine.Layout.RemainingFreeSpace > 0 {
			flexLine.Layout.MainDim += flexLine.Layout.RemainingFreeSpace / float32(flexLine.NumberOfAutoMargins)
		}

		canSkipFlex := !performLayout && sizingModeCrossDim == SizingModeStretchFit
		if canSkipFlex {
			flexLine.Layout.MainDim += child.Style().ComputeMarginForAxis(mainAxis, availableInnerWidth) +
				boundAxisWithinMinAndMax(child, direction, mainAxis, childLayout.computedFlexBasis, mainAxisOwnerSize, ownerWidth).Unwrap()
			flexLine.Layout.CrossDim = availableInnerCrossDim
		} else {
			flexLine.Layout.MainDim += child.DimensionWithMargin(mainAxis, availableInnerWidth)

			if isBaseLayout {
				ascent := calculateBaseline(child) +
					child.Style().ComputeFlexStartMargin(FlexDirectionColumn, direction, availableInnerWidth)
				descent := child.layout.MeasuredDimension(DimensionHeight) +
					child.Style().ComputeMarginForAxis(FlexDirectionColumn, availableInnerWidth) - ascent
				maxAscent = maxOrDefined(maxAscent, ascent)
				maxDescent = maxOrDefined(maxDescent, descent)
			} else {
				flexLine.Layout.CrossDim = maxOrDefined(flexLine.Layout.CrossDim, child.DimensionWithMargin(crossAxis, availableInnerWidth))
			}
		}
	}
	flexLine.Layout.MainDim += trailingPBMain

	if isBaseLayout {
		flexLine.Layout.CrossDim = maxAscent + maxDescent
	}
}
