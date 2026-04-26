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

	lights []tools.Light

	numBands  float32
	thickness float32

	camera rl.Camera3D

	screenWidth  int32
	screenHeight int32
}

var copyright = "(c) Old Rusty Car model by Renafox (https://skfb.ly/LxRy)"

func (c *CelShadingDemo) GetSpec() DemoSpec {
	return DemoSpec{
		Name:    "3D Cel Shading",
		Summary: "A demo showcasing cel shading with an outline effect on a 3D model.",
	}
}

func (c *CelShadingDemo) Init() {
	// car
	c.car = rl.LoadModel("resources/models/old_car_new.glb")

	// cel shader
	c.celShader = rl.LoadShader(
		fmt.Sprintf("resources/shaders/glsl%d/cel.vs", ps.GLSLVersion),
		fmt.Sprintf("resources/shaders/glsl%d/cel.fs", ps.GLSLVersion))
	c.celShader.UpdateLocation(rl.ShaderLocVectorView, rl.GetShaderLocation(c.celShader, "viewPos"))

	c.numBands = 7.0
	c.numBandsLoc = rl.GetShaderLocation(c.celShader, "numBands")
	rl.SetShaderValue(c.celShader, c.numBandsLoc, []float32{c.numBands}, rl.ShaderUniformFloat)

	carMats := c.car.GetMaterials()
	if len(carMats) > 0 {
		carMats[0].Shader = c.celShader
	}

	// outline shader
	c.thickness = 0.003
	c.outlineShader = rl.LoadShader(
		fmt.Sprintf("resources/shaders/glsl%d/outline_hull.vs", ps.GLSLVersion),
		fmt.Sprintf("resources/shaders/glsl%d/outline_hull.fs", ps.GLSLVersion))
	c.outlineThicknessLoc = rl.GetShaderLocation(c.outlineShader, "outlineThickness")
	rl.SetShaderValue(c.outlineShader, c.outlineThicknessLoc, []float32{c.thickness}, rl.ShaderUniformFloat)

	// camera
	c.camera = rl.NewCamera3D(rl.NewVector3(9.0, 4.0, 9.0), rl.NewVector3(0, 0.5, 0), rl.NewVector3(0, 1.0, 0), 45.0, rl.CameraPerspective)

	// lights
	c.lights = make([]tools.Light, 0, 1)

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
		c.Platform.LogIt(ps.AndroidLogDebug, ps.GameTag, fmt.Sprintf("Light count before: %d", len(c.lights)))
		c.lights = append(c.lights, tools.CreateLight(len(c.lights), tools.LightDirectional, def.position, rl.NewVector3(0, 0, 0), def.color, def.intensity, c.celShader))

		c.Platform.LogIt(ps.AndroidLogDebug, ps.GameTag, fmt.Sprintf("Light count after: %d", len(c.lights)))
	}
}

func (c *CelShadingDemo) Update(CurrentWidth int32, CurrentHeight int32) {
	c.screenWidth = CurrentWidth
	c.screenHeight = CurrentHeight

	rl.UpdateCamera(&c.camera, rl.CameraOrbital)

	// cel shader upd
	rl.SetShaderValue(c.celShader, c.celShader.GetLocation(rl.ShaderLocVectorView), []float32{c.camera.Position.X, c.camera.Position.Y, c.camera.Position.Z}, rl.ShaderUniformVec3)
	rl.SetShaderValue(c.celShader, c.numBandsLoc, []float32{c.numBands}, rl.ShaderUniformFloat)

	// outline upd
	rl.SetShaderValue(c.outlineShader, c.outlineThicknessLoc, []float32{c.thickness}, rl.ShaderUniformFloat)

	// lights upd
	for i := range c.lights {
		tools.UpdateLight(c.celShader, c.lights[i])
	}
}

func (c *CelShadingDemo) Draw() {
	// rl.BeginDrawing()
	rl.ClearBackground(rl.White)

	rl.BeginMode3D(c.camera)

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

	// rl.EndDrawing()
}

func (c *CelShadingDemo) DrawUI() {
	insets := c.Platform.GetInsets()

	// insets.Left so that it's matching ;)
	fontSize := tools.UiPxToDp(20)
	tools.DrawTextCenter(copyright, float32(c.screenWidth), fontSize, fontSize, rl.Black)
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
		c.numBands++
	}
	if gui.Button(rl.NewRectangle(float32(startX+buttonWidth+buttonSpacing), float32(y), float32(buttonWidth), float32(buttonHeight)), "Bands DECREASE") {
		c.numBands--
	}
	c.numBands = max(1, c.numBands)

	// thickness increase, decrease
	if gui.Button(rl.NewRectangle(float32(startX+2*(buttonWidth+buttonSpacing)), float32(y), float32(buttonWidth), float32(buttonHeight)), "Thickness INCREASE") {
		c.thickness += 0.001
	}
	if gui.Button(rl.NewRectangle(float32(startX+3*(buttonWidth+buttonSpacing)), float32(y), float32(buttonWidth), float32(buttonHeight)), "Thickness DECREASE") {
		c.thickness -= 0.001
	}
	c.thickness = rl.Clamp(c.thickness, 0.001, 0.01)

	// if we're on android, draw mobile UI for zooming in/out
	if c.Platform.GetOS() == ps.PlatformAndroid {
		c.DrawMobileUI(insets)
	}
}

func (c *CelShadingDemo) DrawMobileUI(insets ps.Insets) {
	var buttonSize int32 = 50
	var spacing int32 = 10

	zoomInRect := rl.NewRectangle(float32(insets.Left), float32(c.screenHeight-insets.Bottom-buttonSize*2-spacing), float32(buttonSize), float32(buttonSize))
	zoomOutRect := rl.NewRectangle(float32(insets.Left), float32(c.screenHeight-insets.Bottom-buttonSize), float32(buttonSize), float32(buttonSize))

	if gui.Button(zoomInRect, gui.IconText(gui.ICON_ARROW_UP_FILL, "+")) {
		rl.CameraMoveToTarget(&c.camera, -0.5)
	}
	if gui.Button(zoomOutRect, gui.IconText(gui.ICON_ARROW_DOWN_FILL, "-")) {
		rl.CameraMoveToTarget(&c.camera, 0.5)
	}
}

func (c *CelShadingDemo) Deinit() {
	rl.UnloadModel(c.car)
	rl.UnloadShader(c.celShader)
	rl.UnloadShader(c.outlineShader)
}
