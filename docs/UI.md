# UI.md

Model
- Holds GameState pointer, viewport size, focus/modal state, keymap.

Components
- StatusBar: resources, epoch, drive level
- MapView: renders viewport using pre-styled glyphs
- LogView: recent events
- Modal: Upgrades menu

Messages
- tea.KeyMsg for input, tea.WindowSizeMsg for resize, TickMsg for periodic updates if needed

Rendering
- Use alt screen; layout with lipgloss width/height; only paint visible viewport

Keymap
- hjkl/arrows: move
- g: harvest
- w: warp
- s: scan
- u: upgrades
- q: quit
