package goda

import "math"

const (
	defaultFlexGrow      float32 = 0.0
	defaultFlexShrink    float32 = 0.0
	webDefaultFlexShrink float32 = 1.0
)

// Style holds all CSS properties for a single Node.
type Style struct {
	direction_      Direction
	flexDirection_  FlexDirection
	justifyContent_ Justify
	justifyItems_   Justify
	justifySelf_    Justify
	alignContent_   Align
	alignItems_     Align
	alignSelf_      Align
	positionType_   PositionType
	flexWrap_       Wrap
	overflow_       Overflow
	display_        Display
	boxSizing_      BoxSizing

	flex_         FloatOptional
	flexGrow_     FloatOptional
	flexShrink_   FloatOptional
	flexBasis_    StyleSizeLength
	margin_       [9]StyleLength
	position_     [9]StyleLength
	padding_      [9]StyleLength
	border_       [9]StyleLength
	gap_          [3]StyleLength
	dimensions_    [2]StyleSizeLength
	minDimensions_ [2]StyleSizeLength
	maxDimensions_ [2]StyleSizeLength
	aspectRatio_  FloatOptional

	gridTemplateColumns_ GridTrackList
	gridTemplateRows_    GridTrackList
	gridAutoColumns_     GridTrackList
	gridAutoRows_        GridTrackList
	gridColumnStart_     GridLine
	gridColumnEnd_       GridLine
	gridRowStart_        GridLine
	gridRowEnd_          GridLine
}

func NewStyle() Style {
	s := Style{
		direction_:      DirectionInherit,
		flexDirection_:   FlexDirectionColumn,
		justifyContent_:  JustifyFlexStart,
		justifyItems_:    JustifyStretch,
		justifySelf_:     JustifyAuto,
		alignContent_:    AlignFlexStart,
		alignItems_:      AlignStretch,
		alignSelf_:       AlignAuto,
		positionType_:    PositionTypeRelative,
		flexWrap_:        WrapNoWrap,
		overflow_:        OverflowVisible,
		display_:         DisplayFlex,
		boxSizing_:       BoxSizingBorderBox,
		flexBasis_:       StyleSizeLengthAuto(),
		dimensions_:      [2]StyleSizeLength{StyleSizeLengthAuto(), StyleSizeLengthAuto()},
	}
	return s
}

func (s *Style) Direction() Direction            { return s.direction_ }
func (s *Style) SetDirection(v Direction)         { s.direction_ = v }
func (s *Style) FlexDirection() FlexDirection     { return s.flexDirection_ }
func (s *Style) SetFlexDirection(v FlexDirection) { s.flexDirection_ = v }
func (s *Style) JustifyContent() Justify          { return s.justifyContent_ }
func (s *Style) SetJustifyContent(v Justify)      { s.justifyContent_ = v }
func (s *Style) JustifyItems() Justify            { return s.justifyItems_ }
func (s *Style) SetJustifyItems(v Justify)        { s.justifyItems_ = v }
func (s *Style) JustifySelf() Justify             { return s.justifySelf_ }
func (s *Style) SetJustifySelf(v Justify)         { s.justifySelf_ = v }
func (s *Style) AlignContent() Align              { return s.alignContent_ }
func (s *Style) SetAlignContent(v Align)          { s.alignContent_ = v }
func (s *Style) AlignItems() Align                { return s.alignItems_ }
func (s *Style) SetAlignItems(v Align)            { s.alignItems_ = v }
func (s *Style) AlignSelf() Align                 { return s.alignSelf_ }
func (s *Style) SetAlignSelf(v Align)             { s.alignSelf_ = v }
func (s *Style) PositionType() PositionType       { return s.positionType_ }
func (s *Style) SetPositionType(v PositionType)   { s.positionType_ = v }
func (s *Style) FlexWrap() Wrap                   { return s.flexWrap_ }
func (s *Style) SetFlexWrap(v Wrap)               { s.flexWrap_ = v }
func (s *Style) Overflow() Overflow               { return s.overflow_ }
func (s *Style) SetOverflow(v Overflow)           { s.overflow_ = v }
func (s *Style) Display() Display                 { return s.display_ }
func (s *Style) SetDisplay(v Display)             { s.display_ = v }
func (s *Style) BoxSizing() BoxSizing             { return s.boxSizing_ }
func (s *Style) SetBoxSizing(v BoxSizing)         { s.boxSizing_ = v }

