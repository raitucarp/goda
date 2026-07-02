package goda

func resolveChildAlignment(node, child *Node) Align {
	align := child.Style().AlignSelf()
	if align == AlignAuto {
		align = node.Style().AlignItems()
	}
	if node.Style().Display() == DisplayFlex && align == AlignBaseline && isColumn(node.Style().FlexDirection()) {
		return AlignFlexStart
	}
	return align
}

func resolveChildJustification(node, child *Node) Justify {
	j := child.Style().JustifySelf()
	if j == JustifyAuto {
		return node.Style().JustifyItems()
	}
	return j
}

func fallbackAlignmentAlign(a Align) Align {
	switch a {
	case AlignSpaceBetween, AlignStretch, AlignSpaceAround, AlignSpaceEvenly:
		return AlignFlexStart
	}
	return a
}

func fallbackAlignmentJustify(j Justify) Justify {
	switch j {
	case JustifySpaceBetween, JustifySpaceAround, JustifySpaceEvenly:
		return JustifyFlexStart
	}
	return j
}
