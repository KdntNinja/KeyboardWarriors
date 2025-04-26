package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
	"image/color"
	"log"
	"time"
)

// GameState represents the current state of the game
type GameState int

const (
	StateTitle GameState = iota
	StatePlaying
	StateGameOver
	StateSongComplete
)

type Game struct {
	notes        []*Note
	keyBindings  map[string]ebiten.Key
	barY         float64
	score        int
	misses       int
	gameState    GameState
	songs        []*Song
	currentSong  int
	songPlayer   *SongPlayer
	lastHitTime  time.Time
	lastHitLane  int
	lastHitScore int // Score of the last hit note
	lastMissTime time.Time
	lastMissLane int
	audioManager *AudioManager
	accuracy     float64
	totalNotes   int
	hitNotes     int
	endTime      time.Time
	titleScreen  *TitleScreen // New title screen component
}

func NewGame(songs []*Song, audioManager *AudioManager) *Game {
	// QWEIOP key bindings
	keyBindings := map[string]ebiten.Key{
		"C": ebiten.KeyQ,
		"D": ebiten.KeyW,
		"E": ebiten.KeyE,
		"F": ebiten.KeyI,
		"G": ebiten.KeyO,
		"A": ebiten.KeyP,
	}

	// Calculate the hit bar position
	hitBarY := 400.0

	// Create a song player with the first song
	var songPlayer *SongPlayer
	if len(songs) > 0 {
		songPlayer = NewSongPlayer(songs[0])
	}

	// Create title screen
	titleScreen := NewTitleScreen(songs)

	return &Game{
		notes:        []*Note{},
		keyBindings:  keyBindings,
		barY:         hitBarY,
		gameState:    StateTitle,
		songs:        songs,
		currentSong:  0,
		songPlayer:   songPlayer,
		audioManager: audioManager,
		titleScreen:  titleScreen,
	}
}

func (g *Game) Update() error {
	switch g.gameState {
	case StateTitle:
		// Press space to start
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.startNewGame()
		}

		// Change song with up/down arrow keys when on the title screen
		songChanged := false
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) && g.currentSong < len(g.songs)-1 {
			g.currentSong++
			songChanged = true
		} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) && g.currentSong > 0 {
			g.currentSong--
			songChanged = true
		}

		// Update title screen and song player if song changed
		if songChanged {
			g.titleScreen.SetCurrentSong(g.currentSong)
			g.songPlayer = NewSongPlayer(g.songs[g.currentSong])
		}

	case StatePlaying:
		// Always update the song player in gameplay mode
		if g.songPlayer != nil {
			notesToAdd := g.songPlayer.Update()

			// Convert song notes to game notes and add them
			if len(notesToAdd) > 0 && g.songPlayer.IsPlaying {
				screenWidth, _ := g.Layout(0, 0)
				for i, songNote := range notesToAdd {
					newNote := CreateNoteFromSong(songNote, screenWidth, g.keyBindings, g.songs[g.currentSong].BPM)
					g.notes = append(g.notes, newNote)
					g.totalNotes++
					// Log note generation for debugging
					if i == 0 {
						log.Printf("Generated note: %s at time %.1f", songNote.Key, songNote.Time)
					}
				}
			}

			// Check if the song is over and all notes have been processed
			if !g.songPlayer.IsPlaying && !g.songPlayer.IsCounting && len(g.notes) == 0 {
				// Song is over and all notes are gone, show song complete screen
				g.gameState = StateSongComplete
				g.endTime = time.Now()

				// Calculate accuracy if there were any notes
				if g.totalNotes > 0 {
					g.accuracy = float64(g.hitNotes) / float64(g.totalNotes) * 100
				}
			}
		}

		// Update notes
		for _, note := range g.notes {
			note.Update()

			// Check for misses - a note is missed only when it has completely passed the hit line
			if note.status == StatusActive && note.HasPassedHitLine(g.barY) {
				note.Miss()
				g.misses++

				// Visual feedback for miss
				g.lastMissTime = time.Now()
				g.lastMissLane = note.lane

				// Game over after 20 misses
				if g.misses >= 20 {
					g.gameState = StateGameOver
					g.endTime = time.Now()
					if g.songPlayer != nil {
						g.songPlayer.Stop()
					}

					// Calculate accuracy
					if g.totalNotes > 0 {
						g.accuracy = float64(g.hitNotes) / float64(g.totalNotes) * 100
					}
				}
			}
		}

		// Check for key presses
		for key, ebitenKey := range g.keyBindings {
			// Handle key presses
			if inpututil.IsKeyJustPressed(ebitenKey) {
				noteHit := false

				// Check for note hits
				for _, note := range g.notes {
					if note.key == key && note.status == StatusActive {
						// Use the improved hitbox check - entire length of note
						if note.IsAtHitLine(g.barY) {
							// Calculate accuracy for scoring
							accuracy := note.GetHitAccuracy(g.barY)

							// Better scoring based on accuracy
							baseScore := 100
							accuracyBonus := int(float64(baseScore) * accuracy)
							totalScore := baseScore + accuracyBonus

							note.Hit()
							g.score += totalScore
							g.hitNotes++
							noteHit = true

							// Flash the lane to indicate hit
							g.lastHitTime = time.Now()
							g.lastHitLane = note.lane
							g.lastHitScore = totalScore

							// Play sound only when a note is actually hit
							if g.audioManager != nil {
								g.audioManager.PlayNote(key)
							}

							break
						}
					}
				}

				// Visual feedback even when no note is hit
				if !noteHit {
					// Just show that the key was pressed but no note was hit
					// No sound is played when no note is hit
				}
			}
		}

		// Clean up off-screen notes and hit/missed notes
		var activeNotes []*Note
		for _, note := range g.notes {
			// Keep notes that are on screen and active
			if note.y > -100 && note.y < 600 && note.status == StatusActive {
				activeNotes = append(activeNotes, note)
			}
		}
		g.notes = activeNotes

	case StateGameOver, StateSongComplete:
		// Return to the title screen after 5 seconds or when space is pressed
		if time.Since(g.endTime) > 5*time.Second || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.gameState = StateTitle
			// No need to stop sounds since they play once and stop automatically
		}
	}

	return nil
}

