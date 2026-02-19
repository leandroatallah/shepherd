package gameplayermethods

import (
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	"github.com/leandroatallah/firefly/internal/game/entity/actors/events"
	gameentitytypes "github.com/leandroatallah/firefly/internal/game/entity/types"
)

type PlayerDeathBehavior struct {
	player gameentitytypes.PlatformerActorEntity
}

func NewPlayerDeathBehavior(p gameentitytypes.PlatformerActorEntity) *PlayerDeathBehavior {
	tm := &PlayerDeathBehavior{
		player: p,
	}
	return tm
}

func (tm *PlayerDeathBehavior) OnDie() {
	tm.player.SetHealth(0)
	// Freeze all actors
	if tm.player.AppContext().ActorManager != nil {
		tm.player.AppContext().ActorManager.ForEach(func(actor actors.ActorEntity) {
			actor.SetImmobile(true)
			actor.SetFreeze(true)
		})
	}

	// Trigger event to reboot scene
	if tm.player.AppContext().EventManager != nil {
		tm.player.AppContext().EventManager.Publish(&events.CharacterDiedEvent{})
	}
}
