package gamescenephases

import (
	"github.com/leandroatallah/firefly/internal/engine/app"
	gameplayer "github.com/leandroatallah/firefly/internal/game/entity/actors/player"
	gameentitytypes "github.com/leandroatallah/firefly/internal/game/entity/types"
)

func createPlayer(ctx *app.AppContext, playerType gameentitytypes.PlayerType) (gameentitytypes.PlatformerActorEntity, error) {
	var f func(*app.AppContext) (gameentitytypes.PlatformerActorEntity, error)

	switch playerType {
	case gameentitytypes.ShepherdPlayerType:
		f = gameplayer.NewShepherdPlayer
	case gameentitytypes.DogPlayerType:
		f = gameplayer.NewDogPlayer
	}

	p, err := f(ctx)
	if err != nil {
		return nil, err
	}

	return p, nil
}
