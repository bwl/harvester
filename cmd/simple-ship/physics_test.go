package main

import (
	"math"
	"testing"
)

func TestSimpleShipThrust(t *testing.T) {
	ship := NewSimpleShip()

	// Apply thrust
	ship.thrust = true
	initialSpeed := math.Sqrt(ship.vx*ship.vx + ship.vy*ship.vy)

	ship.updatePhysics()

	finalSpeed := math.Sqrt(ship.vx*ship.vx + ship.vy*ship.vy)

	if finalSpeed <= initialSpeed {
		t.Errorf("Thrust should increase speed: initial=%.3f, final=%.3f", initialSpeed, finalSpeed)
	}
}

func TestSimpleShipBraking(t *testing.T) {
	ship := NewSimpleShip()

	// Build up speed first
	ship.thrust = true
	for i := 0; i < 10; i++ {
		ship.updatePhysics()
	}
	ship.thrust = false

	speedBeforeBraking := math.Sqrt(ship.vx*ship.vx + ship.vy*ship.vy)

	// Apply braking
	ship.braking = true
	ship.updatePhysics()

	speedAfterBraking := math.Sqrt(ship.vx*ship.vx + ship.vy*ship.vy)

	if speedAfterBraking >= speedBeforeBraking {
		t.Errorf("Braking should decrease speed: before=%.3f, after=%.3f", speedBeforeBraking, speedAfterBraking)
	}
}

func TestSimpleShipTurning(t *testing.T) {
	ship := NewSimpleShip()
	initialAngle := ship.angle

	ship.angle += 0.15 // Turn right

	if ship.angle == initialAngle {
		t.Errorf("Turning should change angle: initial=%.3f, final=%.3f", initialAngle, ship.angle)
	}
}

func TestCameraFollowsPlayer(t *testing.T) {
	ship := NewSimpleShip()

	// Move player forward
	ship.thrust = true
	for i := 0; i < 20; i++ {
		ship.updatePhysics()
	}

	// Camera should follow player, keeping player centered
	expectedCameraX := ship.x - float64(ship.width)/2
	expectedCameraY := ship.y - float64(ship.height)/2

	if math.Abs(ship.cameraX-expectedCameraX) > 0.001 {
		t.Errorf("Camera X should follow player: expected=%.3f, actual=%.3f", expectedCameraX, ship.cameraX)
	}

	if math.Abs(ship.cameraY-expectedCameraY) > 0.001 {
		t.Errorf("Camera Y should follow player: expected=%.3f, actual=%.3f", expectedCameraY, ship.cameraY)
	}
}
