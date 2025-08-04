#!/bin/bash

# Watch for Go files being created or modified in the Harvest of Stars project
# Usage: ./scripts/watch-go-files.sh [directory]

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
WATCH_DIR="${1:-$PROJECT_ROOT}"

echo "Watching for Go files (.go) in: $WATCH_DIR"
echo "Press Ctrl+C to stop"
echo ""

# Check if fswatch is available (macOS preferred)
if command -v fswatch >/dev/null 2>&1; then
    echo "Using fswatch (macOS native file watcher)"
    fswatch -0 -r -e "\.git/" -e "\.saves/" --include "\.go$" \
        --event Created --event Updated --event MovedTo "$WATCH_DIR" | \
    while IFS= read -r -d "" file; do
        if [[ -f "$file" && "$file" =~ \.go$ ]]; then
            rel_path="${file#$PROJECT_ROOT/}"
            timestamp=$(date "+%H:%M:%S")
            
            # Determine if new or modified
            if [[ -f "$file" ]]; then
                # Check if file was recently created (within last 5 seconds)
                if [[ $(find "$file" -newermt "5 seconds ago" 2>/dev/null) ]]; then
                    echo "[$timestamp] NEW GO FILE: $rel_path"
                else
                    echo "[$timestamp] MODIFIED GO FILE: $rel_path"
                fi
                
                # Show package and rough line count
                package_name=$(head -n 10 "$file" | grep -E "^package " | head -n 1 | awk '{print $2}' || echo "unknown")
                line_count=$(wc -l < "$file" 2>/dev/null || echo "unknown")
                file_size=$(stat -f%z "$file" 2>/dev/null || echo "unknown")
                
                echo "           Package: $package_name, Lines: $line_count, Size: ${file_size} bytes"
                
                # Show any new function definitions
                new_funcs=$(grep -E "^func " "$file" | head -n 3 | sed 's/^/           /')
                if [[ -n "$new_funcs" ]]; then
                    echo "           Functions:"
                    echo "$new_funcs"
                fi
                echo ""
            fi
        fi
    done

# Fallback to inotifywait (Linux)
elif command -v inotifywait >/dev/null 2>&1; then
    echo "Using inotifywait (Linux)"
    inotifywait -m -r -e create,modify,moved_to \
        --include '\.go$' \
        --exclude '\.(git|saves)/' \
        --format '[%T] %e: %w%f' \
        --timefmt '%H:%M:%S' \
        "$WATCH_DIR" | \
    while read timestamp event file; do
        if [[ "$file" =~ \.go$ ]]; then
            rel_path="${file#$PROJECT_ROOT/}"
            
            case "$event" in
                CREATE|MOVED_TO)
                    echo "[$timestamp] NEW GO FILE: $rel_path"
                    ;;
                MODIFY)
                    echo "[$timestamp] MODIFIED GO FILE: $rel_path"
                    ;;
            esac
            
            if [[ -f "$file" ]]; then
                package_name=$(head -n 10 "$file" | grep -E "^package " | head -n 1 | awk '{print $2}' || echo "unknown")
                line_count=$(wc -l < "$file" 2>/dev/null || echo "unknown")
                echo "           Package: $package_name, Lines: $line_count"
                echo ""
            fi
        fi
    done

# Fallback to find + polling (cross-platform)
else
    echo "Using find polling (fallback - install fswatch for better performance)"
    echo "Checking every 2 seconds..."
    
    TEMP_DIR=$(mktemp -d)
    BASELINE="$TEMP_DIR/baseline.txt"
    CURRENT="$TEMP_DIR/current.txt"
    
    # Create baseline of Go files with modification times
    find "$WATCH_DIR" -name "*.go" \
        ! -path "*/.git/*" \
        ! -path "*/.saves/*" \
        -exec stat -f "%m %N" {} \; 2>/dev/null > "$BASELINE" || \
    find "$WATCH_DIR" -name "*.go" \
        ! -path "*/.git/*" \
        ! -path "*/.saves/*" \
        -exec stat -c "%Y %n" {} \; > "$BASELINE"
    
    while true; do
        sleep 2
        
        # Get current state
        find "$WATCH_DIR" -name "*.go" \
            ! -path "*/.git/*" \
            ! -path "*/.saves/*" \
            -exec stat -f "%m %N" {} \; 2>/dev/null > "$CURRENT" || \
        find "$WATCH_DIR" -name "*.go" \
            ! -path "*/.git/*" \
            ! -path "*/.saves/*" \
            -exec stat -c "%Y %n" {} \; > "$CURRENT"
        
        # Find new files
        new_files=$(comm -13 <(awk '{print $2}' "$BASELINE" | sort) <(awk '{print $2}' "$CURRENT" | sort))
        
        # Find modified files
        modified_files=$(comm -13 <(sort "$BASELINE") <(sort "$CURRENT") | awk '{print $2}')
        
        timestamp=$(date "+%H:%M:%S")
        
        # Report new files
        if [[ -n "$new_files" ]]; then
            echo "$new_files" | while IFS= read -r file; do
                if [[ -n "$file" && -f "$file" ]]; then
                    rel_path="${file#$PROJECT_ROOT/}"
                    echo "[$timestamp] NEW GO FILE: $rel_path"
                    
                    package_name=$(head -n 10 "$file" | grep -E "^package " | head -n 1 | awk '{print $2}' || echo "unknown")
                    line_count=$(wc -l < "$file" 2>/dev/null || echo "unknown")
                    echo "           Package: $package_name, Lines: $line_count"
                    echo ""
                fi
            done
        fi
        
        # Report modified files
        if [[ -n "$modified_files" ]]; then
            echo "$modified_files" | while IFS= read -r file; do
                if [[ -n "$file" && -f "$file" ]]; then
                    rel_path="${file#$PROJECT_ROOT/}"
                    echo "[$timestamp] MODIFIED GO FILE: $rel_path"
                    
                    package_name=$(head -n 10 "$file" | grep -E "^package " | head -n 1 | awk '{print $2}' || echo "unknown")
                    line_count=$(wc -l < "$file" 2>/dev/null || echo "unknown")
                    echo "           Package: $package_name, Lines: $line_count"
                    echo ""
                fi
            done
        fi
        
        # Update baseline
        cp "$CURRENT" "$BASELINE"
    done
    
    # Cleanup
    rm -rf "$TEMP_DIR"
fi