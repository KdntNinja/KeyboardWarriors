# Keyboard Warrior

A rhythm-based music game built with Go and Ebiten. Hit keys in time with the falling notes to score points.

## How to Play

1. **Start the Game**: Run the game with `go run .` or the built executable.
2. **Select a Song**: Use the UP/DOWN arrow keys to choose a song from the list.
3. **Start Playing**: Press SPACE to begin.
4. **Hit Notes**: Press the corresponding keys (QWEIOP) when notes reach the white hit line.
   - Q = C (Red notes)
   - W = D (Orange notes)
   - E = E (Yellow notes)
   - I = F (Green notes)
   - O = G (Blue notes)
   - P = A (Purple notes)
5. **Score Points**: More accurate hits earn higher scores.
6. **Hold Notes**: For longer notes, keep the key pressed until the note passes.

## Adding Custom Songs

You can add your own songs to the game by creating JSON files in the `songs` directory.

### Song File Format

Create a new JSON file in the `songs` directory with the following format:
