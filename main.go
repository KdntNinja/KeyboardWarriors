package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func main() {
	// Set up logging
	log.SetFlags(log.LstdFlags)
	log.Println("Starting Keyboard Warrior...")

	// Load songs from the songs directory
	songs, err := LoadSongsFromDirectory("songs")
	if err != nil {
		log.Printf("Error loading songs: %v", err)
		songs = []*Song{} // Initialize with empty slice if error
	}

	log.Printf("Loaded %d songs", len(songs))

	// Initialize audio
	log.Println("Initializing audio system...")
	audioManager, err := NewAudioManager()
	if err != nil {
		log.Printf("Error initializing audio: %v", err)
		// Continue without audio if there's an error
	}

	// Create the game
	game := NewGame(songs, audioManager)

	// Set up Ebiten
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Keyboard Warrior - Music Game")
	ebiten.SetWindowResizable(true)

	// Run the game
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
