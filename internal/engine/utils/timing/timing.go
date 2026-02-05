package timing

import "time"

// TPS is the Ticks Per Second of the game engine.
// Ebitengine defaults to 60 TPS.
const TPS = 60

// FromDuration converts a time.Duration to the number of frames (ticks) based on TPS.
func FromDuration(d time.Duration) int {
	return int(d.Seconds() * TPS)
}
