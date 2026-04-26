package demos

type Demo interface {
	Init()                                          // Asset loading - called once after InitWindow.
	Update(CurrentWidth int32, CurrentHeight int32) // Game logic update.
	Draw()                                          // Draw logic.
	Deinit()                                        // Asset unloading - called once before CloseWindow or when game is switched.
}
