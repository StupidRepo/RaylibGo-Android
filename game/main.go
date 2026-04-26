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

var games = []demos.Demo{
	&demos.CelShadingDemo{
		Platform: &platform,
	},
	&demos.TwoDDemo{
		Platform: &platform,
	},
}
var game = games[0]

func main() {
	initRaylib()
	resize()
	windowShouldClose := false

	// Init stage
	platform.LogIt(ps.AndroidLogInfo, ps.GameTag, "Initializing game...")
	game.Init()

	frameCounter := 0

	for !windowShouldClose {
		if (platform.GetOS() == ps.PlatformAndroid && rl.IsKeyDown(rl.KeyBack)) || rl.WindowShouldClose() {
			windowShouldClose = true
		}
		resize()

		frameCounter++
		if frameCounter >= 400 {
			frameCounter = 0

			// switch to the next game in the list, looping back round if we reach the end
			currentGameIndex := 0
			for i, g := range games {
				if g == game {
					currentGameIndex = i
					break
				}
			}
			nextGameIndex := (currentGameIndex + 1) % len(games)
			switchGame(nextGameIndex)
		}

		game.Update(CurrentWidth, CurrentHeight)
		game.Draw()
	}

	game.Deinit()

	rl.CloseWindow()
	os.Exit(0)
}

func switchGame(index int) {
	if index >= 0 && index < len(games) {
		game.Deinit()

		game = games[index]
		game.Init()
	}
}
