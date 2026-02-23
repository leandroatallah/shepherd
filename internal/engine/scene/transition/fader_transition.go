package transition

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	"github.com/leandroatallah/firefly/internal/engine/utils/timing"
)

type Fader struct {
	BaseTransition
	alpha           float64
	waitFrames      int
	holdDuration    time.Duration // wait with black screen BEFORE callback
	fadeSpeed       float64
	visibleDuration time.Duration // wait with black screen AFTER callback (before fade in)
}

func NewFader(holdDuration, visibleDuration time.Duration) *Fader {
	return &Fader{
		holdDuration:    holdDuration,
		visibleDuration: visibleDuration,
		fadeSpeed:       15, // default: ~17 frames for full fade (255/15)
	}
}

// Transition methods
func (f *Fader) Update() {
	if !f.active {
		return
	}

	if f.exiting {
		f.alpha += f.fadeSpeed
		if f.alpha >= 255 {
			f.alpha = 255
			f.exiting = false
			// If no hold duration, call callback immediately
			if f.holdDuration <= 0 {
				if f.onExitCb != nil {
					f.onExitCb()
				}
				// If no visible wait, start fade in immediately
				if f.visibleDuration <= 0 {
					f.starting = true
					return
				}
				f.waitFrames = 0
				return
			}
			f.waitFrames = 0
		}
		return
	}

	// Hold phase: wait with black screen BEFORE callback
	if f.holdDuration > 0 && !f.starting {
		f.waitFrames++
		if timing.ToDuration(f.waitFrames) >= f.holdDuration {
			// Hold complete, call callback and start visible wait
			if f.onExitCb != nil {
				f.onExitCb()
			}
			if f.visibleDuration > 0 {
				f.waitFrames = 0
				return
			}
			// No visible wait, start fade in immediately
			f.starting = true
		}
		return
	}

	// Visible phase: wait with black screen AFTER callback (before fade in)
	if f.visibleDuration > 0 && !f.starting {
		f.waitFrames++
		if timing.ToDuration(f.waitFrames) >= f.visibleDuration {
			f.starting = true
		}
		return
	}

	if f.starting {
		f.alpha -= f.fadeSpeed
		if f.alpha <= 0 {
			f.alpha = 0
			f.starting = false
			f.active = false
		}
		return
	}
}

func (f *Fader) Draw(screen *ebiten.Image) {
	if !f.active {
		return
	}
	c := color.RGBA{A: uint8(f.alpha)}
	img := ebiten.NewImage(config.Get().ScreenWidth, config.Get().ScreenHeight)
	img.Fill(c)
	screen.DrawImage(img, nil)
}

func (f *Fader) StartTransition(cb func()) {
	f.fadeOut(func() {
		f.fadeIn(cb)
	})
}

func (f *Fader) EndTransition(cb func()) {}

// Custom methods
func (f *Fader) fadeOut(cb func()) {
	if f.active {
		return
	}
	f.active = true
	f.exiting = true
	f.alpha = 0
	f.onExitCb = cb
}

func (f *Fader) fadeIn(cb func()) {
	f.starting = true
	cb()
}
