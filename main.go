package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	// Seed random number generator
	rand.NewSource(time.Now().UnixNano())

	// Ensure audio directory exists
	audioDir := "audio"
	if _, err := os.Stat(audioDir); os.IsNotExist(err) {
		err := os.MkdirAll(audioDir, 0755)
		if err != nil {
			log.Printf("Warning: Could not create audio directory: %v", err)
		}
	}

	// Load songs from JSON files
	var songs []*Song

	simpleSong, err := LoadSongFromFile("songs/simple_song.json")
	if err != nil {
		log.Printf("Error loading simple_song.json: %v", err)
	} else {
		songs = append(songs, simpleSong)
	}

	twinkleSong, err := LoadSongFromFile("songs/twinkle_star.json")
	if err != nil {
		log.Printf("Error loading twinkle_star.json: %v", err)
	} else {
		songs = append(songs, twinkleSong)
	}

	// If no songs were loaded, create a default song
	if len(songs) == 0 {
		defaultSong := &Song{
			Title:    "Default Song",
			Artist:   "System",
			BPM:      120,
			Duration: 20.0,
			Notes: []SongNote{
				{Key: "C", Lane: 0, Time: 1.0},
				{Key: "D", Lane: 1, Time: 2.0},
				{Key: "E", Lane: 2, Time: 3.0},
				{Key: "F", Lane: 3, Time: 4.0},
				{Key: "G", Lane: 4, Time: 5.0},
				{Key: "A", Lane: 5, Time: 6.0},
			},
		}
		songs = append(songs, defaultSong)
	}

	// Initialize audio manager
	audioManager := NewAudioManager()
	if err := audioManager.Initialize(); err != nil {
		log.Printf("Warning: Failed to initialize audio: %v", err)
	}

	// Create and run the game with the loaded songs
	game := NewGame(songs, audioManager)
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Keyboard Warrior - JSON Music")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

	// Clean up audio resources
	audioManager.Close()
}
