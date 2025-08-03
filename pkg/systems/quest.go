package systems

import (
	"harvester/pkg/components"
	"harvester/pkg/ecs"
)

type QuestSystem struct{}

func (QuestSystem) Update(dt float64, w *ecs.World) {
	ctx := ecs.GetWorldContext(w)
	if ctx.CurrentLayer != ecs.LayerPlanetSurface {
		return
	}
	var playerInv components.Inventory
	found := false
	ecs.View2Of[components.Player, components.Inventory](w).Each(func(t ecs.Tuple2[components.Player, components.Inventory]) {
		playerInv = *t.B
		found = true
	})
	if !found {
		return
	}
	collected := playerInv.Items["trade_contract"]
	ctx.QuestProgress.ContractsCollected = collected
	if collected >= ctx.QuestProgress.ContractsNeeded {
		ctx.QuestProgress.RoyalCharterComplete = true
	}
	ecs.SetWorldContext(w, ctx)
}
