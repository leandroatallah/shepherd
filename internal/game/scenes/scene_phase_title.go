package gamescene

import (
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/assets/font"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	"github.com/leandroatallah/firefly/internal/engine/render/screenutil"
	"github.com/leandroatallah/firefly/internal/engine/scene"
	"github.com/leandroatallah/firefly/internal/engine/scene/transition"
)

const (
	titleMusic = "assets/audio/Goblins_Den_Regular.ogg"
)

type PhaseTitleScene struct {
	scene.BaseScene

	fontText     *font.FontText
	title        string
	showTitle    bool
	musicStarted bool
}

func NewPhaseTitleScene(ctx *app.AppContext) *PhaseTitleScene {
	fontText, err := font.NewFontText(config.Get().SmallFontFace)
	if err != nil {
		log.Fatal(err)
	}
	s := &PhaseTitleScene{
		fontText: fontText,
	}
	s.SetAppContext(ctx)
	return s
}

func (s *PhaseTitleScene) OnStart() {
	phase, err := s.AppContext().PhaseManager.GetCurrentPhase()
	if err != nil {
		log.Printf("PhaseTitleScene: failed to get current phase: %v", err)
	}
	s.title = phase.Title

	if am := s.AppContext().AudioManager; am != nil {
		am.FadeOutAll(500 * time.Millisecond)
	}

	s.Schedule(550*time.Millisecond, func() {
		s.showTitle = true
		if am := s.AppContext().AudioManager; am != nil && !s.musicStarted {
			am.SetVolume(1.0)
			am.PlayMusic(titleMusic)
			s.musicStarted = true
		}
	})
}

func (s *PhaseTitleScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	if s.showTitle {
		text := s.title
		if text == "" {
			text = "Phase"
		}
		screenutil.DrawCenteredText(screen, s.fontText, text, 16, color.White)
	}
}

func (s *PhaseTitleScene) Update() error {
	if err := s.BaseScene.Update(); err != nil {
		return err
	}
	if s.showTitle && inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		s.AppContext().CompleteCurrentPhase(transition.NewFader(0, config.Get().FadeVisibleDuration), true)
	}
	return nil
}

func (s *PhaseTitleScene) OnFinish() {}
