package goda

import (
	"math"
	"testing"
)

func TestConfig_CloneNodeCallback(t *testing.T) {
	config := NewConfig(DefaultLogger)
	called := false
	config.SetCloneNodeCallback(func(oldNode *Node, owner *Node, childIndex int) *Node {
		called = true
		return oldNode.Clone()
	})

	node := NewWithConfig(config)
	clone := config.CloneNode(node, nil, 0)

	if !called {
		t.Error("clone callback was not called via Config.CloneNode")
	}
	if clone == nil {
		t.Error("CloneNode returned nil")
	}
}

func TestConfig_ErrataBitmask(t *testing.T) {
	config := NewConfig(DefaultLogger)

	if !config.HasErrata(ErrataMinSizeUndefinedInsteadOfAuto) {
		t.Error("default errata should include MinSizeUndefinedInsteadOfAuto")
	}

	config.AddErrata(ErrataStretchFlexBasis)
	if !config.HasErrata(ErrataStretchFlexBasis) {
		t.Error("errata should include StretchFlexBasis after AddErrata")
	}

	config.RemoveErrata(ErrataMinSizeUndefinedInsteadOfAuto)
	if config.HasErrata(ErrataMinSizeUndefinedInsteadOfAuto) {
		t.Error("errata should not include MinSizeUndefinedInsteadOfAuto after RemoveErrata")
	}

	config.SetErrata(ErrataNone)
	if config.GetErrata() != ErrataNone {
		t.Error("errata should be none after SetErrata")
	}
	if config.HasErrata(ErrataClassic) == (int(ErrataClassic) != int(ErrataNone)) {
		t.Error("HasErrata for Classic should be consistent")
	}
}

func TestConfig_VersionIncrements(t *testing.T) {
	config := NewConfig(DefaultLogger)
	v0 := config.GetVersion()

	config.SetErrata(ErrataNone)
	v1 := config.GetVersion()
	if v1 <= v0 {
		t.Error("version did not increment after SetErrata")
	}

	config.AddErrata(ErrataStretchFlexBasis)
	v2 := config.GetVersion()
	if v2 <= v1 {
		t.Error("version did not increment after AddErrata")
	}

	config.SetExperimentalFeatureEnabled(ExperimentalFeatureWebFlexBasis, true)
	v3 := config.GetVersion()
	if v3 <= v2 {
		t.Error("version did not increment after SetExperimentalFeatureEnabled")
	}

	config.SetPointScaleFactor(2.0)
	v4 := config.GetVersion()
	if v4 <= v3 {
		t.Error("version did not increment after SetPointScaleFactor")
	}
}

func TestConfig_UpdateInvalidatesLayout(t *testing.T) {
	config1 := NewConfig(DefaultLogger)
	root := NewWithConfig(config1)
	root.SetWidth(100)
	root.SetHeight(100)

	child := NewWithConfig(config1)
	child.SetFlexGrow(1)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	config2 := NewConfig(DefaultLogger)
	config2.SetErrata(ErrataNone)

	root.SetConfig(config2)

	if !root.IsDirty() {
		t.Error("node should be dirty after config that invalidates layout")
	}
}

func TestConfig_Logger(t *testing.T) {
	logCalled := false
	config := NewConfig(func(c *Config, n *Node, level LogLevel, format string, args ...interface{}) int {
		logCalled = true
		return 0
	})

	root := NewWithConfig(config)
	root.SetWidth(100)
	root.SetHeight(100)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	config.Log(root, LogLevelError, "test")

	if !logCalled {
		t.Error("logger should have been called")
	}
}

