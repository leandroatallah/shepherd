package sequences

import (
	"image"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/contracts/animation"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/movement"
	physicsmovement "github.com/leandroatallah/firefly/internal/engine/physics/movement"
	"github.com/leandroatallah/firefly/internal/engine/physics/space"
)

// MockActor implements actors.ActorEntity for testing
type mockActor struct {
	id             string
	pos            image.Rectangle
	speed          int
	maxSpeed       int
	movementModel  physicsmovement.MovementModel
	movementState  movement.MovementStateEnum
	scripted       bool
	moveLeftForce  int
	moveRightForce int
}

func (m *mockActor) ID() string                                           { return m.id }
func (m *mockActor) SetID(id string)                                      { m.id = id }
func (m *mockActor) Position() image.Rectangle                            { return m.pos }
func (m *mockActor) SetPosition(x, y int)                                 { m.pos = image.Rect(x, y, x+10, y+10) }
func (m *mockActor) SetPosition16(x16, y16 int)                           { m.SetPosition(x16/16, y16/16) }
func (m *mockActor) GetPosition16() (int, int)                            { return m.pos.Min.X * 16, m.pos.Min.Y * 16 }
func (m *mockActor) GetPositionMin() (int, int)                           { return m.pos.Min.X, m.pos.Min.Y }
func (m *mockActor) GetShape() body.Shape                                 { return m }
func (m *mockActor) Width() int                                           { return m.pos.Dx() }
func (m *mockActor) Height() int                                          { return m.pos.Dy() }
func (m *mockActor) Speed() int                                           { return m.speed }
func (m *mockActor) MaxSpeed() int                                        { return m.maxSpeed }
func (m *mockActor) SetSpeed(s int) error                                 { m.speed = s; return nil }
func (m *mockActor) SetMaxSpeed(s int) error                              { m.maxSpeed = s; return nil }
func (m *mockActor) MovementModel() physicsmovement.MovementModel         { return m.movementModel }
func (m *mockActor) SetMovementModel(model physicsmovement.MovementModel) { m.movementModel = model }
func (m *mockActor) OnMoveLeft(force int)                                 { m.moveLeftForce = force }
func (m *mockActor) OnMoveRight(force int)                                { m.moveRightForce = force }
func (m *mockActor) SetMovementState(state movement.MovementStateEnum, target body.MovableCollidable, options ...movement.MovementStateOption) {
	m.movementState = state
}
func (m *mockActor) GetCharacter() *actors.Character {
	return &actors.Character{}
}

// Other interface methods (stubs)
func (m *mockActor) Image() *ebiten.Image                                 { return nil }
func (m *mockActor) ImageOptions() *ebiten.DrawImageOptions               { return nil }
func (m *mockActor) UpdateImageOptions()                                  {}
func (m *mockActor) BlockMovement()                                       {}
func (m *mockActor) UnblockMovement()                                     {}
func (m *mockActor) IsMovementBlocked() bool                              { return false }
func (m *mockActor) State() actors.ActorStateEnum                         { return 0 }
func (m *mockActor) SetState(state actors.ActorState)                     {}
func (m *mockActor) SwitchMovementState(state movement.MovementStateEnum) {}
func (m *mockActor) MovementState() movement.MovementState                { return nil }
func (m *mockActor) NewState(state actors.ActorStateEnum) (actors.ActorState, error) {
	return nil, nil
}
func (m *mockActor) Hurt(damage int)                                     {}
func (m *mockActor) Owner() interface{}                                  { return nil }
func (m *mockActor) SetOwner(interface{})                                {}
func (m *mockActor) LastOwner() interface{}                              { return nil }
func (m *mockActor) Update(space body.BodiesSpace) error                 { return nil }
func (m *mockActor) Health() int                                         { return 100 }
func (m *mockActor) MaxHealth() int                                      { return 100 }
func (m *mockActor) SetHealth(h int)                                     {}
func (m *mockActor) SetMaxHealth(h int)                                  {}
func (m *mockActor) LoseHealth(d int)                                    {}
func (m *mockActor) RestoreHealth(h int)                                 {}
func (m *mockActor) Invulnerable() bool                                  { return false }
func (m *mockActor) SetInvulnerability(v bool)                           {}
func (m *mockActor) GetTouchable() body.Touchable                        { return m }
func (m *mockActor) OnTouch(other body.Collidable)                       {}
func (m *mockActor) OnBlock(other body.Collidable)                       {}
func (m *mockActor) DrawCollisionBox(s *ebiten.Image, p image.Rectangle) {}
func (m *mockActor) CollisionPosition() []image.Rectangle                { return []image.Rectangle{m.pos} }
func (m *mockActor) CollisionShapes() []body.Collidable                  { return nil }
func (m *mockActor) IsObstructive() bool                                 { return true }
func (m *mockActor) SetIsObstructive(v bool)                             {}
func (m *mockActor) AddCollision(list ...body.Collidable)                {}
func (m *mockActor) ClearCollisions()                                    {}
func (m *mockActor) SetTouchable(t body.Touchable)                       {}
func (m *mockActor) ApplyValidPosition(d int, ax bool, sp body.BodiesSpace) (int, int, bool) {
	return 0, 0, false
}
func (m *mockActor) MoveX(d int)                                      {}
func (m *mockActor) MoveY(d int)                                      {}
func (m *mockActor) OnMoveUpLeft(d int)                               {}
func (m *mockActor) OnMoveDownLeft(d int)                             {}
func (m *mockActor) OnMoveUpRight(d int)                              {}
func (m *mockActor) OnMoveDownRight(d int)                            {}
func (m *mockActor) OnMoveUp(d int)                                   {}
func (m *mockActor) OnMoveDown(d int)                                 {}
func (m *mockActor) Velocity() (int, int)                             { return 0, 0 }
func (m *mockActor) SetVelocity(vx, vy int)                           {}
func (m *mockActor) Acceleration() (ax, ay int)                       { return 0, 0 }
func (m *mockActor) SetAcceleration(ax, ay int)                       {}
func (m *mockActor) Immobile() bool                                   { return false }
func (m *mockActor) SetImmobile(i bool)                               {}
func (m *mockActor) SetFreeze(f bool)                                 {}
func (m *mockActor) Freeze() bool                                     { return false }
func (m *mockActor) FaceDirection() animation.FacingDirectionEnum     { return 0 }
func (m *mockActor) SetFaceDirection(v animation.FacingDirectionEnum) {}
func (m *mockActor) IsIdle() bool                                     { return true }
func (m *mockActor) IsWalking() bool                                  { return false }
func (m *mockActor) IsFalling() bool                                  { return false }
func (m *mockActor) IsGoingUp() bool                                  { return false }
func (m *mockActor) CheckMovementDirectionX()                         {}
func (m *mockActor) TryJump(f int)                                    {}
func (m *mockActor) SetJumpForceMultiplier(mu float64)                {}
func (m *mockActor) JumpForceMultiplier() float64                     { return 1.0 }
func (m *mockActor) SetHorizontalInertia(i float64)                   {}
func (m *mockActor) HorizontalInertia() float64                       { return 1.0 }
func (m *mockActor) IsAlive() bool                                    { return true }
func (m *mockActor) Die()                                             {}

