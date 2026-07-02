package goda

import (
	"strconv"
	"strings"
)

// ApplyStyle applies CSS-like properties from a map.
// Keys use kebab-case (e.g. "flex-direction", "justify-content") or camelCase.
// Values are parsed as CSS values: numbers, percentages ("50%"), or keywords.
// Unknown properties are silently ignored. Returns the receiver for chaining.
//
// Example:
//
//	node.ApplyStyle(map[string]string{
//	    "display":       "flex",
//	    "flex-direction": "row",
//	    "width":         "800",
//	    "height":        "600",
//	    "padding":       "16",
//	    "gap":           "8",
//	})
func (n *Node) ApplyStyle(props map[string]string) *Node {
	for key, val := range props {
		n.applyStyleProperty(key, val)
	}
	return n
}

// ParseStyle parses a CSS-like string into a map of property-value pairs.
// Declarations are separated by ";" or newlines. Keys and values are split by ":".
// Lines starting with "//" or "/*" are treated as comments. Only supported
// properties are included in the result.
//
// Example:
//
//	props := goda.ParseStyle(`
//	    display: flex;
//	    flex-direction: row;
//	    width: 800;
//	    height: 600;
//	    padding: 16;
//	    gap: 8;
//	`)
//	node.ApplyStyle(props)
func ParseStyle(css string) map[string]string {
	props := make(map[string]string)
	lines := strings.FieldsFunc(css, func(r rune) bool { return r == ';' || r == '\n' })
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "/*") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		if key == "" || val == "" {
			continue
		}
		if !supportedCSSProperties[normalizeKey(key)] {
			continue
		}
		props[key] = val
	}
	return props
}

// ApplyStyleString parses a CSS-like string and applies the properties.
// This is a convenience combining ParseStyle and ApplyStyle. Returns the
// receiver for chaining.
//
// Example:
//
//	node.ApplyStyleString(`
//	    display: flex;
//	    flex-direction: row;
//	    width: 800;
//	    height: 600;
//	    padding: 16;
//	    gap: 8;
//	`)
func (n *Node) ApplyStyleString(css string) *Node {
	return n.ApplyStyle(ParseStyle(css))
}

// supportedCSSProperties is the set of normalized keys accepted by ApplyStyle.
var supportedCSSProperties = map[string]bool{
	"direction": true, "flexdirection": true, "flexwrap": true,
	"justifycontent": true, "justifyitems": true, "justifyself": true,
	"justify": true,
	"aligncontent": true, "alignitems": true, "alignself": true,
	"position": true, "overflow": true, "display": true, "boxsizing": true,
	"flex": true, "flexgrow": true, "flexshrink": true, "flexbasis": true,
	"width": true, "height": true, "minwidth": true, "maxwidth": true,
	"minheight": true, "maxheight": true,
	"margin": true, "margintop": true, "marginright": true, "marginbottom": true, "marginleft": true,
	"marginhorizontal": true, "marginvertical": true,
	"padding": true, "paddingtop": true, "paddingright": true, "paddingbottom": true, "paddingleft": true,
	"paddinghorizontal": true, "paddingvertical": true,
	"border": true, "bordertop": true, "borderright": true, "borderbottom": true, "borderleft": true,
	"gap": true, "columngap": true, "rowgap": true,
	"aspectratio": true,
}

