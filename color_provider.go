package main

type ColorProvider struct {
	Brightness float64
}

func NewColorProvider(brightness float64) *ColorProvider {
	return &ColorProvider{Brightness: brightness}
}

func (colorProvider *ColorProvider) getColor(
	airQualityIndex float32) NeopixelColor {
	switch {
	case airQualityIndex <= 25:
		return colorProvider.getAdjustedColor(0, 50, 0)
	case airQualityIndex <= 50:
		return colorProvider.getAdjustedColor(40, 50, 0)
	case airQualityIndex <= 75:
		return colorProvider.getAdjustedColor(70, 50, 0)
	case airQualityIndex <= 90:
		return colorProvider.getAdjustedColor(80, 40, 0)
	case airQualityIndex <= 100:
		return colorProvider.getAdjustedColor(90, 0, 0)
	default:
		return colorProvider.getAdjustedColor(100, 0, 20)
	}
}

func (colorProvider *ColorProvider) getAdjustedColor(
	red uint8, green uint8, blue uint8) NeopixelColor {
	return NeopixelColor{
		Red:   uint8(colorProvider.Brightness * float64(red)),
		Green: uint8(colorProvider.Brightness * float64(green)),
		Blue:  uint8(colorProvider.Brightness * float64(blue)),
	}
}
