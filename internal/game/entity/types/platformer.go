package gameentitytypes

import (
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/contracts/context"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	bodyphysics "github.com/leandroatallah/firefly/internal/engine/physics/body"
	"github.com/leandroatallah/firefly/internal/engine/render/sprites"
)

type AlivePlayer interface {
	Hurt(damage int)
}

type PlatformerActorEntity interface {
	actors.ActorEntity
	context.ContextProvider

	OnDie()
}

type PlatformerCharacter struct {
	*actors.Character
	app.AppContextHolder

	coinCount        int
	movementBlockers int
}

func NewPlatformerCharacter(s sprites.SpriteMap, bodyRect *bodyphysics.Rect) *PlatformerCharacter {
	c := actors.NewCharacter(s, bodyRect)
	pf := &PlatformerCharacter{
		Character: c,
	}
	c.SetOwner(pf)
	return pf
}
