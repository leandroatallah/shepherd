package sequences

import (
	"log"

	"github.com/leandroatallah/firefly/internal/engine/app"
)

type SpawnTextCommand struct {
	TargetID string  `json:"target_id,omitempty"`
	Text     string  `json:"text"`
	Duration int     `json:"duration"`
	Type     string  `json:"type"`
	X, Y     float64
}

func (c *SpawnTextCommand) Init(appContext any) {
	ctx := appContext.(*app.AppContext)

	if c.Type == "screen" {
		ctx.VFX.SpawnFloatingText(c.Text, c.X, c.Y, c.Duration)
		return
	}

	if c.TargetID == "" {
		log.Printf("SpawnTextCommand: target_id required for overhead text")
		return
	}

	actor, found := ctx.ActorManager.Find(c.TargetID)
	if !found {
		log.Printf("SpawnTextCommand: actor not found: %s", c.TargetID)
		return
	}

	ctx.VFX.SpawnFloatingTextAbove(actor, c.Text, c.Duration)
}

func (c *SpawnTextCommand) Update() bool {
	return true
}