// startNewGame resets game state and starts playing
func (g *Game) startNewGame() {
	g.notes = []*Note{}
	g.score = 0
	g.misses = 0
	g.totalNotes = 0
	g.hitNotes = 0
	g.accuracy = 0
	g.lastHitScore = 0
	g.gameState = StatePlaying

	// Start the song player with countdown
	if g.songPlayer != nil {
		g.songPlayer.Start()
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screenWidth, screenHeight := g.Layout(0, 0)

	// Background
	screen.Fill(color.Black)

	switch g.gameState {
	case StateTitle:
		// Draw the improved title screen
		g.titleScreen.Draw(screen, screenWidth, screenHeight)

	case StatePlaying:
		// Draw lane separators and lane highlights
		laneWidth := float64(screenWidth) / 6

		// Draw countdown if active
		if g.songPlayer != nil && g.songPlayer.IsCounting {
			countdown := g.songPlayer.GetCountdownSeconds()
			countdownText := fmt.Sprintf("%d", countdown)

			// Center of screen
			centerX := screenWidth / 2
			fontColor := color.RGBA{R: 255, G: 255, B: 255, A: 255}

			// Add "GET READY!" text above
			readyText := "GET READY!"
			text.Draw(screen, readyText, basicfont.Face7x13,
				centerX-len(readyText)*4,
				screenHeight/2-40,
				fontColor)

			// Draw the countdown number
			text.Draw(screen, countdownText, basicfont.Face7x13,
				centerX-4,
				screenHeight/2+20,
				fontColor)
		}

		for i := 0; i < 6; i++ {
			// Draw hit/miss feedback
			// Hit feedback (flashes blue/green/gold based on score)
			if i == g.lastHitLane && time.Since(g.lastHitTime) < 200*time.Millisecond {
				hitFeedback := ebiten.NewImage(int(laneWidth), 480)
				alpha := 255 - uint8(time.Since(g.lastHitTime).Milliseconds())

				// Default hit color (blue)
				hitColor := color.RGBA{R: 100, G: 180, B: 255, A: alpha}

				// Perfect hit color (gold) for scores over 150
				if g.lastHitScore > 150 {
					hitColor = color.RGBA{R: 255, G: 215, B: 0, A: alpha}
				} else if g.lastHitScore > 120 {
					// Good hit color (green) for scores over 120
					hitColor = color.RGBA{R: 50, G: 205, B: 50, A: alpha}
				}

				hitFeedback.Fill(hitColor)
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(i)*laneWidth, 0)
				screen.DrawImage(hitFeedback, op)
			}

			// Miss feedback (flashes red)
			if i == g.lastMissLane && time.Since(g.lastMissTime) < 200*time.Millisecond {
				missFeedback := ebiten.NewImage(int(laneWidth), 480)
				alpha := 255 - uint8(time.Since(g.lastMissTime).Milliseconds())
				missFeedback.Fill(color.RGBA{R: 200, G: 0, B: 0, A: alpha})
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(i)*laneWidth, 0)
				screen.DrawImage(missFeedback, op)
			}

			// Draw lane separators
			if i > 0 {
				x := float64(i) * laneWidth
				lineImg := ebiten.NewImage(1, 480)
				lineImg.Fill(color.RGBA{R: 50, G: 50, B: 50, A: 255})
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(x, 0)
				screen.DrawImage(lineImg, op)
			}
		}

		// Draw the hit bar
		barImage := ebiten.NewImage(screenWidth, 2)
		barImage.Fill(color.White)
		barOp := &ebiten.DrawImageOptions{}
		barOp.GeoM.Translate(0, g.barY)
		screen.DrawImage(barImage, barOp)

		// Draw notes - only draw active notes
		for _, note := range g.notes {
			// Only draw active notes, skip hit or missed notes
			if note.status == StatusActive {
				// Create appropriate sized note image based on hold time
				noteImage := ebiten.NewImage(40, int(note.height))
				noteImage.Fill(GetNoteColor(note.key))

				// For hold notes, add a border
				if note.holdTime > 0 {
					// Draw a border for hold notes
					borderColor := color.RGBA{R: 255, G: 255, B: 255, A: 180}
					borderImage := ebiten.NewImage(40, int(note.height))
					borderImage.Fill(borderColor)

					// Draw border as a 2px frame
					noteInnerImage := ebiten.NewImage(36, int(note.height)-4)
					noteInnerImage.Fill(GetNoteColor(note.key))

					innerOp := &ebiten.DrawImageOptions{}
					innerOp.GeoM.Translate(2, 2)
					borderImage.DrawImage(noteInnerImage, innerOp)

					op := &ebiten.DrawImageOptions{}
					op.GeoM.Translate(note.x, note.y)
					screen.DrawImage(borderImage, op)
				} else {
					// Regular note without hold
					op := &ebiten.DrawImageOptions{}
					op.GeoM.Translate(note.x, note.y)
					screen.DrawImage(noteImage, op)
				}

				// Draw the note letter on top of the note
				text.Draw(screen, note.key, basicfont.Face7x13, int(note.x)+15, int(note.y)+15, color.Black)
			}
		}

		// Draw game info
		text.Draw(screen, fmt.Sprintf("Score: %d", g.score), basicfont.Face7x13, 10, 20, color.White)
		text.Draw(screen, fmt.Sprintf("Misses: %d/20", g.misses), basicfont.Face7x13, 10, 40, color.White)

		if g.songPlayer != nil && g.currentSong < len(g.songs) {
			songName := fmt.Sprintf("Playing: %s", g.songs[g.currentSong].Title)
			text.Draw(screen, songName, basicfont.Face7x13, 10, 60, color.White)

			// Show time remaining
			timeRemaining := g.songs[g.currentSong].Duration - g.songPlayer.ElapsedTime
			if timeRemaining < 0 {
				timeRemaining = 0
			}
			timeText := fmt.Sprintf("Time: %.1f", timeRemaining)
			text.Draw(screen, timeText, basicfont.Face7x13, screenWidth-100, 20, color.White)
		}

	case StateGameOver:
		// Draw game over screen
		var song *Song
		if len(g.songs) > 0 && g.currentSong < len(g.songs) {
			song = g.songs[g.currentSong]
		}
		DrawEndScreen(screen, screenWidth, "GAME OVER", g.score, g.hitNotes, g.totalNotes, g.accuracy, song)

	case StateSongComplete:
		// Draw song complete screen with stats
		var song *Song
		if len(g.songs) > 0 && g.currentSong < len(g.songs) {
			song = g.songs[g.currentSong]
		}
		DrawEndScreen(screen, screenWidth, "SONG COMPLETE!", g.score, g.hitNotes, g.totalNotes, g.accuracy, song)
	}
}

