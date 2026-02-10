package gamenpcs

import (
	"log"

	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/npcs"
	gameplayer "github.com/leandroatallah/firefly/internal/game/entity/actors/player"
	gameentitytypes "github.com/leandroatallah/firefly/internal/game/entity/types"
)

const (
	SheepNpcType npcs.NpcType = "SHEEP"
	DogNpcType   npcs.NpcType = "DOG"
)

func InitNpcMap(ctx *app.AppContext) npcs.NpcMap[gameentitytypes.PlatformerActorEntity] {
	npcMap := map[npcs.NpcType]func(x, y int, id string) gameentitytypes.PlatformerActorEntity{
		SheepNpcType: func(x, y int, id string) gameentitytypes.PlatformerActorEntity {
			npc, err := NewSheep(ctx, x, y, id)
			if err != nil {
				log.Fatal(err)
			}
			return npc
		},
		DogNpcType: func(x, y int, id string) gameentitytypes.PlatformerActorEntity {
			npc, err := gameplayer.NewDogPlayer(ctx)
			if err != nil {
				log.Fatal(err)
			}
			npc.SetPosition(x, y)
			npc.SetID(id)
			return npc
		},
	}
	return npcMap
}