func (s *Style) Flex() FloatOptional              { return s.flex_ }
func (s *Style) SetFlex(v FloatOptional)           { s.flex_ = v }
func (s *Style) FlexGrow() FloatOptional           { return s.flexGrow_ }
func (s *Style) SetFlexGrow(v FloatOptional)       { s.flexGrow_ = v }
func (s *Style) FlexShrink() FloatOptional         { return s.flexShrink_ }
func (s *Style) SetFlexShrink(v FloatOptional)     { s.flexShrink_ = v }
func (s *Style) FlexBasis() StyleSizeLength        { return s.flexBasis_ }
func (s *Style) SetFlexBasis(v StyleSizeLength)    { s.flexBasis_ = v }

func (s *Style) Margin(edge Edge) StyleLength       { return s.margin_[edge] }
func (s *Style) SetMargin(edge Edge, v StyleLength) { s.margin_[edge] = v }
func (s *Style) Position(edge Edge) StyleLength       { return s.position_[edge] }
func (s *Style) SetPosition(edge Edge, v StyleLength) { s.position_[edge] = v }
func (s *Style) Padding(edge Edge) StyleLength       { return s.padding_[edge] }
func (s *Style) SetPadding(edge Edge, v StyleLength) { s.padding_[edge] = v }
func (s *Style) Border(edge Edge) StyleLength        { return s.border_[edge] }
func (s *Style) SetBorder(edge Edge, v StyleLength)  { s.border_[edge] = v }
func (s *Style) Gap(gutter Gutter) StyleLength        { return s.gap_[gutter] }
func (s *Style) SetGap(gutter Gutter, v StyleLength)  { s.gap_[gutter] = v }

func (s *Style) Dimension(axis Dimension) StyleSizeLength        { return s.dimensions_[axis] }
func (s *Style) SetDimension(axis Dimension, v StyleSizeLength)   { s.dimensions_[axis] = v }
func (s *Style) MinDimension(axis Dimension) StyleSizeLength      { return s.minDimensions_[axis] }
func (s *Style) SetMinDimension(axis Dimension, v StyleSizeLength) { s.minDimensions_[axis] = v }
func (s *Style) MaxDimension(axis Dimension) StyleSizeLength      { return s.maxDimensions_[axis] }
func (s *Style) SetMaxDimension(axis Dimension, v StyleSizeLength) { s.maxDimensions_[axis] = v }

func (s *Style) AspectRatio() FloatOptional { return s.aspectRatio_ }
func (s *Style) SetAspectRatio(v FloatOptional) {
	if v.IsDefined() {
		val := v.Unwrap()
		if val == 0.0 || math.IsInf(float64(val), 0) {
			s.aspectRatio_ = FloatOptional{}
			return
		}
	}
	s.aspectRatio_ = v
}

// Grid container properties
func (s *Style) GridTemplateColumns() GridTrackList     { return s.gridTemplateColumns_ }
func (s *Style) SetGridTemplateColumns(v GridTrackList)  { s.gridTemplateColumns_ = v }
func (s *Style) ResizeGridTemplateColumns(count int) {
	if len(s.gridTemplateColumns_) != count {
		newList := make(GridTrackList, count)
		copy(newList, s.gridTemplateColumns_)
		s.gridTemplateColumns_ = newList
	}
}
func (s *Style) SetGridTemplateColumnAt(index int, v GridTrackSize) { s.gridTemplateColumns_[index] = v }

func (s *Style) GridTemplateRows() GridTrackList     { return s.gridTemplateRows_ }
func (s *Style) SetGridTemplateRows(v GridTrackList)  { s.gridTemplateRows_ = v }
func (s *Style) ResizeGridTemplateRows(count int) {
	if len(s.gridTemplateRows_) != count {
		newList := make(GridTrackList, count)
		copy(newList, s.gridTemplateRows_)
		s.gridTemplateRows_ = newList
	}
}
func (s *Style) SetGridTemplateRowAt(index int, v GridTrackSize) { s.gridTemplateRows_[index] = v }

