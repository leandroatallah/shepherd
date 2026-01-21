package gamescenephases

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func (s *PhasesScene) CamDebug() {
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		s.Camera().Kamera().Angle += 0.02
	}
	if ebiten.IsKeyPressed(ebiten.KeyF) {
		s.Camera().Kamera().Angle -= 0.02
	}

	if ebiten.IsKeyPressed(ebiten.KeyBackspace) {
		s.Camera().Kamera().Reset()
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) { // zoom out
		s.Camera().Kamera().ZoomFactor /= 1.02
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) { // zoom in
		s.Camera().Kamera().ZoomFactor *= 1.02
	}

	if ebiten.IsKeyPressed(ebiten.KeyJ) {
		s.Camera().Kamera().CenterOffsetX *= 1.02
	}
	if ebiten.IsKeyPressed(ebiten.KeyL) {
		s.Camera().Kamera().CenterOffsetX /= 1.02
	}
	if ebiten.IsKeyPressed(ebiten.KeyK) {
		s.Camera().Kamera().CenterOffsetY /= 1.02
	}
	if ebiten.IsKeyPressed(ebiten.KeyI) {
		s.Camera().Kamera().CenterOffsetY *= 1.02
	}
}
