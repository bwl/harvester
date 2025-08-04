import SwiftUI

@main
struct HarvesterApp: App {
    var body: some Scene {
        WindowGroup {
            GameView()
                .frame(minWidth: 1024, minHeight: 768)
                .background(Color.black)
        }
        .windowResizability(.contentSize)
    }
}