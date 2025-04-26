package main

import (
	"encoding/json"
	"os"
	"time"
)

// SongNote represents a note in a song
type SongNote struct {
	Key      string  `json:"key"`      // Musical note (C, D, E, etc.)
	Lane     int     `json:"lane"`     // Which lane (0-5 for 6 keys)
	Time     float64 `json:"time"`     // Time in seconds when note should appear
	Duration float64 `json:"duration"` // Optional: How long the note should be held (for future use)
}

// Song represents a complete song with multiple notes
type Song struct {
	Title    string     `json:"title"`
	Artist   string     `json:"artist"`
	BPM      int        `json:"bpm"`
	Notes    []SongNote `json:"notes"`
	Duration float64    `json:"duration"` // Total song duration in seconds
}

// SongPlayer manages playback of a song
type SongPlayer struct {
	Song           *Song
	StartTime      time.Time
	CurrentIndex   int
	IsPlaying      bool
	ElapsedTime    float64
	NotesGenerated map[int]bool // Track which notes have been generated
}

// LoadSongFromFile loads a song from a JSON file
func LoadSongFromFile(filename string) (*Song, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var song Song
	err = json.Unmarshal(data, &song)
	if err != nil {
		return nil, err
	}

	return &song, nil
}

// NewSongPlayer creates a new song player for the given song
func NewSongPlayer(song *Song) *SongPlayer {
	return &SongPlayer{
		Song:           song,
		CurrentIndex:   0,
		IsPlaying:      false,
		NotesGenerated: make(map[int]bool),
	}
}

// Start begins playback of the song
func (sp *SongPlayer) Start() {
	sp.StartTime = time.Now()
	sp.IsPlaying = true
	sp.CurrentIndex = 0
	sp.ElapsedTime = 0
	sp.NotesGenerated = make(map[int]bool)
}

// Stop stops playback of the song
func (sp *SongPlayer) Stop() {
	sp.IsPlaying = false
}

// Update updates the song player state and returns any notes that should be generated
func (sp *SongPlayer) Update() []SongNote {
	if !sp.IsPlaying {
		return nil
	}

	sp.ElapsedTime = time.Since(sp.StartTime).Seconds()

	// Check if the song is over
	if sp.ElapsedTime >= sp.Song.Duration {
		sp.IsPlaying = false
		return nil
	}

	// Look ahead by 2 seconds to generate notes in advance
	lookaheadTime := sp.ElapsedTime + 2.0

	var notesToGenerate []SongNote

	// Check for notes that should be generated
	for i, note := range sp.Song.Notes {
		// Skip notes we've already generated
		if sp.NotesGenerated[i] {
			continue
		}

		// If this note's time is within our lookahead window
		if note.Time <= lookaheadTime {
			notesToGenerate = append(notesToGenerate, note)
			sp.NotesGenerated[i] = true
		}
	}

	return notesToGenerate
}
