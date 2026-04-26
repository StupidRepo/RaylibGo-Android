package demos

import (
	ps "RaylibGoGame/platformspecifics"
	"RaylibGoGame/tools"
	"fmt"

	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type CelShadingDemo struct {
	Platform *ps.Platform

	car rl.Model

	celShader   rl.Shader
	numBandsLoc int32

	outlineShader       rl.Shader
	outlineThicknessLoc int32

	screenWidth  int32
	screenHeight int32
}

var copyright = "(c) Old Rusty Car model by Renafox (https://skfb.ly/LxRy)"

var lights = make([]tools.Light, 0)

var numBands float32 = 3.0
var thickness float32 = 0.005

var camera = rl.Camera3D{
	Position: rl.NewVector3(9.0, 4.0, 9.0),
	Target:   rl.NewVector3(0, 0.5, 0),
	Up:       rl.NewVector3(0, 1.0, 0),

	Fovy:       45.0,
	Projection: rl.CameraPerspective,
}

func (c *CelShadingDemo) Init() {
	// car
	c.car = rl.LoadModel("resources/models/old_car_new.glb")

	// cel shader
	c.celShader = rl.LoadShader(
		fmt.Sprintf("resources/shaders/glsl%d/cel.vs", ps.GLSLVersion),
		fmt.Sprintf("resources/shaders/glsl%d/cel.fs", ps.GLSLVersion))
	c.celShader.UpdateLocation(rl.ShaderLocVectorView, rl.GetShaderLocation(c.celShader, "viewPos"))

	c.numBandsLoc = rl.GetShaderLocation(c.celShader, "numBands")
	rl.SetShaderValue(c.celShader, c.numBandsLoc, []float32{numBands}, rl.ShaderUniformFloat)

	carMats := c.car.GetMaterials()
	if len(carMats) > 0 {
		carMats[0].Shader = c.celShader
	}

	// outline shader
	c.outlineShader = rl.LoadShader(
		fmt.Sprintf("resources/shaders/glsl%d/outline_hull.vs", ps.GLSLVersion),
		fmt.Sprintf("resources/shaders/glsl%d/outline_hull.fs", ps.GLSLVersion))
	c.outlineThicknessLoc = rl.GetShaderLocation(c.outlineShader, "outlineThickness")
	rl.SetShaderValue(c.outlineShader, c.outlineThicknessLoc, []float32{thickness}, rl.ShaderUniformFloat)

	// lights
	lightDefs := []struct {
		position  rl.Vector3
		color     rl.Color
		intensity float32
	}{
		{rl.NewVector3(50, 50, 50), rl.White, 4.0},
		//{rl.NewVector3(2.0, 1.0, 1.0), rl.Green, 3.3},
		//{rl.NewVector3(-2.0, 1.0, 1.0), rl.Red, 8.3},
		//{rl.NewVector3(1.0, 1.0, -2.0), rl.Blue, 2.0},
	}
	for _, def := range lightDefs {
		lights = append(lights, tools.CreateLight(len(lights), tools.LightDirectional, def.position, rl.NewVector3(0, 0, 0), def.color, def.intensity, c.celShader))
	}
}

func (c *CelShadingDemo) Update(CurrentWidth int32, CurrentHeight int32) {
	c.screenWidth = CurrentWidth
	c.screenHeight = CurrentHeight

	rl.UpdateCamera(&camera, rl.CameraOrbital)
	c.Platform.LogIt(ps.AndroidLogDebug, ps.GameTag, fmt.Sprintf("Camera: %+v", camera))

	// cel shader upd
	rl.SetShaderValue(c.celShader, c.celShader.GetLocation(rl.ShaderLocVectorView), []float32{camera.Position.X, camera.Position.Y, camera.Position.Z}, rl.ShaderUniformVec3)
	rl.SetShaderValue(c.celShader, c.numBandsLoc, []float32{numBands}, rl.ShaderUniformFloat)

	// outline upd
	rl.SetShaderValue(c.outlineShader, c.outlineThicknessLoc, []float32{thickness}, rl.ShaderUniformFloat)

	// lights upd
	for i := range lights {
		tools.UpdateLight(c.celShader, lights[i])
	}
}

func (c *CelShadingDemo) Draw() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.White)

	rl.BeginMode3D(camera)

	// outline start
	rl.SetCullFace(0) // CULL_FACE_FRONT

	carMats := c.car.GetMaterials()
	if len(carMats) > 0 {
		carMats[0].Shader = c.outlineShader
	}

	rl.DrawModel(c.car, rl.Vector3Zero(), 0.75, rl.Black)

	if len(carMats) > 0 {
		carMats[0].Shader = c.celShader
	}

	rl.SetCullFace(1) // CULL_FACE_BACK
	// outline end

	rl.DrawModel(c.car, rl.Vector3Zero(), 0.75, rl.Gold)
	rl.DrawGrid(10, 10.0)

	rl.EndMode3D()

	c.DrawUI()

	rl.EndDrawing()
}

