package gameplayer

import (
	"fmt"

	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	physicsmovement "github.com/leandroatallah/firefly/internal/engine/physics/movement"
	"github.com/leandroatallah/firefly/internal/engine/physics/skill"
	gameplayermethods "github.com/leandroatallah/firefly/internal/game/entity/actors/methods"
	gamestates "github.com/leandroatallah/firefly/internal/game/entity/actors/states"
	gameentitytypes "github.com/leandroatallah/firefly/internal/game/entity/types"
)

type DogPlayer struct {
	gameentitytypes.PlatformerCharacter

	*gameplayermethods.PlayerDeathBehavior
}

func NewDogPlayer(ctx *app.AppContext) (gameentitytypes.PlatformerActorEntity, error) {
	spriteData, statData, err := actors.ParseJsonPlayer("internal/game/entity/actors/player/dog.json")
	if err != nil {
		return nil, err
	}

	character, err := CreateAnimatedCharacter(ctx, spriteData)
	if err != nil {
		return nil, err
	}
	character.AddSkill(skill.NewJumpSkill())
	character.AddSkill(skill.NewHorizontalMovementSkill())

	player := &DogPlayer{
		PlatformerCharacter: *character,
	}
	player.SetOwner(player)

	if err = SetPlayerBodies(player, spriteData); err != nil {
		return nil, fmt.Errorf("SetPlayerBodies: %w", err)
	}
	if err = SetPlayerStats(player, statData); err != nil {
		return nil, fmt.Errorf("SetPlayerStats: %w", err)
	}
	// Pass player itself
	if err = SetMovementModel(player, physicsmovement.Platform); err != nil {
		return nil, fmt.Errorf("SetMovementModel: %w", err)
	}

	character.StateCollisionManager.RefreshCollisions()

	player.PlayerDeathBehavior = gameplayermethods.NewPlayerDeathBehavior(player)

	return player, nil
}

func (p *DogPlayer) GetCharacter() *actors.Character {
	return &p.Character
}

func (p *DogPlayer) Hurt(damage int) {
	state, err := p.NewState(gamestates.Dying)
	if err != nil {
		return
	}
	p.SetState(state)
}
