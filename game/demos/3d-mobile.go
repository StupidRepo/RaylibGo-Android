package demos

import (
	"fmt"
	"math"

	ps "RaylibGoGame/platformspecifics"
	"RaylibGoGame/tools"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	gravity      float32 = 32.0
	maxSpeed     float32 = 20.0
	crouchSpeed  float32 = 5.0
	jumpForce    float32 = 12.0
	maxAccel     float32 = 300.0
	friction     float32 = 0.95
	airDrag      float32 = 0.98
	control      float32 = 15.0
	crouchHeight float32 = 0.0
	standHeight  float32 = 1.0
	bottomHeight float32 = 0.5

	cameraFovDefault float32 = 60.0
	cameraFovWalk    float32 = 55.0
)

type body struct {
	position   rl.Vector3
	velocity   rl.Vector3
	dir        rl.Vector3
	isGrounded bool
}

type ThreeDMobileDemo struct {
	Platform *ps.Platform

	screenWidth  int32
	screenHeight int32

	camera       rl.Camera3D
	player       body
	lookRotation rl.Vector2
	headTimer    float32
	walkLerp     float32
	headLerp     float32
	lean         rl.Vector2

	moveJoystick     *tools.Joystick
	lastLookTouchPos rl.Vector2
}

func (d *ThreeDMobileDemo) GetSpec() DemoSpec {
	return DemoSpec{
		Name:    "3D Mobile First-Person Game",
		Summary: "A simple 3D demo optimized for mobile devices.",
	}
}

func (d *ThreeDMobileDemo) Init() {
	d.screenWidth, d.screenHeight = d.Platform.GetWindowSize()

	d.camera = rl.Camera3D{
		Fovy:       cameraFovDefault,
		Projection: rl.CameraPerspective,
	}
	d.headLerp = standHeight
	d.updateCameraBasePosition()
	d.updateCameraFPS()

	joystickRadius := tools.UiPxToDp(125)
	joystickHandleRadius := tools.UiPxToDp(25)
	joystickDeadzone := tools.UiPxToDp(250)

	d.moveJoystick = tools.NewJoystick(
		joystickDeadzone,
		float32(d.screenHeight)-joystickDeadzone,
		joystickRadius,
		joystickHandleRadius,
	)
}

func (d *ThreeDMobileDemo) Update(currentWidth, currentHeight int32) {
	d.screenWidth = currentWidth
	d.screenHeight = currentHeight

	joystickDeadzone := tools.UiPxToDp(250)
	d.moveJoystick.Position = rl.NewVector2(joystickDeadzone, float32(d.screenHeight)-joystickDeadzone)

	d.moveJoystick.Update()

	d.updateLookRotation()
	moveVec := d.moveJoystick.GetValue()

	// Add keyboard input for desktop
	if rl.IsKeyDown(rl.KeyD) {
		moveVec.X += 1.0
	}
	if rl.IsKeyDown(rl.KeyA) {
		moveVec.X -= 1.0
	}
	if rl.IsKeyDown(rl.KeyW) {
		moveVec.Y += 1.0
	}
	if rl.IsKeyDown(rl.KeyS) {
		moveVec.Y -= 1.0
	}

	if rl.Vector2Length(moveVec) > 1.0 {
		moveVec = rl.Vector2Normalize(moveVec)
	}
	sideway, forward := moveVec.X, moveVec.Y

	crouching := rl.IsKeyDown(rl.KeyLeftControl)
	d.updateBody(
		d.lookRotation.X,
		sideway,
		forward,
		rl.IsKeyPressed(rl.KeySpace),
		crouching,
	)

	delta := rl.GetFrameTime()
	d.updateHeadHeight(delta, crouching)
	d.updateCameraBasePosition()
	d.updateWalkEffects(delta, sideway, forward)
	d.updateLean(delta, sideway, forward)
	d.updateCameraFPS()
}

func (d *ThreeDMobileDemo) Draw() {
	rl.ClearBackground(rl.RayWhite)
	rl.BeginMode3D(d.camera)
	d.drawLevel()
	rl.EndMode3D()

	if d.Platform.GetOS() == ps.PlatformMobile {
		d.moveJoystick.Draw()
	}

	insets := d.Platform.GetInsets()
	fontSize := tools.UiPxToDp(24)

	fps := rl.GetFPS()
	color := tools.GradientColor(rl.Red, rl.Green, 30, 144, float32(fps))
	tools.DrawTextCenter(fmt.Sprintf("%d FPS", fps), float32(d.screenWidth), float32(insets.Top)+fontSize, fontSize, color)
}

