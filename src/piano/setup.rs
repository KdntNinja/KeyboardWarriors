use bevy::prelude::*;
use crate::piano::components::*;

// System to set up the piano keys and camera
pub fn setup(mut commands: Commands, _asset_server: Res<AssetServer>) {
    // 2D Camera
    commands.spawn(Camera2dBundle::default());

    // Define the piano keys (C4 octave with A-G notes)
    // White keys: C, D, E, F, G, A, B (bound to A, S, D, F, G, H, J)
    // Black keys: C#, D#, F#, G#, A# (bound to W, E, R, T, Y)
    let white_keys = [
        ("C4", 261.63, KeyCode::A),
        ("D4", 293.66, KeyCode::S),
        ("E4", 329.63, KeyCode::D),
        ("F4", 349.23, KeyCode::F),
        ("G4", 392.00, KeyCode::G),
        ("A4", 440.00, KeyCode::H),
        ("B4", 493.88, KeyCode::J),
        ("C5", 523.25, KeyCode::K),
    ];

    let black_keys = [
        ("C#4", 277.18, KeyCode::W),
        ("D#4", 311.13, KeyCode::E),
        ("F#4", 369.99, KeyCode::R),
        ("G#4", 415.30, KeyCode::T),
        ("A#4", 466.16, KeyCode::Y),
        ("C#5", 554.37, KeyCode::U),
    ];

    // Spawn white keys
    for (i, (note_name, frequency, keyboard_key)) in white_keys.iter().enumerate() {
        let x_position = (i as f32 * WHITE_KEY_WIDTH) - (white_keys.len() as f32 * WHITE_KEY_WIDTH / 2.0) + WHITE_KEY_WIDTH / 2.0;
        
        commands.spawn((
            SpriteBundle {
                sprite: Sprite {
                    color: Color::hex("FFFFFF").unwrap(),
                    custom_size: Some(Vec2::new(WHITE_KEY_WIDTH - 2.0, WHITE_KEY_HEIGHT)),
                    ..default()
                },
                transform: Transform::from_xyz(x_position, PIANO_Y_POSITION, 0.0),
                ..default()
            },
            PianoKey {
                key_type: PianoKeyType::White,
                note_name: note_name.to_string(),
                frequency: *frequency,
                keyboard_key: *keyboard_key,
                is_pressed: false,
            },
        ));
        
        // Add note label
        commands.spawn(Text2dBundle {
            text: Text {
                sections: vec![TextSection {
                    value: format!("{}\n{:?}", note_name, keyboard_key).replace("Key", ""),
                    style: TextStyle {
                        font: default(),
                        font_size: 12.0,
                        color: Color::BLACK,
                    },
                }],
                alignment: TextAlignment::Center,
                ..default()
            },
            transform: Transform::from_xyz(x_position, PIANO_Y_POSITION - WHITE_KEY_HEIGHT/2.0 + 16.0, 1.0),
            ..default()
        });
    }

    // Spawn black keys
    let black_key_positions = [0, 1, 3, 4, 5]; // Positions relative to white keys (no black key between E-F and B-C)
    
    for (i, (note_name, frequency, keyboard_key)) in black_keys.iter().enumerate() {
        if i >= black_key_positions.len() {
            continue; // Skip if we don't have a position for this black key
        }
        
        let white_key_index = black_key_positions[i];
        let x_position = ((white_key_index as f32 * WHITE_KEY_WIDTH) + (WHITE_KEY_WIDTH / 2.0)) 
                        - (white_keys.len() as f32 * WHITE_KEY_WIDTH / 2.0) + WHITE_KEY_WIDTH / 2.0;
        
        commands.spawn((
            SpriteBundle {
                sprite: Sprite {
                    color: Color::hex("1A1A1A").unwrap(),
                    custom_size: Some(Vec2::new(BLACK_KEY_WIDTH, BLACK_KEY_HEIGHT)),
                    ..default()
                },
                transform: Transform::from_xyz(x_position, PIANO_Y_POSITION + (WHITE_KEY_HEIGHT - BLACK_KEY_HEIGHT) / 2.0, 1.0),
                ..default()
            },
            PianoKey {
                key_type: PianoKeyType::Black,
                note_name: note_name.to_string(),
                frequency: *frequency,
                keyboard_key: *keyboard_key,
                is_pressed: false,
            },
        ));
        
        // Add note label
        commands.spawn(Text2dBundle {
            text: Text {
                sections: vec![TextSection {
                    value: format!("{}\n{:?}", note_name, keyboard_key).replace("Key", ""),
                    style: TextStyle {
                        font: default(),
                        font_size: 10.0,
                        color: Color::WHITE,
                    },
                }],
                alignment: TextAlignment::Center,
                ..default()
            },
            transform: Transform::from_xyz(x_position, PIANO_Y_POSITION + (WHITE_KEY_HEIGHT - BLACK_KEY_HEIGHT) / 2.0 - BLACK_KEY_HEIGHT/2.0 + 16.0, 2.0),
            ..default()
        });
    }
}