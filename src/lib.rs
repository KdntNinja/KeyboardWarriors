use bevy::{
    prelude::*,
    utils::Duration,
};
use serde::{Deserialize, Serialize};

// Game state management
#[derive(Debug, Clone, Eq, PartialEq, Hash, States, Default)]
pub enum GameState {
    #[default]
    MainMenu,
    Playing,
    GameOver,
}

// Note representation for the game
#[derive(Debug, Clone, Component)]
pub struct Note {
    pub key: NoteKey,
    pub lane: usize,
    pub spawn_time: f32,
    pub hit_time: f32,
}

// Piano keys that can be played
#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash, Deserialize, Serialize)]
pub enum NoteKey {
    C,
    D,
    E,
    F,
    G,
    A,
}

impl NoteKey {
    pub fn get_audio_path(&self) -> &'static str {
        match self {
            NoteKey::C => "audio/note_C.mp3",
            NoteKey::D => "audio/note_D.mp3",
            NoteKey::E => "audio/note_E.mp3",
            NoteKey::F => "audio/note_F.mp3",
            NoteKey::G => "audio/note_G.mp3",
            NoteKey::A => "audio/note_A.mp3",
        }
    }
    
    pub fn from_key_code(key_code: KeyCode) -> Option<Self> {
        match key_code {
            KeyCode::A => Some(NoteKey::C),
            KeyCode::S => Some(NoteKey::D),
            KeyCode::D => Some(NoteKey::E),
            KeyCode::F => Some(NoteKey::F),
            KeyCode::G => Some(NoteKey::G),
            KeyCode::H => Some(NoteKey::A),
            _ => None,
        }
    }
}

// Song data structure
#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct SongData {
    pub title: String,
    pub bpm: f32,
    pub notes: Vec<SongNote>,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct SongNote {
    pub key: NoteKey,
    pub time: f32, // Time in seconds when the note should be hit
}

// Game resources
#[derive(Resource)]
pub struct GameResources {
    pub score: u32,
    pub combo: u32,
    pub max_combo: u32,
    pub song: Option<SongData>,
    pub song_timer: Timer,
    pub note_speed: f32,
}

impl Default for GameResources {
    fn default() -> Self {
        Self {
            score: 0,
            combo: 0,
            max_combo: 0,
            song: None,
            song_timer: Timer::new(Duration::from_secs(0), TimerMode::Once),
            note_speed: 300.0, // Pixels per second
        }
    }
}

// UI Components
#[derive(Component)]
pub struct ScoreText;

#[derive(Component)]
pub struct ComboText;

#[derive(Component)]
pub struct MainMenu;

#[derive(Component)]
pub struct PlayButton;

#[derive(Component)]
pub struct NoteSprite;

// Game plugin that contains core functionality
pub struct KeyboardWarriorsPlugin;

impl Plugin for KeyboardWarriorsPlugin {
    fn build(&self, app: &mut App) {
        app
            .add_state::<GameState>() // Changed from init_state to add_state
            .init_resource::<GameResources>()
            .add_systems(Startup, setup)
            .add_systems(OnEnter(GameState::MainMenu), setup_main_menu)
            .add_systems(OnExit(GameState::MainMenu), cleanup_main_menu)
            .add_systems(OnEnter(GameState::Playing), setup_game)
            .add_systems(OnExit(GameState::Playing), cleanup_game)
            .add_systems(Update, (
                menu_interaction.run_if(in_state(GameState::MainMenu)),
                spawn_notes.run_if(in_state(GameState::Playing)),
                move_notes.run_if(in_state(GameState::Playing)),
                handle_key_presses.run_if(in_state(GameState::Playing)),
                update_score_ui.run_if(in_state(GameState::Playing)),
                check_game_over.run_if(in_state(GameState::Playing)),
            ));
    }
}

// Setup the game
fn setup(mut commands: Commands) {
    commands.spawn(Camera2dBundle::default());
}

