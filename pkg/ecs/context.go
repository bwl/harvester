package ecs

type GameLayer int

const (
	LayerSpace GameLayer = iota
	LayerPlanetSurface
	LayerPlanetDeep
)

type WorldContext struct {
	CurrentLayer  GameLayer
	PlanetID      int
	Depth         int
	BiomeType     int
	QuestProgress QuestProgress
}

type QuestProgress struct {
	RoyalCharterComplete bool
	ContractsCollected   int
	ContractsNeeded      int
}

// ContextKey is stored in the World to allow systems to fetch context without global state.
// We keep it minimal and typed for simplicity.

type contextKey struct{}

var contextEntity Entity = 0

func ensureContextEntity(w *World) Entity {
	if contextEntity == 0 {
		contextEntity = w.Create()
	}
	return contextEntity
}

func SetWorldContext(w *World, ctx WorldContext) {
	e := ensureContextEntity(w)
	Add(w, e, ctx)
}

func GetWorldContext(w *World) WorldContext {
	e := ensureContextEntity(w)
	ctx, ok := Get[WorldContext](w, e)
	if !ok {
		return WorldContext{}
	}
	return ctx
}

// SystemRegistry lists systems that are conditionally run per layer.
type SystemRegistry struct {
	SpaceSystems     []System
	SurfaceSystems   []System
	DeepSystems      []System
	UniversalSystems []System
}

// SchedulerWithContext runs only systems active for the current context.
type SchedulerWithContext struct {
	Registry SystemRegistry
}

func NewSchedulerWithContext(reg SystemRegistry) *SchedulerWithContext {
	return &SchedulerWithContext{Registry: reg}
}

func (s *SchedulerWithContext) Update(dt float64, w *World) {
	ctx := GetWorldContext(w)
	for _, sys := range s.Registry.UniversalSystems {
		sys.Update(dt, w)
	}
	switch ctx.CurrentLayer {
	case LayerSpace:
		for _, sys := range s.Registry.SpaceSystems {
			sys.Update(dt, w)
		}
	case LayerPlanetSurface:
		for _, sys := range s.Registry.SurfaceSystems {
			sys.Update(dt, w)
		}
	case LayerPlanetDeep:
		for _, sys := range s.Registry.DeepSystems {
			sys.Update(dt, w)
		}
	}
}
