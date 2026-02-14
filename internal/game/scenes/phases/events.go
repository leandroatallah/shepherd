package gamescenephases

import (
	"time"

	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/event"
	"github.com/leandroatallah/firefly/internal/engine/scene/transition"
	"github.com/leandroatallah/firefly/internal/game/events"
	scenestypes "github.com/leandroatallah/firefly/internal/game/scenes/types"
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

	// Story Transition 1
	// TODO: Check if should use event schedule
	ctx.EventManager.Subscribe(events.StoryTransitionTwoType, func(e event.Event) {
		scene.Schedule(4*time.Second+500*time.Millisecond, func() {
			// Advance to next phase (Phase 3) before navigating
			scene.AppContext().PhaseManager.AdvanceToNextPhase()
			scene.AppContext().SceneManager.NavigateTo(scenestypes.SceneStory, transition.NewFader(), true)
		})
	})
	ctx.EventManager.Subscribe(events.StoryTransitionFourType, func(e event.Event) {
		scene.Schedule(3*time.Second, func() {
			// Advance to next phase (Phase 3) before navigating
			scene.AppContext().PhaseManager.AdvanceToNextPhase()
			scene.AppContext().SceneManager.NavigateTo(scenestypes.SceneStory, transition.NewFader(), true)
		})
	})
}
