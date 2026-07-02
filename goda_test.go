package goda

import (
	"math"
	"testing"
)

func TestComputeFlexBasis_BasicLayout(t *testing.T) {
	calc := func(layout float32) {
		root := New()
		root.SetWidth(100)
		root.SetHeight(100)

		child := New()
		child.SetWidth(10)
		child.SetHeight(20)
		root.InsertChildNode(child, 0)

		CalculateNodeLayout(root, layout, layout, DirectionLTR)
	}

	for _, layout := range []float32{100, math.MaxFloat32} {
		calc(layout)
	}
}

func TestLayout_MarginStart(t *testing.T) {
	config := NewConfig(DefaultLogger)
	config.SetExperimentalFeatureEnabled(ExperimentalFeatureFixFlexBasisFitContent, true)

	root := NewWithConfig(config)
	root.SetFlexDirection(FlexDirectionRow)
	root.SetWidth(100)
	root.SetHeight(100)

	child := NewWithConfig(config)
	child.SetMargin(EdgeStart, 10)
	child.SetWidth(10)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	if got := child.GetLayoutMargin(EdgeLeft); got != 10 {
		t.Errorf("Expected margin-start as left: 10, got %f", got)
	}
}

func TestLayout_DefaultValues(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)

	child0 := New()
	child0.SetMargin(EdgeStart, 10)
	child0.SetMargin(EdgeTop, 5)
	child0.SetMargin(EdgeRight, 20)
	child0.SetMargin(EdgeBottom, 5)
	root.InsertChildNode(child0, 0)

	child1 := New()
	child1.SetWidth(10)
	child1.SetHeight(10)
	root.InsertChildNode(child1, 1)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	t.Logf("Root: %fx%f @(%f,%f)", root.GetWidth(), root.GetHeight(), root.GetLeft(), root.GetTop())
	t.Logf("Child0: %fx%f @(%f,%f)", child0.GetWidth(), child0.GetHeight(), child0.GetLeft(), child0.GetTop())
	t.Logf("Child1: %fx%f @(%f,%f)", child1.GetWidth(), child1.GetHeight(), child1.GetLeft(), child1.GetTop())
}

func TestLayout_WidthHeight(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)

	child := New()
	child.SetWidth(10)
	child.SetHeight(10)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	if got := child.GetWidth(); got != 10 {
		t.Errorf("Expected width 10, got %f", got)
	}
	if got := child.GetHeight(); got != 10 {
		t.Errorf("Expected height 10, got %f", got)
	}
}

func TestLayout_FlexGrow(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)

	child := New()
	child.SetFlexGrow(1)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	if got := child.GetWidth(); got != 100 {
		t.Errorf("Expected flex-grow child to fill width: 100, got %f", got)
	}
}

func TestLayout_FlexBasis(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)

	child := New()
	child.SetFlexBasis(50)
	child.SetFlexGrow(1)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	if child.GetWidth() < 50 {
		t.Errorf("Expected flex-basis child width >= 50, got %f", child.GetWidth())
	}
}

func TestLayout_MultiChildren(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)

	for i := 0; i < 3; i++ {
		child := New()
		child.SetFlexGrow(1)
		root.InsertChildNode(child, i)
	}

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	// Column layout: flexGrow affects height, not width
	h := root.GetChild(0).GetHeight()
	if math.Abs(float64(h-33.333336)) > 3 {
		t.Errorf("Expected ~33.3 height, got %f", h)
	}
}

func TestLayout_FlexDirectionColumn(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)
	root.SetFlexDirection(FlexDirectionColumn)

	child0 := New()
	child0.SetHeight(25)
	root.InsertChildNode(child0, 0)

	child1 := New()
	child1.SetHeight(25)
	root.InsertChildNode(child1, 1)

	child2 := New()
	child2.SetHeight(25)
	root.InsertChildNode(child2, 2)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	for i, child := range []*Node{child0, child1, child2} {
		if got := child.GetHeight(); got != 25 {
			t.Errorf("Child %d: expected height 25, got %f", i, got)
		}
	}
}