func (d *ThreeDMobileDemo) Deinit() {}

func (d *ThreeDMobileDemo) updateLookRotation() {
	if d.Platform.GetOS() == ps.PlatformMobile {
		foundLookTouch := false
		for i := int32(0); i < rl.GetTouchPointCount(); i++ {
			touchPos := rl.GetTouchPosition(i)

			if touchPos.X > float32(d.screenWidth)/2 {
				if d.lastLookTouchPos.X != 0 || d.lastLookTouchPos.Y != 0 {
					delta := rl.Vector2Subtract(touchPos, d.lastLookTouchPos)
					d.lookRotation.X -= delta.X * 0.006
					d.lookRotation.Y += delta.Y * 0.006
				}
				d.lastLookTouchPos = touchPos
				foundLookTouch = true
				break
			}
		}

		if !foundLookTouch {
			d.lastLookTouchPos = rl.Vector2Zero()
		}
	} else {
		// On other platforms, process all mouse input.
		mouseDelta := rl.GetMouseDelta()
		d.lookRotation.X -= mouseDelta.X * 0.006
		d.lookRotation.Y += mouseDelta.Y * 0.006
	}
}

func (d *ThreeDMobileDemo) updateHeadHeight(delta float32, crouching bool) {
	target := standHeight
	if crouching {
		target = crouchHeight
	}
	d.headLerp = rl.Lerp(d.headLerp, target, 20.0*delta)
}

func (d *ThreeDMobileDemo) updateCameraBasePosition() {
	d.camera.Position = rl.NewVector3(
		d.player.position.X,
		d.player.position.Y+(bottomHeight+d.headLerp),
		d.player.position.Z,
	)
}

func (d *ThreeDMobileDemo) updateWalkEffects(delta float32, forward, sideway float32) {
	moving := d.player.isGrounded && (forward != 0 || sideway != 0)
	if moving {
		d.headTimer += delta * 3.0
		d.walkLerp = rl.Lerp(d.walkLerp, 1.0, 10.0*delta)
		d.camera.Fovy = rl.Lerp(d.camera.Fovy, cameraFovWalk, 5.0*delta)
		return
	}

	d.walkLerp = rl.Lerp(d.walkLerp, 0.0, 10.0*delta)
	d.camera.Fovy = rl.Lerp(d.camera.Fovy, cameraFovDefault, 5.0*delta)
}

func (d *ThreeDMobileDemo) updateLean(delta float32, sideway, forward float32) {
	d.lean.X = rl.Lerp(d.lean.X, sideway*0.02, 10.0*delta)
	d.lean.Y = rl.Lerp(d.lean.Y, forward*0.015, 10.0*delta)
}

func (d *ThreeDMobileDemo) updateBody(rot float32, side, forward float32, jumpPressed, crouchHold bool) {
	input := rl.NewVector2(side, -forward)
	delta := rl.GetFrameTime()

	if !d.player.isGrounded {
		d.player.velocity.Y -= gravity * delta
	}

	if d.player.isGrounded && jumpPressed {
		d.player.velocity.Y = jumpForce
		d.player.isGrounded = false
	}

	front := rl.NewVector3(
		float32(math.Sin(float64(rot))),
		0.0,
		float32(math.Cos(float64(rot))),
	)
	right := rl.NewVector3(
		float32(math.Cos(float64(-rot))),
		0.0,
		float32(math.Sin(float64(-rot))),
	)

	desiredDir := rl.NewVector3(
		input.X*right.X+input.Y*front.X,
		0.0,
		input.X*right.Z+input.Y*front.Z,
	)
	d.player.dir = rl.Vector3Lerp(d.player.dir, desiredDir, control*delta)

	decel := airDrag
	if d.player.isGrounded {
		decel = friction
	}
	hvel := rl.NewVector3(d.player.velocity.X*decel, 0.0, d.player.velocity.Z*decel)

	if rl.Vector3Length(hvel) < (maxSpeed * 0.01) {
		hvel = rl.NewVector3(0, 0, 0)
	}

	speed := rl.Vector3DotProduct(hvel, d.player.dir)

	maxSpd := maxSpeed
	if crouchHold {
		maxSpd = crouchSpeed
	}

	accel := rl.Clamp(maxSpd-speed, 0.0, maxAccel*delta)
	hvel.X += d.player.dir.X * accel
	hvel.Z += d.player.dir.Z * accel

	d.player.velocity.X = hvel.X
	d.player.velocity.Z = hvel.Z

	d.player.position.X += d.player.velocity.X * delta
	d.player.position.Y += d.player.velocity.Y * delta
	d.player.position.Z += d.player.velocity.Z * delta

	if d.player.position.Y <= 0.0 {
		d.player.position.Y = 0.0
		d.player.velocity.Y = 0.0
		d.player.isGrounded = true
	}
}

