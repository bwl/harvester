package ui

import tea "github.com/charmbracelet/bubbletea"

type InputKind int

const (
	InputNone InputKind = iota
	InputQuit
	InputMoveLeft
	InputMoveRight
	InputMoveUp
	InputMoveDown
	InputEnter
	InputSaveAuto
	InputSaveCompressed
	InputSaveSlot1
	InputSaveSlot2
	InputSaveSlot3
	InputMenuUp
	InputMenuDown
	InputMenuLeft
	InputMenuRight
	InputMenuSelect
	InputMenuBack
	InputDebugToggle
)

type InputAction struct{ Kind InputKind }

type InputHandler interface{ HandleInput(InputAction) tea.Cmd }
