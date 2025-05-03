// Re-export audio modules
pub mod sound_generation;

// Re-export specific items for convenience
pub use sound_generation::{setup_piano_notes, load_audio_files};

// Audio plugin to handle audio functionality
use bevy::prelude::*;

pub struct AudioGenerationPlugin;

impl Plugin for AudioGenerationPlugin {
    fn build(&self, app: &mut App) {
        app.add_systems(Startup, (setup_piano_notes, load_audio_files.after(setup_piano_notes)));
    }
}