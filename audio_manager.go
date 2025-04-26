package main

import (
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// AudioManager handles sound effects for the game
type AudioManager struct {
	audioContext *audio.Context
	noteSounds   map[string]*audio.Player
	mutex        sync.Mutex
}

// NewAudioManager creates a new audio manager
func NewAudioManager() (*AudioManager, error) {
	// Initialize audio with 44.1kHz sample rate
	audioContext := audio.NewContext(44100)

	// Create audio manager
	am := &AudioManager{
		audioContext: audioContext,
		noteSounds:   make(map[string]*audio.Player),
		mutex:        sync.Mutex{},
	}

	// Load note sounds
	err := am.loadSounds()
	if err != nil {
		return nil, err
	}

	log.Println("Audio system initialized")
	return am, nil
}

// loadSounds loads all note sound files
func (am *AudioManager) loadSounds() error {
	// List of notes to load
	notes := []string{"C", "D", "E", "F", "G", "A"}
	successCount := 0

	// Load each note sound
	for _, note := range notes {
		filename := filepath.Join("audio", "note_"+note+".mp3")
		player, err := am.loadSound(filename)
		if err != nil {
			log.Printf("Error loading sound for note %s: %v", note, err)
			continue
		}

		// Store the player
		am.noteSounds[note] = player
		log.Printf("Loaded sound for note: %s from %s", note, filename)
		successCount++
	}

	log.Printf("Successfully loaded %d of %d sound files", successCount, len(notes))
	return nil
}

// loadSound loads a single sound file and returns an audio player
func (am *AudioManager) loadSound(filename string) (*audio.Player, error) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	// Close the file when we're done
	defer file.Close()

	// Decode the MP3
	decoded, err := mp3.DecodeWithSampleRate(44100, file)
	if err != nil {
		return nil, err
	}

	// Create a byte slice from the decoded audio
	audioBytes, err := io.ReadAll(decoded)
	if err != nil {
		return nil, err
	}

	// Create an audio player from the bytes
	player := audio.NewPlayerFromBytes(am.audioContext, audioBytes)
	return player, nil
}

// PlayNote plays the sound for a specific note
func (am *AudioManager) PlayNote(note string) {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	// Check if we have the note sound
	player, exists := am.noteSounds[note]
	if !exists {
		return
	}

	// Rewind and play the sound
	if player != nil {
		player.Rewind()
		player.Play()
	}
}
