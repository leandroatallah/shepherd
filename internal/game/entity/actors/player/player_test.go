package gameplayer

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	"github.com/leandroatallah/firefly/internal/engine/mocks"
	"github.com/leandroatallah/firefly/internal/engine/physics/space"
	_ "github.com/leandroatallah/firefly/internal/game/entity/actors/states"
)

func getModuleRoot() string {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			panic("could not find go.mod")
		}
		dir = parent
	}
}

func TestMain(m *testing.M) {
	err := os.Chdir(getModuleRoot())
	if err != nil {
		panic(err)
	}

	cfg := &config.AppConfig{
		ScreenWidth:  320,
		ScreenHeight: 224,
	}
	config.Set(cfg)

	os.Exit(m.Run())
}

func TestNewShepherdPlayer(t *testing.T) {
	ctx := &app.AppContext{
		ActorManager: actors.NewManager(),
		Space:        space.NewSpace(),
	}
	
	p, err := NewShepherdPlayer(ctx)
	if err != nil {
		t.Fatalf("failed to create shepherd player: %v", err)
	}

	if p == nil {
		t.Fatal("NewShepherdPlayer returned nil")
	}
}

func TestShepherdPlayer_GrabSheep(t *testing.T) {
	ctx := &app.AppContext{
		ActorManager: actors.NewManager(),
		Space:        space.NewSpace(),
	}
	
	p, err := NewShepherdPlayer(ctx)
	if err != nil {
		t.Fatalf("failed to create shepherd player: %v", err)
	}

	shepherd := p.(*ShepherdPlayer)
	
	mockSheep := &mocks.MockActor{Id: "sheep-1"}
	shepherd.GrabSheep(mockSheep)
	
	if !shepherd.IsCarryingSheep() {
		t.Error("expected IsCarryingSheep to be true after GrabSheep")
	}

	shepherd.DropSheep()
	if shepherd.IsCarryingSheep() {
		t.Error("expected IsCarryingSheep to be false after DropSheep")
	}
}
