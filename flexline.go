package goda

// FlexLine represents a single line of flex items.
type FlexLine struct {
	ItemsInFlow         []*Node
	SizeConsumed        float32
	NumberOfAutoMargins int
	Layout              FlexLineRunningLayout
}

// FlexLineRunningLayout holds transient layout state for a flex line.
type FlexLineRunningLayout struct {
	TotalFlexGrowFactors         float32
	TotalFlexShrinkScaledFactors float32
	RemainingFreeSpace           float32
	MainDim                      float32
	CrossDim                     float32
}

func calculateFlexLine(node *Node, ownerDir Direction, ownerWidth, mainAxisOwnerSize, availableInnerWidth, availableInnerMainDim float32, iterator *LayoutableIterator, lineCount int) FlexLine {
	itemsInFlow := make([]*Node, 0, len(node.children))
	var sizeConsumed, totalFlexGrow, totalFlexShrinkScaled float32
	var numberOfAutoMargins int
	var firstElementInLine *Node

	sizeConsumedIncludingMinConstraint := float32(0)
	direction := node.ResolveDirection(ownerDir)
	mainAxis := resolveDirection(node.Style().FlexDirection(), direction)
	isNodeFlexWrap := node.Style().FlexWrap() != WrapNoWrap
	gap := node.Style().ComputeGapForAxis(mainAxis, availableInnerMainDim)

	for child := iterator.Current(); child != nil; child = iterator.Current() {
		if child.Style().Display() == DisplayNone || child.Style().PositionType() == PositionTypeAbsolute {
			if !iterator.Next() {
				break
			}
			continue
		}

		if firstElementInLine == nil {
			firstElementInLine = child
		}

		if child.Style().FlexStartMarginIsAuto(mainAxis, ownerDir) {
			numberOfAutoMargins++
		}
		if child.Style().FlexEndMarginIsAuto(mainAxis, ownerDir) {
			numberOfAutoMargins++
		}

		child.SetLineIndex(lineCount)
		childMarginMainAxis := child.Style().ComputeMarginForAxis(mainAxis, availableInnerWidth)
		childLeadingGap := float32(0.0)
		if child != firstElementInLine {
			childLeadingGap = gap
		}

		flexBasis := boundAxisWithinMinAndMax(child, direction, mainAxis, child.GetLayout().computedFlexBasis, mainAxisOwnerSize, ownerWidth).Unwrap()

		if sizeConsumedIncludingMinConstraint+flexBasis+childMarginMainAxis+childLeadingGap > availableInnerMainDim &&
			isNodeFlexWrap && len(itemsInFlow) > 0 {
			break
		}

		sizeConsumedIncludingMinConstraint += flexBasis + childMarginMainAxis + childLeadingGap
		sizeConsumed += flexBasis + childMarginMainAxis + childLeadingGap

		if child.IsNodeFlexible() {
			totalFlexGrow += child.ResolveFlexGrow()
			totalFlexShrinkScaled += -child.ResolveFlexShrink() * child.GetLayout().computedFlexBasis.Unwrap()
		}

		itemsInFlow = append(itemsInFlow, child)
		if !iterator.Next() {
			break
		}
	}

	if totalFlexGrow > 0 && totalFlexGrow < 1 {
		totalFlexGrow = 1
	}
	if totalFlexShrinkScaled > 0 && totalFlexShrinkScaled < 1 {
		totalFlexShrinkScaled = 1
	}

	return FlexLine{
		ItemsInFlow:         itemsInFlow,
		SizeConsumed:        sizeConsumed,
		NumberOfAutoMargins: numberOfAutoMargins,
		Layout: FlexLineRunningLayout{
			TotalFlexGrowFactors:         totalFlexGrow,
			TotalFlexShrinkScaledFactors:  totalFlexShrinkScaled,
		},
	}
}
