package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// SongNote represents a note in a song
type SongNote struct {
	Key  string  `json:"key"`
	Lane int     `json:"lane"`
	Time float64 `json:"time"`
	Hold float64 `json:"hold"`
}

// Song represents a playable song
type Song struct {
	Title    string     `json:"title"`
	Artist   string     `json:"artist"`
	BPM      int        `json:"bpm"`
	Duration float64    `json:"duration"`
	Notes    []SongNote `json:"notes"`
}

// CalculateDifficulty returns a difficulty score (0-100) based on note count and BPM
func (s *Song) CalculateDifficulty() int {
	// More notes and higher BPM = higher difficulty
	noteCount := len(s.Notes)

	// Calculate note density (notes per second)
	noteDensity := float64(noteCount) / s.Duration

	// Calculate hold note percentage
	holdNotes := 0
	for _, note := range s.Notes {
		if note.Hold > 0 {
			holdNotes++
		}
	}
	holdPercentage := float64(holdNotes) / float64(noteCount)

	// Calculate base difficulty score
	baseScore := (noteDensity * 20.0) + (float64(s.BPM) / 4.0)

	// Add bonus for hold notes
	holdBonus := holdPercentage * 20.0

	// Calculate final score
	finalScore := baseScore + holdBonus

	// Cap at 0-100
	if finalScore > 100 {
		finalScore = 100
	}
	if finalScore < 0 {
		finalScore = 0
	}

	return int(finalScore)
}

// SongPlayer handles playback timing and note spawning
type SongPlayer struct {
	Song         *Song
	StartTime    time.Time
	CountdownEnd time.Time
	IsPlaying    bool
	IsCounting   bool
	ElapsedTime  float64
	LastNoteTime float64
}

// NewSongPlayer creates a new song player for a song
func NewSongPlayer(song *Song) *SongPlayer {
	return &SongPlayer{
		Song:         song,
		IsPlaying:    false,
		IsCounting:   false,
		ElapsedTime:  0,
		LastNoteTime: -1,
	}
}

// Start begins playing the song after a countdown
func (sp *SongPlayer) Start() {
	// Start with a 3-second countdown
	sp.IsPlaying = false
	sp.IsCounting = true
	sp.ElapsedTime = 0
	sp.LastNoteTime = -1
	sp.StartTime = time.Now()
	sp.CountdownEnd = sp.StartTime.Add(3 * time.Second)
}

// GetCountdownSeconds returns the current countdown in seconds (3,2,1)
func (sp *SongPlayer) GetCountdownSeconds() int {
	if !sp.IsCounting {
		return 0
	}

	remaining := sp.CountdownEnd.Sub(time.Now())
	seconds := int(remaining.Seconds()) + 1

	if seconds < 1 {
		seconds = 1
	}

	return seconds
}

// Stop halts the song playback
func (sp *SongPlayer) Stop() {
	sp.IsPlaying = false
	sp.IsCounting = false
}

// Update updates the song time and spawns notes
// Returns notes that should be spawned
func (sp *SongPlayer) Update() []SongNote {
	if sp.IsCounting {
		// Check if countdown is over
		if time.Now().After(sp.CountdownEnd) {
			sp.IsPlaying = true
			sp.IsCounting = false
			sp.StartTime = time.Now()
		}
		return []SongNote{} // No notes during countdown
	}

	if !sp.IsPlaying {
		return []SongNote{}
	}

	// Calculate elapsed time
	sp.ElapsedTime = time.Since(sp.StartTime).Seconds()

	// Check if song is over
	if sp.ElapsedTime >= sp.Song.Duration {
		sp.IsPlaying = false
		return []SongNote{}
	}

	// Find notes that should be spawned now
	// We want to spawn notes early enough that they reach the hit line at the right time
	var notesToSpawn []SongNote

	for _, note := range sp.Song.Notes {
		// Skip notes we've already spawned
		if sp.LastNoteTime >= note.Time {
			continue
		}

		// Only spawn notes that should appear within NoteApproachTime
		if note.Time <= sp.ElapsedTime+NoteApproachTime {
			notesToSpawn = append(notesToSpawn, note)
			sp.LastNoteTime = note.Time
		}
	}

	return notesToSpawn
}

// LoadSongsFromDirectory loads all song files from a directory
func LoadSongsFromDirectory(dirPath string) ([]*Song, error) {
	var songs []*Song

	// Open the directory
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	// Process each JSON file
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			// Load the song file
			filePath := filepath.Join(dirPath, file.Name())
			song, err := LoadSongFromFile(filePath)
			if err != nil {
				continue // Skip if there's an error
			}

			songs = append(songs, song)
		}
	}

	return songs, nil
}

// LoadSongFromFile loads a song from a JSON file
func LoadSongFromFile(filePath string) (*Song, error) {
	// Read the file
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Parse the JSON
	var song Song
	err = json.Unmarshal(fileData, &song)
	if err != nil {
		return nil, err
	}

	return &song, nil
}
