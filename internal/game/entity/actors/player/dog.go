package gameplayer

import (
	"fmt"

	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/data/jsonutil"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/builder"
	physicsmovement "github.com/leandroatallah/firefly/internal/engine/physics/movement"
	gameplayermethods "github.com/leandroatallah/firefly/internal/game/entity/actors/methods"
	gamestates "github.com/leandroatallah/firefly/internal/game/entity/actors/states"
	gameentitytypes "github.com/leandroatallah/firefly/internal/game/entity/types"
)

type DogPlayer struct {
	*gameentitytypes.PlatformerCharacter

	*gameplayermethods.PlayerDeathBehavior
}

func NewDogPlayer(ctx *app.AppContext) (gameentitytypes.PlatformerActorEntity, error) {
	spriteData, statData, err := jsonutil.ParseSpriteAndStats[actors.StatData]("internal/game/entity/actors/player/dog.json")
	if err != nil {
		return nil, err
	}

	stateMap, err := builder.BuildStateMap(spriteData)
	if err != nil {
		return nil, err
	}

	rect := builder.BodyRectFromSpriteData(spriteData)
	character := gameentitytypes.NewPlatformerCharacter(stateMap, spriteData, rect)
	character.SetAppContext(ctx)
	character.SetStateTransitionHandler(gameplayermethods.StandardStateTransitionLogic)

	player := &DogPlayer{
		PlatformerCharacter: character,
	}
	player.SetOwner(player)
	// Ensure the original character pointer (referenced by physics bodies) also points to the player
	character.SetOwner(player)

	player.SetID("player")

	// FIX: Long Parameter List
	if err = builder.ConfigureCharacter(player, spriteData, statData, stateMap, "PLAYER"); err != nil {
		return nil, err
	}
	model, err := physicsmovement.NewMovementModel(physicsmovement.Platform, player)
	if err != nil {
		return nil, fmt.Errorf("SetMovementModel: %w", err)
	}
	player.SetMovementModel(model)

	character.StateCollisionManager.RefreshCollisions()

	player.PlayerDeathBehavior = gameplayermethods.NewPlayerDeathBehavior(player)

	return player, nil
}

func (p *DogPlayer) GetCharacter() *actors.Character {
	return p.Character
}

func (p *DogPlayer) Hurt(damage int) {
	if p.State() == gamestates.Dying {
		return
	}
	state, err := p.NewState(gamestates.Dying)
	if err != nil {
		return
	}
	p.SetState(state)
}
