package gamestates

import "github.com/leandroatallah/firefly/internal/engine/entity/actors"

// CarryingIdle
type CarryingIdleState struct {
	actors.BaseState
}

func (s *CarryingIdleState) OnStart(currentCount int) {
	s.BaseState.OnStart(currentCount)
}

// CarryingWalking
type CarryingWalkingState struct {
	actors.BaseState
}

func (s *CarryingWalkingState) OnStart(currentCount int) {
	s.BaseState.OnStart(currentCount)
}

// CarryingFalling
type CarryingFallingState struct {
	actors.BaseState
}

func (s *CarryingFallingState) OnStart(currentCount int) {
	s.BaseState.OnStart(currentCount)
}

var (
	CarryingIdle    actors.ActorStateEnum
	CarryingWalking actors.ActorStateEnum
	CarryingFalling actors.ActorStateEnum
)

func init() {
	CarryingIdle = actors.RegisterState("carry_idle", func(b actors.BaseState) actors.ActorState { return &CarryingIdleState{BaseState: b} })
	CarryingWalking = actors.RegisterState("carry_walking", func(b actors.BaseState) actors.ActorState { return &CarryingWalkingState{BaseState: b} })
	CarryingFalling = actors.RegisterState("carry_falling", func(b actors.BaseState) actors.ActorState { return &CarryingFallingState{BaseState: b} })
}
