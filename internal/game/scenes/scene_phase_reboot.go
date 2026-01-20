package gamescene

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	"github.com/leandroatallah/firefly/internal/engine/scene"
	"github.com/leandroatallah/firefly/internal/engine/scene/transition"
)

const NAVIGATE_BACK_DELAY = 10 // frames

type PhaseRebootScene struct {
	scene.BaseScene

	count      int
	redirected bool
}

func NewPhaseRebootScene(context *app.AppContext) *PhaseRebootScene {
	overlay := ebiten.NewImage(config.Get().ScreenWidth, config.Get().ScreenHeight)
	overlay.Fill(color.Black)
	scene := PhaseRebootScene{}
	scene.SetAppContext(context)
	return &scene
}

func (s *PhaseRebootScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{A: 255})
}

func (s *PhaseRebootScene) Update() error {
	s.count++

	if s.count > NAVIGATE_BACK_DELAY && !s.redirected {
		s.AppContext().SceneManager.NavigateBack(transition.NewFader())
		s.redirected = true
	}

	return nil
}

func (s *PhaseRebootScene) OnStart() {}

func (s *PhaseRebootScene) OnFinish() {}
