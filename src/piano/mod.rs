// Re-export piano modules
pub mod components;
pub mod setup;
pub mod systems;

// Re-export specific items for convenience
pub use components::*;
pub use setup::setup;
pub use systems::*;

// Piano plugin to easily add piano functionality
use bevy::prelude::*;
use crate::piano::systems::{spawn_falling_notes, move_falling_notes, NoteSpawnTimer};

pub struct PianoPlugin;

impl Plugin for PianoPlugin {
    fn build(&self, app: &mut App) {
        app
            .init_resource::<SoundHandles>()
            .init_resource::<NoteSpawnTimer>()
            .add_systems(Startup, setup)
            .add_systems(Update, (
                handle_key_input,
                update_key_visuals,
                play_sounds,
                spawn_falling_notes,
                move_falling_notes,
            ));
    }
}