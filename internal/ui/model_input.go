package ui

import "harvester/pkg/systems"

func (m *Model) ApplyAction(a InputAction) {
	switch a.Kind {
	case InputMoveLeft:
		systems.SetPlayerInput(m.world, m.player, "left")
	case InputMoveRight:
		systems.SetPlayerInput(m.world, m.player, "right")
	case InputMoveUp:
		systems.SetPlayerInput(m.world, m.player, "up")
	case InputMoveDown:
		systems.SetPlayerInput(m.world, m.player, "down")
	case InputEnter:
		systems.SetPlayerInput(m.world, m.player, "enter")
	}
}
