// This example demonstrates instanced drawing with low-level rlgl bindings:
//   - LoadVertexBufferElement: upload index data (EBO) to GPU
//   - DrawVertexArrayElementsInstanced: instanced indexed draw
//   - SetVertexAttributeDivisor: per-instance attribute advancement
//   - Camera3D with BeginMode3D/EndMode3D for orbital 3D camera control
//
// 625 quads are rendered in a 25x25 grid with animated Y-offsets (sine wave).
package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 1024
	screenHeight = 768

	cols          = 25
	rows          = 25
	instanceCount = cols * rows
	spacing       = 1.5
)

// GLSL 330 vertex shader with instancing support.
const vertexShaderCode = `#version 330
layout(location = 0) in vec3 vertexPosition;
layout(location = 1) in vec2 vertexTexCoord;
layout(location = 2) in vec3 instanceOffset;
layout(location = 3) in vec4 instanceColor;

uniform mat4 mvp;

out vec4 fragColor;

void main() {
    gl_Position = mvp * vec4(vertexPosition + instanceOffset, 1.0);
    fragColor = instanceColor;
}
`

// GLSL 330 fragment shader.
const fragmentShaderCode = `#version 330
in vec4 fragColor;
out vec4 finalColor;

void main() {
    finalColor = fragColor;
}
`

