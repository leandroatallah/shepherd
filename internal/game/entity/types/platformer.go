package gameentitytypes

import (
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	bodyphysics "github.com/leandroatallah/firefly/internal/engine/physics/body"
	"github.com/leandroatallah/firefly/internal/engine/render/sprites"
)

type SheepCarrier interface {
	GrabSheep(s body.MovableCollidableTouchable)
	IsCarryingSheep() bool
	DropSheep()
}

type CoinCollector interface {
	AddCoinCount(amount int)
	CoinCount() int
}

type PlatformerActorEntity interface {
	actors.ActorEntity
}

type PlatformerCharacter struct {
	actors.Character
	app.AppContextHolder

	coinCount        int
	movementBlockers int
}

func NewPlatformerCharacter(s sprites.SpriteMap, bodyRect *bodyphysics.Rect) *PlatformerCharacter {
	c := actors.NewCharacter(s, bodyRect)
	return &PlatformerCharacter{
		Character: *c,
	}
}
