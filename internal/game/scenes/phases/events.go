package gamescenephases

import (
	"github.com/leandroatallah/firefly/internal/engine/app"
	actorevents "github.com/leandroatallah/firefly/internal/engine/entity/actors/events"
	"github.com/leandroatallah/firefly/internal/engine/event"
	"github.com/leandroatallah/firefly/internal/game/entity/actors/events"
)

func subscribeEvents(ctx *app.AppContext, scene *PhasesScene) {
	// Common events
	ctx.EventManager.Subscribe(events.CharacterDiedEventType, func(e event.Event) {
		scene.Reboot()
	})
	ctx.EventManager.Subscribe(actorevents.ActorJumpedType, func(e event.Event) {
		if scene.vfxManager == nil {
			return
		}
		if evt, ok := e.(*actorevents.ActorJumpedEvent); ok {
			yOffset := 1.0
			scene.vfxManager.SpawnJumpPuff(evt.X, evt.Y+yOffset, 1)
		}
	})
	ctx.EventManager.Subscribe(actorevents.ActorLandedType, func(e event.Event) {
		if scene.vfxManager == nil {
			return
		}
		if evt, ok := e.(*actorevents.ActorLandedEvent); ok {
			yOffset := 1.0
			scene.vfxManager.SpawnLandingPuff(evt.X, evt.Y+yOffset, 1)
		}
	})
}
