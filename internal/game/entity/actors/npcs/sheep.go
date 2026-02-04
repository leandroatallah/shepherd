package gamenpcs

import (
	"log"

	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/movement"
	physicsmovement "github.com/leandroatallah/firefly/internal/engine/physics/movement"
	gameplayermethods "github.com/leandroatallah/firefly/internal/game/entity/actors/methods"
	gamemovement "github.com/leandroatallah/firefly/internal/game/entity/actors/movement"
	gamestates "github.com/leandroatallah/firefly/internal/game/entity/actors/states"
	gameentitytypes "github.com/leandroatallah/firefly/internal/game/entity/types"
)

type Sheep struct {
	*gameentitytypes.PlatformerCharacter
	*gameplayermethods.PlayerDeathBehavior
}

// TODO: Use composition to reduce repeated actions in different places
func NewSheep(ctx *app.AppContext, x, y int, id string) (*Sheep, error) {
	spriteData, statData, err := actors.ParseJsonPlayer("internal/game/entity/actors/npcs/sheep.json")
	if err != nil {
		log.Fatal(err)
	}

	character, err := CreateAnimatedCharacter(ctx, spriteData)
	if err != nil {
		log.Fatal(err)
	}

	character.SetPosition(x, y)
	sheep := &Sheep{PlatformerCharacter: character}
	// Set the owner on the embedded character so LastOwner() works correctly
	sheep.SetOwner(sheep)

	if err = SetNpcStats(sheep, statData); err != nil {
		return nil, err
	}
	if err = SetNpcBodies(sheep, spriteData, id); err != nil {
		return nil, err
	}

	model, err := physicsmovement.NewMovementModel(physicsmovement.Platform, nil)
	if err != nil {
		return nil, err
	}
	sheep.SetMovementModel(model)
	sheep.SetTouchable(sheep)
	sheep.Character.SetMovementState(movement.Idle, nil)

	sheep.Character.SetStateTransitionHandler(gameplayermethods.StandardStateTransitionLogic)

	sheep.PlayerDeathBehavior = gameplayermethods.NewPlayerDeathBehavior(sheep)

	return sheep, nil
}

func (s *Sheep) SetTarget(target body.MovableCollidable) {
	s.Character.SetMovementState(gamemovement.Wander, target)
}

// Character Methods
func (s *Sheep) Update(space body.BodiesSpace) error {
	return s.Character.Update(space)
}

func (s *Sheep) GetCharacter() *actors.Character {
	return s.Character
}

func (s *Sheep) OnTouch(other body.Collidable) {
	if s.State() == gamestates.Dying {
		return
	}

	player, found := s.AppContext().ActorManager.GetPlayer()
	if !found {
		return
	}

	if other.ID() != player.ID() {
		return
	}

	sheepCarrier, ok := player.(gameentitytypes.SheepCarrier)
	if ok && !sheepCarrier.IsCarryingSheep() {
		sheepCarrier.GrabSheep(s)
	}
}

func (s *Sheep) Hurt(damage int) {
	if s.State() == gamestates.Dying {
		return
	}
	state, err := s.NewState(gamestates.Dying)
	if err != nil {
		return
	}
	s.SetState(state)
}