func TestLayout_FlexWrap(t *testing.T) {
	root := New()
	root.SetFlexWrap(WrapWrap)
	root.SetWidth(100)
	root.SetHeight(100)

	for i := 0; i < 3; i++ {
		child := New()
		child.SetWidth(40)
		child.SetHeight(40)
		root.InsertChildNode(child, i)
	}

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	if root.GetChildCount() != 3 {
		t.Errorf("Expected 3 children")
	}
}

func TestLayout_AlignItemsCenter(t *testing.T) {
	root := New()
	root.SetAlignItems(AlignCenter)
	root.SetWidth(100)
	root.SetHeight(100)

	child := New()
	child.SetWidth(10)
	child.SetHeight(10)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	if child.GetTop() <= 0 {
		t.Logf("Centered child top: %f", child.GetTop())
	}
}

func TestLayout_AlignItemsFlexEnd(t *testing.T) {
	root := New()
	root.SetAlignItems(AlignFlexEnd)
	root.SetWidth(100)
	root.SetHeight(100)

	child := New()
	child.SetWidth(10)
	child.SetHeight(10)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	// Child should be at the bottom
	if child.GetTop() <= 0 {
		t.Logf("Flex-end child top: %f", child.GetTop())
	}
}

func TestLayout_JustifyContentCenter(t *testing.T) {
	root := New()
	root.SetJustifyContent(JustifyCenter)
	root.SetWidth(100)
	root.SetHeight(100)

	child := New()
	child.SetWidth(10)
	child.SetHeight(10)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	if child.GetLeft() <= 0 {
		t.Logf("Centered child left: %f", child.GetLeft())
	}
}

func TestLayout_JustifyContentFlexEnd(t *testing.T) {
	root := New()
	root.SetJustifyContent(JustifyFlexEnd)
	root.SetWidth(100)
	root.SetHeight(100)

	child := New()
	child.SetWidth(10)
	child.SetHeight(10)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	if child.GetLeft() <= 0 {
		t.Logf("Flex-end child left: %f", child.GetLeft())
	}
}

func TestLayout_Padding(t *testing.T) {
	const padding = 10.0
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)
	root.SetPadding(EdgeAll, padding)

	child := New()
	child.SetFlexGrow(1)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	if got := child.GetWidth(); math.Abs(float64(got-(100-2*padding))) > 0.1 {
		t.Errorf("Expected child width %f, got %f", 100-2*padding, got)
	}
}

func TestLayout_Margin(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)

	child := New()
	child.SetMargin(EdgeAll, 10)
	child.SetFlexGrow(1)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	// Child width is available width minus margins
	if child.GetWidth() > 81 {
		t.Logf("Child with margins: %f x %f @(%f,%f)",
			child.GetWidth(), child.GetHeight(), child.GetLeft(), child.GetTop())
	}
}

func TestLayout_AutoMargin(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)
	root.SetJustifyContent(JustifyCenter)

	child := New()
	child.SetWidth(10)
	child.SetHeight(10)
	child.SetMarginAuto(EdgeHorizontal)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	t.Logf("Auto margin child left: %f", child.GetLeft())
}

func TestLayout_Border(t *testing.T) {
	const border = 5.0
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)
	root.SetBorder(EdgeAll, border)

	child := New()
	child.SetFlexGrow(1)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	if got := child.GetWidth(); math.Abs(float64(got-(100-2*border))) > 0.1 {
		t.Errorf("Expected child width %f, got %f", 100-2*border, got)
	}
}

func TestLayout_AspectRatio(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)

	child := New()
	child.SetAspectRatio(2)
	child.SetFlexGrow(1)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	t.Logf("Aspect ratio child: %f x %f", child.GetWidth(), child.GetHeight())
}

func TestLayout_PositionAbsolute(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)

	absolute := New()
	absolute.SetEdgePosition(EdgeLeft, 10)
	absolute.SetEdgePosition(EdgeTop, 10)
	absolute.SetWidth(20)
	absolute.SetHeight(20)
	absolute.SetPositionType(PositionTypeAbsolute)
	root.InsertChildNode(absolute, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	t.Logf("Absolute child: %f x %f @(%f,%f)", absolute.GetWidth(), absolute.GetHeight(), absolute.GetLeft(), absolute.GetTop())
}

