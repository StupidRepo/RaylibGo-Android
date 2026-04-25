//go:build linux || android || windows || darwin

package main

import (
	ps "RaylibGoGame/platformspecifics"
	"fmt"
	"math"
	"os"
	"runtime"

	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	CurrentHeight int32
	CurrentWidth  int32
)

const (
	copyright = "(c) Old Rusty Car model by Renafox (https://skfb.ly/LxRy)"
)

func init() {
	// log that we are init and print the current OS
	platform.LogIt(ps.AndroidLogInfo, ps.GameTag, fmt.Sprintf("Current OS: %s", runtime.GOOS))
	rl.SetMain(main)
}

const (
	LightDirectional = iota
	LightPoint
	LightSpot

	MaxLights = 4
)

type Light struct {
	Type    int32
	Enabled bool

	Position rl.Vector3
	Target   rl.Vector3

	Color     [4]float32
	Intensity float32

	TypeLoc      int32
	EnabledLoc   int32
	PositionLoc  int32
	TargetLoc    int32
	ColorLoc     int32
	IntensityLoc int32
}

var lights [MaxLights]Light
var lightCount int

var platform = ps.Platform{}

type shaderUniformLocs struct {
	viewPos           int32
	metallicValue     int32
	roughnessValue    int32
	emissiveIntensity int32
	emissiveColor     int32
	textureTiling     int32
}

type materialSetup struct {
	albedoColor   rl.Color
	metalness     float32
	roughness     float32
	occlusion     float32
	emissiveColor rl.Color

	albedoTex    string
	metalnessTex string
	normalTex    string
	emissiveTex  string
}

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

func setShaderInt(shader rl.Shader, loc int32, value int32) {
	packed := math.Float32frombits(uint32(value))
	rl.SetShaderValue(shader, loc, []float32{packed}, rl.ShaderUniformInt)
}

func colorToVec3(c rl.Color) []float32 {
	return []float32{float32(c.R) / 255.0, float32(c.G) / 255.0, float32(c.B) / 255.0}
}

func colorToVec4(c rl.Color) []float32 {
	return []float32{float32(c.R) / 255.0, float32(c.G) / 255.0, float32(c.B) / 255.0, float32(c.A) / 255.0}
}

func setupModelMaterial(model rl.Model, shader rl.Shader, cfg materialSetup) {
	mats := model.GetMaterials()
	if len(mats) == 0 {
		return
	}

	mat := &mats[0]
	mat.Shader = shader
	mat.GetMap(rl.MapAlbedo).Color = cfg.albedoColor
	mat.GetMap(rl.MapMetalness).Value = cfg.metalness
	mat.GetMap(rl.MapRoughness).Value = cfg.roughness
	mat.GetMap(rl.MapOcclusion).Value = cfg.occlusion
	mat.GetMap(rl.MapEmission).Color = cfg.emissiveColor

	if cfg.albedoTex != "" {
		mat.GetMap(rl.MapAlbedo).Texture = rl.LoadTexture(cfg.albedoTex)
	}
	if cfg.metalnessTex != "" {
		mat.GetMap(rl.MapMetalness).Texture = rl.LoadTexture(cfg.metalnessTex)
	}
	if cfg.normalTex != "" {
		mat.GetMap(rl.MapNormal).Texture = rl.LoadTexture(cfg.normalTex)
	}
	if cfg.emissiveTex != "" {
		mat.GetMap(rl.MapEmission).Texture = rl.LoadTexture(cfg.emissiveTex)
	}
}

func applyMaterialUniforms(shader rl.Shader, locs shaderUniformLocs, mat *rl.Material, tiling rl.Vector2) {
	rl.SetShaderValue(shader, locs.textureTiling, []float32{tiling.X, tiling.Y}, rl.ShaderUniformVec2)
	rl.SetShaderValue(shader, locs.emissiveColor, colorToVec4(mat.GetMap(rl.MapEmission).Color), rl.ShaderUniformVec4)
	rl.SetShaderValue(shader, locs.metallicValue, []float32{mat.GetMap(rl.MapMetalness).Value}, rl.ShaderUniformFloat)
	rl.SetShaderValue(shader, locs.roughnessValue, []float32{mat.GetMap(rl.MapRoughness).Value}, rl.ShaderUniformFloat)
}

func unloadModelMaterial(model rl.Model) {
	mats := model.GetMaterials()
	if len(mats) > 0 {
		mats[0].Shader = rl.Shader{}
		rl.UnloadMaterial(mats[0])
		mats[0].Maps = nil
	}
	rl.UnloadModel(model)
}