func (s *Style) GridAutoColumns() GridTrackList     { return s.gridAutoColumns_ }
func (s *Style) SetGridAutoColumns(v GridTrackList)  { s.gridAutoColumns_ = v }
func (s *Style) ResizeGridAutoColumns(count int) {
	if len(s.gridAutoColumns_) != count {
		newList := make(GridTrackList, count)
		copy(newList, s.gridAutoColumns_)
		s.gridAutoColumns_ = newList
	}
}
func (s *Style) SetGridAutoColumnAt(index int, v GridTrackSize) { s.gridAutoColumns_[index] = v }

func (s *Style) GridAutoRows() GridTrackList     { return s.gridAutoRows_ }
func (s *Style) SetGridAutoRows(v GridTrackList)  { s.gridAutoRows_ = v }
func (s *Style) ResizeGridAutoRows(count int) {
	if len(s.gridAutoRows_) != count {
		newList := make(GridTrackList, count)
		copy(newList, s.gridAutoRows_)
		s.gridAutoRows_ = newList
	}
}
func (s *Style) SetGridAutoRowAt(index int, v GridTrackSize) { s.gridAutoRows_[index] = v }

// Grid item properties
func (s *Style) GridColumnStart() GridLine       { return s.gridColumnStart_ }
func (s *Style) SetGridColumnStart(v GridLine)    { s.gridColumnStart_ = v }
func (s *Style) GridColumnEnd() GridLine         { return s.gridColumnEnd_ }
func (s *Style) SetGridColumnEnd(v GridLine)      { s.gridColumnEnd_ = v }
func (s *Style) GridRowStart() GridLine          { return s.gridRowStart_ }
func (s *Style) SetGridRowStart(v GridLine)       { s.gridRowStart_ = v }
func (s *Style) GridRowEnd() GridLine            { return s.gridRowEnd_ }
func (s *Style) SetGridRowEnd(v GridLine)         { s.gridRowEnd_ = v }

// ResolvedMinDimension returns the resolved minimum size for the given axis.
func (s *Style) ResolvedMinDimension(direction Direction, axis Dimension, referenceLength, ownerWidth float32) FloatOptional {
	sl := s.minDimensions_[axis]
	if sl.IsUndefined() {
		return FloatOptional{}
	}
	value := sl.Resolve(referenceLength)
	if s.boxSizing_ == BoxSizingBorderBox || !value.IsDefined() {
		return value
	}
	dimPB := NewFloatOptional(s.ComputePaddingAndBorderForDimension(direction, axis, ownerWidth))
	return NewFloatOptional(value.Unwrap() + dimPB.Unwrap())
}

// ResolvedMaxDimension returns the resolved maximum size for the given axis.
func (s *Style) ResolvedMaxDimension(direction Direction, axis Dimension, referenceLength, ownerWidth float32) FloatOptional {
	sl := s.maxDimensions_[axis]
	if sl.IsUndefined() {
		return FloatOptional{}
	}
	value := sl.Resolve(referenceLength)
	if s.boxSizing_ == BoxSizingBorderBox || !value.IsDefined() {
		return value
	}
	dimPB := NewFloatOptional(s.ComputePaddingAndBorderForDimension(direction, axis, ownerWidth))
	return NewFloatOptional(value.Unwrap() + dimPB.Unwrap())
}

// Edge computation helpers
func (s *Style) computeLeftEdge(edges *[9]StyleLength, dir Direction) StyleLength {
	if dir == DirectionLTR && edges[EdgeStart].IsDefined() {
		return edges[EdgeStart]
	} else if dir == DirectionRTL && edges[EdgeEnd].IsDefined() {
		return edges[EdgeEnd]
	} else if edges[EdgeLeft].IsDefined() {
		return edges[EdgeLeft]
	} else if edges[EdgeHorizontal].IsDefined() {
		return edges[EdgeHorizontal]
	}
	return edges[EdgeAll]
}

func (s *Style) computeTopEdge(edges *[9]StyleLength) StyleLength {
	if edges[EdgeTop].IsDefined() {
		return edges[EdgeTop]
	} else if edges[EdgeVertical].IsDefined() {
		return edges[EdgeVertical]
	}
	return edges[EdgeAll]
}