func (d *ThreeDMobileDemo) updateCameraFPS() {
	up := rl.NewVector3(0.0, 1.0, 0.0)
	targetOffset := rl.NewVector3(0.0, 0.0, -1.0)

	yaw := rl.Vector3RotateByAxisAngle(targetOffset, up, d.lookRotation.X)

	maxAngleUp := rl.Vector3Angle(up, yaw) - 0.001
	if -d.lookRotation.Y > maxAngleUp {
		d.lookRotation.Y = -maxAngleUp
	}

	maxAngleDown := -rl.Vector3Angle(rl.Vector3Negate(up), yaw) + 0.001
	if -d.lookRotation.Y < maxAngleDown {
		d.lookRotation.Y = -maxAngleDown
	}

	right := rl.Vector3Normalize(rl.Vector3CrossProduct(yaw, up))

	pitchAngle := rl.Clamp(
		-d.lookRotation.Y-d.lean.Y,
		-rl.Pi/2+0.0001,
		rl.Pi/2-0.0001,
	)
	pitch := rl.Vector3RotateByAxisAngle(yaw, right, pitchAngle)

	headSin := float32(math.Sin(float64(d.headTimer * rl.Pi)))
	headCos := float32(math.Cos(float64(d.headTimer * rl.Pi)))

	const stepRotation float32 = 0.01
	d.camera.Up = rl.Vector3RotateByAxisAngle(up, pitch, headSin*stepRotation+d.lean.X)

	const (
		bobSide float32 = 0.125
		bobUp   float32 = 0.3
	)
	bobbing := rl.Vector3Scale(right, headSin*bobSide)
	bobbing.Y = float32(math.Abs(float64(headCos * bobUp)))

	d.camera.Position = rl.Vector3Add(d.camera.Position, rl.Vector3Scale(bobbing, d.walkLerp))
	d.camera.Target = rl.Vector3Add(d.camera.Position, pitch)
}

func (d *ThreeDMobileDemo) drawLevel() {
	const (
		floorExtent = 25
		tileSize    = float32(5.0)
	)
	tileColor1 := rl.NewColor(150, 200, 200, 255)

	for y := -floorExtent; y < floorExtent; y++ {
		for x := -floorExtent; x < floorExtent; x++ {
			if !isCheckerTile(x, y) {
				continue
			}

			color := rl.LightGray
			if (y&1) == 1 && (x&1) == 1 {
				color = tileColor1
			}

			rl.DrawPlane(
				rl.NewVector3(float32(x)*tileSize, 0.0, float32(y)*tileSize),
				rl.NewVector2(tileSize, tileSize),
				color,
			)
		}
	}

	towerSize := rl.NewVector3(16.0, 32.0, 16.0)
	towerColor := rl.NewColor(150, 200, 200, 255)

	towerPositions := []rl.Vector3{
		rl.NewVector3(16.0, 16.0, 16.0),
		rl.NewVector3(-16.0, 16.0, 16.0),
		rl.NewVector3(-16.0, 16.0, -16.0),
		rl.NewVector3(16.0, 16.0, -16.0),
	}

	for _, pos := range towerPositions {
		rl.DrawCubeV(pos, towerSize, towerColor)
		rl.DrawCubeWiresV(pos, towerSize, rl.DarkBlue)
	}

	rl.DrawSphereEx(
		rl.NewVector3(300.0, 300.0, 0.0),
		100.0,
		24,
		52,
		rl.Orange,
	)
}

func isCheckerTile(x, y int) bool {
	return ((x&1) == 1 && (y&1) == 1) || ((x&1) == 0 && (y&1) == 0)
}
