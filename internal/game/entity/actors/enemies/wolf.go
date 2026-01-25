package gameenemies

import (
	"log"

	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/enemies"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/movement"
	physicsmovement "github.com/leandroatallah/firefly/internal/engine/physics/movement"
	gameentitytypes "github.com/leandroatallah/firefly/internal/game/entity/types"
)

type WolfEnemy struct {
	gameentitytypes.PlatformerCharacter
	count int
}

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
	enemy := &WolfEnemy{PlatformerCharacter: *character}

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

	return enemy, nil
}

func (e *WolfEnemy) SetTarget(target body.MovableCollidable) {
	e.Character.SetMovementState(movement.SideToSide, target, movement.WithWaitBeforeTurn(60))
}

// Character Methods
func (e *WolfEnemy) Update(space body.BodiesSpace) error {
	e.count++
	return e.Character.Update(space)
}

func (e *WolfEnemy) GetCharacter() *actors.Character {
	return &e.Character
}

func (e *WolfEnemy) OnTouch(other body.Collidable) {
	player := e.MovementState().Target()
	if other.ID() == player.ID() {
		player.(gameentitytypes.PlatformerActorEntity).Hurt(1)
	}
}

func (e *WolfEnemy) OnDie() {}