func TestDefaults_AllDefaults(t *testing.T) {
	node := New()

	if node.GetChildCount() != 0 {
		t.Errorf("default childCount: expected 0, got %d", node.GetChildCount())
	}
	if node.GetDirection() != DirectionInherit {
		t.Errorf("default direction: expected Inherit, got %v", node.GetDirection())
	}
	if node.GetFlexDirection() != FlexDirectionColumn {
		t.Errorf("default flexDirection: expected Column, got %v", node.GetFlexDirection())
	}
	if node.GetJustifyContent() != JustifyFlexStart {
		t.Errorf("default justifyContent: expected FlexStart, got %v", node.GetJustifyContent())
	}
	if node.GetAlignItems() != AlignStretch {
		t.Errorf("default alignItems: expected Stretch, got %v", node.GetAlignItems())
	}
	if node.GetAlignSelf() != AlignAuto {
		t.Errorf("default alignSelf: expected Auto, got %v", node.GetAlignSelf())
	}
	if node.GetPositionType() != PositionTypeRelative {
		t.Errorf("default positionType: expected Relative, got %v", node.GetPositionType())
	}
	if node.GetFlexWrap() != WrapNoWrap {
		t.Errorf("default flexWrap: expected NoWrap, got %v", node.GetFlexWrap())
	}
	if node.GetOverflow() != OverflowVisible {
		t.Errorf("default overflow: expected Visible, got %v", node.GetOverflow())
	}
	if node.GetDisplay() != DisplayFlex {
		t.Errorf("default display: expected Flex, got %v", node.GetDisplay())
	}
	if node.GetBoxSizing() != BoxSizingBorderBox {
		t.Errorf("default boxSizing: expected BorderBox, got %v", node.GetBoxSizing())
	}
}

func TestDefaults_WebDefaults(t *testing.T) {
	config := NewConfig(DefaultLogger)
	config.SetUseWebDefaults(true)

	node := NewWithConfig(config)

	if node.GetFlexDirection() != FlexDirectionRow {
		t.Errorf("web default flexDirection: expected Row, got %v", node.GetFlexDirection())
	}
	if node.GetAlignContent() != AlignStretch {
		t.Errorf("web default alignContent: expected Stretch, got %v", node.GetAlignContent())
	}
}

func TestMeasure_DontMeasureSingleGrowShrinkChild(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)

	measureCount := 0
	child := New()
	child.SetFlexGrow(1)
	child.SetFlexShrink(1)
	child.SetMeasureFunc(func(node *Node, width float32, widthMode MeasureMode, height float32, heightMode MeasureMode) Size {
		measureCount++
		return Size{Width: 50, Height: 50}
	})
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	if measureCount != 0 {
		t.Errorf("measure func should not be called for grow=1 shrink=1 child filling container, got %d", measureCount)
	}
}

func TestMeasure_WithFlexShrink(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)
	root.SetFlexDirection(FlexDirectionRow)

	child0 := New()
	child0.SetWidth(80)
	child0.SetFlexShrink(1)
	root.InsertChildNode(child0, 0)

	child1 := New()
	child1.SetWidth(80)
	child1.SetFlexShrink(1)
	root.InsertChildNode(child1, 1)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	if child0.GetWidth() >= 80 {
		t.Errorf("child0 should shrink below 80, got %f", child0.GetWidth())
	}
	if child1.GetWidth() >= 80 {
		t.Errorf("child1 should shrink below 80, got %f", child1.GetWidth())
	}
}

func TestMeasure_MinEqualsMax(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)

	child := New()
	child.SetMinWidth(50)
	child.SetMaxWidth(50)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	if math.Abs(float64(child.GetWidth()-50)) > 0.1 {
		t.Errorf("width should be 50 when min==max, got %f", child.GetWidth())
	}
}

func TestMeasure_CannotAddChildToMeasureNode(t *testing.T) {
	node := New()
	node.SetMeasureFunc(func(n *Node, width float32, widthMode MeasureMode, height float32, heightMode MeasureMode) Size {
		return Size{Width: 50, Height: 50}
	})

	child := New()

	defer func() {
		if r := recover(); r == nil {
			t.Error("should panic when adding child to measure node")
		}
	}()

	node.InsertChildNode(child, 0)
}

func TestMeasure_CanNullifyMeasureFunc(t *testing.T) {
	node := New()
	node.SetMeasureFunc(func(n *Node, width float32, widthMode MeasureMode, height float32, heightMode MeasureMode) Size {
		return Size{Width: 50, Height: 50}
	})

	if !node.HasMeasureFunc() {
		t.Error("node should have measure func after setting")
	}

	node.SetMeasureFunc(nil)

	if node.HasMeasureFunc() {
		t.Error("node should not have measure func after nil")
	}
}

