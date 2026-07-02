package goda

import "math"

func (n *Node) SetDirection(d Direction) *Node {
	if n.style.Direction() != d {
		n.style.SetDirection(d)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) GetDirection() Direction { return n.style.Direction() }

func (n *Node) SetFlexDirection(fd FlexDirection) *Node {
	if n.style.FlexDirection() != fd {
		n.style.SetFlexDirection(fd)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) GetFlexDirection() FlexDirection { return n.style.FlexDirection() }

func (n *Node) SetJustifyContent(j Justify) *Node {
	if n.style.JustifyContent() != j {
		n.style.SetJustifyContent(j)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) GetJustifyContent() Justify { return n.style.JustifyContent() }

func (n *Node) SetJustifyItems(j Justify) *Node {
	if n.style.JustifyItems() != j {
		n.style.SetJustifyItems(j)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) GetJustifyItems() Justify { return n.style.JustifyItems() }

func (n *Node) SetJustifySelf(j Justify) *Node {
	if n.style.JustifySelf() != j {
		n.style.SetJustifySelf(j)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) GetJustifySelf() Justify { return n.style.JustifySelf() }

func (n *Node) SetAlignContent(a Align) *Node {
	if n.style.AlignContent() != a {
		n.style.SetAlignContent(a)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) GetAlignContent() Align { return n.style.AlignContent() }

func (n *Node) SetAlignItems(a Align) *Node {
	if n.style.AlignItems() != a {
		n.style.SetAlignItems(a)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) GetAlignItems() Align { return n.style.AlignItems() }

func (n *Node) SetAlignSelf(a Align) *Node {
	if n.style.AlignSelf() != a {
		n.style.SetAlignSelf(a)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) GetAlignSelf() Align { return n.style.AlignSelf() }

func (n *Node) SetPositionType(p PositionType) *Node {
	if n.style.PositionType() != p {
		n.style.SetPositionType(p)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) GetPositionType() PositionType { return n.style.PositionType() }

func (n *Node) SetFlexWrap(w Wrap) *Node {
	if n.style.FlexWrap() != w {
		n.style.SetFlexWrap(w)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) GetFlexWrap() Wrap { return n.style.FlexWrap() }

func (n *Node) SetOverflow(o Overflow) *Node {
	if n.style.Overflow() != o {
		n.style.SetOverflow(o)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) GetOverflow() Overflow { return n.style.Overflow() }

func (n *Node) SetDisplay(d Display) *Node {
	if n.style.Display() != d {
		n.style.SetDisplay(d)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) GetDisplay() Display { return n.style.Display() }

func (n *Node) SetFlex(f float32) *Node {
	fo := FloatOptional{}
	if !isUndefined(f) {
		fo = NewFloatOptional(f)
	}
	if !n.style.Flex().Equals(fo) {
		n.style.SetFlex(fo)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) GetFlex() float32 {
	f := n.style.Flex()
	if f.IsUndefined() {
		return float32(math.NaN())
	}
	return f.Unwrap()
}

func (n *Node) SetFlexGrow(f float32) *Node {
	fo := NewFloatOptional(f)
	if !n.style.FlexGrow().Equals(fo) {
		n.style.SetFlexGrow(fo)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) GetFlexGrow() float32 {
	f := n.style.FlexGrow()
	if f.IsUndefined() {
		return defaultFlexGrow
	}
	return f.Unwrap()
}

func (n *Node) SetFlexShrink(f float32) *Node {
	fo := NewFloatOptional(f)
	if !n.style.FlexShrink().Equals(fo) {
		n.style.SetFlexShrink(fo)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) GetFlexShrink() float32 {
	f := n.style.FlexShrink()
	if f.IsUndefined() {
		if n.config.UseWebDefaults() {
			return webDefaultFlexShrink
		}
		return defaultFlexShrink
	}
	return f.Unwrap()
}

func (n *Node) SetFlexBasis(f float32) *Node {
	sl := StyleSizeLengthPoints(f)
	if !n.style.FlexBasis().Equals(sl) {
		n.style.SetFlexBasis(sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) SetFlexBasisPercent(f float32) *Node {
	sl := StyleSizeLengthPercent(f)
	if !n.style.FlexBasis().Equals(sl) {
		n.style.SetFlexBasis(sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) SetFlexBasisAuto() *Node {
	sl := StyleSizeLengthAuto()
	if !n.style.FlexBasis().Equals(sl) {
		n.style.SetFlexBasis(sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) SetFlexBasisMaxContent() *Node {
	sl := StyleSizeLengthMaxContent()
	if !n.style.FlexBasis().Equals(sl) {
		n.style.SetFlexBasis(sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) SetFlexBasisFitContent() *Node {
	sl := StyleSizeLengthFitContent()
	if !n.style.FlexBasis().Equals(sl) {
		n.style.SetFlexBasis(sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) SetFlexBasisStretch() *Node {
	sl := StyleSizeLengthOfStretch()
	if !n.style.FlexBasis().Equals(sl) {
		n.style.SetFlexBasis(sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) GetFlexBasis() Value {
	return n.style.FlexBasis().ToValue()
}

func (n *Node) SetEdgePosition(edge Edge, value float32) *Node {
	sl := StyleLengthPoints(value)
	if n.style.Position(edge) != sl {
		n.style.SetPosition(edge, sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) SetEdgePositionPercent(edge Edge, value float32) *Node {
	sl := StyleLengthPercent(value)
	if n.style.Position(edge) != sl {
		n.style.SetPosition(edge, sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) SetEdgePositionAuto(edge Edge) *Node {
	sl := StyleLengthAuto()
	if n.style.Position(edge) != sl {
		n.style.SetPosition(edge, sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) GetEdgePosition(edge Edge) Value {
	return n.style.Position(edge).ToValue()
}

func (n *Node) SetMargin(edge Edge, value float32) *Node {
	sl := StyleLengthPoints(value)
	if n.style.Margin(edge) != sl {
		n.style.SetMargin(edge, sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) SetMarginPercent(edge Edge, value float32) *Node {
	sl := StyleLengthPercent(value)
	if n.style.Margin(edge) != sl {
		n.style.SetMargin(edge, sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) SetMarginAuto(edge Edge) *Node {
	sl := StyleLengthAuto()
	if n.style.Margin(edge) != sl {
		n.style.SetMargin(edge, sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) GetMargin(edge Edge) Value {
	return n.style.Margin(edge).ToValue()
}

func (n *Node) SetPadding(edge Edge, value float32) *Node {
	sl := StyleLengthPoints(value)
	if n.style.Padding(edge) != sl {
		n.style.SetPadding(edge, sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) SetPaddingPercent(edge Edge, value float32) *Node {
	sl := StyleLengthPercent(value)
	if n.style.Padding(edge) != sl {
		n.style.SetPadding(edge, sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) GetPadding(edge Edge) Value {
	return n.style.Padding(edge).ToValue()
}

func (n *Node) SetBorder(edge Edge, value float32) *Node {
	sl := StyleLengthPoints(value)
	if n.style.Border(edge) != sl {
		n.style.SetBorder(edge, sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) GetBorder(edge Edge) float32 {
	b := n.style.Border(edge)
	if b.IsUndefined() || b.IsAuto() {
		return float32(math.NaN())
	}
	return b.ToValue().Value
}

func (n *Node) SetGap(gutter Gutter, value float32) *Node {
	sl := StyleLengthPoints(value)
	if n.style.Gap(gutter) != sl {
		n.style.SetGap(gutter, sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) SetGapPercent(gutter Gutter, value float32) *Node {
	sl := StyleLengthPercent(value)
	if n.style.Gap(gutter) != sl {
		n.style.SetGap(gutter, sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) GetGap(gutter Gutter) Value {
	return n.style.Gap(gutter).ToValue()
}

func (n *Node) SetBoxSizing(b BoxSizing) *Node {
	if n.style.BoxSizing() != b {
		n.style.SetBoxSizing(b)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) GetBoxSizing() BoxSizing { return n.style.BoxSizing() }

func (n *Node) SetWidth(value float32) *Node {
	sl := StyleSizeLengthPoints(value)
	if !n.style.Dimension(DimensionWidth).Equals(sl) {
		n.style.SetDimension(DimensionWidth, sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) SetWidthPercent(value float32) *Node {
	sl := StyleSizeLengthPercent(value)
	if !n.style.Dimension(DimensionWidth).Equals(sl) {
		n.style.SetDimension(DimensionWidth, sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) SetWidthAuto() *Node {
	sl := StyleSizeLengthAuto()
	if !n.style.Dimension(DimensionWidth).Equals(sl) {
		n.style.SetDimension(DimensionWidth, sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) SetWidthMaxContent() *Node {
	sl := StyleSizeLengthMaxContent()
	if !n.style.Dimension(DimensionWidth).Equals(sl) {
		n.style.SetDimension(DimensionWidth, sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) SetWidthFitContent() *Node {
	sl := StyleSizeLengthFitContent()
	if !n.style.Dimension(DimensionWidth).Equals(sl) {
		n.style.SetDimension(DimensionWidth, sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) SetWidthStretch() *Node {
	sl := StyleSizeLengthOfStretch()
	if !n.style.Dimension(DimensionWidth).Equals(sl) {
		n.style.SetDimension(DimensionWidth, sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) GetWidthValue() Value { return n.style.Dimension(DimensionWidth).ToValue() }

func (n *Node) SetHeight(value float32) *Node {
	sl := StyleSizeLengthPoints(value)
	if !n.style.Dimension(DimensionHeight).Equals(sl) {
		n.style.SetDimension(DimensionHeight, sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) SetHeightPercent(value float32) *Node {
	sl := StyleSizeLengthPercent(value)
	if !n.style.Dimension(DimensionHeight).Equals(sl) {
		n.style.SetDimension(DimensionHeight, sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) SetHeightAuto() *Node {
	sl := StyleSizeLengthAuto()
	if !n.style.Dimension(DimensionHeight).Equals(sl) {
		n.style.SetDimension(DimensionHeight, sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) SetHeightMaxContent() *Node {
	sl := StyleSizeLengthMaxContent()
	if !n.style.Dimension(DimensionHeight).Equals(sl) {
		n.style.SetDimension(DimensionHeight, sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) SetHeightFitContent() *Node {
	sl := StyleSizeLengthFitContent()
	if !n.style.Dimension(DimensionHeight).Equals(sl) {
		n.style.SetDimension(DimensionHeight, sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) SetHeightStretch() *Node {
	sl := StyleSizeLengthOfStretch()
	if !n.style.Dimension(DimensionHeight).Equals(sl) {
		n.style.SetDimension(DimensionHeight, sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) GetHeightValue() Value { return n.style.Dimension(DimensionHeight).ToValue() }

func (n *Node) SetMinWidth(value float32) *Node {
	sl := StyleSizeLengthPoints(value)
	if !n.style.MinDimension(DimensionWidth).Equals(sl) {
		n.style.SetMinDimension(DimensionWidth, sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) SetMinWidthPercent(value float32) *Node {
	sl := StyleSizeLengthPercent(value)
	if !n.style.MinDimension(DimensionWidth).Equals(sl) {
		n.style.SetMinDimension(DimensionWidth, sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) SetMinWidthMaxContent() *Node {
	sl := StyleSizeLengthMaxContent()
	if !n.style.MinDimension(DimensionWidth).Equals(sl) {
		n.style.SetMinDimension(DimensionWidth, sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) SetMinWidthFitContent() *Node {
	sl := StyleSizeLengthFitContent()
	if !n.style.MinDimension(DimensionWidth).Equals(sl) {
		n.style.SetMinDimension(DimensionWidth, sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) SetMinWidthStretch() *Node {
	sl := StyleSizeLengthOfStretch()
	if !n.style.MinDimension(DimensionWidth).Equals(sl) {
		n.style.SetMinDimension(DimensionWidth, sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) GetMinWidth() Value { return n.style.MinDimension(DimensionWidth).ToValue() }

func (n *Node) SetMinHeight(value float32) *Node {
	sl := StyleSizeLengthPoints(value)
	if !n.style.MinDimension(DimensionHeight).Equals(sl) {
		n.style.SetMinDimension(DimensionHeight, sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) SetMinHeightPercent(value float32) *Node {
	sl := StyleSizeLengthPercent(value)
	if !n.style.MinDimension(DimensionHeight).Equals(sl) {
		n.style.SetMinDimension(DimensionHeight, sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) SetMinHeightMaxContent() *Node {
	sl := StyleSizeLengthMaxContent()
	if !n.style.MinDimension(DimensionHeight).Equals(sl) {
		n.style.SetMinDimension(DimensionHeight, sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) SetMinHeightFitContent() *Node {
	sl := StyleSizeLengthFitContent()
	if !n.style.MinDimension(DimensionHeight).Equals(sl) {
		n.style.SetMinDimension(DimensionHeight, sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) SetMinHeightStretch() *Node {
	sl := StyleSizeLengthOfStretch()
	if !n.style.MinDimension(DimensionHeight).Equals(sl) {
		n.style.SetMinDimension(DimensionHeight, sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) GetMinHeight() Value { return n.style.MinDimension(DimensionHeight).ToValue() }

func (n *Node) SetMaxWidth(value float32) *Node {
	sl := StyleSizeLengthPoints(value)
	if !n.style.MaxDimension(DimensionWidth).Equals(sl) {
		n.style.SetMaxDimension(DimensionWidth, sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) SetMaxWidthPercent(value float32) *Node {
	sl := StyleSizeLengthPercent(value)
	if !n.style.MaxDimension(DimensionWidth).Equals(sl) {
		n.style.SetMaxDimension(DimensionWidth, sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) SetMaxWidthMaxContent() *Node {
	sl := StyleSizeLengthMaxContent()
	if !n.style.MaxDimension(DimensionWidth).Equals(sl) {
		n.style.SetMaxDimension(DimensionWidth, sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) SetMaxWidthFitContent() *Node {
	sl := StyleSizeLengthFitContent()
	if !n.style.MaxDimension(DimensionWidth).Equals(sl) {
		n.style.SetMaxDimension(DimensionWidth, sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) SetMaxWidthStretch() *Node {
	sl := StyleSizeLengthOfStretch()
	if !n.style.MaxDimension(DimensionWidth).Equals(sl) {
		n.style.SetMaxDimension(DimensionWidth, sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) GetMaxWidth() Value { return n.style.MaxDimension(DimensionWidth).ToValue() }

func (n *Node) SetMaxHeight(value float32) *Node {
	sl := StyleSizeLengthPoints(value)
	if !n.style.MaxDimension(DimensionHeight).Equals(sl) {
		n.style.SetMaxDimension(DimensionHeight, sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) SetMaxHeightPercent(value float32) *Node {
	sl := StyleSizeLengthPercent(value)
	if !n.style.MaxDimension(DimensionHeight).Equals(sl) {
		n.style.SetMaxDimension(DimensionHeight, sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) SetMaxHeightMaxContent() *Node {
	sl := StyleSizeLengthMaxContent()
	if !n.style.MaxDimension(DimensionHeight).Equals(sl) {
		n.style.SetMaxDimension(DimensionHeight, sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) SetMaxHeightFitContent() *Node {
	sl := StyleSizeLengthFitContent()
	if !n.style.MaxDimension(DimensionHeight).Equals(sl) {
		n.style.SetMaxDimension(DimensionHeight, sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) SetMaxHeightStretch() *Node {
	sl := StyleSizeLengthOfStretch()
	if !n.style.MaxDimension(DimensionHeight).Equals(sl) {
		n.style.SetMaxDimension(DimensionHeight, sl)
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) GetMaxHeight() Value { return n.style.MaxDimension(DimensionHeight).ToValue() }

func (n *Node) SetAspectRatio(value float32) *Node {
	fo := FloatOptional{}
	if !isUndefined(value) {
		fo = NewFloatOptional(value)
	}
	n.style.SetAspectRatio(fo)
	n.MarkDirtyAndPropagate()
	return n
}
func (n *Node) GetAspectRatio() float32 {
	f := n.style.AspectRatio()
	if f.IsUndefined() {
		return float32(math.NaN())
	}
	return f.Unwrap()
}