func (n *Node) applyStyleProperty(key, val string) {
	k := normalizeKey(key)
	switch k {
	// Display & direction
	case "direction":
		n.SetDirection(parseDirection(val))
	case "flexdirection":
		n.SetFlexDirection(parseFlexDirection(val))
	case "flexwrap":
		n.SetFlexWrap(parseWrap(val))
	case "justifycontent":
		n.SetJustifyContent(parseJustify(val))
	case "justifyitems":
		n.SetJustifyItems(parseJustify(val))
	case "justifyself":
		n.SetJustifySelf(parseJustify(val))
	case "justify":
		n.SetJustifyContent(parseJustify(val))
	case "aligncontent":
		n.SetAlignContent(parseAlign(val))
	case "alignitems":
		n.SetAlignItems(parseAlign(val))
	case "alignself":
		n.SetAlignSelf(parseAlign(val))
	case "position":
		n.SetPositionType(parsePositionType(val))
	case "overflow":
		n.SetOverflow(parseOverflow(val))
	case "display":
		n.SetDisplay(parseDisplay(val))
	case "boxsizing":
		n.SetBoxSizing(parseBoxSizing(val))

	// Flex factors
	case "flex":
		n.SetFlex(parseFloat(val))
	case "flexgrow":
		n.SetFlexGrow(parseFloat(val))
	case "flexshrink":
		n.SetFlexShrink(parseFloat(val))
	case "flexbasis":
		n.applyFlexBasis(val)

	// Dimensions
	case "width":
		n.applyDimension(val, DimensionWidth, false, false)
	case "height":
		n.applyDimension(val, DimensionHeight, false, false)
	case "minwidth":
		n.applyDimension(val, DimensionWidth, true, false)
	case "maxwidth":
		n.applyDimension(val, DimensionWidth, false, true)
	case "minheight":
		n.applyDimension(val, DimensionHeight, true, false)
	case "maxheight":
		n.applyDimension(val, DimensionHeight, false, true)

	// Spacing
	case "margin":
		n.SetMargin(EdgeAll, n.resolveCSSValue(val))
	case "margintop":
		n.SetMargin(EdgeTop, n.resolveCSSValue(val))
	case "marginright":
		n.SetMargin(EdgeRight, n.resolveCSSValue(val))
	case "marginbottom":
		n.SetMargin(EdgeBottom, n.resolveCSSValue(val))
	case "marginleft":
		n.SetMargin(EdgeLeft, n.resolveCSSValue(val))
	case "marginhorizontal":
		n.SetMargin(EdgeHorizontal, n.resolveCSSValue(val))
	case "marginvertical":
		n.SetMargin(EdgeVertical, n.resolveCSSValue(val))
	case "padding":
		n.SetPadding(EdgeAll, n.resolveCSSValue(val))
	case "paddingtop":
		n.SetPadding(EdgeTop, n.resolveCSSValue(val))
	case "paddingright":
		n.SetPadding(EdgeRight, n.resolveCSSValue(val))
	case "paddingbottom":
		n.SetPadding(EdgeBottom, n.resolveCSSValue(val))
	case "paddingleft":
		n.SetPadding(EdgeLeft, n.resolveCSSValue(val))
	case "paddinghorizontal":
		n.SetPadding(EdgeHorizontal, n.resolveCSSValue(val))
	case "paddingvertical":
		n.SetPadding(EdgeVertical, n.resolveCSSValue(val))

	// Border
	case "border":
		n.SetBorder(EdgeAll, n.resolveCSSValue(val))
	case "bordertop":
		n.SetBorder(EdgeTop, n.resolveCSSValue(val))
	case "borderright":
		n.SetBorder(EdgeRight, n.resolveCSSValue(val))
	case "borderbottom":
		n.SetBorder(EdgeBottom, n.resolveCSSValue(val))
	case "borderleft":
		n.SetBorder(EdgeLeft, n.resolveCSSValue(val))

	// Gap
	case "gap":
		n.SetGap(GutterAll, n.resolveCSSValue(val))
	case "columngap":
		n.SetGap(GutterColumn, n.resolveCSSValue(val))
	case "rowgap":
		n.SetGap(GutterRow, n.resolveCSSValue(val))

	// Aspect ratio
	case "aspectratio":
		n.SetAspectRatio(parseFloat(val))
	}
}

func normalizeKey(key string) string {
	return strings.ToLower(strings.ReplaceAll(key, "-", ""))
}

// --- Value parsers ---

func parseFloat(s string) float32 {
	s = strings.TrimSpace(s)
	s = strings.TrimSuffix(s, "%")
	s = strings.TrimSuffix(s, "px")
	s = strings.TrimSuffix(s, "rem")
	s = strings.TrimSuffix(s, "em")
	s = strings.TrimSpace(s)
	v, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0
	}
	return float32(v)
}

func parseBool(s string) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "true", "1", "yes":
		return true
	}
	return false
}

func isPercent(s string) bool {
	return strings.HasSuffix(strings.TrimSpace(s), "%")
}

func isAuto(s string) bool {
	return strings.ToLower(strings.TrimSpace(s)) == "auto"
}

func isMaxContent(s string) bool {
	return strings.ToLower(strings.TrimSpace(s)) == "max-content"
}

func isFitContent(s string) bool {
	return strings.ToLower(strings.TrimSpace(s)) == "fit-content"
}

func isStretch(s string) bool {
	s = strings.ToLower(strings.TrimSpace(s))
	return s == "stretch"
}

func parseDirection(s string) Direction {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "ltr":
		return DirectionLTR
	case "rtl":
		return DirectionRTL
	case "inherit":
		return DirectionInherit
	}
	return DirectionInherit
}

func parseFlexDirection(s string) FlexDirection {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "row":
		return FlexDirectionRow
	case "row-reverse":
		return FlexDirectionRowReverse
	case "column":
		return FlexDirectionColumn
	case "column-reverse":
		return FlexDirectionColumnReverse
	}
	return FlexDirectionColumn
}

func parseJustify(s string) Justify {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "flex-start":
		return JustifyFlexStart
	case "center":
		return JustifyCenter
	case "flex-end":
		return JustifyFlexEnd
	case "space-between":
		return JustifySpaceBetween
	case "space-around":
		return JustifySpaceAround
	case "space-evenly":
		return JustifySpaceEvenly
	case "stretch":
		return JustifyStretch
	case "start":
		return JustifyStart
	case "end":
		return JustifyEnd
	case "auto":
		return JustifyAuto
	}
	return JustifyFlexStart
}

