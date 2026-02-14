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
			actor.GetCharacter().ClearSkills()
			actor.GetCharacter().SetMovementState(
				movement.Follow,
				s.player,
			)
		}
	case "set_speed":
		var spd int
		if v, ok := behavior.Config["speed"]; ok {
			switch s := v.(type) {
			case float64:
				spd = int(s)
			case int:
				spd = s
			}
		}
		if spd > 0 {
			if movable, ok := b.(body.Movable); ok {
				movable.SetSpeed(spd)
				movable.SetMaxSpeed(spd)
			}
		}
	case "delay":
		return
	}
}
