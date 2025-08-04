package main

/*
#include <stdlib.h>

typedef struct {
    int x, y;
    int glyph;
    int foregroundR, foregroundG, foregroundB;
    int backgroundR, backgroundG, backgroundB;
    int style;
    float alpha;
} CGlyph;

typedef struct {
    CGlyph* glyphs;
    int width, height;
    int count;
} CGlyphMatrix;

extern void updateGameView(CGlyphMatrix matrix);
*/
import "C"
import (
	"harvester/pkg/components"
	"harvester/pkg/ecs"
	"harvester/pkg/engine"
	"harvester/pkg/rendering"
	"harvester/pkg/systems"
	"math"
	"unsafe"
)

// Global game state
type DesktopGame struct {
	x, y      float64 // position
	vx, vy    float64 // velocity
	angle     float64 // rotation
	width     int
	height    int
	thrust    bool
	braking   bool
	cameraX   float64
	cameraY   float64
	world     *ecs.World
	player    ecs.Entity
	render    *systems.Render
	scheduler *ecs.SchedulerWithContext
}

var globalGame *DesktopGame

//export initGame
func initGame(width, height C.int) {
	bs := engine.New(nil)

	globalGame = &DesktopGame{
		x: 40, y: 12,
		angle: -math.Pi / 2,
		width: int(width), height: int(height),
		world:     bs.World,
		player:    bs.Player,
		render:    bs.Render,
		scheduler: bs.Scheduler,
	}

	// Initialize camera
	globalGame.cameraX = globalGame.x - float64(globalGame.width)/2
	globalGame.cameraY = globalGame.y - float64(globalGame.height)/2

	// Set player position in ECS
	if pos, ok := ecs.Get[components.Position](globalGame.world, globalGame.player); ok {
		pos.X = globalGame.x
		pos.Y = globalGame.y
		ecs.Add(globalGame.world, globalGame.player, pos)
	}

	// Make sure we're in space layer
	ctx := ecs.GetWorldContext(globalGame.world)
	ctx.CurrentLayer = ecs.LayerSpace
	ecs.SetWorldContext(globalGame.world, ctx)
}

//export updateGame
func updateGame(dt C.float, thrustInput, brakeInput, leftInput, rightInput C.int) {
	if globalGame == nil {
		return
	}

	// Update input state
	globalGame.thrust = thrustInput != 0
	globalGame.braking = brakeInput != 0

	// Handle rotation
	if leftInput != 0 {
		globalGame.angle -= 0.15
	}
	if rightInput != 0 {
		globalGame.angle += 0.15
	}

	// Update physics (similar to simple-ship)
	globalGame.updatePhysics()

	// Run ECS scheduler
	globalGame.scheduler.Update(float64(dt), globalGame.world)
}

func (g *DesktopGame) updatePhysics() {
	const dt = 0.016
	const thrustPower = 30.0
	const friction = 0.999
	const brakeFriction = 0.92
	const maxSpeed = 80.0

	// Apply thrust
	if g.thrust {
		thrustX := math.Cos(g.angle) * thrustPower * dt
		thrustY := math.Sin(g.angle) * thrustPower * dt
		g.vx += thrustX
		g.vy += thrustY
	}

	// Apply friction
	if g.braking {
		g.vx *= brakeFriction
		g.vy *= brakeFriction
	} else {
		g.vx *= friction
		g.vy *= friction
	}

	// Speed limit
	speed := math.Sqrt(g.vx*g.vx + g.vy*g.vy)
	if speed > maxSpeed {
		g.vx = (g.vx / speed) * maxSpeed
		g.vy = (g.vy / speed) * maxSpeed
	}

	// Stop if very slow
	if speed < 0.1 {
		g.vx = 0
		g.vy = 0
	}

	// Update position
	g.x += g.vx * dt
	g.y += g.vy * dt

	// Update camera
	g.cameraX = g.x - float64(g.width)/2
	g.cameraY = g.y - float64(g.height)/2

	// Sync ECS position
	if pos, ok := ecs.Get[components.Position](g.world, g.player); ok {
		pos.X = g.x
		pos.Y = g.y
		ecs.Add(g.world, g.player, pos)
	}

	// Update camera in ECS
	if cam, ok := ecs.Get[components.Camera](g.world, g.player); ok {
		cam.X = int(g.cameraX)
		cam.Y = int(g.cameraY)
		ecs.Add(g.world, g.player, cam)
	}
}

//export getGlyphMatrix
func getGlyphMatrix() C.CGlyphMatrix {
	if globalGame == nil {
		return C.CGlyphMatrix{glyphs: nil, width: 0, height: 0, count: 0}
	}

	// Build glyph matrix like simple-ship
	glyphs := globalGame.buildGameGlyphs()
	if glyphs == nil {
		return C.CGlyphMatrix{glyphs: nil, width: 0, height: 0, count: 0}
	}

	return convertToC(glyphs)
}

func (g *DesktopGame) buildGameGlyphs() [][]rendering.Glyph {
	if g.width <= 0 || g.height <= 0 {
		return nil
	}

	// Use unified render system
	g.render.Update(0, g.world)
	cam, _ := ecs.Get[components.Camera](g.world, g.player)
	mx0, my0 := cam.X, cam.Y

	glyphs := make([][]rendering.Glyph, g.height)
	for y := 0; y < g.height; y++ {
		row := make([]rendering.Glyph, g.width)
		for x := 0; x < g.width; x++ {
			row[x] = rendering.Glyph{Char: '.'}
		}
		glyphs[y] = row
	}

	// Render all ECS entities
	for _, d := range g.render.Output {
		x := d.X - mx0
		y := d.Y - my0
		if x >= 0 && y >= 0 && x < g.width && y < g.height {
			glyphs[y][x] = rendering.Glyph{Char: rune(d.Glyph)}
		}
	}

	return glyphs
}

func convertToC(grid [][]rendering.Glyph) C.CGlyphMatrix {
	height := len(grid)
	width := 0
	count := 0
	for y := 0; y < height; y++ {
		if len(grid[y]) > width {
			width = len(grid[y])
		}
		count += len(grid[y])
	}
	if count == 0 {
		return C.CGlyphMatrix{glyphs: nil, width: C.int(width), height: C.int(height), count: 0}
	}
	sz := C.size_t(count) * C.size_t(unsafe.Sizeof(C.CGlyph{}))
	ptr := C.malloc(sz)
	glyphs := (*[1 << 30]C.CGlyph)(ptr)[:count:count]
	idx := 0
	for y := 0; y < height; y++ {
		row := grid[y]
		for x := 0; x < len(row); x++ {
			g := row[x]
			glyphs[idx] = C.CGlyph{
				x:           C.int(x),
				y:           C.int(y),
				glyph:       C.int(g.Char),
				foregroundR: C.int(g.Foreground.R),
				foregroundG: C.int(g.Foreground.G),
				foregroundB: C.int(g.Foreground.B),
				backgroundR: C.int(g.Background.R),
				backgroundG: C.int(g.Background.G),
				backgroundB: C.int(g.Background.B),
				style:       C.int(g.Style),
				alpha:       C.float(g.Alpha),
			}
			idx++
		}
	}
	return C.CGlyphMatrix{glyphs: (*C.CGlyph)(ptr), width: C.int(width), height: C.int(height), count: C.int(count)}
}

func main() {}