func (s *Style) computeRightEdge(edges *[9]StyleLength, dir Direction) StyleLength {
	if dir == DirectionLTR && edges[EdgeEnd].IsDefined() {
		return edges[EdgeEnd]
	} else if dir == DirectionRTL && edges[EdgeStart].IsDefined() {
		return edges[EdgeStart]
	} else if edges[EdgeRight].IsDefined() {
		return edges[EdgeRight]
	} else if edges[EdgeHorizontal].IsDefined() {
		return edges[EdgeHorizontal]
	}
	return edges[EdgeAll]
}

func (s *Style) computeBottomEdge(edges *[9]StyleLength) StyleLength {
	if edges[EdgeBottom].IsDefined() {
		return edges[EdgeBottom]
	} else if edges[EdgeVertical].IsDefined() {
		return edges[EdgeVertical]
	}
	return edges[EdgeAll]
}

func (s *Style) computePosition(edge PhysicalEdge, dir Direction) StyleLength {
	switch edge {
	case PhysicalEdgeLeft:
		return s.computeLeftEdge(&s.position_, dir)
	case PhysicalEdgeTop:
		return s.computeTopEdge(&s.position_)
	case PhysicalEdgeRight:
		return s.computeRightEdge(&s.position_, dir)
	case PhysicalEdgeBottom:
		return s.computeBottomEdge(&s.position_)
	default:
		panic("Invalid physical edge")
	}
}

func (s *Style) computeMargin(edge PhysicalEdge, dir Direction) StyleLength {
	switch edge {
	case PhysicalEdgeLeft:
		return s.computeLeftEdge(&s.margin_, dir)
	case PhysicalEdgeTop:
		return s.computeTopEdge(&s.margin_)
	case PhysicalEdgeRight:
		return s.computeRightEdge(&s.margin_, dir)
	case PhysicalEdgeBottom:
		return s.computeBottomEdge(&s.margin_)
	default:
		panic("Invalid physical edge")
	}
}

func (s *Style) computePadding(edge PhysicalEdge, dir Direction) StyleLength {
	switch edge {
	case PhysicalEdgeLeft:
		return s.computeLeftEdge(&s.padding_, dir)
	case PhysicalEdgeTop:
		return s.computeTopEdge(&s.padding_)
	case PhysicalEdgeRight:
		return s.computeRightEdge(&s.padding_, dir)
	case PhysicalEdgeBottom:
		return s.computeBottomEdge(&s.padding_)
	default:
		panic("Invalid physical edge")
	}
}

func (s *Style) computeBorder(edge PhysicalEdge, dir Direction) StyleLength {
	switch edge {
	case PhysicalEdgeLeft:
		return s.computeLeftEdge(&s.border_, dir)
	case PhysicalEdgeTop:
		return s.computeTopEdge(&s.border_)
	case PhysicalEdgeRight:
		return s.computeRightEdge(&s.border_, dir)
	case PhysicalEdgeBottom:
		return s.computeBottomEdge(&s.border_)
	default:
		panic("Invalid physical edge")
	}
}

func (s *Style) HorizontalInsetsDefined() bool {
	return s.position_[EdgeLeft].IsDefined() ||
		s.position_[EdgeRight].IsDefined() ||
		s.position_[EdgeAll].IsDefined() ||
		s.position_[EdgeHorizontal].IsDefined() ||
		s.position_[EdgeStart].IsDefined() ||
		s.position_[EdgeEnd].IsDefined()
}

func (s *Style) VerticalInsetsDefined() bool {
	return s.position_[EdgeTop].IsDefined() ||
		s.position_[EdgeBottom].IsDefined() ||
		s.position_[EdgeAll].IsDefined() ||
		s.position_[EdgeVertical].IsDefined()
}

