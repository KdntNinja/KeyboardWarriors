# Keyboard Warriors

A rhythm-based typing game built with Rust and Bevy where you become a true "keyboard warrior"! Match your typing skills to the beat, launch attacks with your keyboard mastery, and battle through increasingly difficult challenges.

## 🎮 Game Concept

Keyboard Warriors is a rhythm-based typing combat game that combines several unique gameplay elements:

- **Rhythm Typing Combat**: Type words and phrases in time with the music. Perfect timing and accuracy deal maximum damage to your opponents.
- **Keyboard Sound Magic**: Different keys produce different sound effects and magical attacks, with special key combinations creating powerful combo moves.
- **QWERTY Dungeons**: Journey through levels designed around keyboard layouts, where each zone represents a different section of the keyboard with unique challenges.
- **Boss Battles**: Face off against legendary keyboard warriors who test both your typing speed and your strategic thinking.

## 🚀 Technical Details

- **Engine**: Built with [Bevy](https://bevyengine.org/), a data-driven game engine written in Rust
- **Audio**: Dynamic sound system with key-specific notes and rhythm detection
- **Performance**: Optimized for low-latency input handling crucial for a rhythm game

## 📋 Features

- [x] Basic Bevy setup with 2D camera
- [ ] Keyboard input detection and visualization
- [ ] Sound system with note playback
- [ ] Rhythm tracking and scoring system
- [ ] Enemy combat system
- [ ] Level progression
- [ ] Custom keyboard theming
- [ ] Song editor

## 🎵 Music System

The game includes a simple but powerful music system where songs are defined in JSON format. Players can even create their own custom song patterns to play through.

## 🔧 Development

### Prerequisites

- Rust (latest stable)
- Cargo

### Building and Running

```bash
# Clone the repository
git clone https://www.github.com/KdntNinja/KeyboardWarriors.git
cd KeyboardWarriors

# Run in debug mode
cargo run

# Build for release
cargo build --release
```

## 🤝 Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.
