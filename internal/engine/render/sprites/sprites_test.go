package sprites

import (
	"image"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestAnimatedSpriteImageFrameSelection(t *testing.T) {
	// Create a 96x16 image with 3 frames of 32x16 each
	img := ebiten.NewImage(96, 16)
	sprite := &Sprite{Image: img, Loop: true}

	se := NewSpriteEntity(SpriteMap{"idle": sprite})

	rect := image.Rect(0, 0, 32, 16)

	// count/frameRate -> frame index
	frame := se.AnimatedSpriteImage(sprite, rect, 0, 1)
	if frame.Bounds().Dx() != 32 || frame.Bounds().Dy() != 16 {
		t.Fatalf("unexpected frame size: %+v", frame.Bounds())
	}

	frame2 := se.AnimatedSpriteImage(sprite, rect, 1, 1) // next frame
	if frame2.Bounds().Dx() != 32 || frame2.Bounds().Dy() != 16 {
		t.Fatalf("unexpected frame2 size: %+v", frame2.Bounds())
	}

	se.SetFrameRate(2)
	if se.FrameRate() != 2 {
		t.Fatalf("FrameRate not set")
	}

	// Non-looping clamps to last frame
	spriteNL := &Sprite{Image: img, Loop: false}
	_ = se.AnimatedSpriteImage(spriteNL, rect, 999, 1) // should not panic, clamps internally
}