func main() {
	rl.SetConfigFlags(rl.FlagWindowResizable | rl.FlagMsaa4xHint)
	rl.InitWindow(screenWidth, screenHeight, "rlgl instanced quads example")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	// -----------------------------------------------------------------
	// Quad geometry (XZ plane, unit-sized)
	// -----------------------------------------------------------------
	type QuadVertex struct {
		Pos rl.Vector3
		Tex rl.Vector2
	}
	// Texcords are implicitly zero, because we use the default white texture.
	quadVertices := []QuadVertex{
		{Pos: rl.NewVector3(-0.5, 0, -0.5)}, // 0 top-left
		{Pos: rl.NewVector3(0.5, 0, -0.5)},  // 1 top-right
		{Pos: rl.NewVector3(-0.5, 0, 0.5)},  // 2 bottom-left
		{Pos: rl.NewVector3(0.5, 0, 0.5)},   // 3 bottom-right
	}
	quadIndices := []uint16{
		0, 2, 1, // first triangle (CCW from above)
		1, 2, 3, // second triangle
	}

	// -----------------------------------------------------------------
	// Instance data: offsets and colors for 625 instances
	// -----------------------------------------------------------------
	instanceOffsets := make([]rl.Vector3, instanceCount)
	instanceColors := make([]rl.Vector4, instanceCount)

	// Center the grid at the origin
	originX := -float32(cols-1) * spacing / 2.0
	originZ := -float32(rows-1) * spacing / 2.0

	for row := range rows {
		for col := range cols {
			i := row*cols + col
			x := originX + float32(col)*spacing
			z := originZ + float32(row)*spacing
			instanceOffsets[i] = rl.NewVector3(x, 0, z)

			// HSV gradient: hue from column, saturation from row
			hue := float32(col) / float32(cols)
			sat := 0.5 + 0.5*float32(row)/float32(rows-1)
			r, g, b := hsvToRGB(hue, sat, 1.0)
			instanceColors[i] = rl.NewVector4(r, g, b, 1.0)
		}
	}

	// Working buffer for animated offsets
	animOffsets := make([]rl.Vector3, instanceCount)

	// -----------------------------------------------------------------
	// VAO + VBO + EBO setup
	// DONT FORGOT TO CLEANUP YOUR GPU RESOURCES
	// -----------------------------------------------------------------
	vao := rl.LoadVertexArray()
	rl.EnableVertexArray(vao)

	// Quad vertex VBO (static)
	quadVBO := rl.LoadVertexBuffer(quadVertices, false)
	defer rl.UnloadVertexBuffer(quadVBO)
	rl.SetVertexAttributes(quadVertices, []rl.VertexAttributesConfig{
		{Field: "Pos", Attribute: 0},
		{Field: "Tex", Attribute: 1},
	})

	// Instance offset VBO (dynamic, updated each frame)
	offsetVBO := rl.LoadVertexBuffer(animOffsets, true)
	defer rl.UnloadVertexBuffer(offsetVBO)
	rl.SetVertexAttribute(2, 3, rl.Float, false, 0, 0)
	rl.EnableVertexAttribute(2)
	rl.SetVertexAttributeDivisor(2, 1)

	// Instance color VBO (static)
	colorVBO := rl.LoadVertexBuffer(instanceColors, false)
	defer rl.UnloadVertexBuffer(colorVBO)
	rl.SetVertexAttribute(3, 4, rl.Float, false, 0, 0)
	rl.EnableVertexAttribute(3)
	rl.SetVertexAttributeDivisor(3, 1)

	// Index buffer (EBO)
	ebo := rl.LoadVertexBufferElement(quadIndices, false)
	defer rl.UnloadVertexBuffer(ebo)

	rl.DisableVertexArray()

	// -----------------------------------------------------------------
	// Custom shader
	// -----------------------------------------------------------------
	shaderID := rl.LoadShaderCode(vertexShaderCode, fragmentShaderCode)
	defer rl.UnloadShaderProgram(shaderID)
	mvpLoc := rl.GetLocationUniform(shaderID, "mvp")

	// -----------------------------------------------------------------
	// Camera
	// -----------------------------------------------------------------
	camera := rl.Camera3D{
		Position:   rl.NewVector3(30, 25, 30),
		Target:     rl.NewVector3(0, 0, 0),
		Up:         rl.NewVector3(0, 1, 0),
		Fovy:       45.0,
		Projection: rl.CameraPerspective,
	}

	// -----------------------------------------------------------------
	// Render loop
	// -----------------------------------------------------------------
	for !rl.WindowShouldClose() {
		rl.UpdateCamera(&camera, rl.CameraOrbital)

		// Animate instance Y-offsets with sine wave
		t := float64(rl.GetTime())
		for i := range instanceCount {
			base := instanceOffsets[i]
			// Phase based on XZ distance from origin for radial wave
			phase := float64(base.X+base.Z) * 0.3
			y := float32(math.Sin(t*2.0+phase)) * 2.0
			animOffsets[i] = rl.NewVector3(base.X, y, base.Z)
		}
		rl.UpdateVertexBuffer(offsetVBO, animOffsets, 0)

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode3D(camera)

		// Reference grid
		rl.DrawGrid(40, 1.0)

		// Build MVP from current camera matrices
		mvp := rl.MatrixMultiply(rl.GetMatrixModelview(), rl.GetMatrixProjection())

		// Flush raylib's internal batch before custom GL draws
		rl.DrawRenderBatchActive()

		// Custom instanced draw
		rl.EnableShader(shaderID)
		rl.SetUniformMatrix(mvpLoc, mvp)

		rl.EnableVertexArray(vao)
		rl.DrawVertexArrayElementsInstanced(0, 6, nil, int32(instanceCount))
		rl.DisableVertexArray()

		rl.DisableShader()

		// Restore state for raylib's internal renderer
		rl.DrawRenderBatchActive()

		rl.EndMode3D()

		// Text overlays
		rl.DrawText("rlgl Instanced Quads (625 instances)", 10, 10, 20, rl.DarkGray)
		rl.DrawText("Scroll to zoom", 10, 35, 16, rl.Gray)
		rl.DrawFPS(int32(rl.GetScreenWidth())-100, 10)

		rl.EndDrawing()
	}
}

// hsvToRGB converts HSV (all 0-1 range) to RGB (0-1 range).
func hsvToRGB(h, s, v float32) (r, g, b float32) {
	if s == 0 {
		return v, v, v
	}
	h *= 6.0
	if h >= 6.0 {
		h = 0
	}
	i := int(h)
	f := h - float32(i)
	p := v * (1.0 - s)
	q := v * (1.0 - s*f)
	t := v * (1.0 - s*(1.0-f))

	switch i {
	case 0:
		r, g, b = v, t, p
	case 1:
		r, g, b = q, v, p
	case 2:
		r, g, b = p, v, t
	case 3:
		r, g, b = p, q, v
	case 4:
		r, g, b = t, p, v
	default:
		r, g, b = v, p, q
	}
	return
}