// DrawEndScreen draws the game over or song complete screen
func DrawEndScreen(screen *ebiten.Image, width int, message string, score, hits, total int, accuracy float64, song *Song) {
	centerX := width / 2

	// Draw heading message
	text.Draw(screen, message, basicfont.Face7x13, centerX-len(message)*4, 120, color.White)

	// Draw song information
	if song != nil {
		songTitle := fmt.Sprintf("Song: %s", song.Title)
		text.Draw(screen, songTitle, basicfont.Face7x13, centerX-len(songTitle)*3, 170, color.White)

		songArtist := fmt.Sprintf("By: %s", song.Artist)
		text.Draw(screen, songArtist, basicfont.Face7x13, centerX-len(songArtist)*3, 190, color.White)
	}

	// Draw statistics
	scoreText := fmt.Sprintf("Final Score: %d", score)
	text.Draw(screen, scoreText, basicfont.Face7x13, centerX-len(scoreText)*3, 230, color.White)

	hitsText := fmt.Sprintf("Notes Hit: %d/%d", hits, total)
	text.Draw(screen, hitsText, basicfont.Face7x13, centerX-len(hitsText)*3, 250, color.White)

	accuracyText := fmt.Sprintf("Accuracy: %.1f%%", accuracy)
	text.Draw(screen, accuracyText, basicfont.Face7x13, centerX-len(accuracyText)*3, 270, color.White)

	// Draw rank based on accuracy
	var rankText string
	switch {
	case accuracy >= 95:
		rankText = "Rank: S (Amazing!)"
	case accuracy >= 90:
		rankText = "Rank: A (Excellent!)"
	case accuracy >= 80:
		rankText = "Rank: B (Great!)"
	case accuracy >= 70:
		rankText = "Rank: C (Good)"
	case accuracy >= 60:
		rankText = "Rank: D (Fair)"
	default:
		rankText = "Rank: F (Need Practice)"
	}
	text.Draw(screen, rankText, basicfont.Face7x13, centerX-len(rankText)*3, 300, color.RGBA{R: 255, G: 215, B: 0, A: 255})

	// Draw prompt to continue
	continueText := "Press SPACE to continue"
	text.Draw(screen, continueText, basicfont.Face7x13, centerX-len(continueText)*3, 340, color.White)
}

// Layout implements the ebiten.Game interface
func (g *Game) Layout(int, int) (int, int) {
	return 640, 480
}
