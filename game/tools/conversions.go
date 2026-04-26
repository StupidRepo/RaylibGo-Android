package tools

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func SetShaderInt(shader rl.Shader, loc int32, value int32) {
	packed := math.Float32frombits(uint32(value))
	rl.SetShaderValue(shader, loc, []float32{packed}, rl.ShaderUniformInt)
}

func ColorToVec3(c rl.Color) []float32 {
	return []float32{float32(c.R) / 255.0, float32(c.G) / 255.0, float32(c.B) / 255.0}
}

func ColorToVec4(c rl.Color) []float32 {
	return []float32{float32(c.R) / 255.0, float32(c.G) / 255.0, float32(c.B) / 255.0, float32(c.A) / 255.0}
}
