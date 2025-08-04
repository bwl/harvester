import SwiftUI

struct GameView: View {
    @StateObject private var gameBridge = GameBridge()
    @State private var gameInput = GameInput()
    @State private var lastUpdateTime = Date()
    
    let timer = Timer.publish(every: 1.0/60.0, on: .main, in: .common).autoconnect()
    
    var body: some View {
        GeometryReader { geometry in
            GlyphMatrixView(glyphs: gameBridge.glyphMatrix)
                .focusable()
                .onKeyPress { keyPress in
                    handleKeyPress(keyPress)
                    return .handled
                }
                .onReceive(timer) { _ in
                    let now = Date()
                    let deltaTime = Float(now.timeIntervalSince(lastUpdateTime))
                    lastUpdateTime = now
                    
                    gameBridge.update(deltaTime: deltaTime, input: gameInput)
                }
                .onAppear {
                    let width = max(80, Int(geometry.size.width / 12))
                    let height = max(24, Int(geometry.size.height / 20))
                    gameBridge.initialize(width: width, height: height)
                }
        }
        .background(Color.black)
    }
    
    private func handleKeyPress(_ keyPress: KeyPress) -> KeyPress.Result {
        switch keyPress.key {
        case .upArrow, .character("w"), .character("W"):
            gameInput.thrust = keyPress.phase == .down
        case .downArrow, .character("s"), .character("S"):
            gameInput.brake = keyPress.phase == .down
        case .leftArrow, .character("a"), .character("A"):
            gameInput.left = keyPress.phase == .down
        case .rightArrow, .character("d"), .character("D"):
            gameInput.right = keyPress.phase == .down
        default:
            return .ignored
        }
        return .handled
    }
}

struct GlyphMatrixView: View {
    let glyphs: [[GlyphData]]
    
    var body: some View {
        GeometryReader { geometry in
            if !glyphs.isEmpty {
                let cellWidth = geometry.size.width / CGFloat(glyphs[0].count)
                let cellHeight = geometry.size.height / CGFloat(glyphs.count)
                
                Canvas { context, size in
                    for (y, row) in glyphs.enumerated() {
                        for (x, glyph) in row.enumerated() {
                            let rect = CGRect(
                                x: CGFloat(x) * cellWidth,
                                y: CGFloat(y) * cellHeight,
                                width: cellWidth,
                                height: cellHeight
                            )
                            
                            // Draw background
                            context.fill(
                                Path(rect),
                                with: .color(Color(
                                    red: Double(glyph.backgroundColor.r) / 255.0,
                                    green: Double(glyph.backgroundColor.g) / 255.0,
                                    blue: Double(glyph.backgroundColor.b) / 255.0
                                ))
                            )
                            
                            // Draw character
                            if glyph.character != " " {
                                let text = Text(String(glyph.character))
                                    .font(.system(size: min(cellWidth, cellHeight) * 0.8, design: .monospaced))
                                    .foregroundColor(Color(
                                        red: Double(glyph.foregroundColor.r) / 255.0,
                                        green: Double(glyph.foregroundColor.g) / 255.0,
                                        blue: Double(glyph.foregroundColor.b) / 255.0
                                    ))
                                
                                context.draw(text, at: CGPoint(
                                    x: rect.midX,
                                    y: rect.midY
                                ))
                            }
                        }
                    }
                }
            }
        }
    }
}