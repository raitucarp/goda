package goda

import (
	"strconv"
	"strings"
)

// Node is the fundamental unit of the layout tree.
// Each Node has a Style and computed LayoutResults.
type Node struct {
	hasNewLayout               bool
	isReferenceBaseline        bool
	isDirty                    bool
	alwaysFormsContainingBlock bool
	nodeType                   NodeType
	context                    interface{}
	measureFunc                MeasureFunc
	minContentMeasureFunc      MeasureFunc
	minContentWidth            FloatOptional
	minContentHeight           FloatOptional
	baselineFunc               BaselineFunc
	dirtiedFunc                DirtiedFunc
	style                      Style
	layout                     LayoutResults
	lineIndex                  int
	contentsChildrenCount      int
	owner                      *Node
	children                   []*Node
	config                     *Config
	processedDimensions        [2]StyleSizeLength
	fontSizeEstimate           float32
	id                         string
	classes                    []string
}

func NewNode() *Node {
	return NewNodeWithConfig(GetDefaultConfig())
}

func NewNodeWithConfig(config *Config) *Node {
	if config == nil {
		panic("Tried to construct Node with null config")
	}
	n := &Node{
		hasNewLayout:     true,
		isDirty:          true,
		nodeType:         NodeTypeDefault,
		lineIndex:        0,
		config:           config,
		style:            NewStyle(),
		layout:           NewLayoutResults(),
		fontSizeEstimate: 16,
	}
	n.processedDimensions[0] = StyleSizeLengthUndefined()
	n.processedDimensions[1] = StyleSizeLengthUndefined()
	if config.UseWebDefaults() {
		n.useWebDefaults()
	}
	return n
}

func (n *Node) Clone() *Node {
	clone := *n
	clone.layout = NewLayoutResults()
	clone.layout.configVersion = n.layout.configVersion
	clone.layout.computedFlexBasis = n.layout.computedFlexBasis
	clone.layout.computedFlexBasisGeneration = n.layout.computedFlexBasisGeneration
	clone.layout.computedAutoMinMainSize = n.layout.computedAutoMinMainSize
	clone.layout.generationCount = n.layout.generationCount
	clone.layout.lastOwnerDirection = n.layout.lastOwnerDirection
	clone.layout.direction_ = n.layout.direction_
	clone.layout.dimensions_ = n.layout.dimensions_
	clone.layout.measuredDimensions_ = n.layout.measuredDimensions_
	clone.layout.rawDimensions_ = n.layout.rawDimensions_
	clone.layout.position_ = n.layout.position_
	clone.layout.margin_ = n.layout.margin_
	clone.layout.border_ = n.layout.border_
	clone.layout.padding_ = n.layout.padding_
	clone.owner = nil
	clone.children = nil
	if n.classes != nil {
		clone.classes = make([]string, len(n.classes))
		copy(clone.classes, n.classes)
	}
	return &clone
}

func (n *Node) useWebDefaults() {
	n.style.SetFlexDirection(FlexDirectionRow)
	n.style.SetAlignContent(AlignStretch)
}

