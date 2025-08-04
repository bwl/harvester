import Foundation

// C structs matching the Go bridge
struct CGlyph {
    let x: Int32
    let y: Int32
    let glyph: Int32
    let foregroundR: Int32
    let foregroundG: Int32
    let foregroundB: Int32
    let backgroundR: Int32
    let backgroundG: Int32
    let backgroundB: Int32
    let style: Int32
    let alpha: Float
}

struct CGlyphMatrix {
    let glyphs: UnsafeMutablePointer<CGlyph>?
    let width: Int32
    let height: Int32
    let count: Int32
}

// Swift wrapper for glyph data
struct GlyphData {
    let x: Int
    let y: Int
    let character: Character
    let foregroundColor: (r: UInt8, g: UInt8, b: UInt8)
    let backgroundColor: (r: UInt8, g: UInt8, b: UInt8)
    let style: Int32
    let alpha: Float
}

// Bridge functions
@_cdecl("initGame")
func initGame(width: Int32, height: Int32)

@_cdecl("updateGame") 
func updateGame(dt: Float, thrust: Int32, brake: Int32, left: Int32, right: Int32)

@_cdecl("getGlyphMatrix")
func getGlyphMatrix() -> CGlyphMatrix

class GameBridge: ObservableObject {
    @Published var glyphMatrix: [[GlyphData]] = []
    @Published var gameWidth: Int = 80
    @Published var gameHeight: Int = 24
    
    private var gameInitialized = false
    
    func initialize(width: Int, height: Int) {
        if !gameInitialized {
            gameWidth = width
            gameHeight = height
            initGame(width: Int32(width), height: Int32(height))
            gameInitialized = true
        }
    }
    
    func update(deltaTime: Float, input: GameInput) {
        updateGame(
            dt: deltaTime,
            thrust: input.thrust ? 1 : 0,
            brake: input.brake ? 1 : 0,
            left: input.left ? 1 : 0,
            right: input.right ? 1 : 0
        )
        
        updateGlyphMatrix()
    }
    
    private func updateGlyphMatrix() {
        let matrix = getGlyphMatrix()
        
        guard matrix.glyphs != nil, matrix.count > 0 else {
            glyphMatrix = Array(repeating: Array(repeating: GlyphData(
                x: 0, y: 0, character: " ",
                foregroundColor: (255, 255, 255),
                backgroundColor: (0, 0, 0),
                style: 0, alpha: 1.0
            ), count: gameWidth), count: gameHeight)
            return
        }
        
        var newMatrix: [[GlyphData]] = Array(repeating: Array(repeating: GlyphData(
            x: 0, y: 0, character: " ",
            foregroundColor: (255, 255, 255),
            backgroundColor: (0, 0, 0),
            style: 0, alpha: 1.0
        ), count: Int(matrix.width)), count: Int(matrix.height))
        
        for i in 0..<Int(matrix.count) {
            let cGlyph = matrix.glyphs![i]
            let x = Int(cGlyph.x)
            let y = Int(cGlyph.y)
            
            if y >= 0 && y < Int(matrix.height) && x >= 0 && x < Int(matrix.width) {
                let char = Character(UnicodeScalar(Int(cGlyph.glyph)) ?? UnicodeScalar(32)!)
                
                newMatrix[y][x] = GlyphData(
                    x: x, y: y, character: char,
                    foregroundColor: (UInt8(cGlyph.foregroundR), UInt8(cGlyph.foregroundG), UInt8(cGlyph.foregroundB)),
                    backgroundColor: (UInt8(cGlyph.backgroundR), UInt8(cGlyph.backgroundG), UInt8(cGlyph.backgroundB)),
                    style: cGlyph.style,
                    alpha: cGlyph.alpha
                )
            }
        }
        
        glyphMatrix = newMatrix
    }
}

struct GameInput {
    var thrust: Bool = false
    var brake: Bool = false
    var left: Bool = false
    var right: Bool = false
}