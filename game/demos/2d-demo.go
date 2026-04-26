package demos

import (
	ps "RaylibGoGame/platformspecifics"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type TwoDDemo struct {
	Platform *ps.Platform

	screenWidth  int32
	screenHeight int32
}

func (d *TwoDDemo) Init() {
	d.screenWidth, d.screenHeight = d.Platform.GetWindowSize()
}

func (d *TwoDDemo) Update(CurrentWidth int32, CurrentHeight int32) {
}

func (d *TwoDDemo) Draw() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)

	rl.DrawText("This is a 2D demo!", d.screenWidth/2-rl.MeasureText("This is a 2D demo!", 20)/2, d.screenHeight/2-10, 20, rl.Black)

	rl.EndDrawing()
}

func (d *TwoDDemo) Deinit() {
}
