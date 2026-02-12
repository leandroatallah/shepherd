package gamespeech

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	"github.com/leandroatallah/firefly/internal/engine/ui/speech"
)

type StorySpeech struct {
	*baseSpeech
}

func NewStorySpeech(fontSource *speech.SpeechFont) *StorySpeech {
	// Create indicator image (a simple white square)
	indicatorImg := ebiten.NewImage(8, 8)
	indicatorImg.Fill(color.White)

	s := &StorySpeech{
		baseSpeech: newBaseSpeech(fontSource),
	}
	s.indicator = indicatorImg
	s.SetID("story")
	s.SetColor(color.White) // Default to white for storytelling
	return s
}

func (s *StorySpeech) Show() {
	s.baseSpeech.Show()
	s.SpeechBase.SetSpellingDelay(0)
}

func (s *StorySpeech) Draw(screen *ebiten.Image, msg string) {
	if !s.Visible() && s.removed {
		return
	}

	// For storytelling, we use more space and no bubble
	w := config.Get().ScreenWidth - minMargin*2
	h := 64
	x := float64(minMargin)
	var y float64

	if s.GetPosition() == "top" {
		y = float64(minMargin)
	} else {
		// Default to bottom
		y = float64(config.Get().ScreenHeight) - float64(h) - float64(minMargin)
	}

	s.DrawText(screen, msg, x, y, w, h)

	s.DrawIndicator(screen, x, y, w, h)
}
