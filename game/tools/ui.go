package tools

import rl "github.com/gen2brain/raylib-go/raylib"

func DrawTextCenter(text string, x float32, y float32, fontSize float32, color rl.Color) {
	spacing := UiPxToDp(2)

	size := rl.MeasureTextEx(rl.GetFontDefault(), text, fontSize, spacing)
	pos := rl.NewVector2((x-size.X)/2, (y-size.Y)/2)

	rl.DrawTextEx(rl.GetFontDefault(), text, pos, fontSize, spacing, color)
}

func GradientColor(startColor, endColor rl.Color, min float32, max float32, t float32) rl.Color {
	if t < min {
		return startColor
	}
	if t > max {
		return endColor
	}

	ratio := (t - min) / (max - min)

	return rl.NewColor(
		uint8(float32(startColor.R)*(1-ratio)+float32(endColor.R)*ratio),
		uint8(float32(startColor.G)*(1-ratio)+float32(endColor.G)*ratio),
		uint8(float32(startColor.B)*(1-ratio)+float32(endColor.B)*ratio),
		uint8(float32(startColor.A)*(1-ratio)+float32(endColor.A)*ratio),
	)
}
