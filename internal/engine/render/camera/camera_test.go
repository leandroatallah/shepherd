package camera

import (
	"image"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	bodyphysics "github.com/leandroatallah/firefly/internal/engine/physics/body"
)

func TestControllerFollowTargetAndBounds(t *testing.T) {
	config.Set(&config.AppConfig{ScreenWidth: 320, ScreenHeight: 240})

	ctrl := NewController(0, 0)

	rect := bodyphysics.NewRect(0, 0, 32, 16)
	target := bodyphysics.NewCollidableBodyFromRect(rect)
	target.SetPosition(300, 200) // near bottom-right

	ctrl.SetFollowTarget(target)
	ctrl.SetFollowing(true)

	// Set bounds smaller than target position to force clamping inside camera view
	bounds := image.Rect(0, 0, 320, 240)
	ctrl.SetBounds(&bounds)

	ctrl.Update()

	if ctrl.Bounds() == nil || *ctrl.Bounds() != bounds {
		t.Fatalf("bounds were not set correctly")
	}

	// Smoke test draw flow (no panic)
	dst := ebiten.NewImage(320, 240)
	opts := &ebiten.DrawImageOptions{}
	src := ebiten.NewImage(1, 1)
	ctrl.Draw(src, opts, dst)
}

type fakeCollidable struct {
	*bodyphysics.CollidableBody
	obstructive bool
}

func (f *fakeCollidable) IsObstructive() bool { return f.obstructive }

func TestDrawCollisionBoxColorsDontPanic(t *testing.T) {
	ctrl := NewController(0, 0)

	base := bodyphysics.NewCollidableBodyFromRect(bodyphysics.NewRect(10, 10, 10, 10))
	base.SetID("b")
	x, y := base.GetPositionMin()
	col := bodyphysics.NewCollidableBodyFromRect(base.GetShape())
	col.SetPosition(x, y)
	base.AddCollision(col)
	f := &fakeCollidable{CollidableBody: base, obstructive: true}

	screen := ebiten.NewImage(100, 100)
	ctrl.DrawCollisionBox(screen, f)

	f.obstructive = false
	ctrl.DrawCollisionBox(screen, f)
}