// Main menu setup
fn setup_main_menu(mut commands: Commands) {
    // Menu root
    commands
        .spawn((
            NodeBundle {
                style: Style {
                    width: Val::Percent(100.0),
                    height: Val::Percent(100.0),
                    flex_direction: FlexDirection::Column,
                    align_items: AlignItems::Center,
                    justify_content: JustifyContent::Center,
                    ..default()
                },
                ..default()
            },
            MainMenu,
        ))
        .with_children(|parent| {
            // Title
            parent.spawn(
                TextBundle::from_section(
                    "Keyboard Warriors",
                    TextStyle {
                        font_size: 80.0,
                        color: Color::rgb(0.9, 0.9, 0.9),
                        ..default()
                    },
                )
                .with_text_alignment(TextAlignment::Center)
            );

            // Play button
            parent
                .spawn((
                    ButtonBundle {
                        style: Style {
                            width: Val::Px(200.0),
                            height: Val::Px(65.0),
                            margin: UiRect::all(Val::Px(20.0)),
                            justify_content: JustifyContent::Center,
                            align_items: AlignItems::Center,
                            ..default()
                        },
                        background_color: Color::rgb(0.15, 0.15, 0.15).into(),
                        ..default()
                    },
                    PlayButton,
                ))
                .with_children(|parent| {
                    parent.spawn(
                        TextBundle::from_section(
                            "Play",
                            TextStyle {
                                font_size: 40.0,
                                color: Color::rgb(0.9, 0.9, 0.9),
                                ..default()
                            },
                        )
                        .with_text_alignment(TextAlignment::Center)
                    );
                });
            
            // Controls instructions
            parent.spawn(
                TextBundle::from_section(
                    "Controls: A S D F G H keys to play notes",
                    TextStyle {
                        font_size: 24.0,
                        color: Color::rgb(0.7, 0.7, 0.7),
                        ..default()
                    },
                )
                .with_text_alignment(TextAlignment::Center)
            );
        });
}

// Cleanup menu when transitioning to game
fn cleanup_main_menu(
    mut commands: Commands, 
    menu_query: Query<Entity, With<MainMenu>>,
) {
    for entity in menu_query.iter() {
        commands.entity(entity).despawn_recursive();
    }
}

// Load song data from file
fn load_song_data() -> SongData {
    // In a real implementation, this would load from a file
    // For now, we'll return a simple example song
    SongData {
        title: "Example Song".to_string(),
        bpm: 120.0,
        notes: vec![
            SongNote { key: NoteKey::C, time: 1.0 },
            SongNote { key: NoteKey::D, time: 1.5 },
            SongNote { key: NoteKey::E, time: 2.0 },
            SongNote { key: NoteKey::F, time: 2.5 },
            SongNote { key: NoteKey::G, time: 3.0 },
            SongNote { key: NoteKey::A, time: 3.5 },
            SongNote { key: NoteKey::C, time: 4.0 },
            SongNote { key: NoteKey::E, time: 4.5 },
            SongNote { key: NoteKey::G, time: 5.0 },
            // Add more notes here
        ],
    }
}

// Game setup
fn setup_game(
    mut commands: Commands, 
    mut game_resources: ResMut<GameResources>,
) {
    // Load song
    game_resources.song = Some(load_song_data());
    game_resources.song_timer = Timer::new(Duration::from_secs(0), TimerMode::Once);
    game_resources.score = 0;
    game_resources.combo = 0;
    game_resources.max_combo = 0;
    
    // UI Setup
    commands
        .spawn(NodeBundle {
            style: Style {
                width: Val::Percent(100.0),
                height: Val::Percent(100.0),
                flex_direction: FlexDirection::Column,
                ..default()
            },
            ..default()
        })
        .with_children(|parent| {
            // Score text
            parent.spawn((
                TextBundle::from_section(
                    "Score: 0",
                    TextStyle {
                        font_size: 30.0,
                        color: Color::WHITE,
                        ..default()
                    },
                )
                .with_style(Style {
                    position_type: PositionType::Absolute,
                    top: Val::Px(10.0),
                    left: Val::Px(10.0),
                    ..default()
                }),
                ScoreText,
            ));
            
            // Combo text
            parent.spawn((
                TextBundle::from_section(
                    "Combo: 0",
                    TextStyle {
                        font_size: 30.0,
                        color: Color::WHITE,
                        ..default()
                    },
                )
                .with_style(Style {
                    position_type: PositionType::Absolute,
                    top: Val::Px(50.0),
                    left: Val::Px(10.0),
                    ..default()
                }),
                ComboText,
            ));
        });
    
    // Setup piano key zones at bottom of screen
    let lane_width = 100.0;
    for (i, key) in [NoteKey::C, NoteKey::D, NoteKey::E, NoteKey::F, NoteKey::G, NoteKey::A].iter().enumerate() {
        let x_pos = (i as f32 - 2.5) * lane_width;
        
        commands.spawn(SpriteBundle {
            sprite: Sprite {
                color: Color::rgb(0.3, 0.3, 0.3),
                custom_size: Some(Vec2::new(lane_width - 10.0, 50.0)),
                ..default()
            },
            transform: Transform::from_xyz(x_pos, -300.0, 0.0),
            ..default()
        });
        
        // Key label
        commands.spawn(Text2dBundle {
            text: Text::from_section(
                format!("{:?}", key),
                TextStyle {
                    font_size: 30.0,
                    color: Color::WHITE,
                    ..default()
                },
            )
            .with_alignment(TextAlignment::Center),
            transform: Transform::from_xyz(x_pos, -300.0, 1.0),
            ..default()
        });
    }
}

