package gamescenephases

import (
	gameplayer "github.com/leandroatallah/firefly/internal/game/entity/actors/player"
	gameentitytypes "github.com/leandroatallah/firefly/internal/game/entity/types"
)

func createPlayer() (gameentitytypes.PlatformerActorEntity, error) {
	p, err := gameplayer.NewDogPlayer()
	if err != nil {
		return nil, err
	}

	return p, nil
}
