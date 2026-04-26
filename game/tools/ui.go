package tools

import rl "github.com/gen2brain/raylib-go/raylib"

func DrawTextCenter(text string, x float32, y float32, fontSize float32, color rl.Color) {
	spacing := UiPxToDp(2)

	size := rl.MeasureTextEx(rl.GetFontDefault(), text, fontSize, spacing)
	pos := rl.NewVector2((x-size.X)/2, (y-size.Y)/2)

	rl.DrawTextEx(rl.GetFontDefault(), text, pos, fontSize, spacing, color)
}
