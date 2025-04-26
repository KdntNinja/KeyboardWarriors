package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"math"
)

// NoteStatus represents the state of a note
type NoteStatus int

const (
	StatusActive NoteStatus = iota
	StatusHit
	StatusMissed
)

// NoteApproachTime is the time in seconds it takes for a note to reach the hit line
const NoteApproachTime = 2.0

// Note represents a playable note in the game
type Note struct {
	key      string     // The key to press (C, D, E, etc.)
	lane     int        // Which lane the note appears in (0-5)
	x        float64    // X position on screen
	y        float64    // Y position on screen
	height   float64    // Height of the note (larger for hold notes)
	speed    float64    // Speed the note travels
	status   NoteStatus // Status of the note (active, hit, missed)
	holdTime float64    // How long to hold the note (in seconds)
}

// CreateNoteFromSong creates a new game note from a song note
func CreateNoteFromSong(songNote SongNote, screenWidth int, keyBindings map[string]ebiten.Key, bpm int) *Note {
	// Calculate lane width
	laneWidth := float64(screenWidth) / 6

	// Calculate lane from key mapping if not provided
	lane := songNote.Lane
	if lane < 0 {
		// Default mapping: C=0, D=1, E=2, F=3, G=4, A=5
		keyToLane := map[string]int{
			"C": 0,
			"D": 1,
			"E": 2,
			"F": 3,
			"G": 4,
			"A": 5,
		}
		if l, exists := keyToLane[songNote.Key]; exists {
			lane = l
		}
	}

	// Calculate x position based on lane
	x := (float64(lane) * laneWidth) + (laneWidth/2 - 20)

	// Start above the screen
	y := -20.0

	// Calculate height based on hold time
	height := 20.0 // Default height for regular notes
	holdTime := songNote.Hold

	// Make sure hold time is reasonable
	if holdTime < 0 {
		holdTime = 0
	}

	if holdTime > 0 {
		// For hold notes, make the height proportional to the hold time
		// 100 pixels per second of hold time is a good starting point
		height = holdTime * 100

		// Set reasonable max
		if height > 400 {
			height = 400 // Don't let it get too large
		}
	}

	// Calculate speed based on BPM
	// Faster BPM means notes should move faster
	// Base speed is 2.0 pixels per frame at 100 BPM
	baseSpeed := 2.0
	speedMultiplier := float64(bpm) / 100.0
	speed := baseSpeed * speedMultiplier

	// Ensure speed is within reasonable bounds
	if speed < 1.0 {
		speed = 1.0 // Minimum speed
	} else if speed > 4.0 {
		speed = 4.0 // Maximum speed
	}

	// Adjust starting position for hold notes
	if holdTime > 0 {
		y = -20 - height + 20 // Start higher for hold notes
	}

	return &Note{
		key:      songNote.Key,
		lane:     lane,
		x:        x,
		y:        y,
		height:   height,
		speed:    speed,
		status:   StatusActive,
		holdTime: holdTime,
	}
}

// Update updates the note's position
func (n *Note) Update() {
	n.y += n.speed // Dynamic speed based on BPM
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

// IsAtHitLine checks if any part of the note is at the hit line
func (n *Note) IsAtHitLine(hitLineY float64) bool {
	// For regular notes, use the previous window check
	if n.height <= 20 {
		// Allow a window of +/- 30 pixels around the hit line
		return math.Abs(n.y-hitLineY) <= 30
	} else {
		// For hold notes, check if any part of the note is at the hit line
		noteTop := n.y
		noteBottom := n.y + n.height

		// The hold note is hittable if the hit line is between top and bottom
		// or if it's very close to either end
		return (noteTop-30 <= hitLineY && hitLineY <= noteBottom+30)
	}
}

// HasPassedHitLine checks if the note has completely passed the hit line
func (n *Note) HasPassedHitLine(hitLineY float64) bool {
	// For regular notes
	if n.height <= 20 {
		return n.y > hitLineY+30 // Note is fully below hit line
	} else {
		// For hold notes, check if the entire note has passed
		noteBottom := n.y + n.height
		return noteBottom > hitLineY+30
	}
}

// GetHitAccuracy returns a value from 0.0 to 1.0 based on how accurate the hit was
func (n *Note) GetHitAccuracy(hitLineY float64) float64 {
	// For regular notes
	if n.height <= 20 {
		// Calculate distance from perfect hit (closer is better)
		distance := math.Abs(n.y - hitLineY)

		// Convert to a 0-1 scale (max distance is 30 pixels)
		if distance > 30 {
			return 0.0
		}
		return 1.0 - (distance / 30.0)
	} else {
		// For hold notes, check how close to perfect the timing is
		// Distance from note head to hit line
		headDistance := math.Abs(n.y - hitLineY)

		// If note head is very close to hit line, it's a perfect hit
		if headDistance < 10 {
			return 1.0
		}

		// If hit line is within the hold note body, it's a good hit
		noteBottom := n.y + n.height
		if n.y <= hitLineY && hitLineY <= noteBottom {
			return 0.8
		}

		// Otherwise, calculate based on distance
		if headDistance > 30 {
			return 0.5
		}
		return 0.7 - (headDistance / 100.0)
	}
}

// GetNoteColor returns a color for a specific note key
func GetNoteColor(key string) color.RGBA {
	switch key {
	case "C":
		return color.RGBA{R: 255, G: 50, B: 50, A: 255} // Red
	case "D":
		return color.RGBA{R: 255, G: 150, B: 50, A: 255} // Orange
	case "E":
		return color.RGBA{R: 255, G: 255, B: 50, A: 255} // Yellow
	case "F":
		return color.RGBA{R: 50, G: 255, B: 50, A: 255} // Green
	case "G":
		return color.RGBA{R: 50, G: 150, B: 255, A: 255} // Blue
	case "A":
		return color.RGBA{R: 150, G: 50, B: 255, A: 255} // Purple
	default:
		return color.RGBA{R: 200, G: 200, B: 200, A: 255} // Gray
	}
}
