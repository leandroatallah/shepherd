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

	actor := s.GetActor()

	// Try to find the root owner (e.g. ShepherdPlayer) if the actor is just a component (e.g. Character)
	var root interface{} = actor
	if lastOwner := actor.LastOwner(); lastOwner != nil {
		root = lastOwner
	}

	if p, ok := root.(gameentitytypes.PlatformerActorEntity); ok {
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

// CarryingFalling
type CarryingFallingState struct {
	actors.BaseState
}

func (s *CarryingFallingState) OnStart(currentCount int) {
	s.BaseState.OnStart(currentCount)
}

var (
	Dying           actors.ActorStateEnum
	CarryingIdle    actors.ActorStateEnum
	CarryingWalking actors.ActorStateEnum
	CarryingFalling actors.ActorStateEnum
)

func init() {
	Dying = actors.RegisterState("die", func(b actors.BaseState) actors.ActorState { return &DyingState{BaseState: b} })
	CarryingIdle = actors.RegisterState("carry_idle", func(b actors.BaseState) actors.ActorState { return &CarryingIdleState{BaseState: b} })
	CarryingWalking = actors.RegisterState("carry_walking", func(b actors.BaseState) actors.ActorState { return &CarryingWalkingState{BaseState: b} })
	CarryingFalling = actors.RegisterState("carry_falling", func(b actors.BaseState) actors.ActorState { return &CarryingFallingState{BaseState: b} })
}
