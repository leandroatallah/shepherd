package gameenemies

import (
	"log"

	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/data/jsonutil"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/builder"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/movement"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/platformer"
	physicsmovement "github.com/leandroatallah/firefly/internal/engine/physics/movement"
	gamenpcs "github.com/leandroatallah/firefly/internal/game/entity/actors/npcs"
	gameplayer "github.com/leandroatallah/firefly/internal/game/entity/actors/player"
	gamestates "github.com/leandroatallah/firefly/internal/game/entity/actors/states"
)

type BatEnemy struct {
	*platformer.PlatformerCharacter
}

// TODO: Use composition to reduce repeated actions in different places
func NewBatEnemy(ctx *app.AppContext, x, y int, id string) (*BatEnemy, error) {
	spriteData, statData, err := jsonutil.ParseSpriteAndStats[actors.StatData]("internal/game/entity/actors/enemies/bat.json")
	if err != nil {
		log.Fatal(err)
	}

	stateMap, err := builder.BuildStateMap(spriteData)
	if err != nil {
		return nil, err
	}

	rect := builder.BodyRectFromSpriteData(spriteData)
	character := platformer.NewPlatformerCharacter(stateMap, spriteData, rect)
	character.SetAppContext(ctx)
	character.SetPosition(x, y)

	enemy := &BatEnemy{PlatformerCharacter: character}
	// Set the owner on the embedded character so LastOwner() works correctly
	enemy.SetOwner(enemy)

	if err = builder.ConfigureCharacter(enemy, spriteData, statData, stateMap, "ENEMY"); err != nil {
		return nil, err
	}

	model, err := physicsmovement.NewMovementModel(physicsmovement.Platform, nil)
	if err != nil {
		return nil, err
	}
	enemy.SetMovementModel(model)
	enemy.SetTouchable(enemy)
	enemy.SetGravityEnabled(false)
	enemy.Character.SetMovementState(movement.SideToSide, nil, movement.WithIgnoreLedges(true), movement.WithWaitBeforeTurn(60))

	return enemy, nil
}

func (e *BatEnemy) SetTarget(target body.MovableCollidable) {
	e.Character.MovementState().SetTarget(target)
}

// Character Methods
func (e *BatEnemy) Update(space body.BodiesSpace) error {
	return e.Character.Update(space)
}

func (e *BatEnemy) GetCharacter() *actors.Character {
	return e.Character
}

func (e *BatEnemy) OnTouch(other body.Collidable) {
	owner := other.LastOwner()
	switch owner.(type) {
	case *gameplayer.ShepherdPlayer, *gameplayer.DogPlayer, *gamenpcs.Sheep:
		if owner.(platformer.PlatformerActorEntity).State() == gamestates.Dying {
			return
		}

		if alive, ok := owner.(platformer.AlivePlayer); ok {
			alive.Hurt(1)
		}
	}
}

func (e *BatEnemy) OnDie() {}
