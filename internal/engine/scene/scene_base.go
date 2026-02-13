// - Add scene system (menu, playing, paused, game over)
// - Implement scene transitions and lifecycle management
package scene

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/physics/space"
	"github.com/leandroatallah/firefly/internal/engine/utils/timing"
)

type BaseScene struct {
	app.AppContextHolder

	count          int
	space          *space.Space
	IsKeysDisabled bool

	scheduledActions []scheduledAction
}

type scheduledAction struct {
	targetFrame uint64
	action      func()
}

func NewScene() *BaseScene {
	return &BaseScene{}
}

func (s *BaseScene) Draw(screen *ebiten.Image) {}

func (s *BaseScene) Update() error {
	// Update scheduled actions
	for i := 0; i < len(s.scheduledActions); i++ {
		if s.AppContext().FrameCount >= s.scheduledActions[i].targetFrame {
			s.scheduledActions[i].action()
			s.scheduledActions = append(s.scheduledActions[:i], s.scheduledActions[i+1:]...)
			i--
		}
	}
	return nil
}

func (s *BaseScene) Schedule(delay time.Duration, action func()) {
	target := s.AppContext().FrameCount + uint64(timing.FromDuration(delay))
	s.scheduledActions = append(s.scheduledActions, scheduledAction{
		targetFrame: target,
		action:      action,
	})
}

func (s *BaseScene) OnStart() {
	s.AppContext().Space.Clear()
	if s.AppContext().ActorManager != nil {
		s.AppContext().ActorManager.Clear()
	}
}

func (s *BaseScene) OnFinish() {}

func (s *BaseScene) Exit() {}

func (s *BaseScene) AddBoundaries(boundaries ...body.MovableCollidable) {
	space := s.PhysicsSpace()
	for _, o := range boundaries {
		space.AddBody(o)
	}
}

func (s *BaseScene) PhysicsSpace() body.BodiesSpace {
	return s.AppContext().Space
}

func (s *BaseScene) EnableKeys() {
	s.IsKeysDisabled = false
}

func (s *BaseScene) DisableKeys() {
	s.IsKeysDisabled = true
}
