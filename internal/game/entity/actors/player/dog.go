package gameplayer

import (
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/builder"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/platformer"
	gameplayermethods "github.com/leandroatallah/firefly/internal/game/entity/actors/methods"
	gamestates "github.com/leandroatallah/firefly/internal/game/entity/actors/states"
)

type DogPlayer struct {
	*platformer.PlatformerCharacter

	*gameplayermethods.PlayerDeathBehavior
}

// NewDogPlayer creates a new dog player.
func NewDogPlayer(ctx *app.AppContext) (platformer.PlatformerActorEntity, error) {
	character, spriteData, statData, stateMap, err := builder.PreparePlatformer(ctx, "internal/game/entity/actors/player/dog.json")
	if err != nil {
		return nil, err
	}

	character.SetStateTransitionHandler(gameplayermethods.StandardStateTransitionLogic)

	player := &DogPlayer{
		PlatformerCharacter: character,
	}
	player.SetOwner(player)
	// Ensure the original character pointer (referenced by physics bodies) also points to the player
	character.SetOwner(player)

	if err = builder.ConfigureCharacter(player, spriteData, statData, stateMap, "player"); err != nil {
		return nil, err
	}

	if err = builder.ApplyPlatformerPhysics(player, player); err != nil {
		return nil, err
	}

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