func parseAlign(s string) Align {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "auto":
		return AlignAuto
	case "flex-start":
		return AlignFlexStart
	case "center":
		return AlignCenter
	case "flex-end":
		return AlignFlexEnd
	case "stretch":
		return AlignStretch
	case "baseline":
		return AlignBaseline
	case "space-between":
		return AlignSpaceBetween
	case "space-around":
		return AlignSpaceAround
	case "space-evenly":
		return AlignSpaceEvenly
	case "start":
		return AlignStart
	case "end":
		return AlignEnd
	}
	return AlignFlexStart
}

func parseWrap(s string) Wrap {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "wrap":
		return WrapWrap
	case "wrap-reverse":
		return WrapWrapReverse
	case "nowrap", "no-wrap":
		return WrapNoWrap
	}
	return WrapNoWrap
}

func parsePositionType(s string) PositionType {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "absolute":
		return PositionTypeAbsolute
	case "relative":
		return PositionTypeRelative
	case "static":
		return PositionTypeStatic
	}
	return PositionTypeRelative
}

func parseOverflow(s string) Overflow {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "visible":
		return OverflowVisible
	case "hidden":
		return OverflowHidden
	case "scroll":
		return OverflowScroll
	}
	return OverflowVisible
}

func parseDisplay(s string) Display {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "flex":
		return DisplayFlex
	case "none":
		return DisplayNone
	case "contents":
		return DisplayContents
	case "grid":
		return DisplayGrid
	}
	return DisplayFlex
}

func parseBoxSizing(s string) BoxSizing {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "border-box":
		return BoxSizingBorderBox
	case "content-box":
		return BoxSizingContentBox
	}
	return BoxSizingBorderBox
}

func (n *Node) applyFlexBasis(val string) {
	val = strings.TrimSpace(val)
	switch {
	case isAuto(val):
		n.SetFlexBasisAuto()
	case isMaxContent(val):
		n.SetFlexBasisMaxContent()
	case isFitContent(val):
		n.SetFlexBasisFitContent()
	case isStretch(val):
		n.SetFlexBasisStretch()
	case isPercent(val):
		n.SetFlexBasisPercent(parseFloat(val))
	default:
		n.SetFlexBasis(n.resolveCSSValue(val))
	}
}

func (n *Node) applyDimension(val string, dim Dimension, isMin, isMax bool) {
	val = strings.TrimSpace(val)
	switch {
	case isAuto(val):
		if isMax {
			return
		}
		if isMin {
			if dim == DimensionWidth {
				n.SetMinWidth(0)
			} else {
				n.SetMinHeight(0)
			}
		} else {
			if dim == DimensionWidth {
				n.SetWidthAuto()
			} else {
				n.SetHeightAuto()
			}
		}
	case isMaxContent(val):
		if dim == DimensionWidth {
			if isMin {
				n.SetMinWidthMaxContent()
			} else if isMax {
				n.SetMaxWidthMaxContent()
			} else {
				n.SetWidthMaxContent()
			}
		} else {
			if isMin {
				n.SetMinHeightMaxContent()
			} else if isMax {
				n.SetMaxHeightMaxContent()
			} else {
				n.SetHeightMaxContent()
			}
		}
	case isFitContent(val):
		if dim == DimensionWidth {
			if isMin {
				n.SetMinWidthFitContent()
			} else if isMax {
				n.SetMaxWidthFitContent()
			} else {
				n.SetWidthFitContent()
			}
		} else {
			if isMin {
				n.SetMinHeightFitContent()
			} else if isMax {
				n.SetMaxHeightFitContent()
			} else {
				n.SetHeightFitContent()
			}
		}
	case isStretch(val):
		if dim == DimensionWidth {
			if isMin {
				n.SetMinWidthStretch()
			} else if isMax {
				n.SetMaxWidthStretch()
			} else {
				n.SetWidthStretch()
			}
		} else {
			if isMin {
				n.SetMinHeightStretch()
			} else if isMax {
				n.SetMaxHeightStretch()
			} else {
				n.SetHeightStretch()
			}
		}
	case isPercent(val):
		v := parseFloat(val)
		if dim == DimensionWidth {
			if isMin {
				n.SetMinWidthPercent(v)
			} else if isMax {
				n.SetMaxWidthPercent(v)
			} else {
				n.SetWidthPercent(v)
			}
		} else {
			if isMin {
				n.SetMinHeightPercent(v)
			} else if isMax {
				n.SetMaxHeightPercent(v)
			} else {
				n.SetHeightPercent(v)
			}
		}
	default:
		v := n.resolveCSSValue(val)
		if dim == DimensionWidth {
			if isMin {
				n.SetMinWidth(v)
			} else if isMax {
				n.SetMaxWidth(v)
			} else {
				n.SetWidth(v)
			}
		} else {
			if isMin {
				n.SetMinHeight(v)
			} else if isMax {
				n.SetMaxHeight(v)
			} else {
				n.SetHeight(v)
			}
		}
	}
}
