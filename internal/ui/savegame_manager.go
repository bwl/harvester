package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"harvester/pkg/ecs"
)

// SaveGameManager handles all save/load operations
type SaveGameManager struct {
	saveDir string
}

// SaveSlotInfo contains information about a save slot
type SaveSlotInfo struct {
	SlotNum  int
	Exists   bool
	ModTime  time.Time
	Size     int64
	GameInfo string // extracted from save data if possible
}

// NewSaveGameManager creates a new save game manager
func NewSaveGameManager() *SaveGameManager {
	return &SaveGameManager{
		saveDir: ".saves",
	}
}

// HasAutosave checks if an autosave file exists
func (sgm *SaveGameManager) HasAutosave() bool {
	autosavePath := filepath.Join(sgm.saveDir, "autosave.gz")
	_, err := os.Stat(autosavePath)
	return err == nil
}

// GetSaveSlots returns information about all save slots
func (sgm *SaveGameManager) GetSaveSlots() []SaveSlotInfo {
	slots := make([]SaveSlotInfo, 3)
	for i := 1; i <= 3; i++ {
		slotPath := filepath.Join(sgm.saveDir, fmt.Sprintf("slot%d.gz", i))
		info := SaveSlotInfo{
			SlotNum: i,
			Exists:  false,
		}

		if stat, err := os.Stat(slotPath); err == nil {
			info.Exists = true
			info.ModTime = stat.ModTime()
			info.Size = stat.Size()
			info.GameInfo = sgm.extractGameInfo(slotPath)
		}

		slots[i-1] = info
	}
	return slots
}

// LoadAutosave loads the autosave file into the given world
func (sgm *SaveGameManager) LoadAutosave(world *ecs.World) error {
	autosavePath := filepath.Join(sgm.saveDir, "autosave.gz")
	return sgm.loadFromFile(autosavePath, world)
}

// LoadSlot loads a specific save slot into the given world
func (sgm *SaveGameManager) LoadSlot(slotNum int, world *ecs.World) error {
	slotPath := filepath.Join(sgm.saveDir, fmt.Sprintf("slot%d.gz", slotNum))
	return sgm.loadFromFile(slotPath, world)
}

// SaveAutosave saves the world to the autosave file
func (sgm *SaveGameManager) SaveAutosave(world *ecs.World) error {
	autosavePath := filepath.Join(sgm.saveDir, "autosave.gz")
	return sgm.saveToFile(autosavePath, world)
}

// SaveSlot saves the world to a specific save slot
func (sgm *SaveGameManager) SaveSlot(slotNum int, world *ecs.World) error {
	slotPath := filepath.Join(sgm.saveDir, fmt.Sprintf("slot%d.gz", slotNum))
	return sgm.saveToFile(slotPath, world)
}

// loadFromFile loads a save file into the given world
func (sgm *SaveGameManager) loadFromFile(path string, world *ecs.World) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read save file %s: %w", path, err)
	}

	snapshot, err := ecs.DecodeSnapshot(b, ecs.SaveOptions{Compress: true})
	if err != nil {
		return fmt.Errorf("failed to decode save file %s: %w", path, err)
	}

	err = ecs.Load(world, snapshot, nil)
	if err != nil {
		return fmt.Errorf("failed to load save file %s into world: %w", path, err)
	}

	return nil
}

// saveToFile saves the world to a file
func (sgm *SaveGameManager) saveToFile(path string, world *ecs.World) error {
	// Ensure save directory exists
	if err := os.MkdirAll(sgm.saveDir, 0o755); err != nil {
		return fmt.Errorf("failed to create save directory: %w", err)
	}

	snapshot, err := ecs.Save(world, nil)
	if err != nil {
		return fmt.Errorf("failed to create world snapshot: %w", err)
	}

	b, err := ecs.EncodeSnapshot(snapshot, ecs.SaveOptions{Compress: true})
	if err != nil {
		return fmt.Errorf("failed to encode snapshot: %w", err)
	}

	err = os.WriteFile(path, b, 0o644)
	if err != nil {
		return fmt.Errorf("failed to write save file %s: %w", path, err)
	}

	return nil
}

// extractGameInfo extracts basic information from a save file
func (sgm *SaveGameManager) extractGameInfo(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		return "Unknown"
	}

	snapshot, err := ecs.DecodeSnapshot(b, ecs.SaveOptions{Compress: true})
	if err != nil {
		return "Unknown"
	}

	// Look for world context to get layer info
	// This is a simplified extraction - the actual snapshot structure may vary
	dataStr := fmt.Sprintf("%+v", snapshot)
	if strings.Contains(dataStr, "LayerSpace") {
		return "Space"
	} else if strings.Contains(dataStr, "LayerPlanetSurface") {
		return "Planet Surface"
	} else if strings.Contains(dataStr, "LayerPlanetDeep") {
		return "Deep Underground"
	}

	return "Unknown"
}
