package events

const (
	StoryTransitionTwoType  = "transition_story_two"
	StoryTransitionFourType = "transition_story_four"
	PlayerJumpedType        = "player_jumped"
	PlayerLandedType        = "player_landed"
)

type StoryTransitionTwoEvent struct{}

func (e *StoryTransitionTwoEvent) Type() string {
	return StoryTransitionTwoType
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