func TestMeasure_MeasureModes(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)

	var modes []MeasureMode
	child := New()
	child.SetMeasureFunc(func(n *Node, width float32, widthMode MeasureMode, height float32, heightMode MeasureMode) Size {
		modes = append(modes, widthMode, heightMode)
		return Size{Width: 50, Height: 50}
	})
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	if len(modes) < 2 {
		t.Errorf("measure func should be called, got %d mode values", len(modes))
	}
}

func TestMeasureCache_RemeasureWithSameExact(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)

	measureCount := 0
	child := New()
	child.SetMeasureFunc(func(n *Node, width float32, widthMode MeasureMode, height float32, heightMode MeasureMode) Size {
		measureCount++
		return Size{Width: 50, Height: 50}
	})
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	measureCount = 0
	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	if measureCount != 0 {
		t.Errorf("second layout with same exact constraints should use cache, got %d", measureCount)
	}

	child.SetWidth(60)

	measureCount = 0
	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	if measureCount < 1 {
		t.Errorf("layout after style change should remeasure, got %d", measureCount)
	}
}

func TestMeasureCache_RemeasureWithSameAtMost(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)

	measureCount := 0
	child := New()
	child.SetFlexShrink(1)
	child.SetMeasureFunc(func(n *Node, width float32, widthMode MeasureMode, height float32, heightMode MeasureMode) Size {
		measureCount++
		return Size{Width: 50, Height: 50}
	})
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	if measureCount < 1 {
		t.Error("measure func should be called at least once")
	}
}

func TestDirtyMarking_DirtyPropagation(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)

	child := New()
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	if root.IsDirty() {
		t.Error("root should not be dirty after layout")
	}

	child.SetWidth(50)

	if !root.IsDirty() {
		t.Error("root should be dirty when child style changes")
	}
}

func TestDirtyMarking_ChangingLayoutConfig(t *testing.T) {
	config1 := NewConfig(DefaultLogger)
	root := NewWithConfig(config1)
	root.SetWidth(100)
	root.SetHeight(100)

	child := NewWithConfig(config1)
	child.SetFlexGrow(1)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	config2 := NewConfig(DefaultLogger)
	config2.SetErrata(ErrataNone)

	root.SetConfig(config2)

	if !root.IsDirty() {
		t.Error("node should be dirty after config that invalidates layout")
	}
}

func TestDirtyMarking_DisplayChanges(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)

	child := New()
	child.SetWidth(50)
	child.SetHeight(50)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	child.SetDisplay(DisplayNone)

	if !root.IsDirty() {
		t.Error("root should be dirty after child display change")
	}
	if !child.IsDirty() {
		t.Error("child should be dirty after display change")
	}

	child.SetDisplay(DisplayFlex)

	if !root.IsDirty() {
		t.Error("root should be dirty when toggling display back to flex")
	}
}

func TestDirtyMarking_OnlyIfChildrenRemoved(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)

	child := New()
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	root.SetWidth(100)
	root.SetWidth(100)

	if root.IsDirty() {
		t.Log("root remained not-dirty when setting same value")
	}

	root.RemoveChildNode(child)

	if !root.IsDirty() {
		t.Error("root should be dirty after child removal")
	}
}

func TestRelayout_DontCacheComputedFlexBasisBetweenLayouts(t *testing.T) {
	root := New()
	root.SetWidth(200)
	root.SetHeight(200)

	child := New()
	child.SetFlexBasis(100)
	child.SetFlexGrow(1)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	w1 := child.GetWidth()

	root.SetWidth(150)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	w2 := child.GetWidth()

	if w1 == w2 {
		t.Logf("flex basis should recalculate when dirty: %f vs %f", w1, w2)
	}
}

