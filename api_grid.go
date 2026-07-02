package goda

func (n *Node) SetGridColumnStart(v int32) *Node {
	n.style.SetGridColumnStart(GridLineFromInteger(v))
	n.MarkDirtyAndPropagate()
	return n
}
func (n *Node) SetGridColumnStartAuto() *Node {
	n.style.SetGridColumnStart(GridLineAuto())
	n.MarkDirtyAndPropagate()
	return n
}
func (n *Node) SetGridColumnStartSpan(span int32) *Node {
	n.style.SetGridColumnStart(GridLineSpan(span))
	n.MarkDirtyAndPropagate()
	return n
}
func (n *Node) GetGridColumnStart() int32 {
	gl := n.style.GridColumnStart()
	if gl.IsInteger() {
		return gl.Integer
	}
	return 0
}

func (n *Node) SetGridColumnEnd(v int32) *Node {
	n.style.SetGridColumnEnd(GridLineFromInteger(v))
	n.MarkDirtyAndPropagate()
	return n
}
func (n *Node) SetGridColumnEndAuto() *Node {
	n.style.SetGridColumnEnd(GridLineAuto())
	n.MarkDirtyAndPropagate()
	return n
}
func (n *Node) SetGridColumnEndSpan(span int32) *Node {
	n.style.SetGridColumnEnd(GridLineSpan(span))
	n.MarkDirtyAndPropagate()
	return n
}
func (n *Node) GetGridColumnEnd() int32 {
	gl := n.style.GridColumnEnd()
	if gl.IsInteger() {
		return gl.Integer
	}
	return 0
}

func (n *Node) SetGridRowStart(v int32) *Node {
	n.style.SetGridRowStart(GridLineFromInteger(v))
	n.MarkDirtyAndPropagate()
	return n
}
func (n *Node) SetGridRowStartAuto() *Node {
	n.style.SetGridRowStart(GridLineAuto())
	n.MarkDirtyAndPropagate()
	return n
}
func (n *Node) SetGridRowStartSpan(span int32) *Node {
	n.style.SetGridRowStart(GridLineSpan(span))
	n.MarkDirtyAndPropagate()
	return n
}
func (n *Node) GetGridRowStart() int32 {
	gl := n.style.GridRowStart()
	if gl.IsInteger() {
		return gl.Integer
	}
	return 0
}

func (n *Node) SetGridRowEnd(v int32) *Node {
	n.style.SetGridRowEnd(GridLineFromInteger(v))
	n.MarkDirtyAndPropagate()
	return n
}
func (n *Node) SetGridRowEndAuto() *Node {
	n.style.SetGridRowEnd(GridLineAuto())
	n.MarkDirtyAndPropagate()
	return n
}
func (n *Node) SetGridRowEndSpan(span int32) *Node {
	n.style.SetGridRowEnd(GridLineSpan(span))
	n.MarkDirtyAndPropagate()
	return n
}
func (n *Node) GetGridRowEnd() int32 {
	gl := n.style.GridRowEnd()
	if gl.IsInteger() {
		return gl.Integer
	}
	return 0
}

func (n *Node) SetGridTemplateColumnsCount(count int) *Node {
	n.style.ResizeGridTemplateColumns(count)
	n.MarkDirtyAndPropagate()
	return n
}
func (n *Node) SetGridTemplateColumn(index int, trackType GridTrackType, value float32) *Node {
	n.style.SetGridTemplateColumnAt(index, gridTrackSizeFromTypeAndValue(trackType, value))
	n.MarkDirtyAndPropagate()
	return n
}
func (n *Node) SetGridTemplateColumnMinMax(index int, minType GridTrackType, minVal float32, maxType GridTrackType, maxVal float32) *Node {
	n.style.SetGridTemplateColumnAt(index, GridTrackSizeMinmax(
		styleSizeLengthFromTypeAndValue(minType, minVal),
		styleSizeLengthFromTypeAndValue(maxType, maxVal)))
	n.MarkDirtyAndPropagate()
	return n
}

