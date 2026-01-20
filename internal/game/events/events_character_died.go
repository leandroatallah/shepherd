package events

const (
	// CharacterDiedEventType is the event type for when a character dies.
	CharacterDiedEventType = "CharacterDied"
)

// CharacterDiedEvent is dispatched when a character dies.
type CharacterDiedEvent struct{}

// Type returns the event type.
func (e *CharacterDiedEvent) Type() string {
	return CharacterDiedEventType
}
