.PHONY: all clean build-go build-c build-swift run run-c build-app install-app

GAME_LIB = libgame.a
SWIFT_EXECUTABLE = .build/debug/HarvesterDesktop
C_EXECUTABLE = harvester-terminal

all: build-go build-c

# Build Go shared library
build-go:
	@echo "Building Go shared library..."
	go build -buildmode=c-archive -o $(GAME_LIB) .
	@echo "Go library built: $(GAME_LIB)"

# Build C terminal application (works with command line tools only)
build-c: build-go
	@echo "Building C terminal application..."
	gcc -o $(C_EXECUTABLE) terminal/terminal.c $(GAME_LIB) -lm
	@echo "C application built: $(C_EXECUTABLE)"

# Build Swift application (requires full Xcode)
build-swift: build-go
	@echo "Building Swift application..."
	swift build
	@echo "Swift application built: $(SWIFT_EXECUTABLE)"

# Run the C terminal version
run: build-c
	@echo "Running Harvester Desktop (Terminal Mode)..."
	./$(C_EXECUTABLE)

# Run the C terminal version  
run-c: build-c
	@echo "Running Harvester Desktop (Terminal Mode)..."
	./$(C_EXECUTABLE)

# Run the Swift version (if available)
run-swift: build-swift
	@echo "Running Harvester Desktop (Swift)..."
	./$(SWIFT_EXECUTABLE)

# Build macOS app bundle
build-app: build-c
	@echo "Creating macOS app bundle..."
	./create-app.sh
	@echo "App bundle created: Harvester.app"

# Install app to Applications folder
install-app: build-app
	@echo "Installing Harvester.app to /Applications..."
	cp -r Harvester.app /Applications/
	@echo "Installation complete! You can now find Harvester in your Applications folder."

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(GAME_LIB) libgame.h $(C_EXECUTABLE)
	rm -rf .build Harvester.app
	@echo "Clean complete"

# Test Go compilation
test-go:
	@echo "Testing Go compilation..."
	go build -buildmode=c-archive -o test.a .
	@echo "Go compilation test successful"
	rm -f test.a test.h