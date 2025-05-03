use bevy::prelude::*;
use bevy_kira_audio::AudioSource;
use crate::piano::components::{PianoNotes, SoundHandles};

// System to load audio files for piano notes
pub fn setup_piano_notes(
    mut commands: Commands,
    mut sound_handles: ResMut<SoundHandles>,
    asset_server: Res<AssetServer>,
) {
    // Define the piano notes with their frequencies
    let notes = [
        ("C4", 261.63),
        ("D4", 293.66),
        ("E4", 329.63),
        ("F4", 349.23),
        ("G4", 392.00),
        ("A4", 440.00),
        ("B4", 493.88),
        ("C5", 523.25),
    ];

    // Load sound files from src/audio/sounds instead of assets/sounds
    for (note_name, _) in notes.iter() {
        // Use a string directly instead of Path::new with a reference
        let file_path = format!("src/audio/sounds/{}.wav", note_name);
        let handle: Handle<AudioSource> = asset_server.load(&file_path);
        sound_handles.handles.push((note_name.to_string(), handle));
        println!("Loading audio: {} from path {}", note_name, file_path);
    }

    // Save notes to resource for use in sound generation
    commands.insert_resource(PianoNotes {
        notes: notes
            .iter()
            .map(|(name, freq)| (name.to_string(), *freq))
            .collect(),
    });
}

// System to manually load audio files after they've been generated
pub fn load_audio_files(mut sound_handles: ResMut<SoundHandles>, asset_server: Res<AssetServer>) {
    // For each note in our sound_handles
    for (note_name, handle) in sound_handles.handles.iter_mut() {
        // Use a string directly instead of Path::new with a reference
        let file_path = format!("src/audio/sounds/{}.wav", note_name);
        
        // Load the audio file using asset_server
        *handle = asset_server.load(&file_path);
        println!("Loaded audio: {} from path {}", note_name, file_path);
    }
}
