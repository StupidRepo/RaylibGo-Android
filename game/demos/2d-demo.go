package demos

import (
	ps "RaylibGoGame/platformspecifics"
	"RaylibGoGame/tools"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type TwoDDemo struct {
	Platform *ps.Platform

	screenWidth  int32
	screenHeight int32
}

func (d *TwoDDemo) GetSpec() DemoSpec {
	return DemoSpec{
		Name:    "2D Demo",
		Summary: "A simple 2D demo showcasing basic drawing capabilities.",
	}
}

func (d *TwoDDemo) Init() {
	d.screenWidth, d.screenHeight = d.Platform.GetWindowSize()
}

func (d *TwoDDemo) Update(CurrentWidth int32, CurrentHeight int32) {
	d.screenWidth = CurrentWidth
	d.screenHeight = CurrentHeight
}

func (d *TwoDDemo) Draw() {
	// rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)

	fontSize := tools.UiPxToDp(20)

	for i, line := range []string{
		"Welcome to the 2D Demo!",
		"This demo showcases basic 2D drawing capabilities.",
		"Enjoy the colorful shapes and text!",
	} {
		y := fontSize + ((fontSize * 2) * float32(i))
		tools.DrawTextCenter(line, float32(d.screenWidth), y, fontSize, rl.Black)
	}

	rl.DrawRectangleGradientH(d.screenWidth/2-150, d.screenHeight/2, 300, 100, rl.Blue, rl.Green)
	rl.DrawCircleGradient(d.screenWidth/2, d.screenHeight/2+150, 50, rl.Red, rl.Yellow)

	// rl.EndDrawing()
}

func (d *TwoDDemo) Deinit() {
}
