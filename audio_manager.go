package main

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"log"
	"os"
	"path/filepath"
)

const (
	sampleRate = 44100
)

// AudioManager handles loading and playing audio files
type AudioManager struct {
	initialized  bool
	audioContext *audio.Context
	noteSounds   map[string][]byte // Store audio data for each note
}

// NewAudioManager creates a new audio manager
func NewAudioManager() *AudioManager {
	return &AudioManager{
		initialized: false,
		noteSounds:  make(map[string][]byte),
	}
}

// Initialize sets up the audio context and loads sound files
func (am *AudioManager) Initialize() error {
	log.Println("Initializing audio system...")

	// Create audio context
	am.audioContext = audio.NewContext(sampleRate)

	// Load note sound files from the audio directory
	notes := []string{"C", "D", "E", "F", "G", "A"}

	for _, note := range notes {
		filename := filepath.Join("audio", "note_"+note+".mp3")

		// Check if file exists
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			log.Printf("Warning: Sound file %s not found", filename)
			continue
		}

		// Load the audio file data
		data, err := os.ReadFile(filename)
		if err != nil {
			log.Printf("Error reading sound file %s: %v", filename, err)
			continue
		}

		// Store the audio data
		am.noteSounds[note] = data
		log.Printf("Loaded sound for note: %s from %s", note, filename)
	}

	am.initialized = true
	log.Println("Audio system initialized")
	return nil
}

// PlayNote plays the sound for the given note
func (am *AudioManager) PlayNote(note string) {
	if !am.initialized {
		return
	}

	// Get the audio data for this note
	data, found := am.noteSounds[note]
	if !found {
		log.Printf("No sound loaded for note: %s", note)
		return
	}

	// Create a reader for the audio data
	reader := bytes.NewReader(data)

	// Decode the MP3 data
	decoded, err := mp3.DecodeWithSampleRate(sampleRate, reader)
	if err != nil {
		log.Printf("Error decoding MP3 data for note %s: %v", note, err)
		return
	}

	// Create a new player for this sound
	player, err := am.audioContext.NewPlayer(decoded)
	if err != nil {
		log.Printf("Error creating player for note %s: %v", note, err)
		return
	}

	// Play the sound once
	player.Play()

	// For a 1-second audio file, we can just let it play without explicitly stopping it
	// The player will be garbage collected after it finishes

	log.Printf("Playing note: %s", note)
}

// Close cleans up the audio resources
func (am *AudioManager) Close() {
	if !am.initialized {
		return
	}

	log.Println("Audio system closed")
}