// Cleanup game when exiting
fn cleanup_game(
    mut commands: Commands,
    note_query: Query<Entity, With<Note>>,
    ui_query: Query<Entity, Or<(With<ScoreText>, With<ComboText>)>>,
) {
    // Remove all notes
    for entity in note_query.iter() {
        commands.entity(entity).despawn_recursive();
    }
    
    // Remove UI elements
    for entity in ui_query.iter() {
        commands.entity(entity).despawn_recursive();
    }
}

// Handle menu interactions
fn menu_interaction(
    mut next_state: ResMut<NextState<GameState>>,
    mut interaction_query: Query<
        (&Interaction, &mut BackgroundColor),
        (Changed<Interaction>, With<PlayButton>),
    >,
) {
    for (interaction, mut color) in &mut interaction_query {
        match *interaction {
            Interaction::Pressed => {
                next_state.set(GameState::Playing);
            }
            Interaction::Hovered => {
                *color = Color::rgb(0.25, 0.25, 0.25).into();
            }
            Interaction::None => {
                *color = Color::rgb(0.15, 0.15, 0.15).into();
            }
        }
    }
}

// Spawn new notes based on song data
fn spawn_notes(
    mut commands: Commands,
    time: Res<Time>,
    game_resources: Res<GameResources>,
) {
    // Create a copy of the song data to avoid borrowing issues
    let song_clone = match &game_resources.song {
        Some(song) => song.clone(),
        None => return,
    };
    
    // Get the current time
    let current_time = game_resources.song_timer.elapsed_secs();
    
    // Look ahead time - when notes should be spawned before they need to be hit
    let look_ahead = 2.0; // seconds
    
    // Find notes that need to be spawned
    for note in song_clone.notes.iter() {
        // Check if note should be spawned now
        if note.time - current_time <= look_ahead && note.time - current_time > look_ahead - time.delta_seconds() {
            // Calculate which lane this note belongs in
            let lane = match note.key {
                NoteKey::C => 0,
                NoteKey::D => 1,
                NoteKey::E => 2,
                NoteKey::F => 3,
                NoteKey::G => 4,
                NoteKey::A => 5,
            };
            
            let lane_width = 100.0;
            let x_pos = (lane as f32 - 2.5) * lane_width;
            
            // Calculate color based on note
            let color = match note.key {
                NoteKey::C => Color::rgb(1.0, 0.0, 0.0), // Red
                NoteKey::D => Color::rgb(1.0, 0.5, 0.0), // Orange
                NoteKey::E => Color::rgb(1.0, 1.0, 0.0), // Yellow
                NoteKey::F => Color::rgb(0.0, 1.0, 0.0), // Green
                NoteKey::G => Color::rgb(0.0, 0.0, 1.0), // Blue
                NoteKey::A => Color::rgb(0.5, 0.0, 1.0), // Purple
            };
            
            // Spawn the note entity
            commands.spawn((
                SpriteBundle {
                    sprite: Sprite {
                        color,
                        custom_size: Some(Vec2::new(lane_width - 10.0, 40.0)),
                        ..default()
                    },
                    transform: Transform::from_xyz(x_pos, 400.0, 0.0), // Start at top
                    ..default()
                },
                Note {
                    key: note.key,
                    lane,
                    spawn_time: current_time,
                    hit_time: note.time,
                },
                NoteSprite,
            ));
        }
    }
}