type mockMovementModel struct {
	isScripted bool
}

func (m *mockMovementModel) Update(b body.MovableCollidable, s body.BodiesSpace) error { return nil }
func (m *mockMovementModel) SetIsScripted(is bool)                                     { m.isScripted = is }

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
	actor := &mockActor{id: "test_actor", speed: 5, movementModel: &mockMovementModel{}}
	actor.SetPosition(0, 0)
	am.Register(actor)

	cmd := &MoveActorCommand{
		TargetID: "test_actor",
		EndX:     100,
		Speed:    10,
	}

	cmd.Init(ctx)

	if !actor.movementModel.(*mockMovementModel).isScripted {
		t.Error("expected actor to be set to scripted mode")
	}

	// Move right
	finished := cmd.Update()
	if finished {
		t.Error("command should not be finished yet")
	}
	if actor.moveRightForce != 10 {
		t.Errorf("expected move right force 10, got %d", actor.moveRightForce)
	}

	// Reach destination
	actor.SetPosition(100, 0)
	finished = cmd.Update()
	if !finished {
		t.Error("command should be finished when destination reached")
	}
	if actor.movementModel.(*mockMovementModel).isScripted {
		t.Error("expected actor scripted mode to be disabled after completion")
	}
}

func TestSetSpeedCommand(t *testing.T) {
	ctx, am := setupTestContext()
	actor := &mockActor{id: "test_actor", speed: 5, maxSpeed: 5}
	am.Register(actor)

	cmd := &SetSpeedCommand{
		TargetID: "test_actor",
		Speed:    20,
	}

	cmd.Init(ctx)

	if actor.speed != 20 || actor.maxSpeed != 20 {
		t.Errorf("expected speed 20, got speed %d, maxSpeed %d", actor.speed, actor.maxSpeed)
	}
}

func TestFollowCommands(t *testing.T) {
	ctx, am := setupTestContext()
	player := &mockActor{id: "player"}
	am.Register(player)

	npc := &mockActor{id: "npc"}
	am.Register(npc)

	// Test Follow
	followCmd := &FollowPlayerCommand{TargetID: "npc"}
	followCmd.Init(ctx)
	if npc.movementState != Follow {
		t.Errorf("expected npc to be in Follow state, got %v", npc.movementState)
	}

	// Test Stop Following
	stopCmd := &StopFollowingCommand{TargetID: "npc"}
	stopCmd.Init(ctx)
	if npc.movementState != Idle {
		t.Errorf("expected npc to be in Idle state, got %v", npc.movementState)
	}
}

func TestRemoveActorCommand(t *testing.T) {
	ctx, am := setupTestContext()
	actor := &mockActor{id: "to_remove"}
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