func (s *Style) IsFlexStartPositionDefined(axis FlexDirection, dir Direction) bool {
	return s.computePosition(flexStartEdge(axis), dir).IsDefined()
}
func (s *Style) IsFlexStartPositionAuto(axis FlexDirection, dir Direction) bool {
	return s.computePosition(flexStartEdge(axis), dir).IsAuto()
}
func (s *Style) IsInlineStartPositionDefined(axis FlexDirection, dir Direction) bool {
	return s.computePosition(inlineStartEdge(axis, dir), dir).IsDefined()
}
func (s *Style) IsInlineStartPositionAuto(axis FlexDirection, dir Direction) bool {
	return s.computePosition(inlineStartEdge(axis, dir), dir).IsAuto()
}
func (s *Style) IsFlexEndPositionDefined(axis FlexDirection, dir Direction) bool {
	return s.computePosition(flexEndEdge(axis), dir).IsDefined()
}
func (s *Style) IsFlexEndPositionAuto(axis FlexDirection, dir Direction) bool {
	return s.computePosition(flexEndEdge(axis), dir).IsAuto()
}
func (s *Style) IsInlineEndPositionDefined(axis FlexDirection, dir Direction) bool {
	return s.computePosition(inlineEndEdge(axis, dir), dir).IsDefined()
}
func (s *Style) IsInlineEndPositionAuto(axis FlexDirection, dir Direction) bool {
	return s.computePosition(inlineEndEdge(axis, dir), dir).IsAuto()
}

func (s *Style) ComputeFlexStartPosition(axis FlexDirection, dir Direction, axisSize float32) float32 {
	return s.computePosition(flexStartEdge(axis), dir).Resolve(axisSize).UnwrapOrDefault(0.0)
}
func (s *Style) ComputeInlineStartPosition(axis FlexDirection, dir Direction, axisSize float32) float32 {
	return s.computePosition(inlineStartEdge(axis, dir), dir).Resolve(axisSize).UnwrapOrDefault(0.0)
}
func (s *Style) ComputeFlexEndPosition(axis FlexDirection, dir Direction, axisSize float32) float32 {
	return s.computePosition(flexEndEdge(axis), dir).Resolve(axisSize).UnwrapOrDefault(0.0)
}
func (s *Style) ComputeInlineEndPosition(axis FlexDirection, dir Direction, axisSize float32) float32 {
	return s.computePosition(inlineEndEdge(axis, dir), dir).Resolve(axisSize).UnwrapOrDefault(0.0)
}

// Margin helpers
func (s *Style) ComputeFlexStartMargin(axis FlexDirection, dir Direction, widthSize float32) float32 {
	return s.computeMargin(flexStartEdge(axis), dir).Resolve(widthSize).UnwrapOrDefault(0.0)
}
func (s *Style) ComputeInlineStartMargin(axis FlexDirection, dir Direction, widthSize float32) float32 {
	return s.computeMargin(inlineStartEdge(axis, dir), dir).Resolve(widthSize).UnwrapOrDefault(0.0)
}
func (s *Style) ComputeFlexEndMargin(axis FlexDirection, dir Direction, widthSize float32) float32 {
	return s.computeMargin(flexEndEdge(axis), dir).Resolve(widthSize).UnwrapOrDefault(0.0)
}
func (s *Style) ComputeInlineEndMargin(axis FlexDirection, dir Direction, widthSize float32) float32 {
	return s.computeMargin(inlineEndEdge(axis, dir), dir).Resolve(widthSize).UnwrapOrDefault(0.0)
}

// Border helpers
func (s *Style) ComputeFlexStartBorder(axis FlexDirection, dir Direction) float32 {
	v := s.computeBorder(flexStartEdge(axis), dir).Resolve(0).Unwrap()
	return maxOrDefined(v, 0.0)
}
func (s *Style) ComputeInlineStartBorder(axis FlexDirection, dir Direction) float32 {
	v := s.computeBorder(inlineStartEdge(axis, dir), dir).Resolve(0).Unwrap()
	return maxOrDefined(v, 0.0)
}
func (s *Style) ComputeFlexEndBorder(axis FlexDirection, dir Direction) float32 {
	v := s.computeBorder(flexEndEdge(axis), dir).Resolve(0).Unwrap()
	return maxOrDefined(v, 0.0)
}
func (s *Style) ComputeInlineEndBorder(axis FlexDirection, dir Direction) float32 {
	v := s.computeBorder(inlineEndEdge(axis, dir), dir).Resolve(0).Unwrap()
	return maxOrDefined(v, 0.0)
}

