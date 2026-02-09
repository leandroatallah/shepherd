package items

import (
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
)

type ItemType string

type ItemMap[T Item] map[ItemType]func(x, y int, id string) T

type Item interface {
	body.MovableCollidable
	body.Drawable

	Update(space body.BodiesSpace) error
	IsRemoved() bool
	SetRemoved(value bool)
}