func (c *CelShadingDemo) DrawUI() {
	insets := c.Platform.GetInsets()

	// insets.Left so that it's matching ;)
	rl.DrawText(copyright, c.screenWidth-insets.Left-rl.MeasureText(copyright, 20), insets.Top, 20, rl.SkyBlue)
	rl.DrawFPS(insets.Left, insets.Top)

	// 2 gui buttons that are in the screen center and above bottom bar (insets.bottom)
	var buttonWidth int32 = 150
	var buttonHeight int32 = 40
	var buttonSpacing int32 = 10
	var buttonCount int32 = 4

	var totalWidth = buttonWidth*buttonCount + buttonSpacing*(buttonCount-1)
	var startX = (c.screenWidth - totalWidth) / 2
	var y = c.screenHeight - insets.Bottom - buttonHeight - 10

	// bands increase, decrease
	if gui.Button(rl.NewRectangle(float32(startX), float32(y), float32(buttonWidth), float32(buttonHeight)), "Bands INCREASE") {
		numBands++
	}
	if gui.Button(rl.NewRectangle(float32(startX+buttonWidth+buttonSpacing), float32(y), float32(buttonWidth), float32(buttonHeight)), "Bands DECREASE") {
		if numBands > 1 {
			numBands--
		}
	}

	// thickness increase, decrease
	if gui.Button(rl.NewRectangle(float32(startX+2*(buttonWidth+buttonSpacing)), float32(y), float32(buttonWidth), float32(buttonHeight)), "Thickness INCREASE") {
		thickness += 0.001
	}
	if gui.Button(rl.NewRectangle(float32(startX+3*(buttonWidth+buttonSpacing)), float32(y), float32(buttonWidth), float32(buttonHeight)), "Thickness DECREASE") {
		if thickness > 0.001 {
			thickness -= 0.001
		}
	}

	// if we're on android, draw mobile UI for zooming in/out
	//if c.Platform.GetOS() == ps.PlatformAndroid {
	c.DrawMobileUI(insets)
	//}
}

func (c *CelShadingDemo) DrawMobileUI(insets ps.Insets) {
	var buttonSize int32 = 50
	var spacing int32 = 10

	zoomInRect := rl.NewRectangle(float32(insets.Left), float32(c.screenHeight-insets.Bottom-buttonSize*2-spacing), float32(buttonSize), float32(buttonSize))
	zoomOutRect := rl.NewRectangle(float32(insets.Left), float32(c.screenHeight-insets.Bottom-buttonSize), float32(buttonSize), float32(buttonSize))

	if gui.Button(zoomInRect, gui.IconText(gui.ICON_ARROW_UP_FILL, "+")) {
		rl.CameraMoveToTarget(&camera, -0.5)
	}
	if gui.Button(zoomOutRect, gui.IconText(gui.ICON_ARROW_DOWN_FILL, "-")) {
		rl.CameraMoveToTarget(&camera, 0.5)
	}
}

func (c *CelShadingDemo) Deinit() {
	rl.UnloadModel(c.car)
	rl.UnloadShader(c.celShader)
	rl.UnloadShader(c.outlineShader)
}
