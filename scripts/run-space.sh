#!/bin/bash

# Simple script to run the space screen directly
echo "Running space screen test..."
echo "Player will start centered on screen"
echo "Controls: WASD to move, 'q' or Ctrl+C to quit"
echo "Press any key to continue..."
read -n 1

cd "$(dirname "$0")/.."
go run ./cmd/space