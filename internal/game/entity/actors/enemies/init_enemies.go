package gameenemies

import (
	"log"

	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/enemies"
	gameentitytypes "github.com/leandroatallah/firefly/internal/game/entity/types"
)

const (
	WolfEnemyType  enemies.EnemyType = "WOLF"
	BatEnemyType   enemies.EnemyType = "BAT"
	SwarmEnemyType enemies.EnemyType = "SWARM"
)

func InitEnemyMap(ctx *app.AppContext) enemies.EnemyMap[gameentitytypes.PlatformerActorEntity] {
	enemyMap := map[enemies.EnemyType]func(x, y int, id string) gameentitytypes.PlatformerActorEntity{
		WolfEnemyType: func(x, y int, id string) gameentitytypes.PlatformerActorEntity {
			enemy, err := NewWolfEnemy(ctx, x, y, id)
			if err != nil {
				log.Fatal(err)
			}
			player, _ := ctx.ActorManager.GetPlayer()
			enemy.SetTarget(player)
			return enemy
		},
		BatEnemyType: func(x, y int, id string) gameentitytypes.PlatformerActorEntity {
			enemy, err := NewBatEnemy(ctx, x, y, id)
			if err != nil {
				log.Fatal(err)
			}
			return enemy
		},
		SwarmEnemyType: func(x, y int, id string) gameentitytypes.PlatformerActorEntity {
			enemy, err := NewSwarmEnemy(ctx, x, y, id)
			if err != nil {
				log.Fatal(err)
			}
			return enemy
		},
	}
	return enemyMap
}