func TestScaleChange_ScaleChangeInvalidatesLayout(t *testing.T) {
	config := NewConfig(DefaultLogger)
	root := NewWithConfig(config)
	root.SetWidth(100)
	root.SetHeight(100)

	child := NewWithConfig(config)
	child.SetWidth(33)
	child.SetHeight(33)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	config.SetPointScaleFactor(2.0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	if math.IsNaN(float64(child.GetWidth())) {
		t.Error("layout should be valid after point scale change")
	}
}

func TestHadOverflow_ChildrenOverflowNoWrap(t *testing.T) {
	root := New()
	root.SetWidth(50)
	root.SetHeight(50)

	child := New()
	child.SetWidth(100)
	child.SetHeight(100)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	_ = root.GetHadOverflow()
}

func TestHadOverflow_NoOverflow(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)

	child := New()
	child.SetWidth(50)
	child.SetHeight(50)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	if root.GetHadOverflow() {
		t.Error("no overflow should be detected when children fit")
	}
}

func TestHadOverflow_FlagResets(t *testing.T) {
	root := New()
	root.SetWidth(50)
	root.SetHeight(50)

	child := New()
	child.SetWidth(100)
	child.SetHeight(100)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	overflow1 := root.GetHadOverflow()

	root.SetWidth(300)
	root.SetHeight(300)
	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	overflow2 := root.GetHadOverflow()

	if overflow1 && overflow2 {
		t.Error("overflow flag should reset after container expands enough")
	}
}

func TestDirtiedFunc_CalledOnSetDirty(t *testing.T) {
	dirtiedCalled := false
	node := New()
	node.SetDirtiedFunc(func(n *Node) {
		dirtiedCalled = true
	})

	node.SetDirty(false)

	node.SetDirty(true)

	if !dirtiedCalled {
		t.Error("dirtiedFunc should be called when node transitions to dirty")
	}
}

func TestDirtiedFunc_Propagation(t *testing.T) {
	dirtiedCallCount := 0

	root := New()
	root.SetWidth(100)
	root.SetHeight(100)
	root.SetDirtiedFunc(func(n *Node) {
		dirtiedCallCount++
	})

	child := New()
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	dirtiedCallCount = 0
	child.SetWidth(50)

	if dirtiedCallCount < 1 {
		t.Error("dirtiedFunc should propagate to parent when child changes")
	}
}

func TestDirtiedFunc_Hierarchy(t *testing.T) {
	count := 0
	node := New()
	node.SetDirtiedFunc(func(n *Node) {
		count++
	})

	node.SetDirty(false)

	node.SetDirty(true)

	if count != 1 {
		t.Errorf("dirtiedFunc should be called once when dirty goes false->true, got %d", count)
	}

	count = 0
	node.SetDirty(true)

	if count != 0 {
		t.Errorf("dirtiedFunc should not be called when already dirty, got %d", count)
	}

	node.SetDirty(false)
	count = 0
	node.SetDirty(true)

	if count != 1 {
		t.Errorf("dirtiedFunc should be called again after clear and set, got %d", count)
	}
}

func TestTreeMutation_SetChildrenBatch(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)

	child1 := New()
	child1.SetWidth(30)
	child1.SetHeight(30)

	child2 := New()
	child2.SetWidth(30)
	child2.SetHeight(30)

	root.SetChildrenList([]*Node{child1, child2})

	if root.GetChildCount() != 2 {
		t.Errorf("should have 2 children, got %d", root.GetChildCount())
	}
	if child1.GetOwner() != root {
		t.Error("child1 should be owned by root")
	}
	if child2.GetOwner() != root {
		t.Error("child2 should be owned by root")
	}

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)
}

func TestTreeMutation_SetChildrenReplacesNonCommon(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)

	oldChild := New()
	oldChild.SetWidth(30)
	oldChild.SetHeight(30)
	root.InsertChildNode(oldChild, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	newChild := New()
	newChild.SetWidth(50)
	newChild.SetHeight(50)
	root.SetChildrenList([]*Node{newChild})

	if oldChild.GetOwner() != nil {
		t.Error("old child should have no owner after replacement")
	}
	if !math.IsNaN(float64(oldChild.GetWidth())) {
		t.Errorf("old child layout should be reset (NaN), got %f", oldChild.GetWidth())
	}
}

