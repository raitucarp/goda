package goda

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

// ── Lexer ───────────────────────────────────────────────────────────────────

var godaLexerDef = lexer.MustSimple([]lexer.SimpleRule{
	{"Comment", `//[^\n]*|/\*[\s\S]*?\*/`},
	{"Ident", `[a-zA-Z_][a-zA-Z0-9_-]*`},
	{"Number", `\d+(\.\d+)?(%|rem|em|px)?`},
	{"Punct", `[.#{}\[\]:;,]`},
	{"Whitespace", `[ \t\n\r]+`},
})

// ── Grammar ─────────────────────────────────────────────────────────────────

type sourceAST struct {
	Items []*sourceItem `@@*`
}

type sourceItem struct {
	ClassDef *classDef `  @@`
	NodeDef  *nodeDef  `| @@`
}

type classDef struct {
	Name  string   `"." @Ident`
	Decls []*decl  `"{" @@* "}"`
}

type nodeDef struct {
	ID      string   `"#" @Ident`
	Classes []string `("[" @Ident ("," @Ident)* "]")?`
	Body    *nodeBody `"{" @@ "}"`
}

type nodeBody struct {
	Items []*bodyItem `@@*`
}

type bodyItem struct {
	Decl    *decl    `  @@`
	NodeDef *nodeDef `| @@`
}

type decl struct {
	Key   string `@Ident ":"`
	Value string `@(Ident | Number) ";"?`
}

var godaParser = participle.MustBuild[sourceAST](
	participle.Lexer(godaLexerDef),
	participle.Elide("Comment", "Whitespace"),
	participle.UseLookahead(5),
)

// ── RenderFrom ──────────────────────────────────────────────────────────────

// RenderFrom parses an extended CSS / QML-like string and returns the root nodes.
// Class definitions (e.g. ".myClass { ... }") define reusable style blocks.
// Node definitions (e.g. "#myId[class1, class2] { ... }") create nodes with
// optional class references whose styles are applied as defaults.
//
// Children are nested inside braces:
//
//	#root {
//	  width: 800;
//	  #child {
//	    flex: 1;
//	  }
//	}
//
// Comments (// and /* */) are supported anywhere.
func RenderFrom(source string) ([]*Node, error) {
	ast, err := godaParser.ParseString("", source)
	if err != nil {
		return nil, err
	}

	classes := make(map[string]map[string]string)
	for _, item := range ast.Items {
		if item.ClassDef != nil {
			props := make(map[string]string)
			for _, d := range item.ClassDef.Decls {
				props[d.Key] = d.Value
			}
			classes[item.ClassDef.Name] = props
		}
	}

	var roots []*Node
	for _, item := range ast.Items {
		if item.NodeDef != nil {
			node := nodeDefToNode(item.NodeDef, classes)
			roots = append(roots, node)
		}
	}
	return roots, nil
}

func nodeDefToNode(def *nodeDef, classes map[string]map[string]string) *Node {
	node := New(def.ID)
	if len(def.Classes) > 0 {
		node.SetClasses(def.Classes)
		for _, cls := range def.Classes {
			if props, ok := classes[cls]; ok {
				node.ApplyStyle(props)
			}
		}
	}
	if def.Body != nil {
		for _, item := range def.Body.Items {
			if item.Decl != nil {
				node.ApplyStyle(map[string]string{item.Decl.Key: item.Decl.Value})
			} else if item.NodeDef != nil {
				child := nodeDefToNode(item.NodeDef, classes)
				node.InsertChildNode(child, node.GetChildCount())
			}
		}
	}
	return node
}

// ── ExportAs ────────────────────────────────────────────────────────────────

// ExportAs serializes the node tree into the same extended CSS format
// that RenderFrom can parse. Only non-default style properties are included.
func (n *Node) ExportAs() string {
	var b strings.Builder
	n.exportAs(&b, 0)
	return b.String()
}

