package gamenpcs

import (
	"log"

	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/npcs"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/platformer"
	gameplayer "github.com/leandroatallah/firefly/internal/game/entity/actors/player"
)

const (
	SheepNpcType npcs.NpcType = "SHEEP"
	DogNpcType   npcs.NpcType = "DOG"
)

func InitNpcMap(ctx *app.AppContext) npcs.NpcMap[platformer.PlatformerActorEntity] {
	npcMap := map[npcs.NpcType]func(x, y int, id string) platformer.PlatformerActorEntity{
		SheepNpcType: func(x, y int, id string) platformer.PlatformerActorEntity {
			npc, err := NewSheep(ctx, x, y, id)
			if err != nil {
				log.Fatal(err)
			}
			return npc
		},
		DogNpcType: func(x, y int, id string) platformer.PlatformerActorEntity {
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
