package gamescenephases

import (
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/event"
	"github.com/leandroatallah/firefly/internal/game/events"
)

func subscribeEvents(ctx *app.AppContext, scene *PhasesScene) {
	// Common events
	ctx.EventManager.Subscribe(events.CharacterDiedEventType, func(e event.Event) {
		scene.Reboot()
	})
	ctx.EventManager.Subscribe(events.PlayerJumpedType, func(e event.Event) {
		if scene.vfxManager == nil {
			return
		}
		if evt, ok := e.(*events.PlayerJumpedEvent); ok {
			yOffset := 1.0
			scene.vfxManager.SpawnJumpPuff(evt.X, evt.Y+yOffset, 1)
		}
	})
	ctx.EventManager.Subscribe(events.PlayerLandedType, func(e event.Event) {
		if scene.vfxManager == nil {
			return
		}
		if evt, ok := e.(*events.PlayerLandedEvent); ok {
			yOffset := 1.0
			scene.vfxManager.SpawnLandingPuff(evt.X, evt.Y+yOffset, 1)
		}
	})
}