func (n *Node) exportAs(b *strings.Builder, indent int) {
	pad := strings.Repeat("  ", indent)

	b.WriteString(pad)
	b.WriteByte('#')
	b.WriteString(n.id)
	if len(n.classes) > 0 {
		b.WriteByte('[')
		for i, c := range n.classes {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(c)
		}
		b.WriteByte(']')
	}
	b.WriteString(" {\n")

	defaultStyle := NewStyle()
	props := styleToCSSProps(&n.style, &defaultStyle)
	for _, entry := range sortedProps(props) {
		b.WriteString(pad)
		b.WriteString("  ")
		b.WriteString(entry.key)
		b.WriteString(": ")
		b.WriteString(entry.val)
		b.WriteString(";\n")
	}

	for _, child := range n.children {
		child.exportAs(b, indent+1)
	}

	b.WriteString(pad)
	b.WriteString("}\n")
}

type propEntry struct {
	key, val string
}

func sortedProps(m map[string]string) []propEntry {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sortStrings(keys)
	out := make([]propEntry, len(keys))
	for i, k := range keys {
		out[i] = propEntry{k, m[k]}
	}
	return out
}

func sortStrings(a []string) {
	for i := 0; i < len(a); i++ {
		for j := i + 1; j < len(a); j++ {
			if a[i] > a[j] {
				a[i], a[j] = a[j], a[i]
			}
		}
	}
}

