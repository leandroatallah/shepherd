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

// shepherdStateTransitionLogic provides custom state handling for the ShepherdPlayer,
// specifically for managing the "carrying" states.
func shepherdStateTransitionLogic(c *actors.Character) bool {
	state := c.State()

	// NOTE: These carrying states (Carrying, CarryingIdle, CarryingWalking, CarryingFalling)
	// must be defined in the `gamestates` package. This logic also assumes they are implemented
	// and registered for the character, allowing for different animations and physics,
	// such as slower movement speed.

	isCarryingState := state == gamestates.CarryingIdle ||
		state == gamestates.CarryingWalking ||
		state == gamestates.CarryingFalling

	if !isCarryingState {
		return false // Let the engine handle other states
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

	// State machine for when the character is carrying something.
	switch {
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

	// Set the custom state transition logic for the player
	character.SetStateTransitionHandler(shepherdStateTransitionLogic)

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

func (p *ShepherdPlayer) Update(space body.BodiesSpace) error {
	return p.Character.Update(space)
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
	state, err := p.NewState(gamestates.CarryingIdle)
	if err != nil {
		return
	}
	p.SetState(state)
	// log.Println(p.AppContext().Space.Debug())
	p.AppContext().Space.QueueForRemoval(s)
}

func (p *ShepherdPlayer) IsCarryingSheep() bool {
	state := p.State()
	return state == gamestates.CarryingIdle ||
		state == gamestates.CarryingWalking ||
		state == gamestates.CarryingFalling
}

func (p *ShepherdPlayer) DropSheep() {
	log.Println("ShepherdPlayer: Dropping sheep")
	state, err := p.NewState(actors.Idle)
	if err != nil {
		return
	}
	p.SetState(state)
}
