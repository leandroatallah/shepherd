package gamenpcs

import (
	"log"

	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/npcs"
	gameentitytypes "github.com/leandroatallah/firefly/internal/game/entity/types"
)

const (
	SheepNpcType npcs.NpcType = "SHEEP"
)

func InitNpcMap(ctx *app.AppContext) npcs.NpcMap[gameentitytypes.PlatformerActorEntity] {
	npcMap := map[npcs.NpcType]func(x, y int, id string) gameentitytypes.PlatformerActorEntity{
		SheepNpcType: func(x, y int, id string) gameentitytypes.PlatformerActorEntity {
			npc, err := NewSheep(x, y, id)
			if err != nil {
				log.Fatal(err)
			}
			return npc
		},
	}
	return npcMap
}
