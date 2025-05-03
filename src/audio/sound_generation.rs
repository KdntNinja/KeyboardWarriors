use bevy::prelude::*;
use bevy_kira_audio::AudioSource;
// Removed unused AudioControl import
use crate::piano::components::{PianoNotes, SoundHandles};
use std::f32::consts::PI;
use std::fs::File;
use std::io::{BufWriter, Write};
use std::path::Path;

// Sound generation constants
pub const SAMPLE_RATE: u32 = 44100;
pub const AMPLITUDE: f32 = 0.2;

// System to generate wave data for piano notes
pub fn setup_piano_notes(
    mut commands: Commands,
    mut sound_handles: ResMut<SoundHandles>,
    _asset_server: Res<AssetServer>, // Prefixed with an underscore to suppress warning
) {
    // Define the piano notes with their frequencies
    let notes = [
        ("C4", 261.63),
        ("C#4", 277.18),
        ("D4", 293.66),
        ("D#4", 311.13),
        ("E4", 329.63),
        ("F4", 349.23),
        ("F#4", 369.99),
        ("G4", 392.00),
        ("G#4", 415.30),
        ("A4", 440.00),
        ("A#4", 466.16),
        ("B4", 493.88),
        ("C5", 523.25),
        ("C#5", 554.37),
    ];

    // Create the assets directory if it doesn't exist
    let assets_path = Path::new("assets/sounds");
    std::fs::create_dir_all(assets_path).unwrap_or_else(|e| {
        println!("Warning: Could not create assets directory: {}", e);
    });

    // Generate WAV files for each note
    for (note_name, frequency) in &notes {
        // Generate the audio data
        let duration = 1.5; // seconds
        let sample_rate = SAMPLE_RATE;
        let num_samples = (sample_rate as f32 * duration) as usize;
        let mut audio_data = Vec::with_capacity(num_samples * 2); // Stereo samples

        for i in 0..num_samples {
            let t = i as f32 / sample_rate as f32;
            let decay = (-3.0 * t).exp();

            // Create a rich piano-like tone with harmonics
            let sample = (2.0 * PI * frequency * t).sin() * 0.6 * decay
                + (2.0 * PI * frequency * 2.0 * t).sin() * 0.3 * decay
                + (2.0 * PI * frequency * 3.0 * t).sin() * 0.1 * decay;

            // Scale and convert to 16-bit PCM
            let pcm_sample = (sample * AMPLITUDE * 32767.0) as i16;

            // Add stereo samples
            audio_data.push(pcm_sample);
            audio_data.push(pcm_sample);
        }

        // Create a WAV file for this note
        let file_path = assets_path.join(format!("{}.wav", note_name.replace("#", "s")));
        if !file_path.exists() {
            let file = File::create(&file_path).expect("Could not create audio file");
            let mut writer = BufWriter::new(file);

            // Write WAV header
            writer.write_all(b"RIFF").unwrap();
            writer
                .write_all(&((36 + audio_data.len() * 2) as u32).to_le_bytes())
                .unwrap(); // File size - 8
            writer.write_all(b"WAVE").unwrap();

            // "fmt " chunk
            writer.write_all(b"fmt ").unwrap();
            writer.write_all(&16u32.to_le_bytes()).unwrap(); // Chunk size
            writer.write_all(&1u16.to_le_bytes()).unwrap(); // Audio format (PCM)
            writer.write_all(&2u16.to_le_bytes()).unwrap(); // Num channels (stereo)
            writer.write_all(&sample_rate.to_le_bytes()).unwrap(); // Sample rate
            writer.write_all(&(sample_rate * 4).to_le_bytes()).unwrap(); // Byte rate
            writer.write_all(&4u16.to_le_bytes()).unwrap(); // Block align
            writer.write_all(&16u16.to_le_bytes()).unwrap(); // Bits per sample

            // "data" chunk
            writer.write_all(b"data").unwrap();
            writer
                .write_all(&((audio_data.len() * 2) as u32).to_le_bytes())
                .unwrap(); // Chunk size

            // Write audio data
            for sample in audio_data {
                writer.write_all(&sample.to_le_bytes()).unwrap();
            }

            println!("Generated sound file: {:?}", file_path);
        }

        // Store the note name for use in the audio plugin's manual loading
        sound_handles
            .handles
            .push((note_name.to_string(), Handle::<AudioSource>::default()));
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
        // Use Bevy's asset_server to load the file path directly
        let file_name = format!("sounds/{}.wav", note_name.replace("#", "s"));

        // Load the audio file using asset_server
        *handle = asset_server.load(&file_name);
        println!("Loaded audio: {} from path {}", note_name, file_name);
    }
}