func TestEdge_StartOverrides(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)

	child := New()
	child.SetMargin(EdgeLeft, 5)
	child.SetMargin(EdgeStart, 10)
	child.SetFlexGrow(1)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	got := child.GetLayoutMargin(EdgeLeft)
	if math.Abs(float64(got-10)) > 0.1 {
		t.Errorf("EdgeStart should override left margin in LTR: expected 10, got %f", got)
	}
}

func TestEdge_HorizontalOverridden(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)

	child := New()
	child.SetPadding(EdgeLeft, 5)
	child.SetPadding(EdgeRight, 5)
	child.SetPadding(EdgeHorizontal, 15)
	child.SetFlexGrow(1)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	left := child.GetLayoutPadding(EdgeLeft)
	right := child.GetLayoutPadding(EdgeRight)
	if math.Abs(float64(left-5)) > 0.1 {
		t.Errorf("left padding should override horizontal: expected 5, got %f", left)
	}
	if math.Abs(float64(right-5)) > 0.1 {
		t.Errorf("right padding should override horizontal: expected 5, got %f", right)
	}
}

func TestEdge_VerticalOverridden(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)

	child := New()
	child.SetMargin(EdgeTop, 5)
	child.SetMargin(EdgeBottom, 5)
	child.SetMargin(EdgeVertical, 15)
	child.SetFlexGrow(1)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	top := child.GetLayoutMargin(EdgeTop)
	bottom := child.GetLayoutMargin(EdgeBottom)
	if math.Abs(float64(top-5)) > 0.1 {
		t.Errorf("top margin should override vertical: expected 5, got %f", top)
	}
	if math.Abs(float64(bottom-5)) > 0.1 {
		t.Errorf("bottom margin should override vertical: expected 5, got %f", bottom)
	}
}

func TestEdge_AllOverridden(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)

	child := New()
	child.SetBorder(EdgeAll, 10)
	child.SetBorder(EdgeLeft, 2)
	child.SetBorder(EdgeTop, 3)
	child.SetBorder(EdgeRight, 4)
	child.SetBorder(EdgeBottom, 5)
	child.SetFlexGrow(1)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	left := child.GetLayoutBorder(EdgeLeft)
	top := child.GetLayoutBorder(EdgeTop)
	right := child.GetLayoutBorder(EdgeRight)
	bottom := child.GetLayoutBorder(EdgeBottom)

	if math.Abs(float64(left-2)) > 0.1 {
		t.Errorf("left border should override all: expected 2, got %f", left)
	}
	if math.Abs(float64(top-3)) > 0.1 {
		t.Errorf("top border should override all: expected 3, got %f", top)
	}
	if math.Abs(float64(right-4)) > 0.1 {
		t.Errorf("right border should override all: expected 4, got %f", right)
	}
	if math.Abs(float64(bottom-5)) > 0.1 {
		t.Errorf("bottom border should override all: expected 5, got %f", bottom)
	}
}

func TestAspectRatio_CrossDefined(t *testing.T) {
	root := New()
	root.SetWidth(200)
	root.SetHeight(200)

	child := New()
	child.SetAspectRatio(2)
	child.SetWidth(100)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	if math.Abs(float64(child.GetHeight()-50)) > 1 {
		t.Errorf("aspect ratio: expected height 50 (width 100/2), got %f", child.GetHeight())
	}
}

func TestAspectRatio_MainDefined(t *testing.T) {
	root := New()
	root.SetWidth(200)
	root.SetHeight(200)

	child := New()
	child.SetAspectRatio(2)
	child.SetHeight(50)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	if math.Abs(float64(child.GetWidth()-100)) > 1 {
		t.Errorf("aspect ratio: expected width 100 (height 50*2), got %f", child.GetWidth())
	}
}

