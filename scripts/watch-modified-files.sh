#!/bin/bash

# Watch for modified files in the Harvest of Stars project
# Usage: ./scripts/watch-modified-files.sh [directory]

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
WATCH_DIR="${1:-$PROJECT_ROOT}"

echo "Watching for modified files in: $WATCH_DIR"
echo "Press Ctrl+C to stop"
echo ""

# Check if fswatch is available (macOS preferred)
if command -v fswatch >/dev/null 2>&1; then
    echo "Using fswatch (macOS native file watcher)"
    fswatch -0 -r -e "\.git/" -e "\.saves/" -e "node_modules/" -e "__pycache__/" \
        --event Updated --event AttributeModified "$WATCH_DIR" | \
    while IFS= read -r -d "" file; do
        # Skip directories and hidden files
        if [[ -f "$file" && ! "$(basename "$file")" =~ ^\. ]]; then
            rel_path="${file#$PROJECT_ROOT/}"
            timestamp=$(date "+%H:%M:%S")
            echo "[$timestamp] MODIFIED: $rel_path"
            
            # Show file extension and size
            ext="${file##*.}"
            file_size=$(stat -f%z "$file" 2>/dev/null || echo "unknown")
            echo "           Extension: .$ext, Size: ${file_size} bytes"
            echo ""
        fi
    done

# Fallback to inotifywait (Linux)
elif command -v inotifywait >/dev/null 2>&1; then
    echo "Using inotifywait (Linux)"
    inotifywait -m -r -e modify,attrib \
        --exclude '\.(git|saves)/' \
        --format '[%T] MODIFIED: %w%f' \
        --timefmt '%H:%M:%S' \
        "$WATCH_DIR"

# Fallback to find + polling (cross-platform)
else
    echo "Using find polling (fallback - install fswatch for better performance)"
    echo "Checking every 1 second..."
    
    # Create baseline with modification times
    TEMP_DIR=$(mktemp -d)
    BASELINE="$TEMP_DIR/baseline.txt"
    CURRENT="$TEMP_DIR/current.txt"
    
    find "$WATCH_DIR" -type f \
        ! -path "*/.git/*" \
        ! -path "*/.saves/*" \
        ! -path "*/node_modules/*" \
        ! -path "*/__pycache__/*" \
        ! -name ".*" \
        -exec stat -f "%m %N" {} \; 2>/dev/null > "$BASELINE" || \
    find "$WATCH_DIR" -type f \
        ! -path "*/.git/*" \
        ! -path "*/.saves/*" \
        ! -path "*/node_modules/*" \
        ! -path "*/__pycache__/*" \
        ! -name ".*" \
        -exec stat -c "%Y %n" {} \; > "$BASELINE"
    
    while true; do
        sleep 1
        
        find "$WATCH_DIR" -type f \
            ! -path "*/.git/*" \
            ! -path "*/.saves/*" \
            ! -path "*/node_modules/*" \
            ! -path "*/__pycache__/*" \
            ! -name ".*" \
            -exec stat -f "%m %N" {} \; 2>/dev/null > "$CURRENT" || \
        find "$WATCH_DIR" -type f \
            ! -path "*/.git/*" \
            ! -path "*/.saves/*" \
            ! -path "*/node_modules/*" \
            ! -path "*/__pycache__/*" \
            ! -name ".*" \
            -exec stat -c "%Y %n" {} \; > "$CURRENT"
        
        # Find modified files (different timestamps)
        changed_files=$(comm -13 <(sort "$BASELINE") <(sort "$CURRENT") | cut -d' ' -f2-)
        
        if [[ -n "$changed_files" ]]; then
            timestamp=$(date "+%H:%M:%S")
            echo "$changed_files" | while IFS= read -r file; do
                if [[ -n "$file" && -f "$file" ]]; then
                    rel_path="${file#$PROJECT_ROOT/}"
                    echo "[$timestamp] MODIFIED: $rel_path"
                    
                    # Show file info
                    ext="${file##*.}"
                    file_size=$(stat -c%s "$file" 2>/dev/null || stat -f%z "$file" 2>/dev/null || echo "unknown")
                    echo "           Extension: .$ext, Size: ${file_size} bytes"
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