func (n *Node) GetContext() interface{}                     { return n.context }
func (n *Node) SetContext(ctx interface{})                   { n.context = ctx }
func (n *Node) GetID() string                                 { return n.id }
func (n *Node) SetID(id string)                               { n.id = id }
func (n *Node) GetClasses() []string {
	if n.classes == nil {
		return nil
	}
	out := make([]string, len(n.classes))
	copy(out, n.classes)
	return out
}
func (n *Node) SetClasses(classes []string) {
	n.classes = nil
	if classes != nil {
		n.classes = make([]string, len(classes))
		copy(n.classes, classes)
	}
}
func (n *Node) AddClass(class string) {
	for _, c := range n.classes {
		if c == class {
			return
		}
	}
	n.classes = append(n.classes, class)
}
func (n *Node) HasClass(class string) bool {
	for _, c := range n.classes {
		if c == class {
			return true
		}
	}
	return false
}
func (n *Node) AlwaysFormsContainingBlock() bool             { return n.alwaysFormsContainingBlock }
func (n *Node) SetAlwaysFormsContainingBlock(v bool)         { n.alwaysFormsContainingBlock = v }
func (n *Node) GetHasNewLayout() bool                        { return n.hasNewLayout }
func (n *Node) SetHasNewLayout(v bool)                       { n.hasNewLayout = v }
func (n *Node) GetNodeType() NodeType                        { return n.nodeType }
func (n *Node) SetNodeType(v NodeType)                       { n.nodeType = v }
func (n *Node) HasMeasureFunc() bool                         { return n.measureFunc != nil }
func (n *Node) HasMinContentMeasureFunc() bool               { return n.minContentMeasureFunc != nil }
func (n *Node) HasBaselineFunc() bool                        { return n.baselineFunc != nil }
func (n *Node) GetMinContentWidth() FloatOptional            { return n.minContentWidth }
func (n *Node) SetMinContentWidth(v FloatOptional)           { n.minContentWidth = v }
func (n *Node) GetMinContentHeight() FloatOptional           { return n.minContentHeight }
func (n *Node) SetMinContentHeight(v FloatOptional)          { n.minContentHeight = v }
func (n *Node) GetDirtiedFunc() DirtiedFunc                  { return n.dirtiedFunc }
func (n *Node) SetDirtiedFunc(f DirtiedFunc)                 { n.dirtiedFunc = f }
func (n *Node) Style() *Style                                { return &n.style }
func (n *Node) GetStyle() Style                              { return n.style }
func (n *Node) SetStyle(s Style)                             { n.style = s }
func (n *Node) GetLayout() *LayoutResults                     { return &n.layout }
func (n *Node) GetLayoutVal() LayoutResults                   { return n.layout }
func (n *Node) SetLayout(l LayoutResults)                     { n.layout = l }
func (n *Node) GetLineIndex() int                             { return n.lineIndex }
func (n *Node) SetLineIndex(i int)                            { n.lineIndex = i }
func (n *Node) IsReferenceBaseline() bool                     { return n.isReferenceBaseline }
func (n *Node) SetIsReferenceBaseline(v bool)                 { n.isReferenceBaseline = v }
func (n *Node) GetOwner() *Node                               { return n.owner }
func (n *Node) SetOwner(owner *Node)                          { n.owner = owner }
func (n *Node) GetChildren() []*Node                          { return n.children }
func (n *Node) GetChild(index int) *Node                      { return n.children[index] }
func (n *Node) GetChildCount() int                            { return len(n.children) }
func (n *Node) GetConfig() *Config                            { return n.config }
func (n *Node) IsDirty() bool                                 { return n.isDirty }
func (n *Node) GetProcessedDimension(dim Dimension) StyleSizeLength {
	return n.processedDimensions[dim]
}
func (n *Node) HasContentsChildren() bool   { return n.contentsChildrenCount != 0 }
func (n *Node) HasErrata(errata Errata) bool { return n.config.HasErrata(errata) }

func (n *Node) GetResolvedDimension(dir Direction, dim Dimension, referenceLength, ownerWidth float32) FloatOptional {
	value := n.processedDimensions[dim].Resolve(referenceLength)
	if n.style.BoxSizing() == BoxSizingBorderBox {
		return value
	}
	dimPB := NewFloatOptional(n.style.ComputePaddingAndBorderForDimension(dir, dim, ownerWidth))
	result := value.Unwrap()
	if dimPB.IsDefined() {
		result += dimPB.Unwrap()
	}
	return NewFloatOptional(result)
}

func (n *Node) HasDefiniteLength(dim Dimension, ownerSize float32) bool {
	usedValue := n.processedDimensions[dim].Resolve(ownerSize)
	return usedValue.IsDefined() && usedValue.Unwrap() >= 0.0
}

func (n *Node) Measure(availableWidth float32, widthMode MeasureMode, availableHeight float32, heightMode MeasureMode) Size {
	size := n.measureFunc(n, availableWidth, widthMode, availableHeight, heightMode)
	if isUndefined(size.Height) || size.Height < 0 ||
		isUndefined(size.Width) || size.Width < 0 {
		return Size{
			Width:  maxOrDefined(0.0, size.Width),
			Height: maxOrDefined(0.0, size.Height),
		}
	}
	return size
}

func (n *Node) MeasureMinContent(availableWidth float32, widthMode MeasureMode, availableHeight float32, heightMode MeasureMode) Size {
	size := n.minContentMeasureFunc(n, availableWidth, widthMode, availableHeight, heightMode)
	if isUndefined(size.Height) || size.Height < 0 ||
		isUndefined(size.Width) || size.Width < 0 {
		return Size{
			Width:  maxOrDefined(0.0, size.Width),
			Height: maxOrDefined(0.0, size.Height),
		}
	}
	return size
}

