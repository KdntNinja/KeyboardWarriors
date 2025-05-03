use bevy::prelude::*;
use bevy_kira_audio::AudioSource;

// Constants for piano keys
pub const WINDOW_WIDTH: f32 = 800.0;
pub const WINDOW_HEIGHT: f32 = 600.0;
pub const WHITE_KEY_WIDTH: f32 = WINDOW_WIDTH / 8.0; // Adjust to fit 8 white keys perfectly across the screen (C4 to C5)
pub const WHITE_KEY_HEIGHT: f32 = 150.0;
pub const PIANO_Y_POSITION: f32 = -WINDOW_HEIGHT / 2.0 + WHITE_KEY_HEIGHT / 2.0; // Position piano at the bottom edge

// Constants for falling notes
pub const NOTE_SPEED: f32 = 200.0; // Pixels per second

// Define piano key types - we only use White now but keeping the enum for future extensibility
#[derive(Clone, Copy, PartialEq, Eq, Debug)]
pub enum PianoKeyType {
    White,
    Black, // Keep for future use
}

// Define piano key component
#[derive(Component)]
pub struct PianoKey {
    pub key_type: PianoKeyType,
    pub note_name: String,
    pub frequency: f32,
    pub keyboard_key: KeyCode,
    pub is_pressed: bool,
}

// Component for falling notes
#[derive(Component)]
pub struct FallingNote {
    pub note_name: String,
    pub keyboard_key: KeyCode,
}

// System resources
#[derive(Resource)]
pub struct PianoNotes {
    pub notes: Vec<(String, f32)>, // (note_name, frequency)
}

// For sound playback
#[derive(Resource)]
pub struct SoundHandles {
    pub handles: Vec<(String, Handle<AudioSource>)>, // (note_name, audio_handle)
}

impl Default for SoundHandles {
    fn default() -> Self {
        Self {
            handles: Vec::new(),
        }
    }
}
