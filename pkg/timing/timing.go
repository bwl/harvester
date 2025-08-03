package timing

import (
	"github.com/charmbracelet/harmonica"
	"time"
)

// Animation constants
const (
	// TargetFPS is the target frame rate for all animations
	TargetFPS = 60
)

// Animation variables
var (
	// HarmonicaFPS is the harmonica FPS configuration used throughout the application
	HarmonicaFPS = harmonica.FPS(TargetFPS)
)

// GlobalTimer provides centralized timing for the entire application
type GlobalTimer struct {
	startTime   time.Time
	currentTick uint64
	lastUpdate  time.Time
	deltaTime   float64
	frameCount  uint64
	paused      bool
}

// NewGlobalTimer creates a new global timer
func NewGlobalTimer() *GlobalTimer {
	now := time.Now()
	return &GlobalTimer{
		startTime:   now,
		lastUpdate:  now,
		currentTick: 0,
		frameCount:  0,
		paused:      false,
	}
}

// Update advances the timer by one frame
func (t *GlobalTimer) Update() {
	if t.paused {
		return
	}

	now := time.Now()
	t.deltaTime = now.Sub(t.lastUpdate).Seconds()
	t.lastUpdate = now
	t.currentTick++
	t.frameCount++
}

// Tick returns the current tick count
func (t *GlobalTimer) Tick() uint64 {
	return t.currentTick
}

// DeltaTime returns the time since last update in seconds
func (t *GlobalTimer) DeltaTime() float64 {
	return t.deltaTime
}

// FrameCount returns total frames processed
func (t *GlobalTimer) FrameCount() uint64 {
	return t.frameCount
}

// ElapsedTime returns total time since timer creation
func (t *GlobalTimer) ElapsedTime() time.Duration {
	return time.Since(t.startTime)
}

// Pause pauses the timer
func (t *GlobalTimer) Pause() {
	t.paused = true
}

// Resume resumes the timer
func (t *GlobalTimer) Resume() {
	if t.paused {
		t.paused = false
		t.lastUpdate = time.Now() // Reset to avoid large delta
	}
}

// IsPaused returns whether the timer is paused
func (t *GlobalTimer) IsPaused() bool {
	return t.paused
}

// Global timer instance
var globalTimer *GlobalTimer

// InitGlobalTimer initializes the global timer
func InitGlobalTimer() {
	globalTimer = NewGlobalTimer()
}

// GetGlobalTimer returns the global timer instance
func GetGlobalTimer() *GlobalTimer {
	if globalTimer == nil {
		InitGlobalTimer()
	}
	return globalTimer
}

// Convenience functions for accessing global timer
func Tick() uint64 {
	return GetGlobalTimer().Tick()
}

func DeltaTime() float64 {
	return GetGlobalTimer().DeltaTime()
}

func FrameCount() uint64 {
	return GetGlobalTimer().FrameCount()
}

func UpdateGlobalTimer() {
	GetGlobalTimer().Update()
}

// Animation utilities

// AnimationState tracks frame-based animations
type AnimationState struct {
	startFrame   uint64
	duration     uint64
	currentFrame uint64
	loop         bool
	finished     bool
}

// NewAnimation creates a new animation state
func NewAnimation(duration uint64, loop bool) *AnimationState {
	return &AnimationState{
		startFrame: Tick(),
		duration:   duration,
		loop:       loop,
		finished:   false,
	}
}

// Update updates the animation state
func (a *AnimationState) Update() {
	if a.finished && !a.loop {
		return
	}

	currentTick := Tick()
	elapsed := currentTick - a.startFrame

	if elapsed >= a.duration {
		if a.loop {
			// Restart animation
			a.startFrame = currentTick
			a.currentFrame = 0
		} else {
			// Finish animation
			a.currentFrame = a.duration
			a.finished = true
		}
	} else {
		a.currentFrame = elapsed
	}
}

// Progress returns animation progress (0.0 to 1.0)
func (a *AnimationState) Progress() float64 {
	if a.duration == 0 {
		return 1.0
	}
	return float64(a.currentFrame) / float64(a.duration)
}

// Frame returns the current animation frame
func (a *AnimationState) Frame() uint64 {
	return a.currentFrame
}

// IsFinished returns whether the animation is complete
func (a *AnimationState) IsFinished() bool {
	return a.finished
}

// Reset restarts the animation
func (a *AnimationState) Reset() {
	a.startFrame = Tick()
	a.currentFrame = 0
	a.finished = false
}

// Easing functions for smooth animations

// EaseInOut applies ease-in-out easing to a 0-1 progress value
func EaseInOut(t float64) float64 {
	if t < 0.5 {
		return 2 * t * t
	}
	return -1 + (4-2*t)*t
}

// EaseIn applies ease-in easing to a 0-1 progress value
func EaseIn(t float64) float64 {
	return t * t
}

// EaseOut applies ease-out easing to a 0-1 progress value
func EaseOut(t float64) float64 {
	return t * (2 - t)
}

// Linear returns the input unchanged (linear progression)
func Linear(t float64) float64 {
	return t
}
