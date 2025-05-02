#!/bin/bash
# Run script for KeyboardWarriors native build
# Ensures proper Wayland support on Fedora Linux

clear

# Set required environment variables for Wayland
export WINIT_UNIX_BACKEND=wayland
export SDL_VIDEODRIVER=wayland

# Check if -r flag is passed to run in release mode
if [ "$1" == "-r" ]; then
    echo "Running in release mode..."
    cargo run --release
else
    echo "Running in debug mode..."
    cargo run
fi