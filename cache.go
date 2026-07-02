package goda

func canUseCachedMeasurement(
	widthMode SizingMode, availableWidth float32,
	heightMode SizingMode, availableHeight float32,
	lastWidthMode SizingMode, lastAvailableWidth float32,
	lastHeightMode SizingMode, lastAvailableHeight float32,
	lastComputedWidth, lastComputedHeight float32,
	marginRow, marginColumn float32,
	config *Config) bool {

	if (isDefined(lastComputedHeight) && lastComputedHeight < 0) ||
		(isDefined(lastComputedWidth) && lastComputedWidth < 0) {
		return false
	}

	pointScaleFactor := config.GetPointScaleFactor()
	useRounded := pointScaleFactor != 0

	effWidth := availableWidth
	effHeight := availableHeight
	effLastWidth := lastAvailableWidth
	effLastHeight := lastAvailableHeight

	if useRounded {
		effWidth = roundValueToPixelGrid(float64(availableWidth), float64(pointScaleFactor), false, false)
		effHeight = roundValueToPixelGrid(float64(availableHeight), float64(pointScaleFactor), false, false)
		effLastWidth = roundValueToPixelGrid(float64(lastAvailableWidth), float64(pointScaleFactor), false, false)
		effLastHeight = roundValueToPixelGrid(float64(lastAvailableHeight), float64(pointScaleFactor), false, false)
	}

	hasSameWidthSpec := lastWidthMode == widthMode && inexactEqualsFloat(effLastWidth, effWidth)
	hasSameHeightSpec := lastHeightMode == heightMode && inexactEqualsFloat(effLastHeight, effHeight)

	widthCompatible := hasSameWidthSpec ||
		sizeIsExactAndMatches(widthMode, availableWidth-marginRow, lastComputedWidth) ||
		oldSizeIsMaxContent(widthMode, availableWidth-marginRow, lastWidthMode, lastComputedWidth) ||
		newSizeIsStricter(widthMode, availableWidth-marginRow, lastWidthMode, lastAvailableWidth, lastComputedWidth)

	heightCompatible := hasSameHeightSpec ||
		sizeIsExactAndMatches(heightMode, availableHeight-marginColumn, lastComputedHeight) ||
		oldSizeIsMaxContent(heightMode, availableHeight-marginColumn, lastHeightMode, lastComputedHeight) ||
		newSizeIsStricter(heightMode, availableHeight-marginColumn, lastHeightMode, lastAvailableHeight, lastComputedHeight)

	return widthCompatible && heightCompatible
}

func sizeIsExactAndMatches(mode SizingMode, size, lastComputed float32) bool {
	return mode == SizingModeStretchFit && inexactEqualsFloat(size, lastComputed)
}

func oldSizeIsMaxContent(mode SizingMode, size float32, lastMode SizingMode, lastComputed float32) bool {
	return mode == SizingModeFitContent && lastMode == SizingModeMaxContent &&
		(size >= lastComputed || inexactEqualsFloat(size, lastComputed))
}

func newSizeIsStricter(mode SizingMode, size float32, lastMode SizingMode, lastSize, lastComputed float32) bool {
	return lastMode == SizingModeFitContent && mode == SizingModeFitContent &&
		isDefined(lastSize) && isDefined(size) && isDefined(lastComputed) &&
		lastSize > size && (lastComputed <= size || inexactEqualsFloat(size, lastComputed))
}
