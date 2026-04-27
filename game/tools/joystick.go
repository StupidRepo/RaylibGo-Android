package tools

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Joystick represents a virtual joystick.
type Joystick struct {
	Position     rl.Vector2
	Radius       float32
	HandleRadius float32
	Value        rl.Vector2
	touchID      int32
	active       bool
}

// NewJoystick creates a new joystick.
func NewJoystick(x, y, radius, handleRadius float32) *Joystick {
	return &Joystick{
		Position:     rl.NewVector2(x, y),
		Radius:       radius,
		HandleRadius: handleRadius,
		touchID:      -1,
	}
}

// Update updates the joystick's state based on touch input.
func (j *Joystick) Update() {
	foundTouch := false
	for i := int32(0); i < rl.GetTouchPointCount(); i++ {
		touchID := rl.GetTouchPointId(i)
		touchPos := rl.GetTouchPosition(i)

		if j.active {
			if j.touchID == touchID {
				j.updateValue(touchPos)
				foundTouch = true
			}
		} else {
			if rl.CheckCollisionPointCircle(touchPos, j.Position, j.Radius) {
				j.active = true
				j.touchID = touchID
				j.updateValue(touchPos)
				foundTouch = true
			}
		}
	}

	if j.active && !foundTouch {
		j.reset()
	}
}

func (j *Joystick) updateValue(touchPos rl.Vector2) {
	delta := rl.Vector2Subtract(touchPos, j.Position)
	distance := rl.Vector2Length(delta)

	if distance > j.Radius {
		j.Value = rl.Vector2Scale(rl.Vector2Normalize(delta), 1.0)
	} else {
		j.Value = rl.Vector2Scale(delta, 1.0/j.Radius)
	}
}

func (j *Joystick) reset() {
	j.active = false
	j.touchID = -1
	j.Value = rl.Vector2Zero()
}

// Draw renders the joystick on the screen.
func (j *Joystick) Draw() {
	// Draw base
	rl.DrawCircleV(j.Position, j.Radius, rl.NewColor(0, 0, 0, 100))

	// Draw handle
	handlePos := rl.Vector2Add(j.Position, rl.Vector2Scale(j.Value, j.Radius))
	rl.DrawCircleV(handlePos, j.HandleRadius, rl.NewColor(0, 0, 0, 150))
}

// GetValue returns the joystick's current value.
func (j *Joystick) GetValue() rl.Vector2 {
	// Invert Y-axis to match typical game coordinate systems (up is positive)
	return rl.NewVector2(j.Value.X, -j.Value.Y)
}

// IsActive returns true if the joystick is currently active.
func (j *Joystick) IsActive() bool {
	return j.active
}
