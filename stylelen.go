package goda

import "math"

// StyleLength represents a CSS length value (for margins, padding, borders, positions).
// It supports point, percent, auto, and undefined units.
type StyleLength struct {
	value FloatOptional
	unit  Unit
}

func StyleLengthUndefined() StyleLength {
	return StyleLength{FloatOptional{}, UnitUndefined}
}

func StyleLengthPoints(v float32) StyleLength {
	if isUndefined(v) || math.IsInf(float64(v), 0) {
		return StyleLengthUndefined()
	}
	return StyleLength{NewFloatOptional(v), UnitPoint}
}

func StyleLengthPercent(v float32) StyleLength {
	if isUndefined(v) || math.IsInf(float64(v), 0) {
		return StyleLengthUndefined()
	}
	return StyleLength{NewFloatOptional(v), UnitPercent}
}

func StyleLengthAuto() StyleLength {
	return StyleLength{FloatOptional{}, UnitAuto}
}

func (l StyleLength) IsAuto() bool      { return l.unit == UnitAuto }
func (l StyleLength) IsUndefined() bool { return l.unit == UnitUndefined }
func (l StyleLength) IsPoints() bool    { return l.unit == UnitPoint }
func (l StyleLength) IsPercent() bool   { return l.unit == UnitPercent }
func (l StyleLength) IsDefined() bool   { return !l.IsUndefined() }
func (l StyleLength) Value() FloatOptional { return l.value }

func (l StyleLength) Resolve(referenceLength float32) FloatOptional {
	switch l.unit {
	case UnitPoint:
		return l.value
	case UnitPercent:
		return NewFloatOptional(l.value.Unwrap() * referenceLength * 0.01)
	default:
		return FloatOptional{}
	}
}

func (l StyleLength) ToValue() Value {
	return Value{l.value.Unwrap(), l.unit}
}

// StyleSizeLength represents a CSS size value for dimensions (width, height, flex-basis).
// It supports all units including max-content, fit-content, and stretch.
type StyleSizeLength struct {
	value FloatOptional
	unit  Unit
}

func StyleSizeLengthUndefined() StyleSizeLength {
	return StyleSizeLength{FloatOptional{}, UnitUndefined}
}

func StyleSizeLengthPoints(v float32) StyleSizeLength {
	if isUndefined(v) || math.IsInf(float64(v), 0) {
		return StyleSizeLengthUndefined()
	}
	return StyleSizeLength{NewFloatOptional(v), UnitPoint}
}

func StyleSizeLengthPercent(v float32) StyleSizeLength {
	if isUndefined(v) || math.IsInf(float64(v), 0) {
		return StyleSizeLengthUndefined()
	}
	return StyleSizeLength{NewFloatOptional(v), UnitPercent}
}

func StyleSizeLengthStretch(fraction float32) StyleSizeLength {
	if isUndefined(fraction) || math.IsInf(float64(fraction), 0) {
		return StyleSizeLengthUndefined()
	}
	return StyleSizeLength{NewFloatOptional(fraction), UnitStretch}
}

func StyleSizeLengthAuto() StyleSizeLength {
	return StyleSizeLength{FloatOptional{}, UnitAuto}
}

func StyleSizeLengthMaxContent() StyleSizeLength {
	return StyleSizeLength{FloatOptional{}, UnitMaxContent}
}

func StyleSizeLengthFitContent() StyleSizeLength {
	return StyleSizeLength{FloatOptional{}, UnitFitContent}
}

func StyleSizeLengthOfStretch() StyleSizeLength {
	return StyleSizeLength{FloatOptional{}, UnitStretch}
}

func (l StyleSizeLength) IsAuto() bool       { return l.unit == UnitAuto }
func (l StyleSizeLength) IsMaxContent() bool { return l.unit == UnitMaxContent }
func (l StyleSizeLength) IsFitContent() bool { return l.unit == UnitFitContent }
func (l StyleSizeLength) IsStretch() bool    { return l.unit == UnitStretch }
func (l StyleSizeLength) IsUndefined() bool  { return l.unit == UnitUndefined }
func (l StyleSizeLength) IsDefined() bool    { return !l.IsUndefined() }
func (l StyleSizeLength) IsPoints() bool     { return l.unit == UnitPoint }
func (l StyleSizeLength) IsPercent() bool    { return l.unit == UnitPercent }
func (l StyleSizeLength) Value() FloatOptional { return l.value }

func (l StyleSizeLength) Resolve(referenceLength float32) FloatOptional {
	switch l.unit {
	case UnitPoint:
		return l.value
	case UnitPercent:
		return NewFloatOptional(l.value.Unwrap() * referenceLength * 0.01)
	default:
		return FloatOptional{}
	}
}

func (l StyleSizeLength) ToValue() Value {
	return Value{l.value.Unwrap(), l.unit}
}

func (l StyleSizeLength) Equals(other StyleSizeLength) bool {
	return l.value.Equals(other.value) && l.unit == other.unit
}
