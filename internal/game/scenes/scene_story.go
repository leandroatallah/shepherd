package gamescene

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/app"
	sequencestype "github.com/leandroatallah/firefly/internal/engine/contracts/sequences"
	"github.com/leandroatallah/firefly/internal/engine/scene"
	"github.com/leandroatallah/firefly/internal/engine/scene/transition"
	"github.com/leandroatallah/firefly/internal/engine/sequences"
	scenestypes "github.com/leandroatallah/firefly/internal/game/scenes/types"
)

type StoryScene struct {
	scene.BaseScene

	sequencePlayer    sequencestype.Player
	isPlayingSequence bool
	isRedirecting     bool
	shouldRedirect    bool

	count int
}

func NewStoryScene(ctx *app.AppContext) *StoryScene {
	scene := StoryScene{
		sequencePlayer: sequences.NewSequencePlayer(ctx),
	}
	scene.SetAppContext(ctx)
	return &scene
}

func (s *StoryScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	s.AppContext().DialogueManager.Draw(screen)
}

func (s *StoryScene) Update() error {
	s.count++

	if err := s.BaseScene.Update(); err != nil {
		return err
	}

	if s.shouldRedirect {
		s.NextScene()
	}

	if s.isRedirecting {
		return nil
	}

	if s.isPlayingSequence && s.sequencePlayer.IsOver() {
		s.isPlayingSequence = false
		s.shouldRedirect = true
	}

	s.sequencePlayer.Update()

	return nil
}

func (s *StoryScene) NextScene() {
	s.isRedirecting = true
	s.shouldRedirect = false
	s.AppContext().SceneManager.NavigateTo(scenestypes.ScenePhases, transition.NewFader(), true)
}

func (s *StoryScene) OnStart() {
	// TODO: Use a constant
	s.AppContext().DialogueManager.SetActiveSpeech("story")

	sequence, err := sequences.NewSequenceFromJSON("assets/sequences/story-1.json")
	if err != nil {
		log.Fatal(err)
	}
	s.sequencePlayer.Play(sequence)
	s.isPlayingSequence = true
}

func (s *StoryScene) OnFinish() {}
