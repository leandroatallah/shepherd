package gameentitytypes

import (
	"image"

	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/contracts/context"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	bodyphysics "github.com/leandroatallah/firefly/internal/engine/physics/body"
	"github.com/leandroatallah/firefly/internal/engine/physics/skill"
	"github.com/leandroatallah/firefly/internal/engine/render/sprites"
	"github.com/leandroatallah/firefly/internal/game/events"
)

type AlivePlayer interface {
	Hurt(damage int)
}

type PlatformerActorEntity interface {
	actors.ActorEntity
	context.ContextProvider

	OnDie()
	SetOnJump(func(image.Point))
}

type PlatformerCharacter struct {
	*actors.Character
	app.AppContextHolder

	coinCount        int
	movementBlockers int

	OnJump func(image.Point)
	OnLand func(image.Point)
}

func (p *PlatformerCharacter) SetOnJump(f func(image.Point)) {
	p.OnJump = f
}

func (p *PlatformerCharacter) SetOnLand(f func(image.Point)) {
	p.OnLand = f
}

func (p *PlatformerCharacter) AddSkill(s skill.Skill) {
	p.Character.AddSkill(s)

	if js, ok := s.(*skill.JumpSkill); ok {
		js.OnJump = func(b body.MovableCollidable) {
			if p.OnJump != nil {
				rect := b.Position()
				// Bottom center
				pos := image.Point{X: rect.Min.X + rect.Dx()/2, Y: rect.Max.Y}
				p.OnJump(pos)
			}
		}
	}
}

func NewPlatformerCharacter(s sprites.SpriteMap, bodyRect *bodyphysics.Rect) *PlatformerCharacter {
	c := actors.NewCharacter(s, bodyRect)
	pf := &PlatformerCharacter{
		Character: c,
	}

	pf.Character.OnStateChange = func(oldState, newState actors.ActorStateEnum) {
		isLanding := (oldState == actors.Falling && newState == actors.Landing)

		if !isLanding {
			carryFalling, ok1 := actors.GetStateEnum("carry_falling")
			carryLanding, ok2 := actors.GetStateEnum("carry_landing")
			if ok1 && ok2 && oldState == carryFalling && newState == carryLanding {
				isLanding = true
			}
		}

		if isLanding {
			if pf.OnLand != nil {
				rect := pf.Position()
				// Bottom center
				pos := image.Point{X: rect.Min.X + rect.Dx()/2, Y: rect.Max.Y}
				pf.OnLand(pos)
			}
		}
	}

	pf.SetOnJump(func(pos image.Point) {
		if pf.AppContext() != nil {
			yOffset := 1.0
			pf.AppContext().EventManager.Publish(&events.PlayerJumpedEvent{
				X: float64(pos.X),
				Y: float64(pos.Y) + yOffset,
			})
		}
	})
	pf.SetOnLand(func(pos image.Point) {
		if pf.AppContext() != nil {
			yOffset := 1.0
			pf.AppContext().EventManager.Publish(&events.PlayerLandedEvent{
				X: float64(pos.X),
				Y: float64(pos.Y) + yOffset,
			})
		}
	})
	c.SetOwner(pf)
	return pf
}
