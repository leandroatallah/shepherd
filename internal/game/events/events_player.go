package events

const (
	PlayerJumpedType = "player_jumped"
	PlayerLandedType = "player_landed"
)

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
