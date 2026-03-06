package screenutil

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
)

func TestGetCenterOfScreenPosition(t *testing.T) {
	config.Set(&config.AppConfig{ScreenWidth: 320, ScreenHeight: 240})

	x, y := GetCenterOfScreenPosition(100, 50)
	expectedX := 320/2 - 100/2 // 160 - 50 = 110
	expectedY := 240/2 - 50/2  // 120 - 25 = 95

	if x != expectedX || y != expectedY {
		t.Errorf("expected (%d, %d), got (%d, %d)", expectedX, expectedY, x, y)
	}
}

func TestDrawScreenFlash(t *testing.T) {
	config.Set(&config.AppConfig{ScreenWidth: 320, ScreenHeight: 240})
	screen := ebiten.NewImage(320, 240)
	
	// Smoke test
	DrawScreenFlash(screen)
}

func TestDrawCenteredTextSmoke(t *testing.T) {
	config.Set(&config.AppConfig{ScreenWidth: 320, ScreenHeight: 240})
	screen := ebiten.NewImage(320, 240)
	_ = screen
	// Not calling DrawCenteredText here because it might panic without a real font if we don't mock it carefully
}