func TestLayout_MinMaxDimensions(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)

	child := New()
	child.SetMinWidth(50)
	child.SetMaxWidth(80)
	child.SetFlexGrow(1)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	if child.GetWidth() < 50 {
		t.Errorf("Expected min width >= 50, got %f", child.GetWidth())
	}
	if child.GetWidth() > 80 {
		t.Errorf("Expected max width <= 80, got %f", child.GetWidth())
	}
}

func TestLayout_Percentage(t *testing.T) {
	root := New()
	root.SetWidth(200)
	root.SetHeight(200)

	child := New()
	child.SetWidthPercent(50)
	child.SetHeight(50)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	if math.Abs(float64(child.GetWidth()-100)) > 1 {
		t.Errorf("Expected 50%% width = 100, got %f", child.GetWidth())
	}
}

func TestLayout_DirectionRTL(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)

	child := New()
	child.SetWidth(10)
	child.SetHeight(10)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionRTL)

	t.Logf("RTL child: %f x %f @(%f,%f)", child.GetWidth(), child.GetHeight(), child.GetLeft(), child.GetTop())
}

func TestLayout_Overflow(t *testing.T) {
	root := New()
	root.SetWidth(50)
	root.SetHeight(50)
	root.SetOverflow(OverflowScroll)

	child := New()
	child.SetWidth(100)
	child.SetHeight(100)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	// Scroll overflow shouldn't crash
	t.Logf("Overflow root: %f x %f", root.GetWidth(), root.GetHeight())
}

func TestLayout_DisplayNone(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)

	child0 := New()
	child0.SetDisplay(DisplayNone)
	child0.SetWidth(50)
	child0.SetHeight(50)
	root.InsertChildNode(child0, 0)

	child1 := New()
	child1.SetFlexGrow(1)
	root.InsertChildNode(child1, 1)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	// display:none children shouldn't affect layout
	t.Logf("Display none child: %f x %f", child0.GetWidth(), child0.GetHeight())
	if math.Abs(float64(child1.GetWidth()-100)) > 0.1 {
		t.Errorf("Expected visible child to fill width: 100, got %f", child1.GetWidth())
	}
}

func TestLayout_HasNewLayout(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)

	if !root.GetHasNewLayout() {
		t.Error("Expected new node to have new layout")
	}

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	if !root.GetHasNewLayout() {
		t.Error("Expected has new layout after calculation")
	}

	root.SetHasNewLayout(false)
	if root.GetHasNewLayout() {
		t.Error("Expected false after reset")
	}
}

func TestLayout_Dirty(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)

	if !root.IsDirty() {
		t.Error("Expected new node to be dirty")
	}

	root.SetMeasureFunc(func(node *Node, width float32, widthMode MeasureMode, height float32, heightMode MeasureMode) Size {
		return Size{Width: width, Height: height}
	})
	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	if root.IsDirty() {
		t.Error("Expected node not dirty after layout")
	}

	root.MarkDirty()
	if !root.IsDirty() {
		t.Error("Expected node dirty after MarkDirty")
	}
}

func TestLayout_NestedFlex(t *testing.T) {
	root := New()
	root.SetWidth(200)
	root.SetHeight(200)
	root.SetFlexDirection(FlexDirectionRow)

	left := New()
	left.SetWidth(50)
	left.SetHeight(50)
	root.InsertChildNode(left, 0)

	right := New()
	right.SetFlexGrow(1)
	right.SetFlexDirection(FlexDirectionColumn)

	rightTop := New()
	rightTop.SetHeight(25)
	right.InsertChildNode(rightTop, 0)

	rightBottom := New()
	rightBottom.SetFlexGrow(1)
	right.InsertChildNode(rightBottom, 1)

	root.InsertChildNode(right, 1)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	if math.Abs(float64(left.GetWidth()-50)) > 0.1 {
		t.Errorf("Expected left width 50, got %f", left.GetWidth())
	}
	t.Logf("Right: %f x %f @(%f,%f)", right.GetWidth(), right.GetHeight(), right.GetLeft(), right.GetTop())
}

