package gameentitytypes

import "github.com/leandroatallah/firefly/internal/engine/contracts/body"

type SheepCarrier interface {
	GrabSheep(s body.MovableCollidableTouchable)
	IsCarryingSheep() bool
	DropSheep()
}
