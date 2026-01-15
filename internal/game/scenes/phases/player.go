package gamescenephases

import (
	"github.com/leandroatallah/firefly/internal/engine/app"
	gameplayer "github.com/leandroatallah/firefly/internal/game/entity/actors/player"
	gameentitytypes "github.com/leandroatallah/firefly/internal/game/entity/types"
)

func createPlayer(ctx *app.AppContext) (gameentitytypes.PlatformerActorEntity, error) {
	p, err := gameplayer.NewDogPlayer(ctx)
	if err != nil {
		return nil, err
	}

	return p, nil
}
