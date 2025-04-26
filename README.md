## Keyboard Warrior - Music Game

A rhythm game where notes are loaded from JSON song files with detailed end screens.

### How to Play

1. Use arrow keys to browse through available songs on the title screen
2. Press SPACE to start playing the selected song
3. Use Q, W, E, I, O, P keys to hit the notes when they reach the white line
4. The game ends after 20 misses or when the song is complete
5. View your performance statistics on the end screen
6. Press SPACE to return to the title screen

### Features

- Title screen with song selection
- Gameplay with visual feedback
- Detailed end screens showing:
    - Final score
    - Notes hit count
    - Accuracy percentage
    - Performance rank (S, A, B, C, D, F)
- Different screens for game over and song completion

### Game States

- Title Screen: Select songs and start a game
- Playing: Hit notes as they come down lanes
- Game Over: Displayed when player misses too many notes
- Song Complete: Shown when a song finishes successfully

### Song JSON Format

Songs are stored in JSON files in the `songs` directory. Each song file follows this format:
