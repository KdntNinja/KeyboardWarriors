package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

// NoteStatus represents the current status of a note
type NoteStatus int

const (
	StatusActive NoteStatus = iota
	StatusHit
	StatusMissed
)

// Note represents a musical note in the game
type Note struct {
	key        string     // Musical note (C, D, E, etc.)
	keyBinding ebiten.Key // Keyboard key to press
	x, y       float64
	status     NoteStatus
	lane       int // Which lane the note is in (0-5)
}

// Update moves the note down the screen
func (n *Note) Update() {
	n.y += 2.0 // Fixed speed
}

// Hit marks the note as hit and removes it
func (n *Note) Hit() {
	n.status = StatusHit
	n.y = -5000 // Move far offscreen so it will be cleaned up
}

// Miss marks the note as missed and removes it
func (n *Note) Miss() {
	n.status = StatusMissed
	n.y = -5000 // Move far offscreen so it will be cleaned up
}

// CreateNoteFromSong creates a new note from song data
func CreateNoteFromSong(songNote SongNote, screenWidth int, keyBindings map[string]ebiten.Key) *Note {
	// Calculate x position based on lane
	keyWidth := float64(screenWidth) / 6
	x := float64(songNote.Lane)*keyWidth + keyWidth/2 - 20

	return &Note{
		key:        songNote.Key,
		keyBinding: keyBindings[songNote.Key],
		x:          x,
		y:          -50,
		status:     StatusActive,
		lane:       songNote.Lane,
	}
}

// GetNoteColor returns a color for the note based on its key
func GetNoteColor(key string) color.RGBA {
	switch key {
	case "C":
		return color.RGBA{R: 220, G: 60, B: 60, A: 255} // Softer red
	case "D":
		return color.RGBA{R: 220, G: 140, B: 40, A: 255} // Amber/orange
	case "E":
		return color.RGBA{R: 220, G: 220, B: 60, A: 255} // Softer yellow
	case "F":
		return color.RGBA{R: 60, G: 180, B: 100, A: 255} // Softer green
	case "G":
		return color.RGBA{R: 60, G: 100, B: 220, A: 255} // Softer blue
	case "A":
		return color.RGBA{R: 140, G: 60, B: 200, A: 255} // Violet/purple
	default:
		return color.RGBA{R: 200, G: 200, B: 200, A: 255} // Light gray
	}
}
