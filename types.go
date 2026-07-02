package goda

import "math"

// Undefined is the NaN sentinel value used throughout the layout engine
// to represent unset or "auto" dimensions.
var Undefined = float32(math.NaN())

// IsUndefinedFloat returns true if v is the NaN sentinel value.
func IsUndefinedFloat(v float32) bool { return math.IsNaN(float64(v)) }

func isUndefined(v float32) bool { return math.IsNaN(float64(v)) }
func isDefined(v float32) bool   { return !isUndefined(v) }

func maxOrDefined(a, b float32) float32 {
	if isDefined(a) && isDefined(b) {
		return max(a, b)
	}
	if isUndefined(a) {
		return b
	}
	return a
}

func minOrDefined(a, b float32) float32 {
	if isDefined(a) && isDefined(b) {
		return min(a, b)
	}
	if isUndefined(a) {
		return b
	}
	return a
}

// Align represents the CSS align-items/align-self/align-content values.
type Align int

const (
	AlignAuto Align = iota
	AlignFlexStart
	AlignCenter
	AlignFlexEnd
	AlignStretch
	AlignBaseline
	AlignSpaceBetween
	AlignSpaceAround
	AlignSpaceEvenly
	AlignStart
	AlignEnd
)

func (a Align) String() string {
	switch a {
	case AlignAuto:
		return "auto"
	case AlignFlexStart:
		return "flex-start"
	case AlignCenter:
		return "center"
	case AlignFlexEnd:
		return "flex-end"
	case AlignStretch:
		return "stretch"
	case AlignBaseline:
		return "baseline"
	case AlignSpaceBetween:
		return "space-between"
	case AlignSpaceAround:
		return "space-around"
	case AlignSpaceEvenly:
		return "space-evenly"
	case AlignStart:
		return "start"
	case AlignEnd:
		return "end"
	}
	return "unknown"
}

// BoxSizing represents the CSS box-sizing property.
type BoxSizing int

const (
	BoxSizingBorderBox  BoxSizing = iota
	BoxSizingContentBox
)

func (b BoxSizing) String() string {
	switch b {
	case BoxSizingBorderBox:
		return "border-box"
	case BoxSizingContentBox:
		return "content-box"
	}
	return "unknown"
}

// Dimension represents width or height axis.
type Dimension int

const (
	DimensionWidth  Dimension = iota
	DimensionHeight
)

func (d Dimension) String() string {
	switch d {
	case DimensionWidth:
		return "width"
	case DimensionHeight:
		return "height"
	}
	return "unknown"
}

// Direction represents the text direction (LTR/RTL).
type Direction int

const (
	DirectionInherit Direction = iota
	DirectionLTR
	DirectionRTL
)

func (d Direction) String() string {
	switch d {
	case DirectionInherit:
		return "inherit"
	case DirectionLTR:
		return "ltr"
	case DirectionRTL:
		return "rtl"
	}
	return "unknown"
}

// Display represents the CSS display property.
type Display int

const (
	DisplayFlex     Display = iota
	DisplayNone
	DisplayContents
	DisplayGrid
)

func (d Display) String() string {
	switch d {
	case DisplayFlex:
		return "flex"
	case DisplayNone:
		return "none"
	case DisplayContents:
		return "contents"
	case DisplayGrid:
		return "grid"
	}
	return "unknown"
}

// Edge represents a CSS edge (left, top, right, bottom, start, end, etc.).
type Edge int

const (
	EdgeLeft       Edge = iota
	EdgeTop
	EdgeRight
	EdgeBottom
	EdgeStart
	EdgeEnd
	EdgeHorizontal
	EdgeVertical
	EdgeAll
)

func (e Edge) String() string {
	switch e {
	case EdgeLeft:
		return "left"
	case EdgeTop:
		return "top"
	case EdgeRight:
		return "right"
	case EdgeBottom:
		return "bottom"
	case EdgeStart:
		return "start"
	case EdgeEnd:
		return "end"
	case EdgeHorizontal:
		return "horizontal"
	case EdgeVertical:
		return "vertical"
	case EdgeAll:
		return "all"
	}
	return "unknown"
}