// Move notes down the screen
fn move_notes(
    mut commands: Commands,
    time: Res<Time>,
    game_resources: Res<GameResources>,
    mut note_query: Query<(Entity, &mut Transform), With<NoteSprite>>,
) {
    let note_speed = game_resources.note_speed;
    
    for (entity, mut transform) in note_query.iter_mut() {
        // Move the note down
        transform.translation.y -= note_speed * time.delta_seconds();
        
        // Remove notes that are off-screen
        if transform.translation.y < -400.0 {
            commands.entity(entity).despawn();
        }
    }
}

// Handle keyboard input to hit notes
fn handle_key_presses(
    mut commands: Commands,
    keyboard_input: Res<Input<KeyCode>>, // Changed from ButtonInput to Input
    mut game_resources: ResMut<GameResources>,
    note_query: Query<(Entity, &Transform, &Note), With<NoteSprite>>,
) {
    // Check for key presses
    for key_code in keyboard_input.get_just_pressed() {
        if let Some(note_key) = NoteKey::from_key_code(*key_code) {
            // Find the lane for this key
            let lane = match note_key {
                NoteKey::C => 0,
                NoteKey::D => 1,
                NoteKey::E => 2,
                NoteKey::F => 3,
                NoteKey::G => 4,
                NoteKey::A => 5,
            };
            
            // Look for notes in this lane
            let hit_zone_y = -300.0; // Where the note should be hit
            let hit_window = 50.0; // How close the note needs to be
            
            let mut closest_note: Option<(Entity, f32)> = None;
            
            // Find closest note in the correct lane
            for (entity, transform, note) in note_query.iter() {
                if note.lane == lane {
                    let distance = (transform.translation.y - hit_zone_y).abs();
                    if distance < hit_window {
                        if let Some((_, best_distance)) = closest_note {
                            if distance < best_distance {
                                closest_note = Some((entity, distance));
                            }
                        } else {
                            closest_note = Some((entity, distance));
                        }
                    }
                }
            }
            
            // If we found a note to hit
            if let Some((entity, distance)) = closest_note {
                // Score based on accuracy
                let accuracy = 1.0 - (distance / hit_window);
                let score_gain = (accuracy * 100.0).round() as u32;
                
                game_resources.score += score_gain;
                game_resources.combo += 1;
                
                if game_resources.combo > game_resources.max_combo {
                    game_resources.max_combo = game_resources.combo;
                }
                
                // Remove the hit note
                commands.entity(entity).despawn();
                
                // Here we would play the note sound
                // audio.play(asset_server.load(note_key.get_audio_path()));
            } else {
                // Missed - break combo
                game_resources.combo = 0;
            }
        }
    }
}

// Update the score and combo display
fn update_score_ui(
    game_resources: Res<GameResources>,
    mut score_query: Query<&mut Text, With<ScoreText>>,
    mut combo_query: Query<&mut Text, (With<ComboText>, Without<ScoreText>)>,
) {
    // Update score text
    if let Ok(mut text) = score_query.get_single_mut() {
        text.sections[0].value = format!("Score: {}", game_resources.score);
    }
    
    // Update combo text
    if let Ok(mut text) = combo_query.get_single_mut() {
        text.sections[0].value = format!("Combo: {}", game_resources.combo);
    }
}

// Check if the game is over
fn check_game_over(
    game_resources: Res<GameResources>,
    note_query: Query<&Note>,
    mut next_state: ResMut<NextState<GameState>>,
) {
    // Get current song time
    let current_time = game_resources.song_timer.elapsed_secs();
    
    // Get song end time
    if let Some(song) = &game_resources.song {
        let last_note_time = song.notes.iter().map(|note| note.time).max_by(|a, b| a.partial_cmp(b).unwrap()).unwrap_or(0.0);
        
        // If song is over and no notes remain
        if current_time > last_note_time + 3.0 && note_query.iter().count() == 0 {
            next_state.set(GameState::GameOver);
        }
    }
}