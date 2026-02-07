package gamestates

import (
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	gameentitytypes "github.com/leandroatallah/firefly/internal/game/entity/types"
)

// Dying
type DyingState struct {
	actors.BaseState
}

func (s *DyingState) OnStart(currentCount int) {
	s.BaseState.OnStart(currentCount)

	if p, ok := s.GetRootOwner().(gameentitytypes.PlatformerActorEntity); ok {
		p.OnDie()
	}
}

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

// CarryingJump
type CarryingJumpState struct {
	actors.BaseState
}

func (s *CarryingJumpState) OnStart(currentCount int) {
	s.BaseState.OnStart(currentCount)

	if j, ok := s.GetRootOwner().(actors.Jumpable); ok {
		j.OnJump()
	}
}

// CarryingFalling
type CarryingFallingState struct {
	actors.BaseState
}

func (s *CarryingFallingState) OnStart(currentCount int) {
	s.BaseState.OnStart(currentCount)

	if f, ok := s.GetRootOwner().(actors.Fallable); ok {
		f.OnFall()
	}
}

// CarryingLanding
type CarryingLandingState struct {
	actors.BaseState
}

func (s *CarryingLandingState) OnStart(currentCount int) {
	s.BaseState.OnStart(currentCount)

	if l, ok := s.GetRootOwner().(actors.Landable); ok {
		l.OnLand()
	}
}

var (
	Dying           actors.ActorStateEnum
	CarryingIdle    actors.ActorStateEnum
	CarryingWalking actors.ActorStateEnum
	CarryingJump    actors.ActorStateEnum
	CarryingFalling actors.ActorStateEnum
	CarryingLanding actors.ActorStateEnum
)

func init() {
	Dying = actors.RegisterState("die", func(b actors.BaseState) actors.ActorState { return &DyingState{BaseState: b} })
	CarryingIdle = actors.RegisterState("carry_idle", func(b actors.BaseState) actors.ActorState { return &CarryingIdleState{BaseState: b} })
	CarryingWalking = actors.RegisterState("carry_walking", func(b actors.BaseState) actors.ActorState { return &CarryingWalkingState{BaseState: b} })
	CarryingJump = actors.RegisterState("carry_jump", func(b actors.BaseState) actors.ActorState { return &CarryingJumpState{BaseState: b} })
	CarryingFalling = actors.RegisterState("carry_falling", func(b actors.BaseState) actors.ActorState { return &CarryingFallingState{BaseState: b} })
	CarryingLanding = actors.RegisterState("carry_landing", func(b actors.BaseState) actors.ActorState { return &CarryingLandingState{BaseState: b} })
}