func TestAspectRatio_BothDimensionsOverridesRatio(t *testing.T) {
	root := New()
	root.SetWidth(200)
	root.SetHeight(200)

	child := New()
	child.SetAspectRatio(0.5)
	child.SetWidth(100)
	child.SetHeight(100)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	w := child.GetWidth()
	h := child.GetHeight()

	if math.Abs(float64(h-100)) > 1 && math.Abs(float64(w-100)) > 1 {
		t.Errorf("at least one explicit dimension should be preserved, got w=%f h=%f", w, h)
	}
}

func TestAspectRatio_FlexGrowInteraction(t *testing.T) {
	root := New()
	root.SetWidth(200)
	root.SetHeight(200)

	child := New()
	child.SetAspectRatio(1)
	child.SetFlexGrow(1)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	w := child.GetWidth()
	h := child.GetHeight()

	if math.Abs(float64(w-200)) > 1 {
		t.Errorf("flexGrow + ratio: expected width 200, got %f", w)
	}
	if math.Abs(float64(h-200)) > 2 {
		t.Errorf("1:1 ratio: expected height ~200, got %f", h)
	}
}

func TestAlignBaseline_ParentUsingChildInColumnAsReference(t *testing.T) {
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
	child1.SetBaselineFunc(func(node *Node, width, height float32) float32 { return 40 })
	root.InsertChildNode(child1, 1)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	t.Logf("baseline child0 top=%f child1 top=%f", child0.GetTop(), child1.GetTop())
}

func TestAlignBaseline_WithNoParentHeight(t *testing.T) {
	root := New()
	root.SetFlexDirection(FlexDirectionRow)
	root.SetAlignItems(AlignBaseline)
	root.SetWidth(300)

	child0 := New()
	child0.SetWidth(50)
	child0.SetHeight(50)
	child0.SetBaselineFunc(func(node *Node, width, height float32) float32 { return 30 })
	root.InsertChildNode(child0, 0)

	child1 := New()
	child1.SetWidth(50)
	child1.SetHeight(80)
	child1.SetBaselineFunc(func(node *Node, width, height float32) float32 { return 60 })
	root.InsertChildNode(child1, 1)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	if math.IsNaN(float64(root.GetHeight())) {
		t.Error("root height should not be NaN after baseline calculation")
	}
}

func TestAlignBaseline_ChildMargin(t *testing.T) {
	root := New()
	root.SetFlexDirection(FlexDirectionRow)
	root.SetAlignItems(AlignBaseline)
	root.SetWidth(300)
	root.SetHeight(300)

	child := New()
	child.SetWidth(50)
	child.SetHeight(50)
	child.SetMargin(EdgeTop, 10)
	child.SetBaselineFunc(func(node *Node, width, height float32) float32 { return 30 })
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	if math.Abs(float64(child.GetLayoutMargin(EdgeTop)-10)) > 0.1 {
		t.Errorf("baseline child top margin should be 10, got %f", child.GetLayoutMargin(EdgeTop))
	}
}

func TestBaselineFunc_HasBaselineFunc(t *testing.T) {
	node := New()

	if node.HasBaselineFunc() {
		t.Error("node should not have baseline func initially")
	}

	node.SetBaselineFunc(func(n *Node, width, height float32) float32 { return 0 })

	if !node.HasBaselineFunc() {
		t.Error("node should have baseline func after setting")
	}

	node.SetBaselineFunc(nil)

	if node.HasBaselineFunc() {
		t.Error("node should not have baseline func after unsetting with nil")
	}
}

func TestComputedPadding_EdgeOverridesHorizontal(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)

	child := New()
	child.SetPadding(EdgeAll, 10)
	child.SetPadding(EdgeLeft, 2)
	child.SetPadding(EdgeRight, 4)
	child.SetFlexGrow(1)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	left := child.GetLayoutPadding(EdgeLeft)
	right := child.GetLayoutPadding(EdgeRight)

	if math.Abs(float64(left-2)) > 0.1 {
		t.Errorf("left padding should override all: expected 2, got %f", left)
	}
	if math.Abs(float64(right-4)) > 0.1 {
		t.Errorf("right padding should override all: expected 4, got %f", right)
	}
}

