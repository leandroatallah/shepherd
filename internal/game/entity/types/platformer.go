package gameentitytypes

import (
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	bodyphysics "github.com/leandroatallah/firefly/internal/engine/physics/body"
	"github.com/leandroatallah/firefly/internal/engine/render/sprites"
)

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

// Overwrite Character NewState to pass PlatformerCharacter as argument.
func (c *PlatformerCharacter) NewState(state actors.ActorStateEnum) (actors.ActorState, error) {
	s, err := actors.NewState(c, state)
	if err != nil {
		return nil, err
	}
	c.SetState(s)
	return s, nil
}
