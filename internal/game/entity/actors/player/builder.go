package gameplayer

import (
	"fmt" // ADDED THIS
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/contracts/animation"
	"github.com/leandroatallah/firefly/internal/engine/data/schemas"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	physicsmovement "github.com/leandroatallah/firefly/internal/engine/physics/movement"
	"github.com/leandroatallah/firefly/internal/game/entity/actors/builder"
	gameentitytypes "github.com/leandroatallah/firefly/internal/game/entity/types"
)

func CreateAnimatedCharacter(ctx *app.AppContext, data schemas.SpriteData) (*gameentitytypes.PlatformerCharacter, error) {
	stateMap := make(map[string]animation.SpriteState)
	for _, stateName := range []string{"idle", "walk", "fall", "hurt", "carry"} {
		enum, ok := actors.GetStateEnum(stateName)
		if !ok {
			return nil, fmt.Errorf("state '%s' not registered", stateName)
		}
		stateMap[stateName] = enum
	}
	return builder.CreateAnimatedCharacter(ctx, data, stateMap)
}

// SetPlayerBodies
func SetPlayerBodies(player gameentitytypes.PlatformerActorEntity, data schemas.SpriteData) error {
	player.SetID("player")

	stateMap := make(map[string]animation.SpriteState)
	for _, stateName := range []string{"idle", "walk", "fall", "hurt", "carry"} {
		enum, ok := actors.GetStateEnum(stateName)
		if !ok {
			return fmt.Errorf("state '%s' not registered", stateName)
		}
		stateMap[stateName] = enum
	}

	return builder.SetCharacterBodies(player, data, stateMap, "PLAYER")
}

func SetPlayerStats(player gameentitytypes.PlatformerActorEntity, data actors.StatData) error {
	return builder.SetCharacterStats(player, data)
}

func SetMovementModel(
	player gameentitytypes.PlatformerActorEntity,
	movementModel physicsmovement.MovementModelEnum,
) error {
	model, err := physicsmovement.NewMovementModel(movementModel, player)
	if err != nil {
		return err
	}
	player.SetMovementModel(model)
	return nil
}
