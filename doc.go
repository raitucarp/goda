// Package goda is a CSS Flexbox layout engine written in pure Go.
// It calculates positions and dimensions for UI elements based on
// CSS Flexbox properties such as flex-direction, justify-content,
// align-items, and more.
//
// # Builder Pattern
//
// All property setters return the receiver (*Node), enabling a fluent
// builder pattern for constructing layout trees:
//
//	root := goda.New().
//	    SetWidth(800).
//	    SetHeight(600).
//	    SetFlexDirection(goda.FlexDirectionRow).
//	    SetJustifyContent(goda.JustifySpaceBetween).
//	    SetAlignItems(goda.AlignCenter).
//	    SetPadding(goda.EdgeAll, 16).
//	    SetGap(goda.GutterAll, 8)
//
//	child := goda.New().
//	    SetWidth(100).
//	    SetHeight(50).
//	    SetFlexGrow(1).
//	    SetMargin(goda.EdgeAll, 8).
//	    SetAlignSelf(goda.AlignCenter)
//
//	root.InsertChildNode(child, 0)
//	goda.CalculateNodeLayout(root, 800, 600, goda.DirectionLTR)
//
//	fmt.Println(child.GetLeft(), child.GetTop())
//	fmt.Println(child.GetWidth(), child.GetHeight())
//
// # Consuming Layout Results
//
// After CalculateNodeLayout, use LayoutOut() to get all computed layout values
// in a single struct designed for GUI library consumption:
//
//	goda.CalculateNodeLayout(root, 800, 600, goda.DirectionLTR)
//	lo := root.LayoutOut()
//
//	// All at once:
//	renderer.DrawBox(lo.Left, lo.Top, lo.Width, lo.Height)
//	renderer.SetMargins(lo.Margin.Top, lo.Margin.Right, lo.Margin.Bottom, lo.Margin.Left)
//
//	// Individual accessors still work:
//	fmt.Printf("Pos:(%f,%f) Size:%fx%f Overflow:%v Dir:%v\n",
//	    lo.Left, lo.Top, lo.Width, lo.Height, lo.HadOverflow, lo.Direction)
//
//	// Child layout:
//	childLo := child.LayoutOut()
//	renderer.DrawBox(childLo.Left, childLo.Top, childLo.Width, childLo.Height)
//
// # CSS String Properties
//
// Use ParseStyle to convert a CSS-like string into a map, or ApplyStyleString
// to parse and apply in one call:
//
//	css := `
//	    display: flex;
//	    flex-direction: row;
//	    width: 800;
//	    height: 600;
//	    padding: 16;
//	    gap: 8;
//	`
//	root := goda.New().ApplyStyleString(css)
//
// // Or parse first, inspect, then apply:
//	props := goda.ParseStyle(css)
//	root.ApplyStyle(props)
//
// Declarations use "key: value" syntax separated by ";" or newlines.
// Lines starting with "//" or "/*" are treated as comments and ignored.
// Unknown CSS properties (e.g. "color", "font-size") are silently skipped.
//
// Length values support px, rem, and em units. rem resolves against the
// root node's font size estimate; em resolves against the node's own estimate
// (default 16 for both). Use SetFontSizeEstimate to customize:
//
//	root := goda.New().SetFontSizeEstimate(14)
//	child := goda.New().
//	    ApplyStyleString("width: 10rem; padding: 2em;").
//	    SetFontSizeEstimate(12) // em=12px here, rem=14px from root
//
// # CSS Map Properties
//
// Use ApplyStyle with a map[string]string to set multiple properties at once:
//
//	root := goda.New().ApplyStyle(map[string]string{
//	    "display":         "flex",
//	    "flex-direction":  "row",
//	    "justify-content": "space-between",
//	    "align-items":     "center",
//	    "width":           "800",
//	    "height":          "600",
//	    "padding":         "16",
//	    "gap":             "8",
//	})
//
// All three APIs chain seamlessly with the builder pattern:
//
//	child := goda.New().
//	    ApplyStyleString("width: 100; height: 50;").
//	    SetFlexGrow(1).
//	    ApplyStyle(map[string]string{"align-self": "center"})
//
// # Supported CSS Properties
//
// Layout:
//
//	display        "flex" | "none" | "contents" | "grid"
//	direction      "ltr" | "rtl" | "inherit"
//	position       "static" | "relative" | "absolute"
//	overflow       "visible" | "hidden" | "scroll"
//	box-sizing     "border-box" | "content-box"
//
// Flex:
//
//	flex-direction  "row" | "row-reverse" | "column" | "column-reverse"
//	flex-wrap       "nowrap" | "wrap" | "wrap-reverse"
//	justify-content "flex-start" | "center" | "flex-end" | "space-between" |
//	                "space-around" | "space-evenly" | "start" | "end"
//	justify-items   same as justify-content + "stretch" | "auto"
//	justify-self    same as justify-content + "stretch" | "auto"
//	align-content   same as align-items
//	align-items     "flex-start" | "center" | "flex-end" | "stretch" |
//	                "baseline" | "start" | "end" | "auto"
//	align-self      same as align-items + "auto"
//
// Flex factors:
//
//	flex        number
//	flex-grow   number
//	flex-shrink number
//	flex-basis  number | "auto" | number% | "max-content" | "fit-content" | "stretch"
//
// Dimensions:
//
//	width      number | "auto" | number% | numberpx | "max-content" | "fit-content" | "stretch"
//	height     same as width
//	min-width  same as width
//	max-width  same as width
//	min-height same as height
//	max-height same as height
//
// Spacing:
//
//	margin             number
//	margin-top         number
//	margin-right       number
//	margin-bottom      number
//	margin-left        number
//	margin-horizontal  number
//	margin-vertical    number
//	padding            number
//	padding-top        number
//	padding-right      number
//	padding-bottom     number
//	padding-left       number
//	padding-horizontal number
//	padding-vertical   number
//
// Border & Gap:
//
//	border         number
//	border-top     number
//	border-right   number
//	border-bottom  number
//	border-left    number
//	gap            number
//	column-gap     number
//	row-gap        number
//
// Other:
//
//	aspect-ratio number
//
// # Complete Example
//
//	config := goda.ConfigNewDefault()
//	config.SetPointScaleFactor(2.0)
//
//	root := goda.NewWithConfig(config).ApplyStyleString(`
//	    width: 800;
//	    height: 600;
//	    flex-direction: column;
//	    justify-content: center;
//	    align-items: stretch;
//	    padding: 20;
//	    gap: 12;
//	`)
//
//	header := goda.New().ApplyStyleString("height: 60;")
//
//	body := goda.New().
//	    SetFlexGrow(1).
//	    SetFlexDirection(goda.FlexDirectionRow).
//	    SetGap(goda.GutterAll, 16)
//
//	sidebar := goda.New().
//	    ApplyStyle(map[string]string{"width": "200"}).
//	    SetFlexShrink(0)
//
//	content := goda.New().
//	    SetFlexGrow(1).
//	    SetMinWidth(300)
//
//	root.InsertChildNode(header, 0)
//	root.InsertChildNode(body, 1)
//	body.InsertChildNode(sidebar, 0)
//	body.InsertChildNode(content, 1)
//
//	goda.CalculateNodeLayout(root, 800, 600, goda.DirectionLTR)
package goda
