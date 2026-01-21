package movement

import (
	"image"

	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
)

// SideToSideMovementState defines a movement behavior where an actor moves
// horizontally, changing direction upon detecting a ledge or a wall.
type SideToSideMovementState struct {
	BaseMovementState
	movingRight bool
}

// NewSideToSideMovementState creates a new SideToSideMovementState.
func NewSideToSideMovementState(base BaseMovementState) *SideToSideMovementState {
	return &SideToSideMovementState{
		BaseMovementState: base,
		movingRight:       true, // Start by moving right
	}
}

// Move executes the side-to-side movement logic. It checks for ledges and walls
// to reverse direction and then applies movement.
func (s *SideToSideMovementState) Move(space body.BodiesSpace) {
	if s.actor.Immobile() {
		return
	}

	if s.shouldTurn(space) {
		s.movingRight = !s.movingRight
	}

	if s.movingRight {
		s.actor.OnMoveRight(s.actor.Speed())
	} else {
		s.actor.OnMoveLeft(s.actor.Speed())
	}
}

// shouldTurn checks for conditions that should make the actor reverse direction.
// It returns true if a wall is directly in front of the actor or if there is no
// ground just ahead of it (a ledge).
func (s *SideToSideMovementState) shouldTurn(space body.BodiesSpace) bool {
	if space == nil {
		return false
	}
	actorPos := s.actor.Position()

	// 1. Ledge detection
	var groundCheckPoint image.Point
	if s.movingRight {
		// Check point is at the actor's bottom-right corner, plus one pixel down.
		groundCheckPoint = image.Point{X: actorPos.Max.X, Y: actorPos.Max.Y + 1}
	} else {
		// Check point is at the actor's bottom-left corner, minus one pixel left, plus one pixel down.
		groundCheckPoint = image.Point{X: actorPos.Min.X - 1, Y: actorPos.Max.Y + 1}
	}

	groundCheckRect := image.Rectangle{Min: groundCheckPoint, Max: groundCheckPoint.Add(image.Point{X: 1, Y: 1})}

	hasGround := false
	colliders := space.Query(groundCheckRect)
	for _, c := range colliders {
		if c.IsObstructive() && c.ID() != s.actor.ID() {
			hasGround = true
			break
		}
	}

	if !hasGround {
		return true // Turn at ledge
	}

	// 2. Wall detection
	var wallCheckRect image.Rectangle
	if s.movingRight {
		// Check a 1-pixel-wide vertical slice right in front of the actor.
		wallCheckRect = image.Rect(actorPos.Max.X, actorPos.Min.Y, actorPos.Max.X+1, actorPos.Max.Y)
	} else {
		// Check a 1-pixel-wide vertical slice right in front of the actor.
		wallCheckRect = image.Rect(actorPos.Min.X-1, actorPos.Min.Y, actorPos.Min.X, actorPos.Max.Y)
	}

	colliders = space.Query(wallCheckRect)
	for _, c := range colliders {
		if c.IsObstructive() && c.ID() != s.actor.ID() {
			return true // Turn at wall
		}
	}

	return false
}
