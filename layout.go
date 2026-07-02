package goda

import "math"

type cachedMeasurement struct {
	availableWidth   float32
	availableHeight  float32
	widthSizingMode  SizingMode
	heightSizingMode SizingMode
	computedWidth    float32
	computedHeight   float32
}

const maxCachedMeasurements = 8

// LayoutResults holds the computed layout output for a Node.
type LayoutResults struct {
	computedFlexBasisGeneration uint32
	computedFlexBasis           FloatOptional
	computedAutoMinMainSize     FloatOptional
	generationCount             uint32
	configVersion               uint32
	lastOwnerDirection           Direction
	nextCachedMeasurementsIndex uint32
	cachedMeasurements          [maxCachedMeasurements]cachedMeasurement
	cachedLayout                cachedMeasurement

	direction_  Direction
	hadOverflow bool

	dimensions_         [2]float32
	measuredDimensions_ [2]float32
	rawDimensions_      [2]float32
	position_           [4]float32
	margin_             [4]float32
	border_             [4]float32
	padding_            [4]float32
}

func NewLayoutResults() LayoutResults {
	lr := LayoutResults{
		direction_:             DirectionInherit,
		configVersion:          0,
		lastOwnerDirection:      DirectionInherit,
		computedFlexBasis:      FloatOptional{},
		computedAutoMinMainSize: FloatOptional{},
	}
	lr.cachedLayout.availableWidth = -1
	lr.cachedLayout.availableHeight = -1
	lr.cachedLayout.widthSizingMode = SizingModeMaxContent
	lr.cachedLayout.heightSizingMode = SizingModeMaxContent
	lr.cachedLayout.computedWidth = -1
	lr.cachedLayout.computedHeight = -1
	lr.dimensions_[0] = float32(math.NaN())
	lr.dimensions_[1] = float32(math.NaN())
	lr.measuredDimensions_[0] = float32(math.NaN())
	lr.measuredDimensions_[1] = float32(math.NaN())
	lr.rawDimensions_[0] = float32(math.NaN())
	lr.rawDimensions_[1] = float32(math.NaN())
	return lr
}

func (l *LayoutResults) Direction() Direction                    { return l.direction_ }
func (l *LayoutResults) SetDirection(d Direction)                { l.direction_ = d }
func (l *LayoutResults) HadOverflow() bool                       { return l.hadOverflow }
func (l *LayoutResults) SetHadOverflow(v bool)                   { l.hadOverflow = v }

func (l *LayoutResults) Dimension(axis Dimension) float32        { return l.dimensions_[axis] }
func (l *LayoutResults) SetDimension(axis Dimension, v float32)  { l.dimensions_[axis] = v }

func (l *LayoutResults) MeasuredDimension(axis Dimension) float32 { return l.measuredDimensions_[axis] }
func (l *LayoutResults) SetMeasuredDimension(axis Dimension, v float32) { l.measuredDimensions_[axis] = v }

func (l *LayoutResults) RawDimension(axis Dimension) float32     { return l.rawDimensions_[axis] }
func (l *LayoutResults) SetRawDimension(axis Dimension, v float32) { l.rawDimensions_[axis] = v }

func (l *LayoutResults) Position(edge PhysicalEdge) float32       { return l.position_[edge] }
func (l *LayoutResults) SetPosition(edge PhysicalEdge, v float32) { l.position_[edge] = v }

func (l *LayoutResults) Margin(edge PhysicalEdge) float32         { return l.margin_[edge] }
func (l *LayoutResults) SetMargin(edge PhysicalEdge, v float32)   { l.margin_[edge] = v }

func (l *LayoutResults) Border(edge PhysicalEdge) float32         { return l.border_[edge] }
func (l *LayoutResults) SetBorder(edge PhysicalEdge, v float32)   { l.border_[edge] = v }

func (l *LayoutResults) Padding(edge PhysicalEdge) float32        { return l.padding_[edge] }
func (l *LayoutResults) SetPadding(edge PhysicalEdge, v float32)  { l.padding_[edge] = v }
