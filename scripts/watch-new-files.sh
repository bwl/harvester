#!/bin/bash

# Watch for new files in the Harvest of Stars project
# Usage: ./scripts/watch-new-files.sh [directory]

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
WATCH_DIR="${1:-$PROJECT_ROOT}"

echo "Watching for new files in: $WATCH_DIR"
echo "Press Ctrl+C to stop"
echo ""

# Check if fswatch is available (macOS preferred)
if command -v fswatch >/dev/null 2>&1; then
    echo "Using fswatch (macOS native file watcher)"
    fswatch -0 -r -e "\.git/" -e "\.saves/" -e "node_modules/" -e "__pycache__/" \
        --event Created --event MovedTo "$WATCH_DIR" | \
    while IFS= read -r -d "" file; do
        # Skip directories
        if [[ -f "$file" ]]; then
            rel_path="${file#$PROJECT_ROOT/}"
            timestamp=$(date "+%H:%M:%S")
            echo "[$timestamp] NEW FILE: $rel_path"
            
            # Show file type and size
            file_type=$(file -b "$file" | cut -d',' -f1)
            file_size=$(stat -f%z "$file" 2>/dev/null || echo "unknown")
            echo "           Type: $file_type, Size: ${file_size} bytes"
            echo ""
        fi
    done

# Fallback to inotifywait (Linux)
elif command -v inotifywait >/dev/null 2>&1; then
    echo "Using inotifywait (Linux)"
    inotifywait -m -r -e create,moved_to \
        --exclude '\.(git|saves)/' \
        --format '[%T] NEW FILE: %w%f' \
        --timefmt '%H:%M:%S' \
        "$WATCH_DIR"

# Fallback to find + polling (cross-platform)
else
    echo "Using find polling (fallback - install fswatch for better performance)"
    echo "Checking every 2 seconds..."
    
    # Create baseline file list
    TEMP_DIR=$(mktemp -d)
    BASELINE="$TEMP_DIR/baseline.txt"
    CURRENT="$TEMP_DIR/current.txt"
    
    find "$WATCH_DIR" -type f \
        ! -path "*/.git/*" \
        ! -path "*/.saves/*" \
        ! -path "*/node_modules/*" \
        ! -path "*/__pycache__/*" \
        > "$BASELINE"
    
    while true; do
        sleep 2
        
        find "$WATCH_DIR" -type f \
            ! -path "*/.git/*" \
            ! -path "*/.saves/*" \
            ! -path "*/node_modules/*" \
            ! -path "*/__pycache__/*" \
            > "$CURRENT"
        
        # Find new files
        new_files=$(comm -13 "$BASELINE" "$CURRENT")
        
        if [[ -n "$new_files" ]]; then
            timestamp=$(date "+%H:%M:%S")
            echo "$new_files" | while IFS= read -r file; do
                if [[ -n "$file" ]]; then
                    rel_path="${file#$PROJECT_ROOT/}"
                    echo "[$timestamp] NEW FILE: $rel_path"
                    
                    # Show file info
                    if [[ -f "$file" ]]; then
                        file_type=$(file -b "$file" 2>/dev/null | cut -d',' -f1 || echo "unknown")
                        file_size=$(stat -c%s "$file" 2>/dev/null || stat -f%z "$file" 2>/dev/null || echo "unknown")
                        echo "           Type: $file_type, Size: ${file_size} bytes"
                    fi
                    echo ""
                fi
            done
            
            # Update baseline
            cp "$CURRENT" "$BASELINE"
        fi
    done
    
    # Cleanup
    rm -rf "$TEMP_DIR"
fi