package goda

import "math"

// New creates a new Node with default configuration.
func New() *Node {
	return NewNode()
}

// NewWithConfig creates a new Node with the given configuration.
func NewWithConfig(config *Config) *Node {
	return NewNodeWithConfig(config)
}

// CalculateNodeLayout performs layout on the given node tree.
func CalculateNodeLayout(node *Node, availableWidth, availableHeight float32, ownerDirection Direction) {
	CalculateLayout(node, availableWidth, availableHeight, ownerDirection)
}

// Layout position getters.
func (n *Node) GetLeft() float32        { return n.layout.Position(PhysicalEdgeLeft) }
func (n *Node) GetTop() float32         { return n.layout.Position(PhysicalEdgeTop) }
func (n *Node) GetRight() float32       { return n.layout.Position(PhysicalEdgeRight) }
func (n *Node) GetBottom() float32      { return n.layout.Position(PhysicalEdgeBottom) }
func (n *Node) GetWidth() float32       { return n.layout.Dimension(DimensionWidth) }
func (n *Node) GetHeight() float32      { return n.layout.Dimension(DimensionHeight) }
func (n *Node) GetLayoutDirection() Direction { return n.layout.Direction() }
func (n *Node) GetHadOverflow() bool          { return n.layout.HadOverflow() }
func (n *Node) GetRawWidth() float32          { return n.layout.RawDimension(DimensionWidth) }
func (n *Node) GetRawHeight() float32         { return n.layout.RawDimension(DimensionHeight) }

func (n *Node) GetLayoutMargin(edge Edge) float32 {
	return getResolvedLayoutProperty(n, edge, func(l *LayoutResults, e PhysicalEdge) float32 { return l.Margin(e) })
}
func (n *Node) GetLayoutBorder(edge Edge) float32 {
	return getResolvedLayoutProperty(n, edge, func(l *LayoutResults, e PhysicalEdge) float32 { return l.Border(e) })
}
func (n *Node) GetLayoutPadding(edge Edge) float32 {
	return getResolvedLayoutProperty(n, edge, func(l *LayoutResults, e PhysicalEdge) float32 { return l.Padding(e) })
}

func getResolvedLayoutProperty(node *Node, edge Edge, getter func(*LayoutResults, PhysicalEdge) float32) float32 {
	if edge > EdgeEnd {
		panic("Cannot get layout properties of multi-edge shorthands")
	}
	layout := &node.layout
	if edge == EdgeStart {
		if layout.Direction() == DirectionRTL {
			return getter(layout, PhysicalEdgeRight)
		}
		return getter(layout, PhysicalEdgeLeft)
	}
	if edge == EdgeEnd {
		if layout.Direction() == DirectionRTL {
			return getter(layout, PhysicalEdgeLeft)
		}
		return getter(layout, PhysicalEdgeRight)
	}
	return getter(layout, PhysicalEdge(edge))
}

// RoundValueToPixelGrid rounds a value to the nearest pixel grid boundary.
func RoundValueToPixelGrid(value float64, pointScaleFactor float64, forceCeil, forceFloor bool) float32 {
	return roundValueToPixelGrid(value, pointScaleFactor, forceCeil, forceFloor)
}

// ConfigNew creates a new Config with the given logger.
func ConfigNew(logger LoggerFunc) *Config {
	return NewConfig(logger)
}

// ConfigNewDefault creates a new Config with the default no-op logger.
func ConfigNewDefault() *Config {
	return NewConfig(DefaultLogger)
}

// Convenience methods.
func (n *Node) SetNodeType_Public(nt NodeType) *Node { n.nodeType = nt; return n }
func (n *Node) GetNodeType_Public() NodeType         { return n.nodeType }

func (n *Node) SetIsReferenceBaseline_Public(v bool) *Node {
	if n.isReferenceBaseline != v {
		n.isReferenceBaseline = v
		n.MarkDirtyAndPropagate()
	}
	return n
}
func (n *Node) GetIsReferenceBaseline() bool { return n.isReferenceBaseline }
func (n *Node) GetParent() *Node             { return n.owner }

func (n *Node) SetMinContentWidthFunc(f MeasureFunc) *Node  { n.minContentMeasureFunc = f; return n }
func (n *Node) SetMinContentWidthValue(v float32) *Node      { n.minContentWidth = NewFloatOptional(v); return n }
func (n *Node) SetMinContentHeightValue(v float32) *Node     { n.minContentHeight = NewFloatOptional(v); return n }
func (n *Node) GetMinContentWidthValue() float32      { return n.minContentWidth.UnwrapOrDefault(float32(math.NaN())) }
func (n *Node) GetMinContentHeightValue() float32     { return n.minContentHeight.UnwrapOrDefault(float32(math.NaN())) }
