package gameobstacles

import (
	"testing"

	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	bodyphysics "github.com/leandroatallah/firefly/internal/engine/physics/body"
	gamesetup "github.com/leandroatallah/firefly/internal/game/app"
)

func TestNewWalls(t *testing.T) {
	// Initialize config
	cfg := gamesetup.NewConfig()
	config.Set(cfg)

	tests := []struct {
		name string
		fn   func() *bodyphysics.ObstacleRect
		id   string
	}{
		{"WallTop", NewWallTop, "WALL-TOP"},
		{"WallLeft", NewWallLeft, "WALL-LEFT"},
		{"WallRight", NewWallRight, "WALL-RIGHT"},
		{"WallDown", NewWallDown, "WALL-DOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := tt.fn()
			if o.ID() != tt.id {
				t.Errorf("expected ID %s, got %s", tt.id, o.ID())
			}
		})
	}
}

func TestInitObstacleMap(t *testing.T) {
	ctx := &app.AppContext{}
	m := InitObstacleMap(ctx)
	if len(m) != 4 {
		t.Errorf("expected 4 obstacles, got %d", len(m))
	}

	for k, v := range m {
		o := v()
		if o == nil {
			t.Errorf("obstacle %v returned nil", k)
		}
	}
}
