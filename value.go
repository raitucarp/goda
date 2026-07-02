package goda

// Value is a CSS-style value with a numeric component and a unit.
type Value struct {
	Value float32
	Unit  Unit
}

var (
	ValueZero      = Value{0, UnitPoint}
	ValueUndefined = Value{Undefined, UnitUndefined}
	ValueAuto      = Value{Undefined, UnitAuto}
)

// Size is a simple width/height pair used by measure callbacks.
type Size struct {
	Width  float32
	Height float32
}

// FloatOptional represents an optional float32 value that may be undefined.
type FloatOptional struct {
	val     float32
	defined bool
}

func NewFloatOptional(v float32) FloatOptional {
	if isUndefined(v) {
		return FloatOptional{}
	}
	return FloatOptional{val: v, defined: true}
}

func (f FloatOptional) Unwrap() float32 {
	return f.val
}

func (f FloatOptional) UnwrapOrDefault(defaultValue float32) float32 {
	if !f.defined {
		return defaultValue
	}
	return f.val
}

func (f FloatOptional) IsUndefined() bool {
	return !f.defined
}

func (f FloatOptional) IsDefined() bool {
	return f.defined
}

func (f FloatOptional) Equals(other FloatOptional) bool {
	if !f.defined && !other.defined {
		return true
	}
	if f.defined != other.defined {
		return false
	}
	return f.val == other.val
}
