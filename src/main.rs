use bevy::prelude::*;
use bevy_kira_audio::AudioPlugin;

// Import our custom modules
mod audio;
mod piano;

// Import our plugins
use audio::AudioGenerationPlugin;
use piano::PianoPlugin;

fn main() {
    App::new()
        .insert_resource(ClearColor(Color::rgb(0.1, 0.1, 0.1))) // Dark gray background
        .add_plugins((
            DefaultPlugins.set(WindowPlugin {
                primary_window: Some(Window {
                    title: "Piano App".into(),
                    resolution: (piano::WINDOW_WIDTH, piano::WINDOW_HEIGHT).into(),
                    ..default()
                }),
                ..default()
            }), 
            AudioPlugin,
            PianoPlugin,
            AudioGenerationPlugin,
        ))
        .run();
}