func (n *Node) Baseline(width, height float32) float32 {
	return n.baselineFunc(n, width, height)
}

func (n *Node) DimensionWithMargin(axis FlexDirection, widthSize float32) float32 {
	return n.layout.MeasuredDimension(dimension(axis)) +
		n.style.ComputeMarginForAxis(axis, widthSize)
}

func (n *Node) IsLayoutDimensionDefined(axis FlexDirection) bool {
	value := n.layout.MeasuredDimension(dimension(axis))
	return isDefined(value) && value >= 0.0
}

func (n *Node) SetMeasureFunc(f MeasureFunc) {
	if f == nil {
		n.nodeType = NodeTypeDefault
	} else {
		if len(n.children) > 0 {
			panic("Cannot set measure function: Nodes with measure functions cannot have children.")
		}
		n.nodeType = NodeTypeText
	}
	n.measureFunc = f
}

func (n *Node) SetMinContentMeasureFunc(f MeasureFunc) { n.minContentMeasureFunc = f }
func (n *Node) SetBaselineFunc(f BaselineFunc)          { n.baselineFunc = f }

func (n *Node) SetLayoutLastOwnerDirection(dir Direction)         { n.layout.lastOwnerDirection = dir }
func (n *Node) SetLayoutComputedFlexBasis(fb FloatOptional)       { n.layout.computedFlexBasis = fb }
func (n *Node) SetLayoutComputedFlexBasisGeneration(g uint32)     { n.layout.computedFlexBasisGeneration = g }
func (n *Node) SetLayoutMeasuredDimension(v float32, dim Dimension) { n.layout.SetMeasuredDimension(dim, v) }
func (n *Node) SetLayoutHadOverflow(v bool)                       { n.layout.SetHadOverflow(v) }
func (n *Node) SetLayoutDimension(v float32, dim Dimension) {
	n.layout.SetDimension(dim, v)
	n.layout.SetRawDimension(dim, v)
}
func (n *Node) SetLayoutDirection(dir Direction)              { n.layout.SetDirection(dir) }
func (n *Node) SetLayoutMargin(v float32, edge PhysicalEdge)  { n.layout.SetMargin(edge, v) }
func (n *Node) SetLayoutBorder(v float32, edge PhysicalEdge)  { n.layout.SetBorder(edge, v) }
func (n *Node) SetLayoutPadding(v float32, edge PhysicalEdge) { n.layout.SetPadding(edge, v) }
func (n *Node) SetLayoutPosition(v float32, edge PhysicalEdge) { n.layout.SetPosition(edge, v) }

func (n *Node) SetDirty(isDirty bool) {
	if isDirty == n.isDirty {
		return
	}
	n.isDirty = isDirty
	if isDirty && n.dirtiedFunc != nil {
		n.dirtiedFunc(n)
	}
}

// Child management
func (n *Node) SetChildren(children []*Node) {
	n.children = children
	n.contentsChildrenCount = 0
	for _, child := range children {
		if child.style.Display() == DisplayContents {
			n.contentsChildrenCount++
		}
	}
}

func (n *Node) ReplaceChildAt(child *Node, index int) {
	prev := n.children[index]
	if prev.style.Display() == DisplayContents && child.style.Display() != DisplayContents {
		n.contentsChildrenCount--
	} else if prev.style.Display() != DisplayContents && child.style.Display() == DisplayContents {
		n.contentsChildrenCount++
	}
	n.children[index] = child
}

func (n *Node) ReplaceChild(oldChild, newChild *Node) {
	if oldChild.style.Display() == DisplayContents && newChild.style.Display() != DisplayContents {
		n.contentsChildrenCount--
	} else if oldChild.style.Display() != DisplayContents && newChild.style.Display() == DisplayContents {
		n.contentsChildrenCount++
	}
	for i, c := range n.children {
		if c == oldChild {
			n.children[i] = newChild
			return
		}
	}
}

func (n *Node) InsertChild(child *Node, index int) {
	if child.style.Display() == DisplayContents {
		n.contentsChildrenCount++
	}
	if index >= len(n.children) {
		n.children = append(n.children, child)
	} else {
		n.children = append(n.children[:index+1], n.children[index:]...)
		n.children[index] = child
	}
}

func (n *Node) RemoveChild(child *Node) bool {
	for i, c := range n.children {
		if c == child {
			if child.style.Display() == DisplayContents {
				n.contentsChildrenCount--
			}
			n.children = append(n.children[:i], n.children[i+1:]...)
			return true
		}
	}
	return false
}

