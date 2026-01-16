package gameplayer

import (
	"fmt"
	"log"

	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body" // ADDED THIS
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	physicsmovement "github.com/leandroatallah/firefly/internal/engine/physics/movement"
	"github.com/leandroatallah/firefly/internal/engine/physics/skill"
	gamestates "github.com/leandroatallah/firefly/internal/game/entity/actors/states"
	gameentitytypes "github.com/leandroatallah/firefly/internal/game/entity/types"
)

type ShepherdPlayer struct {
	gameentitytypes.PlatformerCharacter
	gameentitytypes.SheepCarrier
}

func NewShepherdPlayer(ctx *app.AppContext) (gameentitytypes.PlatformerActorEntity, error) {
	spriteData, statData, err := actors.ParseJsonPlayer("internal/game/entity/actors/player/shepherd.json")
	if err != nil {
		return nil, err
	}

	character, err := CreateAnimatedCharacter(ctx, spriteData)
	if err != nil {
		return nil, err
	}
	character.AddSkill(skill.NewJumpSkill())
	character.AddSkill(skill.NewHorizontalMovementSkill())

	player := &ShepherdPlayer{
		PlatformerCharacter: *character,
	}

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

	return player, nil
}

func (p *ShepherdPlayer) GetCharacter() *actors.Character {
	return &p.Character
}

func (p *ShepherdPlayer) Hurt(damage int) {
	state, err := p.NewState(actors.Hurted)
	if err != nil {
		return
	}
	p.SetState(state)
}

func (p *ShepherdPlayer) GrabSheep(s body.MovableCollidableTouchable) {
	log.Printf("Shepherd: %v - grabs: %v", p.ID(), s.ID())
	state, err := p.NewState(gamestates.Carrying)
	if err != nil {
		return
	}
	p.SetState(state)
}
