package movement

import "github.com/leandroatallah/firefly/internal/engine/contracts/body"

// FollowMovementState defines a movement behavior where an actor follows a target
// maintaining a certain distance or with a delay mechanism based on distance.
type FollowMovementState struct {
	BaseMovementState
	startDistance int // Distance at which to start following
	stopDistance  int // Distance at which to stop following
	isMoving      bool
}

// NewFollowMovementState creates a new FollowMovementState.
func NewFollowMovementState(base BaseMovementState) *FollowMovementState {
	return &FollowMovementState{
		BaseMovementState: base,
		startDistance:     50, // Default to a noticeable delay/leash
		stopDistance:      20, // Stop reasonably close
	}
}

// WithFollowDistances sets the start and stop distances for following.
// start: Distance to target to trigger movement (the "leash" length).
// stop: Distance to target to stop movement.
func WithFollowDistances(start, stop int) MovementStateOption {
	return func(ms MovementState) {
		if s, ok := ms.(*FollowMovementState); ok {
			s.startDistance = start
			s.stopDistance = stop
		}
	}
}

func (s *FollowMovementState) Move(space body.BodiesSpace) {
	if s.actor.Immobile() {
		return
	}

	target := s.target
	if target == nil {
		return
	}

	dist := euclideanDistance(s.actor.Position().Min, target.Position().Min)

	if s.isMoving {
		if dist <= s.stopDistance {
			s.isMoving = false
		}
	} else {
		if dist >= s.startDistance {
			s.isMoving = true
		}
	}

	if s.isMoving {
		// Use direct movement calculation
		directions := calculateMovementDirections(s.actor, target, false)
		executeMovement(s.actor, directions)
	}
}