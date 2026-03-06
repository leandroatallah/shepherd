package movement

import (
	"image"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/contracts/animation"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
)

type mockActor struct {
	pos            image.Rectangle
	speed          int
	moveLeftForce  int
	moveRightForce int
	immobile       bool
	obstructive    bool
}

func (m *mockActor) ID() string { return "mock" }
func (m *mockActor) SetID(id string) {}
func (m *mockActor) Position() image.Rectangle { return m.pos }
func (m *mockActor) SetPosition(x, y int) { m.pos = image.Rect(x, y, x+10, y+10) }
func (m *mockActor) SetPosition16(x16, y16 int) {}
func (m *mockActor) GetPosition16() (int, int) { return m.pos.Min.X * 16, m.pos.Min.Y * 16 }
func (m *mockActor) GetPositionMin() (int, int) { return m.pos.Min.X, m.pos.Min.Y }
func (m *mockActor) GetShape() body.Shape { return nil }
func (m *mockActor) Owner() interface{} { return nil }
func (m *mockActor) SetOwner(o interface{}) {}
func (m *mockActor) LastOwner() interface{} { return nil }

func (m *mockActor) MoveX(d int) { m.moveRightForce = d }
func (m *mockActor) MoveY(d int) {}
func (m *mockActor) OnMoveLeft(d int) { m.moveLeftForce = d }
func (m *mockActor) OnMoveUpLeft(d int) { m.moveLeftForce = d }
func (m *mockActor) OnMoveDownLeft(d int) { m.moveLeftForce = d }
func (m *mockActor) OnMoveRight(d int) { m.moveRightForce = d }
func (m *mockActor) OnMoveUpRight(d int) { m.moveRightForce = d }
func (m *mockActor) OnMoveDownRight(d int) { m.moveRightForce = d }
func (m *mockActor) OnMoveUp(d int) { m.moveRightForce = d }
func (m *mockActor) OnMoveDown(d int) { m.moveRightForce = d }

func (m *mockActor) Velocity() (int, int) { return 0, 0 }
func (m *mockActor) SetVelocity(vx, vy int) {}
func (m *mockActor) Acceleration() (ax, ay int) { return 0, 0 }
func (m *mockActor) SetAcceleration(ax, ay int) {}
func (m *mockActor) SetSpeed(s int) error { m.speed = s; return nil }
func (m *mockActor) SetMaxSpeed(s int) error { return nil }
func (m *mockActor) Speed() int { return m.speed }
func (m *mockActor) MaxSpeed() int { return 0 }
func (m *mockActor) Immobile() bool { return m.immobile }
func (m *mockActor) SetImmobile(i bool) { m.immobile = i }
func (m *mockActor) SetFreeze(f bool) {}
func (m *mockActor) Freeze() bool { return false }
func (m *mockActor) FaceDirection() animation.FacingDirectionEnum { return 0 }
func (m *mockActor) SetFaceDirection(v animation.FacingDirectionEnum) {}
func (m *mockActor) IsIdle() bool { return true }
func (m *mockActor) IsWalking() bool { return false }
func (m *mockActor) IsFalling() bool { return false }
func (m *mockActor) IsGoingUp() bool { return false }
func (m *mockActor) CheckMovementDirectionX() {}
func (m *mockActor) TryJump(f int) { m.moveRightForce = f }
func (m *mockActor) SetJumpForceMultiplier(mu float64) {}
func (m *mockActor) JumpForceMultiplier() float64 { return 1.0 }
func (m *mockActor) SetHorizontalInertia(i float64) {}
func (m *mockActor) HorizontalInertia() float64 { return 1.0 }

// Collidable methods
func (m *mockActor) OnTouch(other body.Collidable) {}
func (m *mockActor) OnBlock(other body.Collidable) {}
func (m *mockActor) GetTouchable() body.Touchable { return nil }
func (m *mockActor) DrawCollisionBox(screen *ebiten.Image, position image.Rectangle) {}
func (m *mockActor) CollisionPosition() []image.Rectangle { return nil }
func (m *mockActor) CollisionShapes() []body.Collidable { return nil }
func (m *mockActor) IsObstructive() bool { return m.obstructive }
func (m *mockActor) SetIsObstructive(value bool) { m.obstructive = value }
func (m *mockActor) AddCollision(list ...body.Collidable) {}
func (m *mockActor) ClearCollisions() {}
func (m *mockActor) SetTouchable(t body.Touchable) {}
func (m *mockActor) ApplyValidPosition(d int, x bool, s body.BodiesSpace) (int, int, bool) {
	return m.pos.Min.X, m.pos.Min.Y, false
}
