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

	// Error cases: nil sprite or image
	if se.AnimatedSpriteImage(nil, rect, 0, 1) != nil {
		t.Error("expected nil for nil sprite")
	}
	if se.AnimatedSpriteImage(&Sprite{Image: nil}, rect, 0, 1) != nil {
		t.Error("expected nil for nil image")
	}

	// No width case
	_ = se.AnimatedSpriteImage(sprite, image.Rect(0, 0, 0, 16), 0, 1)
}

func TestSpriteEntityGetters(t *testing.T) {
	sprite := &Sprite{Image: ebiten.NewImage(32, 32), Loop: true}
	sprites := SpriteMap{"idle": sprite}
	se := NewSpriteEntity(sprites)

	if se.GetFirstSprite() != sprite {
		t.Error("expected GetFirstSprite to return the only sprite")
	}

	if se.GetSpriteByState("idle") != sprite {
		t.Error("expected GetSpriteByState to return the correct sprite")
	}

	if len(se.Sprites()) != 1 {
		t.Error("expected Sprites() to return the sprite map")
	}

	seEmpty := NewSpriteEntity(nil)
	if seEmpty.GetFirstSprite() != nil {
		t.Error("expected GetFirstSprite to return nil for empty entity")
	}
}

func TestSpriteAssets(t *testing.T) {
	var sa SpriteAssets
	sa = sa.AddSprite("idle", "path/to/idle.png", true)
	if len(sa) != 1 {
		t.Fatalf("expected 1 sprite in assets, got %d", len(sa))
	}
	if sa["idle"].Path != "path/to/idle.png" || !sa["idle"].Loop {
		t.Error("sprite assets not set correctly")
	}
}

func TestLoadSpritesError(t *testing.T) {
	sa := SpriteAssets{}.AddSprite("idle", "non_existent.png", true)
	_, err := LoadSprites(sa)
	if err == nil {
		t.Error("expected error loading non-existent sprite")
	}
}
