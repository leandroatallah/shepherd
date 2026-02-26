package gamenpcs

import (
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/builder"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/movement"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/platformer"
	gameplayermethods "github.com/leandroatallah/firefly/internal/game/entity/actors/methods"
	gamestates "github.com/leandroatallah/firefly/internal/game/entity/actors/states"
	gameentitytypes "github.com/leandroatallah/firefly/internal/game/entity/types"
)

type Sheep struct {
	*platformer.PlatformerCharacter
	*gameplayermethods.PlayerDeathBehavior
}

// NewSheep creates a new sheep NPC.
func NewSheep(ctx *app.AppContext, x, y int, id string) (*Sheep, error) {
	character, spriteData, statData, stateMap, err := builder.PreparePlatformer(ctx, "internal/game/entity/actors/npcs/sheep.json")
	if err != nil {
		return nil, err
	}

	sheep := &Sheep{PlatformerCharacter: character}
	// Set the owner on the embedded character so LastOwner() works correctly
	sheep.SetOwner(sheep)
	sheep.SetPosition(x, y)

	if err = builder.ConfigureCharacter(sheep, spriteData, statData, stateMap, id); err != nil {
		return nil, err
	}

	if err = builder.ApplyPlatformerPhysics(sheep, nil); err != nil {
		return nil, err
	}

	sheep.Character.SetMovementState(movement.Idle, nil)
	sheep.Character.SetStateTransitionHandler(gameplayermethods.StandardStateTransitionLogic)
	sheep.PlayerDeathBehavior = gameplayermethods.NewPlayerDeathBehavior(sheep)

	return sheep, nil
}

func (s *Sheep) SetTarget(target body.MovableCollidable) {
	s.Character.SetMovementState(movement.Wander, target)
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
