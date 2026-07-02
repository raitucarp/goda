package goda

import (
	"math"
	"strings"
	"testing"
)

func TestNode_ID(t *testing.T) {
	n := New("my_id")
	if n.GetID() != "my_id" {
		t.Errorf("expected id 'my_id', got %q", n.GetID())
	}
}

func TestNode_ID_Default(t *testing.T) {
	n := New()
	if n.GetID() != "" {
		t.Errorf("expected empty id, got %q", n.GetID())
	}
}

func TestNode_SetID(t *testing.T) {
	n := New()
	n.SetID("hello")
	if n.GetID() != "hello" {
		t.Errorf("expected id 'hello', got %q", n.GetID())
	}
}

func TestNode_AddClass(t *testing.T) {
	n := New()
	n.AddClass("card")
	n.AddClass("highlight")
	n.AddClass("card") // duplicate should not add

	classes := n.GetClasses()
	if len(classes) != 2 {
		t.Errorf("expected 2 classes, got %d", len(classes))
	}
	if !n.HasClass("card") {
		t.Error("expected HasClass('card') to be true")
	}
	if !n.HasClass("highlight") {
		t.Error("expected HasClass('highlight') to be true")
	}
	if n.HasClass("nonexistent") {
		t.Error("expected HasClass('nonexistent') to be false")
	}
}

func TestNode_SetClasses(t *testing.T) {
	n := New()
	n.SetClasses([]string{"a", "b"})
	if !n.HasClass("a") || !n.HasClass("b") {
		t.Error("expected classes a and b")
	}
	n.SetClasses(nil)
	if len(n.GetClasses()) != 0 {
		t.Error("expected no classes after setting nil")
	}
}

func TestNode_ClonePreservesIDAndClasses(t *testing.T) {
	n := New("original")
	n.AddClass("foo")
	n.AddClass("bar")
	clone := n.Clone()

	if clone.GetID() != "original" {
		t.Errorf("expected cloned id 'original', got %q", clone.GetID())
	}
	if !clone.HasClass("foo") || !clone.HasClass("bar") {
		t.Error("clone should preserve classes")
	}
}

func TestApplyStyle_JustifyAlias(t *testing.T) {
	n := New().ApplyStyle(map[string]string{
		"justify": "center",
	})
	if n.GetJustifyContent() != JustifyCenter {
		t.Errorf("expected justify-content center, got %v", n.GetJustifyContent())
	}
}

func TestApplyStyleString_JustifyAlias(t *testing.T) {
	n := New().ApplyStyleString("justify: space-between;")
	if n.GetJustifyContent() != JustifySpaceBetween {
		t.Errorf("expected justify-content space-between, got %v", n.GetJustifyContent())
	}
}

