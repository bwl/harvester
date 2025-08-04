package main

import (
	"harvester/pkg/engine"
	"math"
	"testing"
)

func TestPhysicsBasics(t *testing.T) {
	// Create test model
	model := createTestModel()

	// Test 1: Ship should not move without thrust
	for i := 0; i < 10; i++ {
		model.updatePhysics()
	}

	speed := getSpeed(model)
	if speed > 0.1 {
		t.Errorf("Ship should not move without thrust, but speed is %.3f", speed)
	}
}

func TestThrustIncreasesSpeed(t *testing.T) {
	model := createTestModel()

	// Apply thrust for several frames
	model.thrust = true
	initialSpeed := getSpeed(model)

	for i := 0; i < 5; i++ {
		model.updatePhysics()
	}

	finalSpeed := getSpeed(model)
	if finalSpeed <= initialSpeed {
		t.Errorf("Thrust should increase speed: initial=%.3f, final=%.3f", initialSpeed, finalSpeed)
	}
}

func TestBrakingDecreasesSpeed(t *testing.T) {
	model := createTestModel()

	// Build up some speed first
	model.thrust = true
	for i := 0; i < 10; i++ {
		model.updatePhysics()
	}
	model.thrust = false

	speedBeforeBraking := getSpeed(model)

	// Apply braking
	model.braking = true
	for i := 0; i < 5; i++ {
		model.updatePhysics()
	}

	speedAfterBraking := getSpeed(model)

	if speedAfterBraking >= speedBeforeBraking {
		t.Errorf("Braking should decrease speed: before=%.3f, after=%.3f", speedBeforeBraking, speedAfterBraking)
	}
}

func TestReleasingBrakeDoesNotIncreaseSpeed(t *testing.T) {
	model := createTestModel()

	// Build up speed, then brake
	model.thrust = true
	for i := 0; i < 10; i++ {
		model.updatePhysics()
	}
	model.thrust = false

	model.braking = true
	for i := 0; i < 5; i++ {
		model.updatePhysics()
	}

	speedWhileBraking := getSpeed(model)

	// Release brake
	model.braking = false
	model.updatePhysics()

	speedAfterReleasingBrake := getSpeed(model)

	if speedAfterReleasingBrake > speedWhileBraking {
		t.Errorf("Releasing brake should not increase speed: while_braking=%.3f, after_release=%.3f",
			speedWhileBraking, speedAfterReleasingBrake)
	}
}

func TestMaxSpeedLimit(t *testing.T) {
	model := createTestModel()

	// Apply thrust for a long time
	model.thrust = true
	for i := 0; i < 1000; i++ {
		model.updatePhysics()
	}

	speed := getSpeed(model)
	maxSpeed := 80.0

	if speed > maxSpeed+0.1 { // Small tolerance
		t.Errorf("Speed should not exceed max speed: speed=%.3f, max=%.3f", speed, maxSpeed)
	}
}

func TestBrakingEventuallyStops(t *testing.T) {
	model := createTestModel()

	// Build up speed
	model.thrust = true
	for i := 0; i < 20; i++ {
		model.updatePhysics()
	}
	model.thrust = false

	// Brake for a long time
	model.braking = true
	for i := 0; i < 1000; i++ {
		model.updatePhysics()
	}

	speed := getSpeed(model)
	if speed > 0.1 {
		t.Errorf("Long braking should stop the ship, but speed is %.3f", speed)
	}
}

// Helper functions
func createTestModel() *FreshSpaceModel {
	bs := engine.New(nil)
	return &FreshSpaceModel{
		world:  bs.World,
		player: bs.Player,
		width:  80,
		height: 24,
		angle:  0.0,
	}
}

func getSpeed(model *FreshSpaceModel) float64 {
	return math.Sqrt(model.velocity.x*model.velocity.x + model.velocity.y*model.velocity.y)
}
