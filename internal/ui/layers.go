package ui

import "harvester/pkg/ecs"

func layerName(l ecs.GameLayer) string {
	switch l {
	case ecs.LayerSpace:
		return "Space"
	case ecs.LayerPlanetSurface:
		return "Surface"
	case ecs.LayerPlanetDeep:
		return "Deep"
	default:
		return "Unknown"
	}
}

func royalCharterStatus(q ecs.QuestProgress) string {
	if q.RoyalCharterComplete {
		return "complete"
	}
	need := q.ContractsNeeded
	if need == 0 {
		need = 5
	}
	left := need - q.ContractsCollected
	if left < 0 {
		left = 0
	}
	return "contracts:" + itoa(q.ContractsCollected) + "/" + itoa(need)
}
