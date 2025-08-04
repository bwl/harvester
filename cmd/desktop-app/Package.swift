// swift-tools-version:5.7
import PackageDescription

let package = Package(
    name: "HarvesterDesktop",
    platforms: [.macOS(.v12)],
    products: [
        .executable(name: "HarvesterDesktop", targets: ["HarvesterDesktop"])
    ],
    targets: [
        .executableTarget(
            name: "HarvesterDesktop",
            dependencies: [],
            linkerSettings: [
                .linkedLibrary("game"),
                .unsafeFlags(["-L.", "-Wl,-rpath,."])
            ]
        )
    ]
)