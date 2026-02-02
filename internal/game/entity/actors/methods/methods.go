package gameplayermethods

import (
	gameentitytypes "github.com/leandroatallah/firefly/internal/game/entity/types"
	"github.com/leandroatallah/firefly/internal/game/events"
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
	// TODO: All actors need to freeze.
	tm.player.SetImmobile(true)
	tm.player.SetFreeze(true)

	// Trigger event to reboot scene
	if tm.player.AppContext().EventManager != nil {
		tm.player.AppContext().EventManager.Publish(&events.CharacterDiedEvent{})
	}
}
