use bevy::{
    prelude::*,
    window::{PresentMode, WindowResolution},
    DefaultPlugins,
};

use keyboard_warriors_lib::*;
use bevy_kira_audio::prelude::*;
use bevy_kira_audio::AudioSource;

fn main() {
    // Initialize the game
    let mut app = App::new();
    
    // Configure window with Wayland compatibility
    app.add_plugins(DefaultPlugins.set(WindowPlugin {
        primary_window: Some(Window {
            title: "Keyboard Warriors".to_string(),
            resolution: WindowResolution::new(800.0, 600.0),
            present_mode: PresentMode::AutoVsync,
            ..default()
        }),
        ..default()
    }));
    
    // Add audio support
    app.add_plugins(AudioPlugin);
    
    app.add_systems(Startup, setup_audio);
    
    // Add our game plugin
    app.add_plugins(KeyboardWarriorsPlugin);
    
    // Initialize the game
    app.run();
}

// Setup audio systems
fn setup_audio(_commands: Commands, asset_server: Res<AssetServer>, audio: Res<Audio>) {
    // Preload audio files
    let note_paths = [
        "audio/note_C.mp3",
        "audio/note_D.mp3",
        "audio/note_E.mp3", 
        "audio/note_F.mp3",
        "audio/note_G.mp3",
        "audio/note_A.mp3",
    ];
    
    for path in note_paths.iter() {
        // Add explicit type annotation for the audio asset
        let _handle: Handle<AudioSource> = asset_server.load(*path);
    }
    
    // Play a sound to indicate the game has loaded
    let startup_sound: Handle<AudioSource> = asset_server.load("audio/note_C.mp3");
    audio.play(startup_sound).with_volume(0.5);
}
