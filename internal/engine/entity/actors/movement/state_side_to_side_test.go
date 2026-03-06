package movement

import (
	"image"
	"testing"

	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/contracts/tilemaplayer"
)

type mockSpace struct {
	bodies []body.Collidable
}

func (m *mockSpace) Query(rect image.Rectangle) []body.Collidable {
	var result []body.Collidable
	for _, b := range m.bodies {
		if b.Position().Overlaps(rect) {
			result = append(result, b)
		}
	}
	return result
}

func (m *mockSpace) AddBody(b body.Collidable) {}
func (m *mockSpace) RemoveBody(b body.Collidable) {}
func (m *mockSpace) Bodies() []body.Collidable { return m.bodies }
func (m *mockSpace) ResolveCollisions(b body.Collidable) (bool, bool) { return false, false }
func (m *mockSpace) SetTilemapDimensionsProvider(p tilemaplayer.TilemapDimensionsProvider) {}
func (m *mockSpace) GetTilemapDimensionsProvider() tilemaplayer.TilemapDimensionsProvider { return nil }
func (m *mockSpace) QueueForRemoval(body body.Collidable) {}
func (m *mockSpace) ProcessRemovals() {}
func (m *mockSpace) Clear() {}
func (m *mockSpace) Find(id string) body.Collidable { return nil }

func TestSideToSideMovementState_WallDetection(t *testing.T) {
	actor := &mockActor{speed: 5}
	actor.SetPosition(0, 0) // 10x10 size by default in mock
	
	wall := &mockActor{}
	wall.pos = image.Rect(10, 0, 20, 10)
	
	space := &mockSpace{bodies: []body.Collidable{wall}}

	base := NewBaseMovementState(SideToSide, actor, nil)
	state := NewSideToSideMovementState(base)

	state.Move(space)

	if actor.moveLeftForce != 5 {
		t.Errorf("expected moveLeftForce 5 after hitting wall, got %d", actor.moveLeftForce)
	}
}