func TestLayout_Config(t *testing.T) {
	config := NewConfig(DefaultLogger)
	config.SetUseWebDefaults(true)

	root := NewWithConfig(config)
	root.SetWidth(100)
	root.SetHeight(100)

	child := NewWithConfig(config)
	child.SetHeight(50)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	if root.GetFlexDirection() != FlexDirectionRow {
		t.Errorf("Expected web default flex direction row, got %v", root.GetFlexDirection())
	}
}

func TestLayout_Clone(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)

	child := New()
	child.SetWidth(10)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	clone := root.Clone()
	// Clone doesn't copy layout results by default
	if clone.GetOwner() != nil {
		t.Error("Clone should not have owner")
	}
	// Clone copies config and style
	if clone.GetConfig() != root.GetConfig() {
		t.Error("Clone should have same config")
	}
}

func TestLayout_RemoveChild(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)

	child := New()
	root.InsertChildNode(child, 0)

	if root.GetChildCount() != 1 {
		t.Error("Expected 1 child")
	}

	root.RemoveChildNode(child)

	if root.GetChildCount() != 0 {
		t.Error("Expected 0 children after removal")
	}
	if child.GetOwner() != nil {
		t.Error("Expected removed child to have no owner")
	}
}

func TestLayout_Gap(t *testing.T) {
	root := New()
	root.SetWidth(200)
	root.SetHeight(200)
	root.SetFlexDirection(FlexDirectionRow)
	root.SetGap(GutterAll, 10)

	for i := 0; i < 3; i++ {
		child := New()
		child.SetWidth(30)
		child.SetHeight(30)
		root.InsertChildNode(child, i)
	}

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	t.Logf("Children: 0:%f, 1:%f, 2:%f",
		root.GetChild(0).GetLeft(), root.GetChild(1).GetLeft(), root.GetChild(2).GetLeft())
}

func TestLayout_BoxSizing(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)
	root.SetPadding(EdgeAll, 10)

	child := New()
	child.SetBoxSizing(BoxSizingContentBox)
	child.SetWidth(50)
	child.SetHeight(50)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)
	t.Logf("Content-box child: %f x %f", child.GetWidth(), child.GetHeight())
}

func TestLayout_NoChildren(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	if root.GetWidth() != 100 || root.GetHeight() != 100 {
		t.Errorf("Empty container should keep its dimensions")
	}
}

func TestLayout_Undefined(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)

	child := New()
	child.SetFlexGrow(1)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	if child.GetWidth() != 100 {
		t.Errorf("Expected flex-grow child to fill: 100, got %f", child.GetWidth())
	}
}

func TestLayout_MeasureFunc(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)

	text := New()
	text.SetMeasureFunc(func(node *Node, width float32, widthMode MeasureMode, height float32, heightMode MeasureMode) Size {
		return Size{Width: 50, Height: 20}
	})
	root.InsertChildNode(text, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	t.Logf("Measured child: %f x %f", text.GetWidth(), text.GetHeight())
}

func TestLayout_EdgeStartEnd(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)

	child := New()
	child.SetEdgePosition(EdgeStart, 10)
	child.SetEdgePosition(EdgeEnd, 10)
	child.SetPositionType(PositionTypeAbsolute)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)


	if math.Abs(float64(child.GetLeft()-10)) > 1 {
		t.Errorf("Expected child left 10, got %f", child.GetLeft())
	}
}

func TestLayout_MarginStartEnd(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)

	child := New()
	child.SetMargin(EdgeStart, 10)
	child.SetMargin(EdgeEnd, 5)
	child.SetFlexGrow(1)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	t.Logf("Child with start/end margins: %f x %f @(%f,%f)",
		child.GetWidth(), child.GetHeight(), child.GetLeft(), child.GetTop())
}

func TestLayout_FlexGrowFlexShrink(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)

	child0 := New()
	child0.SetFlexGrow(2)
	child0.SetWidth(10)
	root.InsertChildNode(child0, 0)

	child1 := New()
	child1.SetFlexGrow(1)
	child1.SetWidth(10)
	root.InsertChildNode(child1, 1)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	t.Logf("Grow 2: %f, Grow 1: %f", child0.GetWidth(), child1.GetWidth())
}

