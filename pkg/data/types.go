package data

type TileKind int

const (
	Void TileKind = iota
	Space
	Galaxy
	Nebula
	BlackHole
	Wormhole
	Anomaly
)

type UpgradeKind int

const (
	Drive UpgradeKind = iota
	Sensors
	Cargo
	Shield
	Warp
)
