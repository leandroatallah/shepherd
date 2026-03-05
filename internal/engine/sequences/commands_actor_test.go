package sequences

import (
	"testing"

	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/movement"
	"github.com/leandroatallah/firefly/internal/engine/mocks"
	"github.com/leandroatallah/firefly/internal/engine/physics/space"
)

const (
	Input movement.MovementStateEnum = iota
	Idle
	Rand
	Chase
	DumbChase
	Patrol
	Avoid
	SideToSide
	Follow
)

func setupTestContext() (*app.AppContext, *actors.Manager) {
	appContext := &app.AppContext{}
	actorManager := actors.NewManager()
	appContext.ActorManager = actorManager
	appContext.Space = space.NewSpace()
	return appContext, actorManager
}

func TestMoveActorCommand(t *testing.T) {
	ctx, am := setupTestContext()
	actor := &mocks.MockActor{Id: "test_actor", SpeedVal: 5, MovementMdl: &mocks.MockMovementModel{}}
	actor.SetPosition(0, 0)
	am.Register(actor)

	cmd := &MoveActorCommand{
		TargetID: "test_actor",
		EndX:     100,
		Speed:    10,
	}

	cmd.Init(ctx)

	if !actor.MovementMdl.(*mocks.MockMovementModel).IsScriptedVal {
		t.Error("expected actor to be set to scripted mode")
	}

	// Move right
	finished := cmd.Update()
	if finished {
		t.Error("command should not be finished yet")
	}
	if actor.MoveRightForce != 10 {
		t.Errorf("expected move right force 10, got %d", actor.MoveRightForce)
	}

	// Reach destination
	actor.SetPosition(100, 0)
	finished = cmd.Update()
	if !finished {
		t.Error("command should be finished when destination reached")
	}
	if actor.MovementMdl.(*mocks.MockMovementModel).IsScriptedVal {
		t.Error("expected actor scripted mode to be disabled after completion")
	}
}

func TestSetSpeedCommand(t *testing.T) {
	ctx, am := setupTestContext()
	actor := &mocks.MockActor{Id: "test_actor", SpeedVal: 5, MaxSpeedVal: 5}
	am.Register(actor)

	cmd := &SetSpeedCommand{
		TargetID: "test_actor",
		Speed:    20,
	}

	cmd.Init(ctx)

	if actor.SpeedVal != 20 || actor.MaxSpeedVal != 20 {
		t.Errorf("expected speed 20, got speed %d, maxSpeed %d", actor.SpeedVal, actor.MaxSpeedVal)
	}
}

func TestFollowCommands(t *testing.T) {
	ctx, am := setupTestContext()
	player := &mocks.MockActor{Id: "player"}
	am.Register(player)

	npc := &mocks.MockActor{Id: "npc"}
	am.Register(npc)

	// Test Follow
	followCmd := &FollowPlayerCommand{TargetID: "npc"}
	followCmd.Init(ctx)
	if npc.MovementSt != Follow {
		t.Errorf("expected npc to be in Follow state, got %v", npc.MovementSt)
	}

	// Test Stop Following
	stopCmd := &StopFollowingCommand{TargetID: "npc"}
	stopCmd.Init(ctx)
	if npc.MovementSt != Idle {
		t.Errorf("expected npc to be in Idle state, got %v", npc.MovementSt)
	}
}

func TestRemoveActorCommand(t *testing.T) {
	ctx, am := setupTestContext()
	actor := &mocks.MockActor{Id: "to_remove"}
	actor.SetPosition(0, 0)
	am.Register(actor)
	ctx.Space.AddBody(actor)

	cmd := &RemoveActorCommand{TargetID: "to_remove"}
	cmd.Init(ctx)

	if _, found := am.Find("to_remove"); found {
		t.Error("actor should be unregistered from ActorManager")
	}

	// Check if queued for removal in space
	ctx.Space.ProcessRemovals()
	if ctx.Space.Find("to_remove") != nil {
		t.Error("actor should be removed from Space")
	}
}
