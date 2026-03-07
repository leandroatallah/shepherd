package gamenpcs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	"github.com/leandroatallah/firefly/internal/engine/mocks"
)

type MockSheepCarrierActor struct {
	mocks.MockActor
	carrying bool
}

func (m *MockSheepCarrierActor) GrabSheep(sheep body.MovableCollidableTouchable) {
	m.carrying = true
}

func (m *MockSheepCarrierActor) IsCarryingSheep() bool {
	return m.carrying
}

func (m *MockSheepCarrierActor) DropSheep() {
	m.carrying = false
}

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

func TestNewSheep(t *testing.T) {
	ctx := &app.AppContext{
		ActorManager: actors.NewManager(),
	}
	
	s, err := NewSheep(ctx, 100, 100, "sheep-1")
	if err != nil {
		t.Fatalf("failed to create sheep: %v", err)
	}

	if s == nil {
		t.Fatal("NewSheep returned nil")
	}

	if s.ID() != "sheep-1" {
		t.Errorf("expected ID sheep-1, got %s", s.ID())
	}
}

func TestSheep_OnTouch(t *testing.T) {
	actorMgr := actors.NewManager()
	ctx := &app.AppContext{
		ActorManager: actorMgr,
	}
	
	s, err := NewSheep(ctx, 100, 100, "sheep-1")
	if err != nil {
		t.Fatalf("failed to create sheep: %v", err)
	}

	mockPlayer := &MockSheepCarrierActor{
		MockActor: mocks.MockActor{Id: "player"},
	}
	actorMgr.Register(mockPlayer)
	
	s.OnTouch(mockPlayer)
	if !mockPlayer.carrying {
		t.Error("player should be carrying sheep after OnTouch")
	}
}

func TestSheep_Hurt(t *testing.T) {
	ctx := &app.AppContext{
		ActorManager: actors.NewManager(),
	}
	
	s, err := NewSheep(ctx, 100, 100, "sheep-1")
	if err != nil {
		t.Fatalf("failed to create sheep: %v", err)
	}

	s.Hurt(1)
	// Should be in Dying state
}
