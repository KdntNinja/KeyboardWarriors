use bevy::prelude::*;
use bevy_kira_audio::{Audio, AudioControl};
use crate::piano::components::*;
use bevy::time::Timer;
use crate::piano::components::{FallingNote, NOTE_SPEED};
use rand::seq::IteratorRandom;

// Resource to track time for spawning notes
#[derive(Resource, Default)]
pub struct NoteSpawnTimer(pub Timer);

// System to handle keyboard input and update piano key state
pub fn handle_key_input(
    keyboard: Res<Input<KeyCode>>,
    mut piano_keys: Query<&mut PianoKey>,
) {
    for mut piano_key in piano_keys.iter_mut() {
        // Update key press state based on keyboard input
        piano_key.is_pressed = keyboard.pressed(piano_key.keyboard_key);
    }
}

// System to update piano key visuals based on press state
pub fn update_key_visuals(
    mut query: Query<(&PianoKey, &mut Sprite)>,
) {
    for (key, mut sprite) in query.iter_mut() {
        match key.key_type {
            PianoKeyType::White => {
                sprite.color = if key.is_pressed {
                    Color::hex("BBBBFF").unwrap() // Light blue when pressed
                } else {
                    Color::hex("FFFFFF").unwrap() // White when not pressed
                };
            }
            PianoKeyType::Black => {
                sprite.color = if key.is_pressed {
                    Color::hex("333355").unwrap() // Dark blue when pressed
                } else {
                    Color::hex("1A1A1A").unwrap() // Black when not pressed
                };
            }
        }
    }
}

// System to play sounds when keys are pressed
pub fn play_sounds(
    audio: Res<Audio>,
    keyboard: Res<Input<KeyCode>>,
    piano_keys: Query<&PianoKey>,
    sound_handles: Res<SoundHandles>,
) {
    // For each piano key that was just pressed
    for key in piano_keys.iter() {
        if keyboard.just_pressed(key.keyboard_key) {
            // Find the corresponding audio handle
            if let Some((_, handle)) = sound_handles.handles.iter()
                .find(|(name, _)| name == &key.note_name) {
                
                // Play the sound
                audio.play(handle.clone());
                
                println!("Playing note: {} at frequency {}", key.note_name, key.frequency);
            }
        }
    }
}

// System to spawn falling notes
pub fn spawn_falling_notes(
    mut commands: Commands,
    time: Res<Time>,
    mut timer: ResMut<NoteSpawnTimer>,
    piano_keys: Query<&PianoKey>,
) {
    // Update the timer
    if timer.0.tick(time.delta()).just_finished() {
        // Randomly select a piano key to spawn a note for
        if let Some(key) = piano_keys.iter().choose(&mut rand::thread_rng()) {
            commands.spawn((
                FallingNote {
                    note_name: key.note_name.clone(),
                    keyboard_key: key.keyboard_key,
                },
                SpriteBundle {
                    sprite: Sprite {
                        color: Color::hex("FF0000").unwrap(), // Red for falling notes
                        custom_size: Some(Vec2::new(WHITE_KEY_WIDTH - 2.0, 20.0)),
                        ..default()
                    },
                    transform: Transform::from_xyz(0.0, WINDOW_HEIGHT / 2.0, 1.0),
                    ..default()
                },
            ));
        }
    }
}

// System to move falling notes
pub fn move_falling_notes(
    mut commands: Commands,
    mut query: Query<(Entity, &mut Transform, &FallingNote)>,
    time: Res<Time>,
) {
    for (entity, mut transform, _note) in query.iter_mut() {
        // Move the note downward
        transform.translation.y -= NOTE_SPEED * time.delta_seconds();

        // Despawn the note if it goes off the bottom of the screen
        if transform.translation.y < -WINDOW_HEIGHT / 2.0 {
            commands.entity(entity).despawn();
        }
    }
}