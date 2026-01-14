package gameenemies

import (
	"log"

	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/enemies"
	gameentitytypes "github.com/leandroatallah/firefly/internal/game/entity/types"
)

const (
	WolfEnemyType enemies.EnemyType = "WOLF"
)

func InitEnemyMap(ctx *app.AppContext) enemies.EnemyMap[gameentitytypes.PlatformerActorEntity] {
	enemyMap := map[enemies.EnemyType]func(x, y int, id string) gameentitytypes.PlatformerActorEntity{
		WolfEnemyType: func(x, y int, id string) gameentitytypes.PlatformerActorEntity {
			enemy, err := NewWolfEnemy(x, y, id)
			if err != nil {
				log.Fatal(err)
			}
			player, _ := ctx.ActorManager.GetPlayer()
			enemy.SetTarget(player)
			return enemy
		},
	}
	return enemyMap
}
