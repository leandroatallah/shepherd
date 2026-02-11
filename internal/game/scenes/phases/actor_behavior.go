package gamescenephases

import (
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/movement"
	"github.com/leandroatallah/firefly/internal/engine/scene/phases"
	gameentitytypes "github.com/leandroatallah/firefly/internal/game/entity/types"
)

// ApplyActorBehavior applies a defined behavior to an actor
func ApplyActorBehavior(s *PhasesScene, b body.Body, behavior phases.ActorBehavior) {
	switch behavior.Type {
	case "follow_player":
		if actor, ok := b.(gameentitytypes.PlatformerActorEntity); ok {
			// Set Dog player movement to follow the player
			actor.GetCharacter().ClearSkills()
			actor.GetCharacter().SetMovementState(
				movement.Follow,
				s.player,
			)
		}
	}
}
