package goda

const (
	LayoutPassInitial         = iota
	LayoutPassAbsLayout
	LayoutPassStretch
	LayoutPassMultilineStretch
	LayoutPassFlexLayout
	LayoutPassMeasureChild
	LayoutPassAbsMeasureChild
	LayoutPassFlexMeasure
	LayoutPassGridLayout
	LayoutPassCount
)

// LayoutData tracks layout performance counters.
type LayoutData struct {
	Layouts                int
	Measures               int
	MaxMeasureCache        uint32
	CachedLayouts          int
	CachedMeasures         int
	MeasureCallbacks       int
	MeasureCallbackReasons [LayoutPassCount]int
}