func TestRenderFrom_Basic(t *testing.T) {
	roots, err := RenderFrom(`
		#root {
			width: 800;
			height: 600;
		}
	`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(roots) != 1 {
		t.Fatalf("expected 1 root, got %d", len(roots))
	}
	root := roots[0]
	if root.GetID() != "root" {
		t.Errorf("expected id 'root', got %q", root.GetID())
	}
	if v := root.GetWidthValue().Value; v != 800 {
		t.Errorf("expected width 800, got %f", v)
	}
	if v := root.GetHeightValue().Value; v != 600 {
		t.Errorf("expected height 600, got %f", v)
	}
}

func TestRenderFrom_WithClasses(t *testing.T) {
	roots, err := RenderFrom(`
		.myClass {
			flex: 1;
			display: flex;
			justify: space-between;
		}
		#myId[myClass] {
			height: 30;
			width: 500;
		}
	`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(roots) != 1 {
		t.Fatalf("expected 1 root, got %d", len(roots))
	}
	node := roots[0]
	if node.GetID() != "myId" {
		t.Errorf("expected id 'myId', got %q", node.GetID())
	}
	if !node.HasClass("myClass") {
		t.Error("expected to have class 'myClass'")
	}
	if node.GetDisplay() != DisplayFlex {
		t.Error("expected display flex from class")
	}
	if node.GetJustifyContent() != JustifySpaceBetween {
		t.Errorf("expected justify-content space-between from class, got %v", node.GetJustifyContent())
	}
	// Verify flex:1 works via layout
	if node.GetFlex() != 1 {
		t.Errorf("expected flex 1 from class, got %f", node.GetFlex())
	}
}

func TestRenderFrom_NestedChildren(t *testing.T) {
	roots, err := RenderFrom(`
		#root {
			width: 800;
			height: 600;
			#child1 {
				flex: 1;
			}
			#child2 {
				width: 200;
			}
		}
	`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	root := roots[0]
	if root.GetChildCount() != 2 {
		t.Fatalf("expected 2 children, got %d", root.GetChildCount())
	}

	child1 := root.GetChild(0)
	if child1.GetID() != "child1" {
		t.Errorf("expected child1 id, got %q", child1.GetID())
	}
	if child1.GetFlex() != 1 {
		t.Errorf("expected child1 flex 1, got %f", child1.GetFlex())
	}

	child2 := root.GetChild(1)
	if child2.GetID() != "child2" {
		t.Errorf("expected child2 id, got %q", child2.GetID())
	}
	if v := child2.GetWidthValue().Value; v != 200 {
		t.Errorf("expected child2 width 200, got %f", v)
	}
}

func TestRenderFrom_MultipleClasses(t *testing.T) {
	roots, err := RenderFrom(`
		.card {
			padding: 8;
		}
		.highlight {
			border: 2;
		}
		#item[card, highlight] {
			width: 100;
		}
	`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	node := roots[0]
	if !node.HasClass("card") || !node.HasClass("highlight") {
		t.Error("expected both classes")
	}
	if v := node.GetBorder(EdgeAll); math.IsNaN(float64(v)) || v != 2 {
		t.Errorf("expected border 2, got %f", v)
	}
}

func TestRenderFrom_Comments(t *testing.T) {
	roots, err := RenderFrom(`
		// This is a line comment
		#root {
			/* block comment */
			width: 500;
			// inline comment
			height: 400;
		}
	`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	root := roots[0]
	if v := root.GetWidthValue().Value; v != 500 {
		t.Errorf("expected width 500, got %f", v)
	}
	if v := root.GetHeightValue().Value; v != 400 {
		t.Errorf("expected height 400, got %f", v)
	}
}

func TestRenderFrom_ClassWithoutBracket(t *testing.T) {
	roots, err := RenderFrom(`
		#simple {
			width: 100;
		}
	`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if roots[0].GetID() != "simple" {
		t.Errorf("expected id 'simple', got %q", roots[0].GetID())
	}
}

func TestRenderFrom_MultipleRoots(t *testing.T) {
	roots, err := RenderFrom(`
		#first {
			width: 100;
		}
		#second {
			height: 200;
		}
	`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(roots) != 2 {
		t.Fatalf("expected 2 roots, got %d", len(roots))
	}
	if roots[0].GetID() != "first" {
		t.Errorf("expected first root id 'first', got %q", roots[0].GetID())
	}
	if roots[1].GetID() != "second" {
		t.Errorf("expected second root id 'second', got %q", roots[1].GetID())
	}
}

func TestRenderFrom_SelfOverridesClass(t *testing.T) {
	roots, err := RenderFrom(`
		.card {
			width: 100;
			height: 50;
		}
		#item[card] {
			width: 200;
		}
	`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	node := roots[0]
	if v := node.GetWidthValue().Value; v != 200 {
		t.Errorf("self width should override class width: expected 200, got %f", v)
	}
	if v := node.GetHeightValue().Value; v != 50 {
		t.Errorf("height should come from class: expected 50, got %f", v)
	}
}

func TestRenderFrom_LayoutWorks(t *testing.T) {
	roots, err := RenderFrom(`
		#container {
			width: 800;
			height: 600;
			flex-direction: row;
			#sidebar {
				width: 200;
				flex-shrink: 0;
			}
			#content {
				flex: 1;
			}
		}
	`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	root := roots[0]
	CalculateNodeLayout(root, 800, 600, DirectionLTR)

	sidebar := root.GetChild(0)
	content := root.GetChild(1)

	if math.Abs(float64(sidebar.GetWidth()-200)) > 0.1 {
		t.Errorf("expected sidebar width 200, got %f", sidebar.GetWidth())
	}
	if math.Abs(float64(content.GetWidth()-600)) > 0.1 {
		t.Errorf("expected content width 600, got %f", content.GetWidth())
	}
}

func TestExportAs_RoundTrip(t *testing.T) {
	n := New("root")
	n.AddClass("container")
	n.ApplyStyleString("width: 800; height: 600; flex-direction: row; padding: 16;")

	child := New("header")
	child.ApplyStyleString("height: 64; flex-shrink: 0;")
	n.InsertChildNode(child, 0)

	body := New("body")
	body.ApplyStyleString("flex: 1; padding: 8;")
	n.InsertChildNode(body, 1)

	exported := n.ExportAs()
	if exported == "" {
		t.Fatal("ExportAs returned empty string")
	}
	if !strings.Contains(exported, "#root") {
		t.Error("export should contain #root")
	}
	if !strings.Contains(exported, "width:") {
		t.Error("export should contain width property")
	}
	if !strings.Contains(exported, "#header") {
		t.Error("export should contain #header child")
	}
	if !strings.Contains(exported, "#body") {
		t.Error("export should contain #body child")
	}

	roots, err := RenderFrom(exported)
	if err != nil {
		t.Fatalf("round-trip parse error: %v", err)
	}
	if len(roots) != 1 {
		t.Fatalf("round-trip expected 1 root, got %d", len(roots))
	}
	parsed := roots[0]
	if parsed.GetID() != "root" {
		t.Errorf("round-trip id mismatch: expected 'root', got %q", parsed.GetID())
	}
	if v := parsed.GetWidthValue().Value; v != 800 {
		t.Errorf("round-trip width: expected 800, got %f", v)
	}
	if parsed.GetChildCount() != 2 {
		t.Errorf("round-trip child count: expected 2, got %d", parsed.GetChildCount())
	}
	pHeader := parsed.GetChild(0)
	if pHeader.GetID() != "header" {
		t.Errorf("round-trip child0 id: expected 'header', got %q", pHeader.GetID())
	}
	pBody := parsed.GetChild(1)
	if pBody.GetID() != "body" {
		t.Errorf("round-trip child1 id: expected 'body', got %q", pBody.GetID())
	}
}

func TestExportAs_OmitsDefaults(t *testing.T) {
	n := New("test")
	out := n.ExportAs()
	if !strings.Contains(out, "#test") {
		t.Error("export should contain the node id")
	}
	lines := strings.Split(out, "\n")
	propCount := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, ":") && !strings.Contains(trimmed, "{") && !strings.Contains(trimmed, "}") {
			propCount++
		}
	}
	if propCount > 0 {
		t.Errorf("expected no properties for default node, got %d: %s", propCount, out)
	}
}

func TestExportAs_DeepNesting(t *testing.T) {
	n := New("root")
	n.ApplyStyleString("width: 300; height: 200;")
	a := New("a")
	a.ApplyStyleString("width: 100;")
	b := New("b")
	b.ApplyStyleString("height: 50;")
	n.InsertChildNode(a, 0)
	a.InsertChildNode(b, 0)

	out := n.ExportAs()
	roots, err := RenderFrom(out)
	if err != nil {
		t.Fatalf("round-trip error: %v", err)
	}
	parsed := roots[0]
	if parsed.GetChildCount() != 1 {
		t.Fatalf("expected 1 child, got %d", parsed.GetChildCount())
	}
	parsedA := parsed.GetChild(0)
	if parsedA.GetChildCount() != 1 {
		t.Fatalf("expected 1 grandchild, got %d", parsedA.GetChildCount())
	}
	parsedB := parsedA.GetChild(0)
	if parsedB.GetID() != "b" {
		t.Errorf("expected grandchild id 'b', got %q", parsedB.GetID())
	}
}

func TestRenderFrom_StyleVariants(t *testing.T) {
	roots, err := RenderFrom(`
		#test {
			width: 50%;
			height: auto;
			min-width: max-content;
			flex-basis: fit-content;
			position: absolute;
			box-sizing: content-box;
			overflow: hidden;
			flex-wrap: wrap;
			direction: rtl;
			aspect-ratio: 2;
		}
	`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	node := roots[0]

	if unit := node.GetWidthValue().Unit; unit != UnitPercent {
		t.Errorf("width unit should be percent, got %v", unit)
	}
	if unit := node.GetHeightValue().Unit; unit != UnitAuto {
		t.Errorf("height unit should be auto, got %v", unit)
	}
	if node.GetPositionType() != PositionTypeAbsolute {
		t.Error("position should be absolute")
	}
	if node.GetBoxSizing() != BoxSizingContentBox {
		t.Error("box-sizing should be content-box")
	}
	if node.GetOverflow() != OverflowHidden {
		t.Error("overflow should be hidden")
	}
	if node.GetFlexWrap() != WrapWrap {
		t.Error("flex-wrap should be wrap")
	}
	if node.GetDirection() != DirectionRTL {
		t.Error("direction should be rtl")
	}
	if node.GetAspectRatio() != 2 {
		t.Errorf("aspect-ratio should be 2, got %f", node.GetAspectRatio())
	}
}

func TestRenderFrom_Empty(t *testing.T) {
	roots, err := RenderFrom("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(roots) != 0 {
		t.Errorf("expected 0 roots, got %d", len(roots))
	}
}

func TestRenderFrom_OnlyClasses(t *testing.T) {
	roots, err := RenderFrom(`
		.card {
			padding: 8;
		}
		.box {
			margin: 4;
		}
	`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(roots) != 0 {
		t.Errorf("expected 0 roots (only class defs), got %d", len(roots))
	}
}