func (n *Node) RemoveChildAt(index int) {
	if n.children[index].style.Display() == DisplayContents {
		n.contentsChildrenCount--
	}
	n.children = append(n.children[:index], n.children[index+1:]...)
}

func (n *Node) ClearChildren() {
	n.children = nil
}

func (n *Node) SetConfig(config *Config) {
	if config == nil {
		panic("Attempting to set a null config on a Node")
	}
	if config.UseWebDefaults() != n.config.UseWebDefaults() {
		panic("UseWebDefaults may not be changed after constructing a Node")
	}
	if configUpdateInvalidatesLayout(n.config, config) {
		n.MarkDirtyAndPropagate()
		n.layout.configVersion = 0
	} else {
		n.layout.configVersion = config.GetVersion()
	}
	n.config = config
}

func (n *Node) CloneChildrenIfNeeded() {
	for i, child := range n.children {
		if child.GetOwner() != n {
			clone := n.config.CloneNode(child, n, i)
			clone.SetOwner(n)
			n.children[i] = clone
			if clone.style.Display() == DisplayContents {
				clone.CloneChildrenIfNeeded()
			} else if clone.HasContentsChildren() {
				clone.CloneContentsChildrenIfNeeded()
			}
		}
	}
}

func (n *Node) CloneContentsChildrenIfNeeded() {
	for i, child := range n.children {
		if child.style.Display() == DisplayContents && child.GetOwner() != n {
			clone := n.config.CloneNode(child, n, i)
			clone.SetOwner(n)
			n.children[i] = clone
			clone.CloneChildrenIfNeeded()
		}
	}
}

func (n *Node) MarkDirtyAndPropagate() {
	if !n.isDirty {
		n.SetDirty(true)
		n.SetLayoutComputedFlexBasis(FloatOptional{})
		if n.owner != nil {
			n.owner.MarkDirtyAndPropagate()
		}
	}
}

func (n *Node) ResolveFlexGrow() float32 {
	if n.owner == nil {
		return 0.0
	}
	if n.style.FlexGrow().IsDefined() {
		return n.style.FlexGrow().Unwrap()
	}
	if n.style.Flex().IsDefined() && n.style.Flex().Unwrap() > 0.0 {
		return n.style.Flex().Unwrap()
	}
	return defaultFlexGrow
}

func (n *Node) ResolveFlexShrink() float32 {
	if n.owner == nil {
		return 0.0
	}
	if n.style.FlexShrink().IsDefined() {
		return n.style.FlexShrink().Unwrap()
	}
	if !n.config.UseWebDefaults() && n.style.Flex().IsDefined() && n.style.Flex().Unwrap() < 0.0 {
		return -n.style.Flex().Unwrap()
	}
	if n.config.UseWebDefaults() {
		return webDefaultFlexShrink
	}
	return defaultFlexShrink
}

func (n *Node) IsNodeFlexible() bool {
	return n.style.PositionType() != PositionTypeAbsolute &&
		(n.ResolveFlexGrow() != 0 || n.ResolveFlexShrink() != 0)
}

// FontSizeEstimate controls how rem and em units are resolved in CSS strings.
// Default is 16. For em, uses the node's own estimate; for rem, walks up to
// the root node's estimate.
func (n *Node) SetFontSizeEstimate(v float32) *Node { n.fontSizeEstimate = v; return n }
func (n *Node) GetFontSizeEstimate() float32         { return n.fontSizeEstimate }

// resolveCSSValue parses a CSS length string and resolves rem/em units
// against the node's font size estimate.
func (n *Node) resolveCSSValue(val string) float32 {
	val = strings.TrimSpace(val)
	if strings.HasSuffix(val, "rem") {
		root := n
		for root.owner != nil {
			root = root.owner
		}
		v := strings.TrimSpace(strings.TrimSuffix(val, "rem"))
		f, err := strconv.ParseFloat(v, 32)
		if err != nil {
			return 0
		}
		return float32(f) * root.fontSizeEstimate
	}
	if strings.HasSuffix(val, "em") {
		v := strings.TrimSpace(strings.TrimSuffix(val, "em"))
		f, err := strconv.ParseFloat(v, 32)
		if err != nil {
			return 0
		}
		return float32(f) * n.fontSizeEstimate
	}
	return parseFloat(val)
}

