// This example demonstrates the low-level rlgl vertex buffer bindings:
//   - LoadVertexBuffer / LoadVertexBufferElement: upload vertex + index data to GPU
//   - SetVertexAttribute: configure vertex attribute pointers
//   - SetVertexAttributes: configure vertex attribute pointers (interleaved buffers)
//   - SetVertexAttributeDefault: provide default values for unused attributes
//   - DrawVertexArrayElements: indexed draw (quad)
//   - DrawVertexArray: non-indexed draw (triangle)
//   - UpdateVertexBuffer: animate vertex colors each frame
package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

func main() {
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(screenWidth, screenHeight, "rlgl vertex buffer example")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	// ---------------------------------------------------------------
	// Quad (indexed draw) - centered on left half of screen
	// Vertices ordered CCW in NDC (since ortho flips Y, screen-CW = NDC-CW,
	// so we use CCW winding via index order)
	// ---------------------------------------------------------------

	// QuadVertex stores Pos and Texture coordinates (UV)
	type QuadVertex struct {
		Pos rl.Vector3
		Tex rl.Vector2 // texture cords are unused. But we stil define them so we can bind.
	}
	// Texcords are implicitly zero, because we use the default white texture.
	quadVertices := []QuadVertex{
		// top left vertex,
		{Pos: rl.NewVector3(150, 150, 0)}, // 0 top-left
		// top right vertex,
		{Pos: rl.NewVector3(450, 150, 0)}, // 1 top-right
		// bottom right vertex,
		{Pos: rl.NewVector3(450, 450, 0)}, // 2 bottom-right
		// bottom left vertex,
		{Pos: rl.NewVector3(150, 450, 0)}, // 3 bottom-left
	}
	// quad Colors are stored in a seperate buffer because we will be updating them.
	// Keeping it in a seperate buffer allows us to only reupload the updated colors to the GPU.
	quadColors := []rl.Vector4{
		rl.NewVector4(1, 0, 0, 1), // top-left: red
		rl.NewVector4(0, 1, 0, 1), // top-right: green
		rl.NewVector4(0, 0, 1, 1), // bottom-right: blue
		rl.NewVector4(1, 1, 0, 1), // bottom-left: yellow
	}
	// Describe how to draw those 4 vertices. We will draw them 6 times to make a square (quad).
	// Look at the indices above.
	// CCW winding in NDC (with Y-flip ortho: reverse the original CW order)
	quadIndices := []uint16{
		0, 2, 1, // first triangle (CCW in NDC)
		0, 3, 2, // second triangle (CCW in NDC)
	}

	// ---------------------------------------------------------------
	// Quad VAO, VBO, and EBO Creations
	// DONT FORGOT TO CLEANUP YOUR GPU RESOURCES
	// ---------------------------------------------------------------
	// Create quad VAO for binding attributes
	quadVAO := rl.LoadVertexArray()
	rl.EnableVertexArray(quadVAO)

	// Create quad VBO for storing vertices
	quadVBO := rl.LoadVertexBuffer(quadVertices, false)
	defer rl.UnloadVertexBuffer(quadVBO)
	// Bind attributes to quadVAO
	rl.SetVertexAttributes(quadVertices, []rl.VertexAttributesConfig{
		// Position VBO (attribute 0, vec3)
		{Field: "Pos", Attribute: 0},
		// Texcoord VBO (attribute 1, vec2)
		{Field: "Tex", Attribute: 1},
	})
	// Create quad VBO for storing vertex colors. Explanation above.
	quadColorVBO := rl.LoadVertexBuffer(quadColors, true)
	defer rl.UnloadVertexBuffer(quadColorVBO)

	// Color VBO (attribute 3, vec4) - dynamic for animation
	rl.SetVertexAttribute(3, 4, rl.Float, false, 0, 0)
	rl.EnableVertexAttribute(3)

	// Index buffer (EBO) (for indexed drawing of vertices)
	quadIBO := rl.LoadVertexBufferElement(quadIndices, false)
	defer rl.UnloadVertexBuffer(quadIBO)
	rl.DisableVertexArray() // disable quadVAO that was enabled on creation.

	// ---------------------------------------------------------------
	// Triangle (non-indexed draw) - on the right side
	// Vertices in CCW order in NDC (with Y-flip ortho)
	// ---------------------------------------------------------------
	// We dont update the triangle's colors, so keep everything in 1 buffer
	type TriangleVertex struct {
		Pos   rl.Vector3
		Tex   rl.Vector2 // texture cords are unused. But we stil define them so we can bind.
		Color rl.Vector4 // RGBA represented with floats.
	}
	// define how the fields will be bound. Define them up here for clarity.
	triangleVertexAttributesConfig := []rl.VertexAttributesConfig{
		// Position VBO (attribute 0, vec3)
		{Field: "Pos", Attribute: 0},
		// Texcoord VBO (attribute 1, vec2)
		{Field: "Tex", Attribute: 1},
		// Color VBO (attribute 3, vec4)
		{Field: "Color", Attribute: 3},
	}
	// actual VBO data that will be uploaded to the GPU only once.
	triangleVertices := []TriangleVertex{
		// (CCW in NDC with Y-flip)       Magenta
		{Pos: rl.NewVector3(575, 200, 0), Color: rl.NewVector4(1, 0, 1, 1)}, // top
		// left                           Cyan
		{Pos: rl.NewVector3(575, 450, 0), Color: rl.NewVector4(0, 1, 1, 1)}, // left
		// right                          White
		{Pos: rl.NewVector3(700, 450, 0), Color: rl.NewVector4(1, 1, 1, 1)}, // right
	}
	// ---------------------------------------------------------------
	// Triangle VAO, and VBO creation
	// DONT FORGOT TO CLEANUP YOUR GPU RESOURCES
	// ---------------------------------------------------------------
	// Create triangle VAO
	triangleVAO := rl.LoadVertexArray()
	rl.EnableVertexArray(triangleVAO)
	// Upload vertices
	triangleVBO := rl.LoadVertexBuffer(triangleVertices, false)
	defer rl.UnloadVertexBuffer(triangleVBO)
	// Bind the attributes.
	rl.SetVertexAttributes(triangleVertices, triangleVertexAttributesConfig)
	rl.DisableVertexArray()

	// ---------------------------------------------------------------
	// Shader + MVP setup
	// ---------------------------------------------------------------
	defaultShaderID := rl.GetShaderIdDefault()
	mvpLoc := rl.GetLocationUniform(defaultShaderID, "mvp")
	colDiffuseLoc := rl.GetLocationUniform(defaultShaderID, "colDiffuse")
	defaultTexID := rl.GetTextureIdDefault()

	fmt.Printf("DEBUG: shaderID=%d mvpLoc=%d colDiffuseLoc=%d texID=%d\n",
		defaultShaderID, mvpLoc, colDiffuseLoc, defaultTexID)
	fmt.Printf("DEBUG: quadVAO=%d quadVBO=%d triangleVAO=%d triangleVBO=%d\n",
		quadVAO, quadVBO, triangleVAO, triangleVBO)

	// Orthographic projection: screen-space coordinates, origin top-left
	// NOTE: rl.MatrixMultiply(rl.GetMatrixModelview(), rl.GetMatrixProjection()) could be used to work with rl camera
	// See examples/others/rlgl_instanced_quad/main.go
	mvpMatrix := rl.MatrixOrtho(0, screenWidth, screenHeight, 0, -1, 1)

	// White diffuse color so vertex colors pass through unmodified
	whiteDiffuse := []float32{1, 1, 1, 1}

	// Working buffer for animated quad colors
	animColors := make([]rl.Vector4, 4)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		// ---------------------------------------------------------------
		// Animated quad colors
		// Thanks Inigo: https://iquilezles.org/articles/palettes/
		// ---------------------------------------------------------------
		t := float32(rl.GetTime())
		for i := range 4 {
			phase := float32(i) * math.Pi / 2.0
			r := 0.5 + 0.5*float32(math.Sin(float64(t*2.0+phase)))
			g := 0.5 + 0.5*float32(math.Sin(float64(t*2.0+phase+math.Pi*2.0/3.0)))
			b := 0.5 + 0.5*float32(math.Sin(float64(t*2.0+phase+math.Pi*4.0/3.0)))
			animColors[i] = rl.NewVector4(r, g, b, 1.0)
		}
		rl.UpdateVertexBuffer(quadColorVBO, animColors, 0)

		// Flush raylib's internal batch before custom GL draws
		rl.DrawRenderBatchActive()

		// Set up default shader for our custom VAO draws
		rl.EnableShader(defaultShaderID)
		rl.SetUniformMatrix(mvpLoc, mvpMatrix)
		rl.SetUniform(colDiffuseLoc, whiteDiffuse, int32(rl.ShaderUniformVec4), 1)
		rl.ActiveTextureSlot(0)
		rl.EnableTexture(defaultTexID)

		// Draw quad (indexed)
		rl.EnableVertexArray(quadVAO)
		rl.DrawVertexArrayElements(0, 6, nil)
		rl.DisableVertexArray()

		// Draw triangle (non-indexed)
		rl.EnableVertexArray(triangleVAO)
		rl.DrawVertexArray(0, 3)
		rl.DisableVertexArray()

		rl.CheckErrors()

		// Restore state for raylib's internal renderer
		rl.DisableShader()
		rl.DisableTexture()
		rl.DrawRenderBatchActive()

		// Text overlays
		rl.DrawText("rlgl Vertex Buffer Bindings Test", 10, 10, 20, rl.DarkGray)
		rl.DrawFPS(int32(rl.GetScreenWidth())-100, 10)
		rl.EndDrawing()
	}
}
