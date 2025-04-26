package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
	"image/color"
	"math"
	"time"
)

// TitleScreen represents the game's title screen UI
type TitleScreen struct {
	songs       []*Song
	currentSong int
	startTime   time.Time
	selectedBg  *ebiten.Image
}

// NewTitleScreen creates a new title screen
func NewTitleScreen(songs []*Song) *TitleScreen {
	// Create a background image for the selected song
	selectedBg := ebiten.NewImage(220, 50)
	selectedBg.Fill(color.RGBA{R: 40, G: 60, B: 120, A: 255})

	return &TitleScreen{
		songs:       songs,
		currentSong: 0,
		startTime:   time.Now(),
		selectedBg:  selectedBg,
	}
}

// SetCurrentSong updates the currently selected song
func (ts *TitleScreen) SetCurrentSong(index int) {
	if index >= 0 && index < len(ts.songs) {
		ts.currentSong = index
	}
}

// Draw renders the title screen
func (ts *TitleScreen) Draw(screen *ebiten.Image, screenWidth, screenHeight int) {
	// Draw header background
	headerBg := ebiten.NewImage(screenWidth, 120)
	headerBg.Fill(color.RGBA{R: 25, G: 25, B: 50, A: 255})
	screen.DrawImage(headerBg, &ebiten.DrawImageOptions{})

	// Draw decorative line under header
	decorLine := ebiten.NewImage(screenWidth, 2)
	decorLine.Fill(color.RGBA{R: 60, G: 100, B: 200, A: 255})
	lineOp := &ebiten.DrawImageOptions{}
	lineOp.GeoM.Translate(0, 120)
	screen.DrawImage(decorLine, lineOp)

	// Draw title with larger visual impact
	titleText := "KEYBOARD WARRIOR"
	titleX := screenWidth/2 - len(titleText)*4
	titleY := 60

	// Draw shadow for title text
	text.Draw(screen, titleText, basicfont.Face7x13, titleX+1, titleY+1, color.RGBA{R: 40, G: 40, B: 80, A: 255})
	text.Draw(screen, titleText, basicfont.Face7x13, titleX, titleY, color.RGBA{R: 255, G: 255, B: 255, A: 255})

	// Draw animated subtitle with pulsing effect
	pulse := float64(math.Sin(float64(time.Since(ts.startTime).Milliseconds())/500.0)*0.5 + 0.5)
	subtitleColor := color.RGBA{
		R: 255,
		G: 255,
		B: uint8(200 + pulse*55),
		A: 255,
	}
	subtitleText := "MUSIC RHYTHM GAME"
	subtitleX := screenWidth/2 - len(subtitleText)*4
	subtitleY := 90
	text.Draw(screen, subtitleText, basicfont.Face7x13, subtitleX, subtitleY, subtitleColor)

	// Draw song selection section header with background
	sectionBg := ebiten.NewImage(220, 25)
	sectionBg.Fill(color.RGBA{R: 40, G: 40, B: 90, A: 255})
	sectionOp := &ebiten.DrawImageOptions{}
	sectionOp.GeoM.Translate(float64(screenWidth/2-110), 142)
	screen.DrawImage(sectionBg, sectionOp)

	text.Draw(screen, "SELECT SONG", basicfont.Face7x13, screenWidth/2-43, 160, color.White)

	// Draw song list background
	listBg := ebiten.NewImage(400, len(ts.songs)*35+20)
	listBg.Fill(color.RGBA{R: 20, G: 20, B: 40, A: 150})
	listBgOp := &ebiten.DrawImageOptions{}
	listBgOp.GeoM.Translate(float64(screenWidth/2-200), 175)
	screen.DrawImage(listBg, listBgOp)

	// Draw song list
	listY := 200
	listX := screenWidth/2 - 180
	const songSpacing = 35 // Increased spacing between songs

	// Draw selected song highlight
	selectedBgOp := &ebiten.DrawImageOptions{}
	selectedBgOp.GeoM.Translate(float64(screenWidth/2-190), float64(listY+ts.currentSong*songSpacing-15))
	ts.selectedBg.Fill(color.RGBA{R: 40, G: 70, B: 140, A: 255}) // Brighter blue
	screen.DrawImage(ts.selectedBg, selectedBgOp)

	// Draw each song with better spacing
	for i, song := range ts.songs {
		// Determine color based on selection
		songColor := color.RGBA{R: 200, G: 200, B: 200, A: 255}
		if i == ts.currentSong {
			// Animate the selected song text
			pulse := float64(math.Sin(float64(time.Since(ts.startTime).Milliseconds())/300.0)*0.5 + 0.5)
			songColor = color.RGBA{
				R: 255,
				G: 255,
				B: uint8(180 + pulse*75),
				A: 255,
			}
		}

		// Draw song title and artist
		songTitle := fmt.Sprintf("%s - %s", song.Title, song.Artist)
		text.Draw(screen, songTitle, basicfont.Face7x13, listX, listY+i*songSpacing, songColor)

		// Draw difficulty indicator with a visual background
		difficulty := song.CalculateDifficulty()
		difficultyText := "Easy"
		var diffColor color.RGBA
		if difficulty >= 80 {
			difficultyText = "Expert"
			diffColor = color.RGBA{R: 220, G: 50, B: 50, A: 255}
		} else if difficulty >= 60 {
			difficultyText = "Hard"
			diffColor = color.RGBA{R: 220, G: 150, B: 50, A: 255}
		} else if difficulty >= 40 {
			difficultyText = "Medium"
			diffColor = color.RGBA{R: 50, G: 200, B: 50, A: 255}
		} else {
			diffColor = color.RGBA{R: 100, G: 180, B: 255, A: 255}
		}

		// Draw difficulty badge
		diffBadge := ebiten.NewImage(60, 18)
		diffBadge.Fill(diffColor)
		badgeOp := &ebiten.DrawImageOptions{}
		badgeOp.GeoM.Translate(float64(listX+300), float64(listY+i*songSpacing-13))
		screen.DrawImage(diffBadge, badgeOp)

		text.Draw(screen, difficultyText, basicfont.Face7x13, listX+304, listY+i*songSpacing, color.Black)
	}

	// Draw controls section
	controlsY := 320
	if len(ts.songs) > 0 {
		controlsY = listY + len(ts.songs)*songSpacing + 30
	}

	// Draw controls background
	controlsBg := ebiten.NewImage(270, 80)
	controlsBg.Fill(color.RGBA{R: 30, G: 30, B: 60, A: 200})
	controlsBgOp := &ebiten.DrawImageOptions{}
	controlsBgOp.GeoM.Translate(float64(screenWidth/2-135), float64(controlsY-20))
	screen.DrawImage(controlsBg, controlsBgOp)

	// Draw controls title
	text.Draw(screen, "CONTROLS", basicfont.Face7x13, screenWidth/2-33, controlsY, color.RGBA{R: 220, G: 220, B: 220, A: 255})

	// Draw instructions with icons
	instructionsX := screenWidth/2 - 120
	text.Draw(screen, "↑↓       Change Song", basicfont.Face7x13, instructionsX, controlsY+25, color.White)
	text.Draw(screen, "SPACE    Start Game", basicfont.Face7x13, instructionsX, controlsY+45, color.White)
	text.Draw(screen, "QWEIOP   Hit Notes", basicfont.Face7x13, instructionsX, controlsY+65, color.White)

	// Draw selected song details at bottom with background
	if len(ts.songs) > 0 && ts.currentSong < len(ts.songs) {
		song := ts.songs[ts.currentSong]

		// Draw info background
		infoBg := ebiten.NewImage(screenWidth, 30)
		infoBg.Fill(color.RGBA{R: 20, G: 20, B: 40, A: 220})
		infoBgOp := &ebiten.DrawImageOptions{}
		infoBgOp.GeoM.Translate(0, float64(screenHeight-40))
		screen.DrawImage(infoBg, infoBgOp)

		// Draw song details
		infoY := screenHeight - 20
		infoText := fmt.Sprintf("BPM: %d   Notes: %d   Duration: %.1fs",
			song.BPM, len(song.Notes), song.Duration)
		text.Draw(screen, infoText, basicfont.Face7x13, screenWidth/2-100, infoY, color.White)
	}

	// Draw footer
	versionText := "v1.0.0"
	text.Draw(screen, versionText, basicfont.Face7x13, 10, screenHeight-10, color.RGBA{R: 150, G: 150, B: 150, A: 255})
}
