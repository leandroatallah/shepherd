package events

const (
	PlayerReachedFirstPointType = "player_reached_first_point"
)

type PlayerReachedFirstPointEvent struct{}

func (e *PlayerReachedFirstPointEvent) Type() string {
	return PlayerReachedFirstPointType
}
