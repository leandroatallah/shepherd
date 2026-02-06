package events

const (
	PlayerReachedFirstPointType = "player_reached_first_point"
	PlayerJumpedType            = "player_jumped"
	PlayerLandedType            = "player_landed"
)

type PlayerReachedFirstPointEvent struct{}

func (e *PlayerReachedFirstPointEvent) Type() string {
	return PlayerReachedFirstPointType
}

type PlayerJumpedEvent struct {
	X, Y float64
}

func (e *PlayerJumpedEvent) Type() string {
	return PlayerJumpedType
}

type PlayerLandedEvent struct {
	X, Y float64
}

func (e *PlayerLandedEvent) Type() string {
	return PlayerLandedType
}
