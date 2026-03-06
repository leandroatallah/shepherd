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
	
	// Test without following but with center
	ctrl.SetFollowing(false)
	ctrl.SetCenter(50, 50)
	ctrl.Update()
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

func TestCameraNew(t *testing.T) {
	config.Set(&config.AppConfig{ScreenWidth: 320, ScreenHeight: 240})
	cam := NewCamera(10, 10)
	if cam == nil {
		t.Fatal("NewCamera returned nil")
	}
}

func TestControllerGettersSetters(t *testing.T) {
	config.Set(&config.AppConfig{ScreenWidth: 320, ScreenHeight: 240})
	ctrl := NewController(0, 0)

	if ctrl.IsFollowing() {
		t.Error("expected IsFollowing to be false by default")
	}
	ctrl.SetFollowing(true)
	if !ctrl.IsFollowing() {
		t.Error("expected IsFollowing to be true after SetFollowing(true)")
	}

	ctrl.SetCenter(100, 200)
	cx, cy := ctrl.Kamera().Center()
	if cx != 100 || cy != 200 {
		t.Errorf("expected center (100, 200), got (%f, %f)", cx, cy)
	}

	ctrl.SetPositionTopLeft(10, 20)
	cx, cy = ctrl.Kamera().Center()
	if cx != 170 || cy != 140 {
		t.Errorf("expected center (170, 140), got (%f, %f)", cx, cy)
	}
}

func TestControllerTargetAndPosition(t *testing.T) {
	config.Set(&config.AppConfig{ScreenWidth: 320, ScreenHeight: 240})
	ctrl := NewController(0, 0)

	rect := bodyphysics.NewRect(10, 20, 32, 16)
	target := bodyphysics.NewCollidableBodyFromRect(rect)
	target.SetPosition(10, 20)
	ctrl.SetFollowTarget(target)

	if ctrl.Target() != target {
		t.Error("expected Target() to return the follow target")
	}

	pos := ctrl.Position()
	if pos.Min.X != 10 || pos.Min.Y != 20 || pos.Dx() != 32 || pos.Dy() != 16 {
		t.Errorf("expected position (10, 20, 32, 16), got %+v", pos)
	}
}

func TestCamDebugSmoke(t *testing.T) {
	ctrl := NewController(0, 0)
	ctrl.CamDebug()
}
