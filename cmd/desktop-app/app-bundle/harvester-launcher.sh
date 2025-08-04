#!/bin/bash

# Get the directory where this script is located (inside the app bundle)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
APP_DIR="$(dirname "$SCRIPT_DIR")"
EXECUTABLE="$APP_DIR/MacOS/harvester-terminal"

# Check if the executable exists
if [ ! -f "$EXECUTABLE" ]; then
    echo "Error: harvester-terminal not found at $EXECUTABLE"
    exit 1
fi

# Create a temporary script that will run in the new Terminal window
TEMP_SCRIPT=$(mktemp)
cat > "$TEMP_SCRIPT" << EOF
#!/bin/bash
cd "$APP_DIR/MacOS"
echo "🚀 Starting Harvester..."
echo "Controls: W=thrust, S=brake, A/D=turn, Q=quit"
echo ""
./harvester-terminal
echo ""
echo "Thanks for playing Harvester!"
echo "Press any key to close this window..."
read -n 1
EOF

chmod +x "$TEMP_SCRIPT"

# Open a new Terminal window and run our script
osascript -e "tell application \"Terminal\" to do script \"$TEMP_SCRIPT; rm '$TEMP_SCRIPT'\""

# Focus Terminal
osascript -e "tell application \"Terminal\" to activate"