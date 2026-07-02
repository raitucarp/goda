package goda

import "math"

const inexactEpsilon float64 = 0.0001

func inexactEqualsFloat(a, b float32) bool {
	if isDefined(a) && isDefined(b) {
		return math.Abs(float64(a-b)) < inexactEpsilon
	}
	return isUndefined(a) && isUndefined(b)
}

func inexactEqualsDouble(a, b float64) bool {
	if !math.IsNaN(a) && !math.IsNaN(b) {
		return math.Abs(a-b) < inexactEpsilon
	}
	return math.IsNaN(a) && math.IsNaN(b)
}