// Errata is a bitmask of legacy behavior flags.
type Errata int

const (
	ErrataNone                                        Errata = 0
	ErrataStretchFlexBasis                            Errata = 1 << 0
	ErrataAbsolutePositionWithoutInsetsExcludesPadding Errata = 1 << 1
	ErrataAbsolutePercentAgainstInnerSize             Errata = 1 << 2
	ErrataMinSizeUndefinedInsteadOfAuto               Errata = 1 << 3
	ErrataAll                                         Errata = 1<<31 - 1
	ErrataClassic                                     Errata = ErrataAll & ^ErrataMinSizeUndefinedInsteadOfAuto
)

// ExperimentalFeature represents feature flags for optional behavior.
type ExperimentalFeature int

const (
	ExperimentalFeatureWebFlexBasis         ExperimentalFeature = iota
	ExperimentalFeatureFixFlexBasisFitContent
)

// FlexDirection represents the CSS flex-direction property.
type FlexDirection int

const (
	FlexDirectionColumn        FlexDirection = iota
	FlexDirectionColumnReverse
	FlexDirectionRow
	FlexDirectionRowReverse
)

func (f FlexDirection) String() string {
	switch f {
	case FlexDirectionColumn:
		return "column"
	case FlexDirectionColumnReverse:
		return "column-reverse"
	case FlexDirectionRow:
		return "row"
	case FlexDirectionRowReverse:
		return "row-reverse"
	}
	return "unknown"
}

// GridTrackType describes a CSS Grid track sizing function type.
type GridTrackType int

const (
	GridTrackTypeAuto    GridTrackType = iota
	GridTrackTypePoints
	GridTrackTypePercent
	GridTrackTypeFr
	GridTrackTypeMinmax
)

func (g GridTrackType) String() string {
	switch g {
	case GridTrackTypeAuto:
		return "auto"
	case GridTrackTypePoints:
		return "points"
	case GridTrackTypePercent:
		return "percent"
	case GridTrackTypeFr:
		return "fr"
	case GridTrackTypeMinmax:
		return "minmax"
	}
	return "unknown"
}

// Gutter represents a CSS Grid gap axis.
type Gutter int

const (
	GutterColumn Gutter = iota
	GutterRow
	GutterAll
)

func (g Gutter) String() string {
	switch g {
	case GutterColumn:
		return "column"
	case GutterRow:
		return "row"
	case GutterAll:
		return "all"
	}
	return "unknown"
}

// Justify represents CSS justify-content/justify-items/justify-self values.
type Justify int

const (
	JustifyAuto          Justify = iota
	JustifyFlexStart
	JustifyCenter
	JustifyFlexEnd
	JustifySpaceBetween
	JustifySpaceAround
	JustifySpaceEvenly
	JustifyStretch
	JustifyStart
	JustifyEnd
)

func (j Justify) String() string {
	switch j {
	case JustifyAuto:
		return "auto"
	case JustifyFlexStart:
		return "flex-start"
	case JustifyCenter:
		return "center"
	case JustifyFlexEnd:
		return "flex-end"
	case JustifySpaceBetween:
		return "space-between"
	case JustifySpaceAround:
		return "space-around"
	case JustifySpaceEvenly:
		return "space-evenly"
	case JustifyStretch:
		return "stretch"
	case JustifyStart:
		return "start"
	case JustifyEnd:
		return "end"
	}
	return "unknown"
}

// LogLevel represents the severity of a log message.
type LogLevel int

const (
	LogLevelError   LogLevel = iota
	LogLevelWarn
	LogLevelInfo
	LogLevelDebug
	LogLevelVerbose
	LogLevelFatal
)

func (l LogLevel) String() string {
	switch l {
	case LogLevelError:
		return "error"
	case LogLevelWarn:
		return "warn"
	case LogLevelInfo:
		return "info"
	case LogLevelDebug:
		return "debug"
	case LogLevelVerbose:
		return "verbose"
	case LogLevelFatal:
		return "fatal"
	}
	return "unknown"
}

// MeasureMode describes how a measurement constraint is applied.
type MeasureMode int