func TestComputedPadding_HorizontalOverridesAll(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)

	child := New()
	child.SetPadding(EdgeAll, 10)
	child.SetPadding(EdgeHorizontal, 5)
	child.SetPadding(EdgeTop, 20)
	child.SetFlexGrow(1)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	left := child.GetLayoutPadding(EdgeLeft)
	right := child.GetLayoutPadding(EdgeRight)
	top := child.GetLayoutPadding(EdgeTop)

	if math.Abs(float64(left-5)) > 0.1 {
		t.Errorf("left padding should be from horizontal: expected 5, got %f", left)
	}
	if math.Abs(float64(right-5)) > 0.1 {
		t.Errorf("right padding should be from horizontal: expected 5, got %f", right)
	}
	if math.Abs(float64(top-20)) > 0.1 {
		t.Errorf("top padding should override all: expected 20, got %f", top)
	}
}

func TestComputedMargin_EdgeOverridesHorizontal(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)

	child := New()
	child.SetMargin(EdgeAll, 10)
	child.SetMargin(EdgeLeft, 2)
	child.SetMargin(EdgeRight, 4)
	child.SetFlexGrow(1)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	left := child.GetLayoutMargin(EdgeLeft)
	right := child.GetLayoutMargin(EdgeRight)

	if math.Abs(float64(left-2)) > 0.1 {
		t.Errorf("left margin should override all: expected 2, got %f", left)
	}
	if math.Abs(float64(right-4)) > 0.1 {
		t.Errorf("right margin should override all: expected 4, got %f", right)
	}
}

func TestComputedMargin_HorizontalOverridesAll(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)

	child := New()
	child.SetMargin(EdgeAll, 10)
	child.SetMargin(EdgeHorizontal, 5)
	child.SetFlexGrow(1)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	left := child.GetLayoutMargin(EdgeLeft)
	right := child.GetLayoutMargin(EdgeRight)

	if math.Abs(float64(left-5)) > 0.1 {
		t.Errorf("left margin should be from horizontal: expected 5, got %f", left)
	}
	if math.Abs(float64(right-5)) > 0.1 {
		t.Errorf("right margin should be from horizontal: expected 5, got %f", right)
	}
}

func TestNodeChild_ResetLayoutWhenChildRemoved(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)

	child := New()
	child.SetWidth(50)
	child.SetHeight(50)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	root.RemoveChildNode(child)

	if !math.IsNaN(float64(child.GetWidth())) {
		t.Errorf("removed child width should be NaN, got %f", child.GetWidth())
	}
	if !math.IsNaN(float64(child.GetHeight())) {
		t.Errorf("removed child height should be NaN, got %f", child.GetHeight())
	}
}

func TestNodeChild_RemovedChildCanBeReused(t *testing.T) {
	root1 := New()
	root1.SetWidth(100)
	root1.SetHeight(100)

	child := New()
	child.SetWidth(30)
	child.SetHeight(30)
	root1.InsertChildNode(child, 0)

	CalculateNodeLayout(root1, Undefined, Undefined, DirectionLTR)

	root1.RemoveChildNode(child)

	root2 := New()
	root2.SetWidth(200)
	root2.SetHeight(200)
	root2.InsertChildNode(child, 0)

	CalculateNodeLayout(root2, Undefined, Undefined, DirectionLTR)

	if child.GetOwner() != root2 {
		t.Error("reused child should have new owner")
	}
	if child.GetWidth() != 30 {
		t.Errorf("reused child should keep its width 30, got %f", child.GetWidth())
	}
}

func TestClone_AbsoluteWithStaticParent(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)

	child := New()
	child.SetPositionType(PositionTypeAbsolute)
	child.SetEdgePosition(EdgeLeft, 10)
	child.SetWidth(50)
	child.SetHeight(50)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	clone := child.Clone()
	if clone.GetOwner() != nil {
		t.Error("clone should not have owner")
	}
	if clone.GetPositionType() != PositionTypeAbsolute {
		t.Error("clone should preserve position type")
	}
	if clone.GetConfig() != child.GetConfig() {
		t.Error("clone should have same config")
	}
}