// Padding helpers
func (s *Style) ComputeFlexStartPadding(axis FlexDirection, dir Direction, widthSize float32) float32 {
	v := s.computePadding(flexStartEdge(axis), dir).Resolve(widthSize).Unwrap()
	return maxOrDefined(v, 0.0)
}
func (s *Style) ComputeInlineStartPadding(axis FlexDirection, dir Direction, widthSize float32) float32 {
	v := s.computePadding(inlineStartEdge(axis, dir), dir).Resolve(widthSize).Unwrap()
	return maxOrDefined(v, 0.0)
}
func (s *Style) ComputeFlexEndPadding(axis FlexDirection, dir Direction, widthSize float32) float32 {
	v := s.computePadding(flexEndEdge(axis), dir).Resolve(widthSize).Unwrap()
	return maxOrDefined(v, 0.0)
}
func (s *Style) ComputeInlineEndPadding(axis FlexDirection, dir Direction, widthSize float32) float32 {
	v := s.computePadding(inlineEndEdge(axis, dir), dir).Resolve(widthSize).Unwrap()
	return maxOrDefined(v, 0.0)
}

func (s *Style) ComputeInlineStartPaddingAndBorder(axis FlexDirection, dir Direction, widthSize float32) float32 {
	return s.ComputeInlineStartPadding(axis, dir, widthSize) + s.ComputeInlineStartBorder(axis, dir)
}
func (s *Style) ComputeFlexStartPaddingAndBorder(axis FlexDirection, dir Direction, widthSize float32) float32 {
	return s.ComputeFlexStartPadding(axis, dir, widthSize) + s.ComputeFlexStartBorder(axis, dir)
}
func (s *Style) ComputeInlineEndPaddingAndBorder(axis FlexDirection, dir Direction, widthSize float32) float32 {
	return s.ComputeInlineEndPadding(axis, dir, widthSize) + s.ComputeInlineEndBorder(axis, dir)
}
func (s *Style) ComputeFlexEndPaddingAndBorder(axis FlexDirection, dir Direction, widthSize float32) float32 {
	return s.ComputeFlexEndPadding(axis, dir, widthSize) + s.ComputeFlexEndBorder(axis, dir)
}

func (s *Style) ComputePaddingAndBorderForDimension(dir Direction, dim Dimension, widthSize float32) float32 {
	flexDir := FlexDirectionRow
	if dim == DimensionHeight {
		flexDir = FlexDirectionColumn
	}
	return s.ComputeFlexStartPaddingAndBorder(flexDir, dir, widthSize) +
		s.ComputeFlexEndPaddingAndBorder(flexDir, dir, widthSize)
}

func (s *Style) ComputeBorderForAxis(axis FlexDirection) float32 {
	return s.ComputeInlineStartBorder(axis, DirectionLTR) +
		s.ComputeInlineEndBorder(axis, DirectionLTR)
}

func (s *Style) ComputeMarginForAxis(axis FlexDirection, widthSize float32) float32 {
	return s.ComputeInlineStartMargin(axis, DirectionLTR, widthSize) +
		s.ComputeInlineEndMargin(axis, DirectionLTR, widthSize)
}

func (s *Style) ComputeGapForAxis(axis FlexDirection, ownerSize float32) float32 {
	var gap StyleLength
	if isRow(axis) {
		gap = s.computeColumnGap()
	} else {
		gap = s.computeRowGap()
	}
	return maxOrDefined(gap.Resolve(ownerSize).Unwrap(), 0.0)
}

func (s *Style) ComputeGapForDimension(dim Dimension, ownerSize float32) float32 {
	var gap StyleLength
	if dim == DimensionWidth {
		gap = s.computeColumnGap()
	} else {
		gap = s.computeRowGap()
	}
	return maxOrDefined(gap.Resolve(ownerSize).Unwrap(), 0.0)
}

func (s *Style) computeColumnGap() StyleLength {
	if s.gap_[GutterColumn].IsDefined() {
		return s.gap_[GutterColumn]
	}
	return s.gap_[GutterAll]
}

func (s *Style) computeRowGap() StyleLength {
	if s.gap_[GutterRow].IsDefined() {
		return s.gap_[GutterRow]
	}
	return s.gap_[GutterAll]
}