const (
	MeasureModeUndefined MeasureMode = iota
	MeasureModeExactly
	MeasureModeAtMost
)

func (m MeasureMode) String() string {
	switch m {
	case MeasureModeUndefined:
		return "undefined"
	case MeasureModeExactly:
		return "exactly"
	case MeasureModeAtMost:
		return "at-most"
	}
	return "unknown"
}

// NodeType categorizes a node (default container or text leaf).
type NodeType int

const (
	NodeTypeDefault NodeType = iota
	NodeTypeText
)

func (n NodeType) String() string {
	switch n {
	case NodeTypeDefault:
		return "default"
	case NodeTypeText:
		return "text"
	}
	return "unknown"
}

// Overflow represents the CSS overflow property.
type Overflow int

const (
	OverflowVisible Overflow = iota
	OverflowHidden
	OverflowScroll
)

func (o Overflow) String() string {
	switch o {
	case OverflowVisible:
		return "visible"
	case OverflowHidden:
		return "hidden"
	case OverflowScroll:
		return "scroll"
	}
	return "unknown"
}

// PhysicalEdge represents a fixed physical edge (not logical).
type PhysicalEdge int

const (
	PhysicalEdgeLeft   PhysicalEdge = iota
	PhysicalEdgeTop
	PhysicalEdgeRight
	PhysicalEdgeBottom
)

// PositionType represents the CSS position property.
type PositionType int

const (
	PositionTypeStatic   PositionType = iota
	PositionTypeRelative
	PositionTypeAbsolute
)

func (p PositionType) String() string {
	switch p {
	case PositionTypeStatic:
		return "static"
	case PositionTypeRelative:
		return "relative"
	case PositionTypeAbsolute:
		return "absolute"
	}
	return "unknown"
}

// Unit represents a CSS length unit.
type Unit int

const (
	UnitUndefined  Unit = iota
	UnitPoint
	UnitPercent
	UnitAuto
	UnitMaxContent
	UnitFitContent
	UnitStretch
)

func (u Unit) String() string {
	switch u {
	case UnitUndefined:
		return "undefined"
	case UnitPoint:
		return "point"
	case UnitPercent:
		return "percent"
	case UnitAuto:
		return "auto"
	case UnitMaxContent:
		return "max-content"
	case UnitFitContent:
		return "fit-content"
	case UnitStretch:
		return "stretch"
	}
	return "unknown"
}

// Wrap represents the CSS flex-wrap property.
type Wrap int

const (
	WrapNoWrap     Wrap = iota
	WrapWrap
	WrapWrapReverse
)

func (w Wrap) String() string {
	switch w {
	case WrapNoWrap:
		return "no-wrap"
	case WrapWrap:
		return "wrap"
	case WrapWrapReverse:
		return "wrap-reverse"
	}
	return "unknown"
}

// SizingMode controls how dimensions are resolved during layout.
type SizingMode int

const (
	SizingModeStretchFit SizingMode = iota
	SizingModeMaxContent
	SizingModeFitContent
)

func measureMode(m SizingMode) MeasureMode {
	switch m {
	case SizingModeStretchFit:
		return MeasureModeExactly
	case SizingModeMaxContent:
		return MeasureModeUndefined
	case SizingModeFitContent:
		return MeasureModeAtMost
	}
	panic("Invalid SizingMode")
}

func sizingMode(m MeasureMode) SizingMode {
	switch m {
	case MeasureModeExactly:
		return SizingModeStretchFit
	case MeasureModeUndefined:
		return SizingModeMaxContent
	case MeasureModeAtMost:
		return SizingModeFitContent
	}
	panic("Invalid MeasureMode")
}

func isRow(fd FlexDirection) bool {
	return fd == FlexDirectionRow || fd == FlexDirectionRowReverse
}

func isColumn(fd FlexDirection) bool {
	return fd == FlexDirectionColumn || fd == FlexDirectionColumnReverse
}

func dimension(fd FlexDirection) Dimension {
	if isRow(fd) {
		return DimensionWidth
	}
	return DimensionHeight
}

func edgeCount() int { return 9 }

func gutterCount() int { return 3 }