func main() {
	initRaylib()
	resize()
	windowShouldClose := false

	var camera = rl.Camera3D{
		Position: rl.NewVector3(2.0, 2.0, 6.0),
		Target:   rl.NewVector3(0, 0.5, 0),
		Up:       rl.NewVector3(0, 1.0, 0),

		Fovy:       45.0,
		Projection: rl.CameraPerspective,
	}

	platform.LogIt(ps.AndroidLogDebug, ps.GameTag, fmt.Sprintf("Camera initialized: %+v", camera))

	shader := rl.LoadShader(
		fmt.Sprintf("resources/shaders/glsl%d/pbr.vs", ps.GLSLVersion),
		fmt.Sprintf("resources/shaders/glsl%d/pbr.fs", ps.GLSLVersion))

	shader.UpdateLocation(rl.ShaderLocMapAlbedo, rl.GetShaderLocation(shader, "albedoMap"))
	shader.UpdateLocation(rl.ShaderLocMapMetalness, rl.GetShaderLocation(shader, "mraMap"))
	shader.UpdateLocation(rl.ShaderLocMapNormal, rl.GetShaderLocation(shader, "normalMap"))
	shader.UpdateLocation(rl.ShaderLocMapEmission, rl.GetShaderLocation(shader, "emissiveMap"))
	shader.UpdateLocation(rl.ShaderLocColorDiffuse, rl.GetShaderLocation(shader, "albedoColor"))

	shader.UpdateLocation(rl.ShaderLocVectorView, rl.GetShaderLocation(shader, "viewPos"))
	locs := shaderUniformLocs{
		viewPos:           rl.GetShaderLocation(shader, "viewPos"),
		metallicValue:     rl.GetShaderLocation(shader, "metallicValue"),
		roughnessValue:    rl.GetShaderLocation(shader, "roughnessValue"),
		emissiveIntensity: rl.GetShaderLocation(shader, "emissivePower"),
		emissiveColor:     rl.GetShaderLocation(shader, "emissiveColor"),
		textureTiling:     rl.GetShaderLocation(shader, "tiling"),
	}

	platform.LogIt(ps.AndroidLogDebug, ps.GameTag, fmt.Sprintf("Using GL version: %d", ps.GLSLVersion))

	setShaderInt(shader, rl.GetShaderLocation(shader, "numOfLights"), MaxLights)

	ambientIntensity := float32(0.02)
	ambientColor := rl.NewColor(26, 32, 135, 255)
	rl.SetShaderValue(shader, rl.GetShaderLocation(shader, "ambientColor"), colorToVec3(ambientColor), rl.ShaderUniformVec3)
	rl.SetShaderValue(shader, rl.GetShaderLocation(shader, "ambient"), []float32{ambientIntensity}, rl.ShaderUniformFloat)

	car := rl.LoadModel("resources/models/old_car_new.glb")
	setupModelMaterial(car, shader, materialSetup{
		albedoColor:   rl.White,
		metalness:     1,
		roughness:     0,
		occlusion:     1,
		emissiveColor: rl.NewColor(255, 162, 0, 255),
		albedoTex:     "resources/old_car_d.png",
		metalnessTex:  "resources/old_car_mra.png",
		normalTex:     "resources/old_car_n.png",
		emissiveTex:   "resources/old_car_e.png",
	})

	floor := rl.LoadModel("resources/models/plane.glb")
	setupModelMaterial(floor, shader, materialSetup{
		albedoColor:   rl.White,
		metalness:     0.8,
		roughness:     0.1,
		occlusion:     1,
		emissiveColor: rl.Black,
		albedoTex:     "resources/road_a.png",
		metalnessTex:  "resources/road_mra.png",
		normalTex:     "resources/road_n.png",
	})

	carTextureTiling := rl.NewVector2(0.5, 0.5)
	floorTextureTiling := rl.NewVector2(0.5, 0.5)
	carEmissiveIntensity := float32(0.01)

	lightDefs := []struct {
		position  rl.Vector3
		color     rl.Color
		intensity float32
	}{
		{rl.NewVector3(-1.0, 1.0, -2.0), rl.Yellow, 4.0},
		{rl.NewVector3(2.0, 1.0, 1.0), rl.Green, 3.3},
		{rl.NewVector3(-2.0, 1.0, 1.0), rl.Red, 8.3},
		{rl.NewVector3(1.0, 1.0, -2.0), rl.Blue, 2.0},
	}
	for i, def := range lightDefs {
		lights[i] = CreateLight(LightPoint, def.position, rl.NewVector3(0, 0, 0), def.color, def.intensity, shader)
	}

	for _, uniform := range []string{"useTexAlbedo", "useTexNormal", "useTexMRA", "useTexEmissive"} {
		setShaderInt(shader, rl.GetShaderLocation(shader, uniform), 1)
	}

	for !windowShouldClose {
		if (runtime.GOOS == "android" && rl.IsKeyDown(rl.KeyBack)) || rl.WindowShouldClose() {
			windowShouldClose = true
		}
		resize()

		rl.UpdateCamera(&camera, rl.CameraOrbital)
		rl.SetShaderValue(shader, locs.viewPos, []float32{camera.Position.X, camera.Position.Y, camera.Position.Z}, rl.ShaderUniformVec3)

		for _, binding := range []struct {
			key   int32
			index int
		}{
			{rl.KeyOne, 2},
			{rl.KeyTwo, 1},
			{rl.KeyThree, 3},
			{rl.KeyFour, 0},
		} {
			if rl.IsKeyPressed(binding.key) {
				lights[binding.index].Enabled = !lights[binding.index].Enabled
			}
		}

		for i := range lights {
			UpdateLight(shader, lights[i])
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		rl.BeginMode3D(camera)

		floorMats := floor.GetMaterials()
		if len(floorMats) > 0 {
			applyMaterialUniforms(shader, locs, &floorMats[0], floorTextureTiling)
		}

		rl.DrawModel(floor, rl.NewVector3(0, 0, 0), 5.0, rl.White)

		carMats := car.GetMaterials()
		if len(carMats) > 0 {
			applyMaterialUniforms(shader, locs, &carMats[0], carTextureTiling)
			rl.SetShaderValue(shader, locs.emissiveIntensity, []float32{carEmissiveIntensity}, rl.ShaderUniformFloat)
		}

		rl.DrawModel(car, rl.NewVector3(0, 0, 0), 0.25, rl.White)

		for i := range lights {
			lightColor := rl.NewColor(uint8(lights[i].Color[0]*255), uint8(lights[i].Color[1]*255), uint8(lights[i].Color[2]*255), uint8(lights[i].Color[3]*255))
			if lights[i].Enabled {
				rl.DrawSphereEx(lights[i].Position, 0.2, 8, 8, lightColor)
			} else {
				rl.DrawSphereWires(lights[i].Position, 0.2, 8, 8, rl.ColorAlpha(lightColor, 0.3))
			}
		}

		rl.EndMode3D()

		insets := platform.GetInsets()

		// insets.Left so that it's matching ;)
		rl.DrawText(copyright, CurrentWidth-insets.Left-rl.MeasureText(copyright, 20), insets.Top, 20, rl.White)
		rl.DrawFPS(insets.Left, insets.Top)

		for i := range lights {
			state := "on"
			if lights[i].Enabled {
				state = "off"
			}
			button := gui.Button(rl.NewRectangle(float32(int(insets.Left)+((150+10)*i)), float32(CurrentHeight-60), 150, 60), fmt.Sprintf("Turn %s light #%d", state, i+1))
			if button {
				lights[i].Enabled = !lights[i].Enabled
				UpdateLight(shader, lights[i])
			}
		}

		rl.EndDrawing()
	}

	unloadModelMaterial(car)
	unloadModelMaterial(floor)

	rl.UnloadShader(shader)

	rl.CloseWindow()
	os.Exit(0)
}

func CreateLight(typ int32, position rl.Vector3, target rl.Vector3, color rl.Color, intensity float32, shader rl.Shader) Light {
	var light Light

	if lightCount < MaxLights {
		idx := lightCount
		light.Enabled = true
		light.Type = typ
		light.Position = position
		light.Target = target
		light.Color = [4]float32{float32(color.R) / 255.0, float32(color.G) / 255.0, float32(color.B) / 255.0, float32(color.A) / 255.0}
		light.Intensity = intensity

		light.EnabledLoc = rl.GetShaderLocation(shader, fmt.Sprintf("lights[%d].enabled", idx))
		light.TypeLoc = rl.GetShaderLocation(shader, fmt.Sprintf("lights[%d].type", idx))
		light.PositionLoc = rl.GetShaderLocation(shader, fmt.Sprintf("lights[%d].position", idx))
		light.TargetLoc = rl.GetShaderLocation(shader, fmt.Sprintf("lights[%d].target", idx))
		light.ColorLoc = rl.GetShaderLocation(shader, fmt.Sprintf("lights[%d].color", idx))
		light.IntensityLoc = rl.GetShaderLocation(shader, fmt.Sprintf("lights[%d].intensity", idx))

		UpdateLight(shader, light)
		lightCount++
	}

	return light
}

func UpdateLight(shader rl.Shader, light Light) {
	enabledBits := int32(0)
	if light.Enabled {
		enabledBits = 1
	}
	setShaderInt(shader, light.EnabledLoc, enabledBits)
	setShaderInt(shader, light.TypeLoc, light.Type)

	position := []float32{light.Position.X, light.Position.Y, light.Position.Z}
	rl.SetShaderValue(shader, light.PositionLoc, position, rl.ShaderUniformVec3)

	target := []float32{light.Target.X, light.Target.Y, light.Target.Z}
	rl.SetShaderValue(shader, light.TargetLoc, target, rl.ShaderUniformVec3)

	rl.SetShaderValue(shader, light.ColorLoc, light.Color[:], rl.ShaderUniformVec4)
	rl.SetShaderValue(shader, light.IntensityLoc, []float32{light.Intensity}, rl.ShaderUniformFloat)
}
