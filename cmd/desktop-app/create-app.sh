#!/bin/bash

set -e

APP_NAME="Harvester"
BUNDLE_ID="com.harvester.desktop"
VERSION="1.0.0"

echo "🏗️  Building Harvester.app..."

# Clean previous builds
rm -rf "$APP_NAME.app"

# Create app bundle structure
mkdir -p "$APP_NAME.app/Contents/MacOS"
mkdir -p "$APP_NAME.app/Contents/Resources"

# Build the executable
echo "Building executable..."
make build-c

# Copy executable to app bundle
cp harvester-terminal "$APP_NAME.app/Contents/MacOS/"

# Copy Info.plist
cp app-bundle/Info.plist "$APP_NAME.app/Contents/"

# Create the launcher script in MacOS directory
cat > "$APP_NAME.app/Contents/MacOS/Harvester" << 'EOF'
#!/bin/bash

# Get the directory where this script is located (inside the app bundle)
BUNDLE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
EXECUTABLE="$BUNDLE_DIR/MacOS/harvester-terminal"

# Check if the executable exists
if [ ! -f "$EXECUTABLE" ]; then
    osascript -e 'display dialog "Error: Game executable not found!" buttons {"OK"} default button "OK"'
    exit 1
fi

# Create a temporary script that will run in the new Terminal window
TEMP_SCRIPT=$(mktemp)
cat > "$TEMP_SCRIPT" << INNER_EOF
#!/bin/bash
cd "$BUNDLE_DIR/MacOS"
echo "🚀 Starting Harvester..."
echo "Controls: W=thrust, S=brake, A/D=turn, Q=quit"
echo ""
./harvester-terminal
echo ""
echo "Thanks for playing Harvester!"
echo "Press any key to close this window..."
read -n 1
INNER_EOF

chmod +x "$TEMP_SCRIPT"

# Open a new Terminal window and run our script
osascript << APPLESCRIPT
tell application "Terminal"
    do script "$TEMP_SCRIPT; rm '$TEMP_SCRIPT'"
    activate
end tell
APPLESCRIPT
EOF

# Make the launcher executable
chmod +x "$APP_NAME.app/Contents/MacOS/Harvester"

# Create Resources directory (skip icon for now - will use default)
mkdir -p "$APP_NAME.app/Contents/Resources"

echo "✅ $APP_NAME.app created successfully!"
echo ""
echo "To install:"
echo "  cp -r '$APP_NAME.app' /Applications/"
echo ""
echo "To run:"
echo "  open '$APP_NAME.app'"
echo "  # or double-click in Finder"
echo ""
echo "To test locally:"
echo "  open './$APP_NAME.app'"