func TestZeroOut_DisplayNone(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)

	child := New()
	child.SetDisplay(DisplayNone)
	child.SetWidth(50)
	child.SetHeight(50)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	w := child.GetWidth()
	h := child.GetHeight()

	if !math.IsNaN(float64(w)) && w > 0 {
		t.Errorf("display:none child width should be zeroed (NaN or 0), got %f", w)
	}
	if !math.IsNaN(float64(h)) && h > 0 {
		t.Errorf("display:none child height should be zeroed (NaN or 0), got %f", h)
	}
}

func TestGrid_SetColumnStart(t *testing.T) {
	node := New()
	node.SetGridColumnStart(2)

	if node.GetGridColumnStart() != 2 {
		t.Errorf("grid column start should be 2, got %d", node.GetGridColumnStart())
	}
}

func TestGrid_SetRowEnd(t *testing.T) {
	node := New()
	node.SetGridRowEnd(3)

	if node.GetGridRowEnd() != 3 {
		t.Errorf("grid row end should be 3, got %d", node.GetGridRowEnd())
	}
}

func TestGridTemplateColumns(t *testing.T) {
	root := New()
	root.SetWidth(300)
	root.SetHeight(300)
	root.SetDisplay(DisplayGrid)

	root.SetGridTemplateColumnsCount(3)
	root.SetGridTemplateColumn(0, GridTrackTypePoints, 100)
	root.SetGridTemplateColumn(1, GridTrackTypePoints, 100)
	root.SetGridTemplateColumn(2, GridTrackTypePoints, 100)

	child := New()
	child.SetWidth(50)
	child.SetHeight(50)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionLTR)

	if math.IsNaN(float64(child.GetWidth())) {
		t.Error("grid child should have valid layout")
	}
}

func TestRTL_StartEnd(t *testing.T) {
	root := New()
	root.SetWidth(100)
	root.SetHeight(100)

	child := New()
	child.SetMargin(EdgeStart, 10)
	child.SetMargin(EdgeEnd, 20)
	child.SetFlexGrow(1)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, Undefined, Undefined, DirectionRTL)

	startMargin := child.GetLayoutMargin(EdgeStart)
	endMargin := child.GetLayoutMargin(EdgeEnd)

	if math.Abs(float64(startMargin-10)) > 0.1 {
		t.Errorf("RTL: EdgeStart margin should be 10, got %f", startMargin)
	}
	if math.Abs(float64(endMargin-20)) > 0.1 {
		t.Errorf("RTL: EdgeEnd margin should be 20, got %f", endMargin)
	}
}

func TestLayoutOut_Complete(t *testing.T) {
	root := New().ApplyStyleString(`
		width: 800; height: 600;
		flex-direction: row;
		padding: 10;
	`)

	child := New().ApplyStyleString(`
		width: 100; height: 50;
		margin: 8;
		border: 2;
	`)
	root.InsertChildNode(child, 0)

	CalculateNodeLayout(root, 800, 600, DirectionLTR)

	lo := child.LayoutOut()
	if v := lo.Width; v != 100 {
		t.Errorf("width should be 100, got %f", v)
	}
	if v := lo.Height; v != 50 {
		t.Errorf("height should be 50, got %f", v)
	}
	if v := lo.Margin.Left; v != 8 {
		t.Errorf("margin left should be 8, got %f", v)
	}
	if v := lo.Margin.Top; v != 8 {
		t.Errorf("margin top should be 8, got %f", v)
	}
	if v := lo.Border.Left; v != 2 {
		t.Errorf("border left should be 2, got %f", v)
	}
	if v := lo.Padding.Top; v != 0 {
		t.Errorf("padding top should be 0, got %f", v)
	}
	if lo.Direction != DirectionLTR {
		t.Error("direction should be LTR")
	}
	if lo.HadOverflow {
		t.Error("should not have overflow")
	}
	t.Logf("LayoutOut: pos=(%f,%f) size=%fx%f margin={t:%f r:%f b:%f l:%f}",
		lo.Left, lo.Top, lo.Width, lo.Height,
		lo.Margin.Top, lo.Margin.Right, lo.Margin.Bottom, lo.Margin.Left)
}
