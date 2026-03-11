package gameplayer

import (
	"log"

	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/builder"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/platformer"
	gameplayermethods "github.com/leandroatallah/firefly/internal/game/entity/actors/methods"
	gamestates "github.com/leandroatallah/firefly/internal/game/entity/actors/states"
	gameentitytypes "github.com/leandroatallah/firefly/internal/game/entity/types"
)

// shepherdStateTransitionLogic provides custom state handling for the ShepherdPlayer,
// specifically for managing the "carrying" states.
func shepherdStateTransitionLogic(c *actors.Character) bool {
	if gameplayermethods.StandardStateTransitionLogic(c) {
		return true
	}

	setNewState := func(s actors.ActorStateEnum) {
		state, err := c.NewState(s)
		if err != nil {
			// Log the error instead of crashing if a state is not registered.
			log.Printf("Failed to create new state %v: %v", s, err)
			return
		}
		c.SetState(state)
	}

	state := c.State()

	if state == gamestates.Rising && c.IsAnimationFinished() {
		setNewState(actors.Idle)
		return true
	}

	if state == gamestates.Exiting || state == gamestates.Lying || state == gamestates.Rising {
		return true
	}

	isCarryingState := state == gamestates.CarryingIdle ||
		state == gamestates.CarryingWalking ||
		state == gamestates.CarryingFalling ||
		state == gamestates.CarryingLanding ||
		state == gamestates.CarryingJump

	if !isCarryingState {
		return false // Let the engine handle other states
	}

	// State machine for when the character is carrying something.
	switch {
	case state == gamestates.CarryingLanding:
		isAnimationOver := c.IsAnimationFinished()
		if c.IsWalking() {
			setNewState(gamestates.CarryingWalking)
		} else if isAnimationOver {
			setNewState(gamestates.CarryingIdle)
		}
	case state == gamestates.CarryingFalling && !c.IsFalling():
		setNewState(gamestates.CarryingLanding)
	case state != gamestates.CarryingFalling && c.IsFalling():
		setNewState(gamestates.CarryingFalling)
	case state != gamestates.CarryingWalking && c.IsWalking():
		setNewState(gamestates.CarryingWalking)
	case state != gamestates.CarryingIdle && c.IsIdle():
		// This case also handles the initial transition from the base "Carrying"
		// state to the more specific "CarryingIdle" state.
		setNewState(gamestates.CarryingIdle)
	}

	return true // We've handled the state, so the engine shouldn't.
}

type ShepherdPlayer struct {
	*platformer.PlatformerCharacter
	gameentitytypes.SheepCarrier
	baseSpeed int

	*gameplayermethods.PlayerDeathBehavior
}

// NewShepherdPlayer creates a new shepherd player.
func NewShepherdPlayer(ctx *app.AppContext) (platformer.PlatformerActorEntity, error) {
	character, spriteData, statData, stateMap, err := builder.PreparePlatformer(ctx, "internal/game/entity/actors/player/shepherd.json")
	if err != nil {
		return nil, err
	}

	character.SetStateTransitionHandler(shepherdStateTransitionLogic)

	player := &ShepherdPlayer{
		PlatformerCharacter: character,
	}
	// Set the owner on the embedded character so LastOwner() works correctly
	player.SetOwner(player)
	// Ensure the original character pointer (referenced by physics bodies) also points to the player
	character.SetOwner(player)

	if err = builder.ConfigureCharacter(player, spriteData, statData, stateMap, "player"); err != nil {
		return nil, err
	}
	player.baseSpeed = player.Speed()

	if err = builder.ApplyPlatformerPhysics(player, player); err != nil {
		return nil, err
	}

	character.StateCollisionManager.RefreshCollisions()
	player.PlayerDeathBehavior = gameplayermethods.NewPlayerDeathBehavior(player)

	return player, nil
}

func (p *ShepherdPlayer) Update(space body.BodiesSpace) error {
	if p.IsCarryingSheep() {
		p.SetHorizontalInertia(1.0)
		p.SetSpeed(int(float64(p.baseSpeed) * 0.5))
		p.SetJumpForceMultiplier(0.95)
	} else {
		p.SetHorizontalInertia(-1.0)
		p.SetSpeed(p.baseSpeed)
		p.SetJumpForceMultiplier(1.0)
	}
	return p.Character.Update(space)
}

func (p *ShepherdPlayer) GetCharacter() *actors.Character {
	return p.Character
}

func (p *ShepherdPlayer) Hurt(damage int) {
	if p.State() == gamestates.Dying {
		return
	}

	state, err := p.NewState(gamestates.Dying)
	if err != nil {
		return
	}
	p.SetState(state)
}

// SheepCarrier Methods
func (p *ShepherdPlayer) GrabSheep(s body.MovableCollidableTouchable) {
	state, err := p.NewState(gamestates.CarryingIdle)
	if err != nil {
		return
	}
	p.SetState(state)
	p.AppContext().Space.QueueForRemoval(s)
}

func (p *ShepherdPlayer) IsCarryingSheep() bool {
	state := p.State()
	return state == gamestates.CarryingIdle ||
		state == gamestates.CarryingWalking ||
		state == gamestates.CarryingFalling ||
		state == gamestates.CarryingLanding
}

func (p *ShepherdPlayer) DropSheep() {
	state, err := p.NewState(actors.Idle)
	if err != nil {
		return
	}
	p.SetState(state)
}
