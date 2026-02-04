package gameenemies

import (
	"log"

	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/enemies"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/movement"
	physicsmovement "github.com/leandroatallah/firefly/internal/engine/physics/movement"
	gamenpcs "github.com/leandroatallah/firefly/internal/game/entity/actors/npcs"
	gameplayer "github.com/leandroatallah/firefly/internal/game/entity/actors/player"
	gamestates "github.com/leandroatallah/firefly/internal/game/entity/actors/states"
	gameentitytypes "github.com/leandroatallah/firefly/internal/game/entity/types"
)

type WolfEnemy struct {
	*gameentitytypes.PlatformerCharacter
}

// TODO: Use composition to reduce repeated actions in different places
func NewWolfEnemy(ctx *app.AppContext, x, y int, id string) (*WolfEnemy, error) {
	spriteData, statData, err := enemies.ParseJsonEnemy("internal/game/entity/actors/enemies/wolf.json")
	if err != nil {
		log.Fatal(err)
	}

	character, err := CreateAnimatedCharacter(ctx, spriteData)
	if err != nil {
		log.Fatal(err)
	}

	character.SetPosition(x, y)
	enemy := &WolfEnemy{PlatformerCharacter: character}
	// Set the owner on the embedded character so LastOwner() works correctly
	enemy.SetOwner(enemy)

	if err = SetEnemyStats(enemy, statData); err != nil {
		return nil, err
	}
	if err = SetEnemyBodies(enemy, spriteData, id); err != nil {
		return nil, err
	}

	model, err := physicsmovement.NewMovementModel(physicsmovement.Platform, nil)
	if err != nil {
		return nil, err
	}
	enemy.SetMovementModel(model)
	enemy.SetTouchable(enemy)
	enemy.Character.SetMovementState(movement.SideToSide, nil, movement.WithWaitBeforeTurn(60))

	return enemy, nil
}

func (e *WolfEnemy) SetTarget(target body.MovableCollidable) {
	e.Character.MovementState().SetTarget(target)
}

// Character Methods
func (e *WolfEnemy) Update(space body.BodiesSpace) error {
	return e.Character.Update(space)
}

func (e *WolfEnemy) GetCharacter() *actors.Character {
	return e.Character
}

func (e *WolfEnemy) OnTouch(other body.Collidable) {
	owner := other.LastOwner()
	switch owner.(type) {
	case *gameplayer.ShepherdPlayer, *gameplayer.DogPlayer, *gamenpcs.Sheep:
		if owner.(gameentitytypes.PlatformerActorEntity).State() == gamestates.Dying {
			return
		}

		if alive, ok := owner.(gameentitytypes.AlivePlayer); ok {
			alive.Hurt(1)
		}
	}
}

func (e *WolfEnemy) OnDie() {}