func (n *Node) ProcessFlexBasis() StyleSizeLength {
	flexBasis := n.style.FlexBasis()
	if !flexBasis.IsAuto() && !flexBasis.IsUndefined() {
		return flexBasis
	}
	if n.style.Flex().IsDefined() && n.style.Flex().Unwrap() > 0.0 {
		if n.config.UseWebDefaults() {
			return StyleSizeLengthAuto()
		}
		return StyleSizeLengthPoints(0)
	}
	return StyleSizeLengthAuto()
}

func (n *Node) ResolveFlexBasis(dir Direction, flexDir FlexDirection, referenceLength, ownerWidth float32) FloatOptional {
	value := n.ProcessFlexBasis().Resolve(referenceLength)
	if n.style.BoxSizing() == BoxSizingBorderBox {
		return value
	}
	dim := dimension(flexDir)
	dimPB := NewFloatOptional(n.style.ComputePaddingAndBorderForDimension(dir, dim, ownerWidth))
	result := value.Unwrap()
	if dimPB.IsDefined() {
		result += dimPB.Unwrap()
	}
	return NewFloatOptional(result)
}

func (n *Node) ProcessDimensions() {
	for _, dim := range []Dimension{DimensionWidth, DimensionHeight} {
		if n.style.MaxDimension(dim).IsDefined() &&
			n.style.MaxDimension(dim).Equals(n.style.MinDimension(dim)) {
			n.processedDimensions[dim] = n.style.MaxDimension(dim)
		} else {
			n.processedDimensions[dim] = n.style.Dimension(dim)
		}
	}
}

func (n *Node) ResolveDirection(ownerDir Direction) Direction {
	if n.style.Direction() == DirectionInherit {
		if ownerDir != DirectionInherit {
			return ownerDir
		}
		return DirectionLTR
	}
	return n.style.Direction()
}

func (n *Node) RelativePosition(axis FlexDirection, dir Direction, axisSize float32) float32 {
	if n.style.PositionType() == PositionTypeStatic {
		return 0
	}
	if n.style.IsInlineStartPositionDefined(axis, dir) &&
		!n.style.IsInlineStartPositionAuto(axis, dir) {
		return n.style.ComputeInlineStartPosition(axis, dir, axisSize)
	}
	return -1 * n.style.ComputeInlineEndPosition(axis, dir, axisSize)
}

func (n *Node) SetPosition(dir Direction, ownerWidth, ownerHeight float32) {
	directionRespectingRoot := dir
	if n.owner == nil {
		directionRespectingRoot = DirectionLTR
	}
	mainAxis := resolveDirection(n.style.FlexDirection(), directionRespectingRoot)
	crossAxis := resolveCrossDirection(mainAxis, directionRespectingRoot)

	var mainAxisSize float32
	if isRow(mainAxis) {
		mainAxisSize = ownerWidth
	} else {
		mainAxisSize = ownerHeight
	}
	var crossAxisOwnerW float32
	if isRow(mainAxis) {
		crossAxisOwnerW = ownerHeight
	} else {
		crossAxisOwnerW = ownerWidth
	}

	relativePosMain := n.RelativePosition(mainAxis, directionRespectingRoot, mainAxisSize)
	relativePosCross := n.RelativePosition(crossAxis, directionRespectingRoot, crossAxisOwnerW)

	mainLeading := inlineStartEdge(mainAxis, directionRespectingRoot)
	mainTrailing := inlineEndEdge(mainAxis, directionRespectingRoot)
	crossLeading := inlineStartEdge(crossAxis, directionRespectingRoot)
	crossTrailing := inlineEndEdge(crossAxis, directionRespectingRoot)

	n.SetLayoutPosition(n.style.ComputeInlineStartMargin(mainAxis, directionRespectingRoot, ownerWidth)+relativePosMain, mainLeading)
	n.SetLayoutPosition(n.style.ComputeInlineEndMargin(mainAxis, directionRespectingRoot, ownerWidth)+relativePosMain, mainTrailing)
	n.SetLayoutPosition(n.style.ComputeInlineStartMargin(crossAxis, directionRespectingRoot, ownerWidth)+relativePosCross, crossLeading)
	n.SetLayoutPosition(n.style.ComputeInlineEndMargin(crossAxis, directionRespectingRoot, ownerWidth)+relativePosCross, crossTrailing)
}

func (n *Node) GetLayoutChildCount() int {
	if n.contentsChildrenCount == 0 {
		return len(n.children)
	}
	count := 0
	iter := NewLayoutableIterator(n)
	for iter.Next() {
		count++
	}
	return count
}
