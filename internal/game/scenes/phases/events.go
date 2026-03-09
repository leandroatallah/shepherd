package gamescenephases

import (
	"log"

	"github.com/leandroatallah/firefly/internal/engine/app"
	actorevents "github.com/leandroatallah/firefly/internal/engine/entity/actors/events"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/platformer"
	"github.com/leandroatallah/firefly/internal/engine/event"
	"github.com/leandroatallah/firefly/internal/engine/sequences"
	"github.com/leandroatallah/firefly/internal/game/entity/actors/events"
	gamestates "github.com/leandroatallah/firefly/internal/game/entity/actors/states"
)

func subscribeEvents(ctx *app.AppContext, scene *PhasesScene) {
	onLandFunc := func(e event.Event) {
		if ctx.VFX == nil {
			return
		}
		if evt, ok := e.(*actorevents.ActorLandedEvent); ok {
			yOffset := 1.0
			ctx.VFX.SpawnLandingPuff(evt.X, evt.Y+yOffset, 1)
		}
	}

	em := ctx.EventManager

	// Common events
	em.Subscribe(events.CharacterDiedEventType, func(e event.Event) {
		scene.Reboot()
	})
	em.Subscribe(actorevents.ActorJumpedType, func(e event.Event) {
		if ctx.VFX == nil {
			return
		}
		if evt, ok := e.(*actorevents.ActorJumpedEvent); ok {
			yOffset := 1.0
			ctx.VFX.SpawnJumpPuff(evt.X, evt.Y+yOffset, 1)
		}
	})
	landUnsubscribe := em.Subscribe(actorevents.ActorLandedType, onLandFunc)

	// PHASE EVENTS
	// Area 1
	var fallUnsubscribe func()
	fallUnsubscribe = em.Subscribe("cutscene_shepherd_fall", func(e event.Event) {
		// Unsubscribe after first trigger (one-time event)
		if fallUnsubscribe != nil {
			fallUnsubscribe()
		}

		player, found := ctx.ActorManager.GetPlayer()
		if !found {
			return
		}
		p, ok := player.(platformer.PlatformerActorEntity)
		if !ok {
			return
		}

		// Temporarily disable original land handler
		if landUnsubscribe != nil {
			landUnsubscribe()
		}
		// Add one-time land handler for cutscene
		var cutsceneUnsubscribe func()
		cutsceneUnsubscribe = em.Subscribe(actorevents.ActorLandedType, func(e event.Event) {
			s, err := p.NewState(gamestates.Lying)
			if err != nil {
				log.Printf("Failed to create new state %v: %v", s, err)
				return
			}
			p.SetState(s)

			// Dramatic shake when falling in cutscene
			ctx.VFX.AddTrauma(scene.Camera(), 0.8)

			// Cleanup cutscene handler
			if cutsceneUnsubscribe != nil {
				cutsceneUnsubscribe()
			}

			// Restore original land handler
			landUnsubscribe = em.Subscribe(actorevents.ActorLandedType, onLandFunc)
		})
	})

	var riseUnsubscribe func()
	riseUnsubscribe = em.Subscribe("cutscene_shepherd_fall_rise", func(e event.Event) {
		// Unsubscribe after first trigger (one-time event)
		if riseUnsubscribe != nil {
			riseUnsubscribe()
		}

		player, found := ctx.ActorManager.GetPlayer()
		if !found {
			return
		}
		p, ok := player.(platformer.PlatformerActorEntity)
		if !ok {
			return
		}

		s, err := p.NewState(gamestates.Rising)
		if err != nil {
			log.Printf("Failed to create new state %v: %v", s, err)
			return
		}
		p.SetState(s)

		if scene.sequencePlayer == nil {
			scene.sequencePlayer = sequences.NewSequencePlayer(ctx)
		}
		seq, err := sequences.NewSequenceFromJSON("assets/sequences/phase-1-3-cutscene-2.json")
		if err == nil {
			scene.sequencePlayer.Play(seq)
		} else {
			log.Printf("Failed to load cutscene sequence: %v", err)
		}
	})
}
