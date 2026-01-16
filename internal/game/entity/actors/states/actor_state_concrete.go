package gamestates

import "github.com/leandroatallah/firefly/internal/engine/entity/actors"

var Carrying actors.ActorStateEnum

func init() {
	Carrying = actors.RegisterState("carry", func(b actors.BaseState) actors.ActorState { return &CarryState{BaseState: b} })
}

// Carrying
type CarryState struct {
	actors.BaseState
}

func (s *CarryState) OnStart(currentCount int) {
	s.BaseState.OnStart(currentCount)
}
