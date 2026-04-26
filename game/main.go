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
	lastHeight int32
	lastWidth  int32

	didResize bool

	currentHeight int32
	currentWidth  int32
)

func init() {
	// log that we are init and print the current OS
	platform.LogIt(ps.AndroidLogInfo, ps.GameTag, fmt.Sprintf("Current OS: %s", runtime.GOOS))
	rl.SetMain(main)
}

var platform = ps.Platform{}

func initRaylib() {
	rl.SetConfigFlags(rl.FlagWindowResizable | rl.FlagMsaa4xHint | rl.FlagFullscreenMode)

	resize()

	rl.InitWindow(currentWidth, currentHeight, "Android Game")
	rl.InitAudioDevice()
}

func resize() {
	didResize = false
	currentWidth, currentHeight = platform.GetWindowSize()

	if currentWidth != lastWidth || currentHeight != lastHeight {
		didResize = true
		lastWidth = currentWidth
		lastHeight = currentHeight

		platform.LogIt(ps.AndroidLogInfo, ps.GameTag, fmt.Sprintf("Window resized: %d x %d", currentWidth, currentHeight))
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

	platform.LogIt(ps.AndroidLogInfo, ps.GameTag, "Initializing game...")
	game.Init()

	frameCounter := 0

	for !windowShouldClose {
		if (platform.GetOS() == ps.PlatformAndroid && rl.IsKeyDown(rl.KeyBack)) || rl.WindowShouldClose() {
			windowShouldClose = true
		}
		resize()
		if didResize && platform.GetOS() == ps.PlatformAndroid {
			rl.SetWindowSize(int(currentWidth), int(currentHeight))
			platform.LogIt(ps.AndroidLogInfo, ps.GameTag, fmt.Sprintf("Adjusted window size to: %d x %d", currentWidth, currentHeight))
		}

		frameCounter++
		if frameCounter >= 400 {
			frameCounter = 0

			// switch to the next game in the list, looping back round if we reach the end
			cgi := getCurrentGameIndex()
			ngi := (cgi + 1) % len(games)

			platform.LogIt(ps.AndroidLogInfo, ps.GameTag, fmt.Sprintf("Switching demo: %s -> %s.", games[cgi].GetSpec().Name, games[ngi].GetSpec().Name))
			switchGame(ngi)
		}

		game.Update(currentWidth, currentHeight)

		rl.BeginDrawing()
		game.Draw()
		rl.EndDrawing()
	}

	game.Deinit()

	rl.CloseWindow()
	os.Exit(0)
}

func getCurrentGameIndex() int {
	currentGameIndex := 0
	for i, g := range games {
		if g == game {
			currentGameIndex = i
			break
		}
	}

	return currentGameIndex
}

func switchGame(index int) {
	if index >= 0 && index < len(games) {
		game.Deinit()

		game = games[index]
		game.Init()
	}
}