func (n *Node) SetGridTemplateRowsCount(count int) *Node {
	n.style.ResizeGridTemplateRows(count)
	n.MarkDirtyAndPropagate()
	return n
}
func (n *Node) SetGridTemplateRow(index int, trackType GridTrackType, value float32) *Node {
	n.style.SetGridTemplateRowAt(index, gridTrackSizeFromTypeAndValue(trackType, value))
	n.MarkDirtyAndPropagate()
	return n
}
func (n *Node) SetGridTemplateRowMinMax(index int, minType GridTrackType, minVal float32, maxType GridTrackType, maxVal float32) *Node {
	n.style.SetGridTemplateRowAt(index, GridTrackSizeMinmax(
		styleSizeLengthFromTypeAndValue(minType, minVal),
		styleSizeLengthFromTypeAndValue(maxType, maxVal)))
	n.MarkDirtyAndPropagate()
	return n
}

func (n *Node) SetGridAutoColumnsCount(count int) *Node {
	n.style.ResizeGridAutoColumns(count)
	n.MarkDirtyAndPropagate()
	return n
}
func (n *Node) SetGridAutoColumn(index int, trackType GridTrackType, value float32) *Node {
	n.style.SetGridAutoColumnAt(index, gridTrackSizeFromTypeAndValue(trackType, value))
	n.MarkDirtyAndPropagate()
	return n
}
func (n *Node) SetGridAutoColumnMinMax(index int, minType GridTrackType, minVal float32, maxType GridTrackType, maxVal float32) *Node {
	n.style.SetGridAutoColumnAt(index, GridTrackSizeMinmax(
		styleSizeLengthFromTypeAndValue(minType, minVal),
		styleSizeLengthFromTypeAndValue(maxType, maxVal)))
	n.MarkDirtyAndPropagate()
	return n
}

func (n *Node) SetGridAutoRowsCount(count int) *Node {
	n.style.ResizeGridAutoRows(count)
	n.MarkDirtyAndPropagate()
	return n
}
func (n *Node) SetGridAutoRow(index int, trackType GridTrackType, value float32) *Node {
	n.style.SetGridAutoRowAt(index, gridTrackSizeFromTypeAndValue(trackType, value))
	n.MarkDirtyAndPropagate()
	return n
}
func (n *Node) SetGridAutoRowMinMax(index int, minType GridTrackType, minVal float32, maxType GridTrackType, maxVal float32) *Node {
	n.style.SetGridAutoRowAt(index, GridTrackSizeMinmax(
		styleSizeLengthFromTypeAndValue(minType, minVal),
		styleSizeLengthFromTypeAndValue(maxType, maxVal)))
	n.MarkDirtyAndPropagate()
	return n
}

func gridTrackSizeFromTypeAndValue(trackType GridTrackType, value float32) GridTrackSize {
	switch trackType {
	case GridTrackTypePoints:
		return GridTrackSizeLength(value)
	case GridTrackTypePercent:
		return GridTrackSizePercent(value)
	case GridTrackTypeFr:
		return GridTrackSizeFr(value)
	case GridTrackTypeAuto:
		return GridTrackSizeAuto()
	case GridTrackTypeMinmax:
		return GridTrackSizeAuto()
	}
	return GridTrackSizeAuto()
}

func styleSizeLengthFromTypeAndValue(trackType GridTrackType, value float32) StyleSizeLength {
	switch trackType {
	case GridTrackTypePoints:
		return StyleSizeLengthPoints(value)
	case GridTrackTypePercent:
		return StyleSizeLengthPercent(value)
	case GridTrackTypeFr:
		return StyleSizeLengthStretch(value)
	case GridTrackTypeAuto:
		return StyleSizeLengthAuto()
	case GridTrackTypeMinmax:
		return StyleSizeLengthAuto()
	}
	return StyleSizeLengthAuto()
}