func (s *Style) FlexStartMarginIsAuto(axis FlexDirection, dir Direction) bool {
	return s.computeMargin(flexStartEdge(axis), dir).IsAuto()
}
func (s *Style) FlexEndMarginIsAuto(axis FlexDirection, dir Direction) bool {
	return s.computeMargin(flexEndEdge(axis), dir).IsAuto()
}
func (s *Style) InlineStartMarginIsAuto(axis FlexDirection, dir Direction) bool {
	return s.computeMargin(inlineStartEdge(axis, dir), dir).IsAuto()
}
func (s *Style) InlineEndMarginIsAuto(axis FlexDirection, dir Direction) bool {
	return s.computeMargin(inlineEndEdge(axis, dir), dir).IsAuto()
}

// Equals returns true if all style properties match.
func (s *Style) Equals(other *Style) bool {
	if s.direction_ != other.direction_ || s.flexDirection_ != other.flexDirection_ ||
		s.justifyContent_ != other.justifyContent_ || s.justifyItems_ != other.justifyItems_ ||
		s.justifySelf_ != other.justifySelf_ || s.alignContent_ != other.alignContent_ ||
		s.alignItems_ != other.alignItems_ || s.alignSelf_ != other.alignSelf_ ||
		s.positionType_ != other.positionType_ || s.flexWrap_ != other.flexWrap_ ||
		s.overflow_ != other.overflow_ || s.display_ != other.display_ ||
		s.boxSizing_ != other.boxSizing_ {
		return false
	}
	if !s.flex_.Equals(other.flex_) || !s.flexGrow_.Equals(other.flexGrow_) ||
		!s.flexShrink_.Equals(other.flexShrink_) || !s.flexBasis_.Equals(other.flexBasis_) ||
		!s.aspectRatio_.Equals(other.aspectRatio_) {
		return false
	}
	for i := 0; i <= int(EdgeAll); i++ {
		if s.margin_[i] != other.margin_[i] || s.position_[i] != other.position_[i] ||
			s.padding_[i] != other.padding_[i] || s.border_[i] != other.border_[i] {
			return false
		}
	}
	for i := 0; i <= int(GutterAll); i++ {
		if s.gap_[i] != other.gap_[i] {
			return false
		}
	}
	for i := 0; i < 2; i++ {
		if !s.dimensions_[i].Equals(other.dimensions_[i]) || !s.minDimensions_[i].Equals(other.minDimensions_[i]) ||
			!s.maxDimensions_[i].Equals(other.maxDimensions_[i]) {
			return false
		}
	}
	if len(s.gridTemplateColumns_) != len(other.gridTemplateColumns_) {
		return false
	}
	for i := range s.gridTemplateColumns_ {
		if s.gridTemplateColumns_[i] != other.gridTemplateColumns_[i] {
			return false
		}
	}
	if len(s.gridTemplateRows_) != len(other.gridTemplateRows_) {
		return false
	}
	for i := range s.gridTemplateRows_ {
		if s.gridTemplateRows_[i] != other.gridTemplateRows_[i] {
			return false
		}
	}
	if len(s.gridAutoColumns_) != len(other.gridAutoColumns_) {
		return false
	}
	for i := range s.gridAutoColumns_ {
		if s.gridAutoColumns_[i] != other.gridAutoColumns_[i] {
			return false
		}
	}
	if len(s.gridAutoRows_) != len(other.gridAutoRows_) {
		return false
	}
	for i := range s.gridAutoRows_ {
		if s.gridAutoRows_[i] != other.gridAutoRows_[i] {
			return false
		}
	}
	if s.gridColumnStart_ != other.gridColumnStart_ || s.gridColumnEnd_ != other.gridColumnEnd_ ||
		s.gridRowStart_ != other.gridRowStart_ || s.gridRowEnd_ != other.gridRowEnd_ {
		return false
	}
	return true
}

func copyGridTrackList(src GridTrackList) GridTrackList {
	if src == nil {
		return nil
	}
	dst := make(GridTrackList, len(src))
	copy(dst, src)
	return dst
}

func (s *Style) Copy() Style {
	c := *s
	c.gridTemplateColumns_ = copyGridTrackList(s.gridTemplateColumns_)
	c.gridTemplateRows_ = copyGridTrackList(s.gridTemplateRows_)
	c.gridAutoColumns_ = copyGridTrackList(s.gridAutoColumns_)
	c.gridAutoRows_ = copyGridTrackList(s.gridAutoRows_)
	return c
}