func styleToCSSProps(s, d *Style) map[string]string {
	p := make(map[string]string)

	if s.Direction() != d.Direction() {
		p["direction"] = s.Direction().String()
	}
	if s.FlexDirection() != d.FlexDirection() {
		p["flex-direction"] = s.FlexDirection().String()
	}
	if s.JustifyContent() != d.JustifyContent() {
		p["justify-content"] = s.JustifyContent().String()
	}
	if s.JustifyItems() != d.JustifyItems() {
		p["justify-items"] = s.JustifyItems().String()
	}
	if s.JustifySelf() != d.JustifySelf() {
		p["justify-self"] = s.JustifySelf().String()
	}
	if s.AlignContent() != d.AlignContent() {
		p["align-content"] = s.AlignContent().String()
	}
	if s.AlignItems() != d.AlignItems() {
		p["align-items"] = s.AlignItems().String()
	}
	if s.AlignSelf() != d.AlignSelf() {
		p["align-self"] = s.AlignSelf().String()
	}
	if s.PositionType() != d.PositionType() {
		p["position"] = s.PositionType().String()
	}
	if s.FlexWrap() != d.FlexWrap() {
		p["flex-wrap"] = s.FlexWrap().String()
	}
	if s.Overflow() != d.Overflow() {
		p["overflow"] = s.Overflow().String()
	}
	if s.Display() != d.Display() {
		p["display"] = s.Display().String()
	}
	if s.BoxSizing() != d.BoxSizing() {
		p["box-sizing"] = s.BoxSizing().String()
	}

	flex := s.Flex()
	if flex.IsDefined() && !flex.Equals(d.Flex()) {
		p["flex"] = floatToStr(flex.Unwrap())
	}

	fg := s.FlexGrow()
	if fg.IsDefined() && !fg.Equals(d.FlexGrow()) {
		p["flex-grow"] = floatToStr(fg.Unwrap())
	}

	fs := s.FlexShrink()
	if fs.IsDefined() && !fs.Equals(d.FlexShrink()) {
		p["flex-shrink"] = floatToStr(fs.Unwrap())
	}

	if !s.FlexBasis().Equals(d.FlexBasis()) {
		p["flex-basis"] = styleSizeLengthToStr(s.FlexBasis())
	}

	ar := s.AspectRatio()
	if ar.IsDefined() && !ar.Equals(d.AspectRatio()) {
		p["aspect-ratio"] = floatToStr(ar.Unwrap())
	}

	dimNames := []string{"width", "height"}
	for dim := DimensionWidth; dim <= DimensionHeight; dim++ {
		v := s.Dimension(dim)
		if !v.Equals(d.Dimension(dim)) {
			p[dimNames[dim]] = styleSizeLengthToStr(v)
		}
		minV := s.MinDimension(dim)
		if !minV.Equals(d.MinDimension(dim)) {
			p["min-"+dimNames[dim]] = styleSizeLengthToStr(minV)
		}
		maxV := s.MaxDimension(dim)
		if !maxV.Equals(d.MaxDimension(dim)) {
			p["max-"+dimNames[dim]] = styleSizeLengthToStr(maxV)
		}
	}

	edges := []struct {
		edge Edge
		key  string
	}{
		{EdgeLeft, "margin-left"}, {EdgeRight, "margin-right"},
		{EdgeTop, "margin-top"}, {EdgeBottom, "margin-bottom"},
		{EdgeHorizontal, "margin-horizontal"}, {EdgeVertical, "margin-vertical"},
		{EdgeAll, "margin"},
	}
	for _, e := range edges {
		v := s.Margin(e.edge)
		if v != d.Margin(e.edge) {
			p[e.key] = styleLengthToStr(v)
		}
	}

	padEdges := []struct {
		edge Edge
		key  string
	}{
		{EdgeLeft, "padding-left"}, {EdgeRight, "padding-right"},
		{EdgeTop, "padding-top"}, {EdgeBottom, "padding-bottom"},
		{EdgeHorizontal, "padding-horizontal"}, {EdgeVertical, "padding-vertical"},
		{EdgeAll, "padding"},
	}
	for _, e := range padEdges {
		v := s.Padding(e.edge)
		if v != d.Padding(e.edge) {
			p[e.key] = styleLengthToStr(v)
		}
	}

	borderEdges := []struct {
		edge Edge
		key  string
	}{
		{EdgeLeft, "border-left"}, {EdgeRight, "border-right"},
		{EdgeTop, "border-top"}, {EdgeBottom, "border-bottom"},
		{EdgeHorizontal, "border-horizontal"}, {EdgeVertical, "border-vertical"},
		{EdgeAll, "border"},
	}
	for _, e := range borderEdges {
		v := s.Border(e.edge)
		if v != d.Border(e.edge) {
			p[e.key] = styleLengthToStr(v)
		}
	}

	gapKeys := []struct {
		gutter Gutter
		key    string
	}{
		{GutterColumn, "column-gap"},
		{GutterRow, "row-gap"},
		{GutterAll, "gap"},
	}
	for _, g := range gapKeys {
		v := s.Gap(g.gutter)
		if v != d.Gap(g.gutter) {
			p[g.key] = styleLengthToStr(v)
		}
	}

	return p
}

func styleLengthToStr(sl StyleLength) string {
	if sl.IsAuto() {
		return "auto"
	}
	if sl.IsUndefined() {
		return ""
	}
	if sl.IsPercent() {
		return fmt.Sprintf("%g%%", sl.Value().Unwrap())
	}
	return floatToStr(sl.Value().Unwrap())
}

func styleSizeLengthToStr(sl StyleSizeLength) string {
	if sl.IsAuto() {
		return "auto"
	}
	if sl.IsMaxContent() {
		return "max-content"
	}
	if sl.IsFitContent() {
		return "fit-content"
	}
	if sl.IsStretch() {
		return "stretch"
	}
	if sl.IsUndefined() {
		return ""
	}
	if sl.IsPercent() {
		return fmt.Sprintf("%g%%", sl.Value().Unwrap())
	}
	return floatToStr(sl.Value().Unwrap())
}

func floatToStr(f float32) string {
	s := strconv.FormatFloat(float64(f), 'g', -1, 32)
	if strings.Contains(s, ".") {
		s = strings.TrimRight(s, "0")
		s = strings.TrimRight(s, ".")
	}
	return s
}
