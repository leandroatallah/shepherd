package gamescenephases

import (
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	gameenemies "github.com/leandroatallah/firefly/internal/game/entity/actors/enemies"
	gamenpcs "github.com/leandroatallah/firefly/internal/game/entity/actors/npcs"
)

type BodyCounter struct {
	sheep        int
	sheepRescued int
	wolf         int
	wolfKilled   int
}

func (b *BodyCounter) setBodyCounter(space body.BodiesSpace) {
	for _, sb := range space.Bodies() {
		switch sb.(type) {
		case *gamenpcs.Sheep:
			b.sheep++
		case *gameenemies.WolfEnemy:
			b.wolf++
		}
	}
}
