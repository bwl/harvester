package components

// SpringState stores a spring's position, velocity, and target
type SpringState struct{ Pos, Vel, Target float64 }

// SpaceFlightSprings holds spring states for thrust, angle, and velocity axes
type SpaceFlightSprings struct {
	Thrust SpringState
	Angle  SpringState
	VelX   SpringState
	VelY   SpringState
}
