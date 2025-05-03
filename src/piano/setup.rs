use crate::piano::components::*;
use bevy::prelude::*;

// System to set up the piano keys and camera
pub fn setup(mut commands: Commands, _asset_server: Res<AssetServer>) {
    // 2D Camera
    commands.spawn(Camera2dBundle::default());

    // Define the piano keys (C4 octave with A-G notes)
    // White keys: C, D, E, F, G, A, B, C5 (bound to A, S, D, F, G, H, J, K, L)
    let white_keys = [
        ("C4", 261.63, KeyCode::A),
        ("D4", 293.66, KeyCode::S),
        ("E4", 329.63, KeyCode::D),
        ("F4", 349.23, KeyCode::F),
        ("G4", 392.00, KeyCode::H),
        ("A4", 440.00, KeyCode::J),
        ("B4", 493.88, KeyCode::K),
        ("C5", 523.25, KeyCode::L),
    ];

    // Calculate starting x position to center the piano
    // We're using exact window width now to fill the screen perfectly
    let start_x = -WINDOW_WIDTH / 2.0 + WHITE_KEY_WIDTH / 2.0;

    // Spawn white keys
    for (i, (note_name, frequency, keyboard_key)) in white_keys.iter().enumerate() {
        let x_position = start_x + (i as f32 * WHITE_KEY_WIDTH);

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
            transform: Transform::from_xyz(
                x_position,
                PIANO_Y_POSITION - WHITE_KEY_HEIGHT / 2.0 + 16.0,
                1.0,
            ),
            ..default()
        });
    }
}
