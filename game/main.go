//go:build linux || android || windows || darwin

package main

import (
	"RaylibGoGame/demos"
	ps "RaylibGoGame/platformspecifics"
	"fmt"
	"os"
	"runtime"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	CurrentHeight int32
	CurrentWidth  int32
)

func init() {
	// log that we are init and print the current OS
	platform.LogIt(ps.AndroidLogInfo, ps.GameTag, fmt.Sprintf("Current OS: %s", runtime.GOOS))
	rl.SetMain(main)
}

var platform = ps.Platform{}

func initRaylib() {
	rl.SetConfigFlags(rl.FlagVsyncHint | rl.FlagWindowResizable | rl.FlagMsaa4xHint | rl.FlagFullscreenMode)

	resize()

	rl.InitWindow(CurrentWidth, CurrentHeight, "Android Game")
	rl.InitAudioDevice()
}

func resize() {
	CurrentWidth, CurrentHeight = platform.GetWindowSize()

	if rl.IsWindowResized() {
		rl.SetFramebufferWidth(CurrentWidth)
		rl.SetFramebufferHeight(CurrentHeight)
	}
}

var game demos.Demo

func main() {
	initRaylib()
	resize()
	windowShouldClose := false

	game = &demos.CelShadingDemo{
		Platform: &platform,
	}

	// Init stage
	platform.LogIt(ps.AndroidLogInfo, ps.GameTag, "Initializing game...")
	game.Init()

	for !windowShouldClose {
		if (runtime.GOOS == "android" && rl.IsKeyDown(rl.KeyBack)) || rl.WindowShouldClose() {
			windowShouldClose = true
		}
		resize()

		game.Update(CurrentWidth, CurrentHeight)

		game.Draw()
	}

	game.Deinit()

	rl.CloseWindow()
	os.Exit(0)
}
