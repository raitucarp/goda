package goda

func resolveDirection(fd FlexDirection, dir Direction) FlexDirection {
	if dir == DirectionRTL {
		if fd == FlexDirectionRow {
			return FlexDirectionRowReverse
		} else if fd == FlexDirectionRowReverse {
			return FlexDirectionRow
		}
	}
	return fd
}

func resolveCrossDirection(fd FlexDirection, dir Direction) FlexDirection {
	if isColumn(fd) {
		return resolveDirection(FlexDirectionRow, dir)
	}
	return FlexDirectionColumn
}

func flexStartEdge(fd FlexDirection) PhysicalEdge {
	switch fd {
	case FlexDirectionColumn:
		return PhysicalEdgeTop
	case FlexDirectionColumnReverse:
		return PhysicalEdgeBottom
	case FlexDirectionRow:
		return PhysicalEdgeLeft
	case FlexDirectionRowReverse:
		return PhysicalEdgeRight
	}
	panic("Invalid FlexDirection")
}

func flexEndEdge(fd FlexDirection) PhysicalEdge {
	switch fd {
	case FlexDirectionColumn:
		return PhysicalEdgeBottom
	case FlexDirectionColumnReverse:
		return PhysicalEdgeTop
	case FlexDirectionRow:
		return PhysicalEdgeRight
	case FlexDirectionRowReverse:
		return PhysicalEdgeLeft
	}
	panic("Invalid FlexDirection")
}

func inlineStartEdge(fd FlexDirection, dir Direction) PhysicalEdge {
	if isRow(fd) {
		if dir == DirectionRTL {
			return PhysicalEdgeRight
		}
		return PhysicalEdgeLeft
	}
	return PhysicalEdgeTop
}

func inlineEndEdge(fd FlexDirection, dir Direction) PhysicalEdge {
	if isRow(fd) {
		if dir == DirectionRTL {
			return PhysicalEdgeLeft
		}
		return PhysicalEdgeRight
	}
	return PhysicalEdgeBottom
}

func paddingAndBorderForAxis(node *Node, axis FlexDirection, dir Direction, widthSize float32) float32 {
	return node.style.ComputeInlineStartPaddingAndBorder(axis, dir, widthSize) +
		node.style.ComputeInlineEndPaddingAndBorder(axis, dir, widthSize)
}

func boundAxisWithinMinAndMax(node *Node, dir Direction, axis FlexDirection, value FloatOptional, axisSize, widthSize float32) FloatOptional {
	var min, max FloatOptional
	if isColumn(axis) {
		min = node.style.ResolvedMinDimension(dir, DimensionHeight, axisSize, widthSize)
		max = node.style.ResolvedMaxDimension(dir, DimensionHeight, axisSize, widthSize)
	} else if isRow(axis) {
		min = node.style.ResolvedMinDimension(dir, DimensionWidth, axisSize, widthSize)
		max = node.style.ResolvedMaxDimension(dir, DimensionWidth, axisSize, widthSize)
	}
	if max.IsDefined() && max.Unwrap() >= 0 && value.Unwrap() > max.Unwrap() {
		return max
	}
	if min.IsDefined() && min.Unwrap() >= 0 && value.Unwrap() < min.Unwrap() {
		return min
	}
	return value
}

func boundAxis(node *Node, axis FlexDirection, dir Direction, value, axisSize, widthSize float32) float32 {
	bounded := boundAxisWithinMinAndMax(node, dir, axis, NewFloatOptional(value), axisSize, widthSize).Unwrap()
	return maxOrDefined(bounded, paddingAndBorderForAxis(node, axis, dir, widthSize))
}

func boundAxisWithAutoMin(node *Node, axis FlexDirection, dir Direction, value, axisSize, widthSize float32) float32 {
	bounded := boundAxis(node, axis, dir, value, axisSize, widthSize)
	autoMin := node.layout.computedAutoMinMainSize
	if autoMin.IsDefined() && bounded < autoMin.Unwrap() {
		return autoMin.Unwrap()
	}
	return bounded
}
