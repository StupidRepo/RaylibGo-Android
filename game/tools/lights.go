package tools

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	LightDirectional = iota
	LightPoint
	LightSpot

	MaxLights = 1
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

func CreateLight(lightCount int, typ int32, position rl.Vector3, target rl.Vector3, color rl.Color, intensity float32, shader rl.Shader) Light {
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
	}

	return light
}

func UpdateLight(shader rl.Shader, light Light) {
	enabledBits := int32(0)
	if light.Enabled {
		enabledBits = 1
	}
	SetShaderInt(shader, light.EnabledLoc, enabledBits)
	SetShaderInt(shader, light.TypeLoc, light.Type)

	position := []float32{light.Position.X, light.Position.Y, light.Position.Z}
	rl.SetShaderValue(shader, light.PositionLoc, position, rl.ShaderUniformVec3)

	target := []float32{light.Target.X, light.Target.Y, light.Target.Z}
	rl.SetShaderValue(shader, light.TargetLoc, target, rl.ShaderUniformVec3)

	rl.SetShaderValue(shader, light.ColorLoc, light.Color[:], rl.ShaderUniformVec4)
	rl.SetShaderValue(shader, light.IntensityLoc, []float32{light.Intensity}, rl.ShaderUniformFloat)
}