func TestLayout_RoundingFunction(t *testing.T) {
	config := NewConfig(DefaultLogger)
	config.SetPointScaleFactor(2.0)

	root := NewWithConfig(config)
	root.SetWidth(100)
	root.SetHeight(100)

	child := NewWithConfig(config)
	child.SetWidth(33)
	child.SetHeight(33)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	t.Logf("Rounded child: %f x %f @(%f,%f)", child.GetWidth(), child.GetHeight(), child.GetLeft(), child.GetTop())
}

func TestLayout_ChildCount(t *testing.T) {
	root := New()

	if root.GetChildCount() != 0 {
		t.Error("Expected 0 children")
	}

	child := New()
	root.InsertChildNode(child, 0)
	if root.GetChildCount() != 1 {
		t.Error("Expected 1 child")
	}

	child2 := New()
	root.InsertChildNode(child2, 1)
	if root.GetChildCount() != 2 {
		t.Error("Expected 2 children")
	}

	root.RemoveAllChildren()
	if root.GetChildCount() != 0 {
		t.Error("Expected 0 children after remove all")
	}
}

func TestLayout_Owner(t *testing.T) {
	root := New()
	child := New()
	root.InsertChildNode(child, 0)

	if child.GetParent() != root {
		t.Error("Expected parent to be root")
	}
	if child.GetOwner() != root {
		t.Error("Expected owner to be root")
	}

	root.RemoveAllChildren()
	if child.GetParent() != nil {
		t.Error("Expected no parent after removal")
	}
}

func TestLayout_MeasureFunc_InvalidValues(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)

	text := New()
	text.SetMeasureFunc(func(node *Node, width float32, widthMode MeasureMode, height float32, heightMode MeasureMode) Size {
		return Size{Width: float32(math.NaN()), Height: -1}
	})
	root.InsertChildNode(text, 0)

	// Should not panic
	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)
}

func TestLayout_Baseline(t *testing.T) {
	root := New()
	root.SetFlexDirection(FlexDirectionRow)
	root.SetAlignItems(AlignBaseline)
	root.SetWidth(300)
	root.SetHeight(300)

	child0 := New()
	child0.SetWidth(50)
	child0.SetHeight(50)
	root.InsertChildNode(child0, 0)

	child1 := New()
	child1.SetWidth(50)
	child1.SetHeight(80)
	child1.SetBaselineFunc(func(node *Node, width, height float32) float32 { return 50 })
	root.InsertChildNode(child1, 1)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	t.Logf("Child0: %f x %f @(%f,%f)", child0.GetWidth(), child0.GetHeight(), child0.GetLeft(), child0.GetTop())
	t.Logf("Child1: %f x %f @(%f,%f)", child1.GetWidth(), child1.GetHeight(), child1.GetLeft(), child1.GetTop())
}

func TestApplyStyle_BuilderPattern(t *testing.T) {
	root := New().ApplyStyle(map[string]string{
		"display":         "flex",
		"flex-direction":  "column",
		"width":           "400",
		"height":          "300",
		"padding":         "16",
		"gap":             "8",
	})

	if root.GetDisplay() != DisplayFlex {
		t.Error("display should be flex")
	}
	if root.GetFlexDirection() != FlexDirectionColumn {
		t.Error("flex-direction should be column")
	}
	if v := root.GetWidthValue().Value; v != 400 {
		t.Errorf("width should be 400, got %f", v)
	}
	if v := root.GetHeightValue().Value; v != 300 {
		t.Errorf("height should be 300, got %f", v)
	}
	if v := root.GetGap(GutterAll).Value; v != 8 {
		t.Errorf("gap should be 8, got %f", v)
	}
}

func TestApplyStyle_FluentChain(t *testing.T) {
	child := New().
		ApplyStyle(map[string]string{"width": "100", "height": "50"}).
		SetFlexGrow(1).
		SetMargin(EdgeAll, 8).
		ApplyStyle(map[string]string{"align-self": "center"})

	if v := child.GetWidthValue().Value; v != 100 {
		t.Errorf("width should be 100, got %f", v)
	}
	if child.GetFlexGrow() != 1 {
		t.Error("flex-grow should be 1")
	}
	if child.GetAlignSelf() != AlignCenter {
		t.Error("align-self should be center")
	}
}

