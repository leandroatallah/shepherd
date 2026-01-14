package gamenpcs

import (
	"log"

	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/movement"
	physicsmovement "github.com/leandroatallah/firefly/internal/engine/physics/movement"
	gameentitytypes "github.com/leandroatallah/firefly/internal/game/entity/types"
)

type Sheep struct {
	gameentitytypes.PlatformerCharacter
}

func NewSheep(x, y int, id string) (*Sheep, error) {
	spriteData, statData, err := actors.ParseJsonPlayer("internal/game/entity/actors/npcs/sheep.json")
	if err != nil {
		log.Fatal(err)
	}

	character, err := CreateAnimatedCharacter(spriteData)
	if err != nil {
		log.Fatal(err)
	}

	character.SetPosition(x, y)
	sheep := &Sheep{PlatformerCharacter: *character}

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

	return sheep, nil
}

func (s *Sheep) SetTarget(target body.MovableCollidable) {
	s.Character.SetMovementState(movement.Idle, target)
}

// Character Methods
func (s *Sheep) Update(space body.BodiesSpace) error {
	return s.Character.Update(space)
}

func (s *Sheep) GetCharacter() *actors.Character {
	return &s.Character
}

func (s *Sheep) OnTouch(other body.Collidable) {
	sheepCarrier, ok := other.(gameentitytypes.SheepCarrier)
	_ = sheepCarrier
	if ok {
		// The sheep is "collected".
		// FIX: Implement sheep carrier logic
	}
}
