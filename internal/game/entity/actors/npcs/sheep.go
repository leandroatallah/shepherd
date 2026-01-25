package gamenpcs

import (
	"log"

	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	physicsmovement "github.com/leandroatallah/firefly/internal/engine/physics/movement"
	gamemovement "github.com/leandroatallah/firefly/internal/game/entity/actors/movement"
	gamestates "github.com/leandroatallah/firefly/internal/game/entity/actors/states"
	gameentitytypes "github.com/leandroatallah/firefly/internal/game/entity/types"
	"github.com/leandroatallah/firefly/internal/game/events"
)

type Sheep struct {
	gameentitytypes.PlatformerCharacter
}

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
	sheep := &Sheep{PlatformerCharacter: *character}
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
	sheep.Character.SetMovementState(gamemovement.Wander, nil)

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
	return &s.Character
}

func (s *Sheep) OnTouch(other body.Collidable) {
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
	state, err := s.NewState(gamestates.Dying)
	if err != nil {
		return
	}
	s.SetState(state)
}

func (s *Sheep) OnDie() {
	s.SetHealth(0)
	// TODO: All actors need to freeze.
	s.SetImmobile(true)
	s.SetFreeze(true)

	// Trigger event to reboot scene
	if s.AppContext().EventManager != nil {
		s.AppContext().EventManager.Publish(&events.CharacterDiedEvent{})
	}
}
