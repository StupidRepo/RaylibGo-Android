package demos

type Demo interface {
	GetSpec() DemoSpec // Get demo spec

	Init()                                          // Asset loading - called once after InitWindow.
	Update(CurrentWidth int32, CurrentHeight int32) // Game logic update.
	Draw()                                          // Draw logic.
	Deinit()                                        // Asset unloading - called once before CloseWindow or when game is switched.
}

type DemoSpec struct {
	Name    string // Name of the demo.
	Summary string // Summary of the demo.
}
