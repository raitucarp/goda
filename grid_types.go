package goda

// GridLineType describes the type of a CSS Grid line placement.
type GridLineType int

const (
	GridLineTypeAuto    GridLineType = iota
	GridLineTypeInteger
	GridLineTypeSpan
)

// GridLine represents a CSS Grid line placement value.
type GridLine struct {
	Type    GridLineType
	Integer int32
}

func GridLineAuto() GridLine {
	return GridLine{Type: GridLineTypeAuto}
}

func GridLineFromInteger(value int32) GridLine {
	return GridLine{Type: GridLineTypeInteger, Integer: value}
}

func GridLineSpan(value int32) GridLine {
	return GridLine{Type: GridLineTypeSpan, Integer: value}
}

func (g GridLine) IsAuto() bool    { return g.Type == GridLineTypeAuto }
func (g GridLine) IsInteger() bool { return g.Type == GridLineTypeInteger }
func (g GridLine) IsSpan() bool    { return g.Type == GridLineTypeSpan }

// GridTrackSize represents a CSS Grid track sizing function.
type GridTrackSize struct {
	MinSizingFunction  StyleSizeLength
	MaxSizingFunction  StyleSizeLength
	BaseSize           float32
	GrowthLimit        float32
	InfinitelyGrowable bool
}

func GridTrackSizeAuto() GridTrackSize {
	return GridTrackSize{
		MinSizingFunction: StyleSizeLengthAuto(),
		MaxSizingFunction: StyleSizeLengthAuto(),
	}
}

func GridTrackSizeLength(points float32) GridTrackSize {
	l := StyleSizeLengthPoints(points)
	return GridTrackSize{MinSizingFunction: l, MaxSizingFunction: l}
}

func GridTrackSizeFr(fraction float32) GridTrackSize {
	return GridTrackSize{
		MinSizingFunction: StyleSizeLengthAuto(),
		MaxSizingFunction: StyleSizeLengthStretch(fraction),
	}
}

func GridTrackSizePercent(percentage float32) GridTrackSize {
	return GridTrackSize{
		MinSizingFunction: StyleSizeLengthPercent(percentage),
		MaxSizingFunction: StyleSizeLengthPercent(percentage),
	}
}

func GridTrackSizeMinmax(minFn, maxFn StyleSizeLength) GridTrackSize {
	return GridTrackSize{MinSizingFunction: minFn, MaxSizingFunction: maxFn}
}

// GridTrackList is a slice of GridTrackSize representing a track listing.
type GridTrackList []GridTrackSize
