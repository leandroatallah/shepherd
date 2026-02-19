package gameenemies

import (
	"log"

	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/data/jsonutil"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/builder"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/movement"
	physicsmovement "github.com/leandroatallah/firefly/internal/engine/physics/movement"
	gamenpcs "github.com/leandroatallah/firefly/internal/game/entity/actors/npcs"
	gameplayer "github.com/leandroatallah/firefly/internal/game/entity/actors/player"
	gamestates "github.com/leandroatallah/firefly/internal/game/entity/actors/states"
	gameentitytypes "github.com/leandroatallah/firefly/internal/game/entity/types"
)

type SwarmEnemy struct {
	*gameentitytypes.PlatformerCharacter
}

// TODO: Use composition to reduce repeated actions in different places
func NewSwarmEnemy(ctx *app.AppContext, x, y int, id string) (*SwarmEnemy, error) {
	spriteData, statData, err := jsonutil.ParseSpriteAndStats[actors.StatData]("internal/game/entity/actors/enemies/swarm.json")
	if err != nil {
		log.Fatal(err)
	}

	stateMap, err := builder.BuildStateMap(spriteData)
	if err != nil {
		return nil, err
	}

	rect := builder.BodyRectFromSpriteData(spriteData)
	character := gameentitytypes.NewPlatformerCharacter(stateMap, spriteData, rect)
	character.SetAppContext(ctx)
	character.SetPosition(x, y)

	enemy := &SwarmEnemy{PlatformerCharacter: character}
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
	enemy.SetHorizontalInertia(1.0)
	enemy.Character.SetMovementState(
		movement.SideToSide,
		nil,
		movement.WithIgnoreLedges(true),
		movement.WithWaitBeforeTurn(120),
		movement.WithVerticalMovement(true),
	)

	enemy.Character.SetStateTransitionHandler(func(c *actors.Character) bool {
		// Override the default platformer state transition logic.
		// Since the SwarmEnemy is a flying unit and lacks specific "Jumping" or "Falling" animations,
		// we force the state to remain either Idle or Walking based on horizontal movement.
		// This prevents the sprite from disappearing or flickering when the enemy moves vertically
		// (which would normally trigger a Jump/Fall state in the default platformer model).
		desiredState := actors.Idle
		if vx, _ := c.Velocity(); vx != 0 {
			desiredState = actors.Walking
		}

		if c.State() != desiredState {
			s, err := c.NewState(desiredState)
			if err == nil {
				c.SetState(s)
			}
		}
		return true
	})

	return enemy, nil
}

func (e *SwarmEnemy) SetTarget(target body.MovableCollidable) {
	e.Character.MovementState().SetTarget(target)
}

// Character Methods
func (e *SwarmEnemy) Update(space body.BodiesSpace) error {
	return e.Character.Update(space)
}

func (e *SwarmEnemy) GetCharacter() *actors.Character {
	return e.Character
}

func (e *SwarmEnemy) OnTouch(other body.Collidable) {
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

func (e *SwarmEnemy) OnDie() {}
