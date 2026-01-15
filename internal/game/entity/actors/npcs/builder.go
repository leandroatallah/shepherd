package gamenpcs

import (
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/contracts/animation"
	"github.com/leandroatallah/firefly/internal/engine/data/schemas"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	"github.com/leandroatallah/firefly/internal/game/entity/actors/builder"
	gameentitytypes "github.com/leandroatallah/firefly/internal/game/entity/types"
)

func CreateAnimatedCharacter(ctx *app.AppContext, data schemas.SpriteData) (*gameentitytypes.PlatformerCharacter, error) {
	stateMap := map[string]animation.SpriteState{
		"idle": actors.Idle,
		"walk": actors.Walking,
	}
	return builder.CreateAnimatedCharacter(ctx, data, stateMap)
}

func SetNpcBodies(npc gameentitytypes.PlatformerActorEntity, data schemas.SpriteData, id string) error {
	npc.SetID(id)

	stateMap := map[string]animation.SpriteState{
		"idle": actors.Idle,
		"walk": actors.Walking,
	}

	return builder.SetCharacterBodies(npc, data, stateMap, "NPC")
}

func SetNpcStats(npc gameentitytypes.PlatformerActorEntity, data actors.StatData) error {
	return builder.SetCharacterStats(npc, data)
}