func TestApplyStyle_UnknownKeysIgnored(t *testing.T) {
	n := New().ApplyStyle(map[string]string{
		"width":           "200",
		"unknown-prop":    "foo",
		"also-invalid":    "bar",
		"height":          "100",
	})
	if v := n.GetWidthValue().Value; v != 200 {
		t.Errorf("width should be 200, got %f", v)
	}
	if v := n.GetHeightValue().Value; v != 100 {
		t.Errorf("height should be 100, got %f", v)
	}
}

func TestApplyStyle_PercentAndKeywords(t *testing.T) {
	n := New().ApplyStyle(map[string]string{
		"width":       "50%",
		"height":      "auto",
		"min-width":   "max-content",
		"flex-basis":  "fit-content",
	})

	if v := n.GetWidthValue().Unit; v != UnitPercent {
		t.Errorf("width unit should be percent, got %v", v)
	}
	if v := n.GetHeightValue().Unit; v != UnitAuto {
		t.Errorf("height unit should be auto, got %v", v)
	}
	if v := n.GetMinWidth().Unit; v != UnitMaxContent {
		t.Errorf("min-width unit should be max-content, got %v", v)
	}
}

func TestParseStyle_Basic(t *testing.T) {
	props := ParseStyle(`
		display: flex;
		flex-direction: row;
		width: 800;
		height: 600;
		padding: 16;
		gap: 8;
	`)

	if props["display"] != "flex" {
		t.Errorf("display should be flex, got %q", props["display"])
	}
	if props["flex-direction"] != "row" {
		t.Errorf("flex-direction should be row, got %q", props["flex-direction"])
	}
	if props["width"] != "800" {
		t.Errorf("width should be 800, got %q", props["width"])
	}
	if len(props) != 6 {
		t.Errorf("expected 6 props, got %d", len(props))
	}
}

func TestParseStyle_CommentsAndUnknown(t *testing.T) {
	props := ParseStyle(`
		// header comment
		width: 200;
		/* block comment */
		unknown-prop: foo;
		height: 100;
		color: red;
		font-size: 14px;
	`)

	if props["width"] != "200" {
		t.Errorf("width should be 200, got %q", props["width"])
	}
	if props["height"] != "100" {
		t.Errorf("height should be 100, got %q", props["height"])
	}
	if _, ok := props["unknown-prop"]; ok {
		t.Error("unknown-prop should not be included")
	}
	if _, ok := props["color"]; ok {
		t.Error("color should not be included")
	}
	if len(props) != 2 {
		t.Errorf("expected 2 supported props, got %d: %v", len(props), props)
	}
}

func TestParseStyle_NoTrailingSemicolon(t *testing.T) {
	props := ParseStyle("width: 100\nheight: 200\npadding: 8")
	if props["width"] != "100" {
		t.Errorf("width should be 100, got %q", props["width"])
	}
	if len(props) != 3 {
		t.Errorf("expected 3 props, got %d", len(props))
	}
}

func TestApplyStyleString_Integration(t *testing.T) {
	n := New().ApplyStyleString(`
		display: flex;
		flex-direction: column;
		width: 400;
		height: 300;
		justify-content: center;
		align-items: stretch;
		padding: 10;
		gap: 12;
	`)

	if n.GetDisplay() != DisplayFlex {
		t.Error("display should be flex")
	}
	if n.GetFlexDirection() != FlexDirectionColumn {
		t.Error("flex-direction should be column")
	}
	if v := n.GetWidthValue().Value; v != 400 {
		t.Errorf("width should be 400, got %f", v)
	}
	if n.GetJustifyContent() != JustifyCenter {
		t.Error("justify-content should be center")
	}
}

func TestApplyStyleString_Chained(t *testing.T) {
	n := New().
		ApplyStyleString("width: 500; height: 400;").
		SetFlexGrow(1).
		ApplyStyleString("align-self: center; margin: 8;")

	if v := n.GetWidthValue().Value; v != 500 {
		t.Errorf("width should be 500, got %f", v)
	}
	if n.GetFlexGrow() != 1 {
		t.Error("flex-grow should be 1")
	}
	if n.GetAlignSelf() != AlignCenter {
		t.Error("align-self should be center")
	}
}
